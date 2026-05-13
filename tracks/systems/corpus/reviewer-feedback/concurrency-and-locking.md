# Concurrency and locking

Where the locks go, where they don't, and the recurring race patterns.

## Don't hold the lock during IO

> **rosstimothy** (MarkDHarris, MrChristianL): "Does the lock need to
> be held while doing IO to uphold the guarantees mentioned in this
> comment?"

> **fspmarshall** (mcampo84): "Sending on the streamer while the lock
> is held may produce unacceptable contention."

Expected answer: copy what you need under the lock, release it, then
do IO with the copy. The classic pattern: lock, snapshot state, copy
or clone any byte slice, unlock, then write to the network.

## RWMutex vs Mutex

> **espadolini** (chintamanil): "No need for a rwmutex if there's no
> need for multiple readers to block while holding the read lock and
> there's no contention as proven by a profile"

> **chintamanil** (response): "Agreed. Changed to sync.Mutex"

Expected answer: start with `sync.Mutex`. Only reach for `RWMutex`
after profiling shows read contention. Readers under load are not
common in this challenge; the streaming path holds the lock briefly.

## Double-close / write-after-close

> **dboslee** (benmoss): "there are multiple places where this
> channel is closed. is it possible for a write -> broadcast to come
> in after markDone? I think this would cause a panic."

> **tigrato** (mcampo84): "being this publicly exposed, a user can
> call it multiple times which will cause panics because multiple
> channel close operations"

Expected answer: closing is idempotent. Wrap a channel close in a
`sync.Once`, or check a closed-flag under the lock. Writes after
close return an error or are dropped explicitly with a logged
warning, not crash.

## Mark done as part of close

> **MarkDHarris** (corpus, paraphrased from MrChristianL pattern):
> the writer's `Close()` is the operation that flushes pending bytes
> and marks the stream done. A separate `SetDone()` followed by a
> flush has a race window.

Expected answer: one method for "no more writes." It atomically
transitions the store to closed, increments the version counter, and
broadcasts the condition variable. Readers see the version bump and
return `io.EOF` on the next read.

## Lock-free version counter pattern

> **sclevine** (benmoss): "How will you avoid race conditions like
> this? Writer writes and signals / Reader reads to EOF / Writer
> writes and signals / Reader waits / Writer does not write again for
> 10 minutes. Result: reader has output delayed for 10 minutes"

> **benmoss** (response): "I think I can keep a counter on the reader
> and the writer so that readers can check whether writes occurred
> while they were reading instead of going back to waiting on the
> condition."

Expected answer: a monotonic version counter on the store, snapshotted
under the lock when the reader observes the data, and compared on
re-entry before the reader sleeps on the condition variable. The race
above happens because the reader reads-then-decides-to-wait, with the
write landing in between.

## Goroutine leak on disconnect

> **nklaassen** (Zephan92): "what will happen if the client
> disconnects while blocked on the condition variable here? I think
> it could be sort of a goroutine leak"

> **espadolini** (chintamanil): "This will block until there's new
> output (or the job terminates) even if the client disconnects, so a
> long running job with no output will accumulate goroutines all
> blocked on the waitgroup."

Expected answer: the reader's wait takes a context. `context.AfterFunc`
fires the condition variable broadcast on cancel. The reader rechecks
context after waking. goleak in `TestMain` catches regressions.

## Stop must signal readers too

> **chintamanil** (response): "Yes. I am going to Wake blocked
> waiters when context cancels with `go func()` -> `buf.Broadcast()`
> and then call `buf.WaitForChangeCtx(ctx, offset)` method that
> checks context after each wake"

Expected answer: cancellation triggers a broadcast. Readers waking
to a cancelled context return `ctx.Err()`. No bare sleeps in the
wakeup path.

## Don't store contexts

> **nklaassen** (Zephan92): "avoid storing a context in a struct and
> prefer testing the external API"

Expected answer: context is per-call. Job state does not include a
field of type `context.Context`. Internal goroutines get a context
passed to them at start.

## Wait once

> **tigrato** (kkloberdanz): "l.cmd.Wait() can only be called once.
> Should we document that behaviour?"

Expected answer: one goroutine, started at job creation, owns
`cmd.Wait()`. It captures the exit code into the shared status. Stop
does not call Wait. Status does not call Wait.

## Cleanup goroutine references

> **rosstimothy** (GevorgGal): "Should this goroutine be added to the
> wait group created below?"

Expected answer: every goroutine the library spawns has a defined
exit condition and is accounted for in shutdown. WaitGroups,
sync.Once, or a single supervising goroutine are all acceptable
patterns.

## Inconsistent snapshots across methods

> **rosstimothy** (MrChristianL paraphrase via the Snapshot() answer):
> separate getters for status and exit code can return inconsistent
> state if the job transitions between calls.

Expected answer: one `Snapshot()` method, one lock acquisition,
returns all related fields together.

## Process group is not the lock

> **rosstimothy** (kkloberdanz): "By using cgroups.kill, we can ensure
> that all processes that are apart of that cgroup are killed with
> SIGKILL."

Expected answer: cgroup membership is the source of truth for "what
processes belong to this job." Don't try to reimplement it with
process groups or PID tracking.
