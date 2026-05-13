# Process containment: cgroup, kill, child processes

How the child gets into the cgroup, how the cgroup gets cleaned up,
and what happens to descendants.

## Cgroup membership race

> **espadolini** (Zephan92): "It's very likely that even very normal
> things such as shell scripts will be affected by this race, so it's
> not acceptable to just make a note of it. What ways do you envision
> to have the process be contained in the cgroup right as it's
> starting?"

> **tigrato** (kkloberdanz): "What happens if the process forked
> itself before receiving the sigkill? Is there a way of ensuring all
> processes will be terminated?"

Expected answer: `CLONE_INTO_CGROUP` via Go's
`SysProcAttr.UseCgroupFD + CgroupFD`. The child is placed in the
cgroup before its first instruction. Available in Linux 5.7+.

The naive fork-then-write-to-cgroup.procs has a window where the child
can spawn descendants before the cgroup-write completes. Reviewers
will probe for this.

## cgroup.kill

> **kkloberdanz** (response): "By using cgroups.kill, we can ensure
> that all processes that are apart of that cgroup are killed with
> SIGKILL."

> **nklaassen** (rohitsakala): "sound like you will first send SIGKILL
> to the initial process (via context cancellation), then SIGTERM to
> the process group, the write 1 to the cgroup.kill file to send
> SIGKILL to everything. Consider simplifying"

> **rohitsakala** (response): "Yes you are right. I just added it for
> an additional measure, but as you said it is completely unnecessary.
> In fact, writing 1 to the cgroup.kill is sufficient."

Expected answer: write `"1"` to `cgroup.kill`. One operation, atomic,
covers all descendants. Linux 5.14+. State the kernel floor.

## Process group signaling is not enough

> **smallinsky** (kkloberdanz): "Does this approach guarantees that
> all jobs children process will be killed?"

> **rosstimothy** (GevorgGal): "Is signaling the process group safe?
> Does signaling the process group provide any guarantees that child
> processes are terminated?"

> **GevorgGal** (response): "Process groups are best-effort, not
> guaranteed."

> **espadolini** (chintamanil): "You mention that cgroups are going to
> be used for process management and cleanup, if that's the case why
> are you relying on pgid signaling to kill processes?"

Expected answer: process groups are not authoritative. A child can
`setsid` or `setpgid` itself out. cgroups are.

## Pdeathsig is unsafe in Go

> **espadolini** (Zephan92): "pdeathsig considers the end of the
> thread that spawned the child process, so waiting in a different
> goroutine like you're doing has the potential of just killing the
> child process."

> **Zephan92** (response): "Good catch. I've removed Pdeathsig
> entirely to avoid the thread-affinity risks in Go."

Expected answer: do not use `SysProcAttr.Pdeathsig`. The Go runtime
moves goroutines between threads; the thread that spawned the child
can exit while the program is healthy. Use a cgroup sweep on startup
to handle ungraceful worker crashes.

## File descriptor cleanup at start

> **tigrato** (kkloberdanz): "We should close them as soon as the job
> starts. The reason is simple, if you run without cleanup, you
> easily exhaust FDs"

Expected answer: parent closes the cgroup directory fd immediately
after `cmd.Start()`. The kernel has already used it to place the
child; nothing else needs it.

## Don't reap from multiple places

> **espadolini** (Zephan92): "The contract of os.Process demands that
> nothing else reaps the process, this has the potential to do so and
> then the subsequent Kill and Wait might hit a completely unrelated
> process."

Expected answer: exactly one place calls `cmd.Wait()`. Everyone else
observes via shared state.

## Resource-limit configuration must fail

> **tigrato** (kkloberdanz): "Should we allow running the server
> without any limits? Warns are not errors and running the server
> without resource limits seems to be problematic"

> **kkloberdanz** (response): "Good point, I fixed this in [commit]
> which will error out if we are unable to configure cgroups."

Expected answer: if cgroup setup fails at startup, exit non-zero. Do
not run without containment.

## Child processes survive without explicit cleanup

> **sclevine** (benmoss): "What about child processes?"

> **benmoss** (response): "Right now if they inherit stdout/stderr
> cmd.Wait will actually keep blocking until they exit, but if either
> close stdout/stderr or don't inherit it they will be orphaned.
> Added waiting on all child processes!"

Expected answer: the cgroup is the source of truth. `cmd.Wait()`
returns when the leader exits; the cgroup sweep handles any
descendants the leader spawned and then died. Do not rely on
stdout/stderr inheritance for synchronization.

## Cleanup on EBUSY

Not asked verbatim in the corpus but a known pitfall: removing a
cgroup directory immediately after kill can return EBUSY because the
kernel has not yet reaped killed processes. Retry with backoff and a
deadline. See `../../TECHNICAL_PITFALLS.md`.

## Subtree control delegation

Not asked verbatim in the corpus but a real failure mode: writing
`+cpu +memory +io` into `parent/cgroup.subtree_control` succeeds even
when the parent's parent has not delegated those controllers. Read
the file back and verify each requested controller is present;
hard-fail at startup if not. Otherwise the failure shows up later as
a confusing `EINVAL` on `memory.max`.
