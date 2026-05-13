# Output streaming

The single most contested area of the challenge. Reviewers will not
sign off until they understand exactly how new bytes are signaled, how
multiple readers can stream from the start without duplicating memory,
and what happens when a reader is slow or disconnects.

## How does a reader discover new bytes?

> **tigrato** (joshuarubin): "Can you please indicate more details on
> how the library will inform the reader that the stream ended and how
> it will also inform that new data was written to stderr or stdout?"

> **rosstimothy** (joshuarubin): "How will Read be unblocked when more
> output is available?"

> **sclevine** (benmoss): "How will you avoid race conditions like
> this? Writer writes and signals / Reader reads to EOF / Writer writes
> and signals / Reader waits / Writer does not write again for 10
> minutes. Result: reader has output delayed for 10 minutes."

Expected answer: a per-store version counter incremented on every
write, a `sync.Cond` broadcast, and a reader that compares its observed
version against the store's version under the lock. The naive
condition variable without a counter has the lost-wakeup bug above.

## One backing store, many readers

> **tigrato** (joshuarubin): "Why do we need to keep a full copy of
> the job output per reader? If a job outputs 1 GB and we have 5
> readers, we'll take 5 GB of RAM. Is there a way without this
> bottleneck?"

> **Tener** (GevorgGal): "My only concern is efficiency of streaming:
> if we have a large buffer already accumulated, each reader will
> immediately make a full copy of it."

> **rosstimothy** (RichyHBM): "I think you could reduce the complexity
> a bit if this was inverted. Instead of registering an io.Writer the
> API provided a means for callers to get a unique io.Reader that
> tracked their progress through the buffer."

> **nklaassen** (sabernabil12): "output buffered in channels probably
> doesn't work for multiple clients streaming from the start"

Expected answer: one shared backing store (file-backed is the simplest
correct design) and per-reader file descriptors / offsets. Each reader
seeks independently.

## Slow readers

> **zmb3** (rohitsakala): "IIUC, you're proposing an io.Pipe per
> streamer. The thing about a pipe is it synchronizes reads and
> writes. This means if your reader isn't reading fast enough that it
> can block the writer from being able to write. How will you ensure
> that one slow reader doesn't affect the UX of other clients?"

> **zmb3** (razzam21): "I don't expect the buffered channel to change
> much here as it doesn't solve the slow-reader problem, it merely
> delays it."

> **Joerger** (rohitsakala): "If the channel fills up, is the
> additional content discarded? Is there a way to ensure the slow
> reader eventually gets the full stream of content without discarded
> chunks _and_ without blocking other readers?"

> **fspmarshall** (mcampo84): "Sending on the streamer while the lock
> is held may produce unacceptable contention."

Expected answer: a file-backed log means slow readers do not block
fast ones; the OS handles the slow-reader problem via the page cache.
Each reader has its own fd and seeks at its own pace. Per-reader
in-memory queues with a backpressure policy are an acceptable
alternative but must be defended.

## Final-write vs done ordering

> **rosstimothy** (joshuarubin): "Is it possible for the final
> stdout/stderr write to still be in flight when the jobDone channel
> is closed? If so, wouldn't that cause data loss?"

> **tigrato** (mcampo84): "Are we sure that set done won't occur
> before the log output finishes pushing all log lines to readers?"

> **nklaassen** (benmoss): "I think the order of these should be
> switched, no? Or else you could still miss a final write + close
> that happened after r.file.Read returned and before getNotifyState
> got the lock"

Expected answer: `cmd.Wait()` returns only after stdout/stderr copy
goroutines exit, so wait first and mark done second. Even better, mark
done as part of the same operation that flushes the writer, not as a
separate `SetDone()` call. (This is the bug that got one rejected
submission flagged: a writer's `Write` and a separate `Close` race in
parallel with cleanup.)

## Stream teardown on client disconnect

> **nklaassen** (Zephan92): "what will happen if the client
> disconnects while blocked on the condition variable here? I think it
> could be sort of a goroutine leak if many clients stream output of a
> job that runs indefinitely without writing any more output"

> **nklaassen** (benmoss): "I think this api will limit your ability
> to terminate the stream or clean up resources early if/when the
> client disconnects before the logs have been read to EOF"

> **espadolini** (chintamanil): "This will block until there's new
> output (or the job terminates) even if the client disconnects, so a
> long running job with no output will accumulate goroutines all
> blocked on the waitgroup."

> **tigrato** (chintamanil): "If the client disconnected, this will
> hang until new data is written - no guarantee it will ever happen"

> **tigrato** (kkloberdanz): "if the client disconnected, shouldn't we
> send io.EOF? from the client perspective, he sent a 'disconnect, I
> am no longer interested'"

Expected answer: the read path takes a `context.Context`. Cancellation
wakes the condition variable. `context.AfterFunc` is the idiomatic
trigger.

## Output type: bytes, not strings

> **tigrato** (joshuarubin): "Do we have guarantees that the output is
> a utf8 string? Can we make it generic with a slice of bytes?"

> **espadolini** (Zephan92): "I don't think that we can assume that
> the job will output conveniently small lines of text."

> **rosstimothy** (joshuarubin context): "The output of an arbitrary
> Linux process is not guaranteed to be a string terminated by a
> newline character."

Expected answer: protobuf field is `bytes data`, not `string`. The
reader does not assume newline boundaries.

## Per-read cap

> **eriktate** (GevorgGal): "What if len(data) is larger than 32KB?"

Expected answer: cap per-read response at a small, fixed size
(16-32KB). Avoids huge allocations for late-joining readers that have
megabytes of accumulated output. The caller iterates.

## Read API shape

> **tigrato** (chintamanil): "Can we satisfy io.ReadCloser so io.Copy
> works?"

> **nklaassen** (MrChristianL): "suggestion: `func([]byte) error`
> looks almost like `io.Writer.Write`, I'd consider either having
> this function accept an io.Writer or return an io.Reader to make it
> more composable"

Expected answer: implement standard interfaces. The library's
streaming API should be an `io.ReadCloser` or close cousin.
Callback-based output is a footgun and reviewers cut it.

## Ownership of returned slices

> **rosstimothy** (RichyHBM): "If this doesn't return a copy of the
> bytes then it's a bit unclear what the ownership semantics of the
> returned byte slice are. If the mutex is meant to protect the
> buffer, then doesn't this currently bypass that?"

> **timothyb89** (razzam21): "This is racy as appends will mutate
> memory behind the underlying pointer."

Expected answer: if your reader returns a slice, document that it is
either a copy the caller can hold, or a view valid until the next
call. State which.
