# Status and lifecycle

How the library reports state, what statuses exist, what the exit code
is when a user kills a job, and whether Stop is synchronous.

## Stopped vs Failed vs Succeeded

> **tigrato** (joshuarubin): "can you add other status indicating the
> job was manually stopped?"

> **rosstimothy** (joshuarubin): "What will the exit code of a stopped
> job be? Will the library be able to differentiate between the error
> and stopped states?"

> **nklaassen** (MrChristianL): "we do like to see a solution that is
> able to distinguish between jobs that _stopped naturally_ and jobs
> that were _killed by a user of the library_"

> **eriktate** (GevorgGal): "Do we need this extra bool? It seems like
> this could just be another value in JobStatus"

> **rosstimothy** (GevorgGal): "I'm curious, how come you opted for a
> standalone boolean flag to represent this instead of adding a
> JOB_STATUS_STOPPED state to the JobStatus enum?"

Expected answer: states are `RUNNING`, `STOPPING`, `STOPPED`, `EXITED`
(or `COMPLETED`), `FAILED`. A separate `stop_reason` distinguishes
user-stop from kernel-OOM from natural exit. Do not use a parallel
boolean alongside the enum.

## Enum zero values

> **rosstimothy** (MarkDHarris): "Should we make this iota + 1 to
> prevent any zero values from defaulting to admin?"

Expected answer: never let zero be a privileged value. `iota + 1`, or
explicit numeric assignments, or in protobuf a dedicated
`*_UNSPECIFIED = 0`.

## Exit code for a killed job

> **smallinsky** (kkloberdanz): "It look like the exit code is never
> set when the job is killed. Could you add proper handing and make
> sure that exit code matches `128 + <signal_number>` even if the job
> was killed by `func (l *localJob) Stop() error {` API ?"

Expected answer: pull the exit code from `cmd.ProcessState`. If the
process was signaled, the convention `128 + signum` is the
shell-standard encoding; report it consistently. If you go with `-1`
for "killed by us," explain why.

## Stop semantics: blocking or async

> **rosstimothy** (joshuarubin): "Does this imply that StopJob is
> synchronous and will block until the job is terminated?"

> **greedy52** (MarkDHarris): "I am a little confused on
> context.CancelFunc. Is it canceling a context passed in to exec.Cmd?
> What are the things actually happen when Stop is called?"

Expected answer: state it explicitly. Async (Stop signals the cgroup
and returns; observers watch status) is common. Blocking-until-exited
is also valid if documented. What is not valid is leaving the question
open for the reviewer to chase.

## Double Stop

> **nklaassen** (MrChristianL): "could multiple concurrent calls to
> Stop be a problem?"

> **tigrato** (kkloberdanz): "l.cmd.Wait() can only be called once.
> Should we document that behaviour?"

> **tigrato** (mcampo84): "being this publicly exposed, a user can
> call it multiple times which will cause panics because multiple
> channel close operations"

Expected answer: Stop is idempotent. State explicitly in the design
doc that the second call is a no-op. `cmd.Wait` lives in exactly one
place; readers and Stop callers never invoke it directly.

## Status atomicity

> **rosstimothy** (MrChristianL): "Does the lock need to be held while
> doing IO to uphold the guarantees mentioned in this comment?"

Expected answer: a `Snapshot()` method that takes the lock once,
returns status + exit code + stop reason in one atomic copy. Avoids
the bug where the caller asks for status, gets RUNNING, then asks for
exit code and gets a populated value because the job exited between
calls.

## Jobs persist after stop

> **rosstimothy** (nombiezinja): "Jobs should be stopped but not
> removed. For example, I should be able to terminate a running job,
> and then describe or get output of the job."

Expected answer: terminated jobs remain queryable. Removal is an
explicit operation if it exists at all (it does not need to exist for
the challenge).
