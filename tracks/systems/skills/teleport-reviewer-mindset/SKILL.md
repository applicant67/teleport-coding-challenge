---
name: teleport-reviewer-mindset
description: Review a Teleport Job Worker Service submission from the perspective of a Teleport infrastructure engineer. Use when the candidate has produced a design doc, library, or PR and wants a pre-submission review. Outputs the questions a real reviewer is likely to ask, with the expected answers.
tools: Read, Glob, Grep
---

# Teleport reviewer mindset

The Teleport review pool (rosstimothy, tigrato, nklaassen, eriktate,
espadolini, codingllama, jimbishopp, sclevine, zmb3, GavinFrazar,
smallinsky, greedy52, awly, and others) is consistent. Across the
public corpus the same five or six questions get asked on every
submission. Older submissions also see russjones, r0mant, and
alex-kovoy. This skill simulates that review pass on the candidate's
own work.

The output is not a code review of *style*. It is a list of the
questions the reviewer is likely to ask and whether the submission
has a defensible answer.

## How to run it

For each artifact (design doc, library code, gRPC code), walk through
the question set in order. For each question, classify:

- **Answered** - the artifact contains a clear answer
- **Partially answered** - the answer is implied but not stated
- **Unanswered** - the reviewer will ask

Produce a report with one section per question.

## The question set

### 1. How does a reader discover new output?

Reviewers asked this in: joshuarubin, MarkDHarris, benmoss, Zephan92,
chintamanil, sabernabil12.

Acceptable answers:
- `sync.Cond` broadcast on every write, with a monotonic version
  counter readers check under the lock before sleeping
- `chan struct{}` per-reader with non-blocking send by the writer
  (acceptable but more complex)

Unacceptable: bare condition variable without version (lost-wakeup
bug), `io.Pipe` per reader (slow-reader couples writers), bytes
channel fanout (only one reader receives each chunk).

### 2. What happens when a reader is slow?

Reviewers asked this in: joshuarubin, razzam21, rohitsakala,
mcampo84.

Acceptable: file-backed log lets readers seek independently. The OS
handles slow readers via the page cache. Each reader is just
`os.File` with its own offset.

Unacceptable: shared `bytes.Buffer` cloned per reader (memory
explodes). Bounded channel (drops bytes or blocks producer).

### 3. What happens on client disconnect mid-stream?

Reviewers asked this in: Zephan92, benmoss, chintamanil, kkloberdanz.

Acceptable: the reader takes a `context.Context`. Cancellation wakes
the condition variable (via `context.AfterFunc`). The reader rechecks
context after waking.

Unacceptable: the reader blocks forever on the condition. goroutine
leaks accumulate.

### 4. What is the exit code of a killed job?

Reviewers asked this in: joshuarubin, kkloberdanz.

Acceptable: `128 + SIGKILL = 137` from `cmd.ProcessState`. Or `-1`
with an explanation. State which.

Unacceptable: zero, leaving it unset, or returning the same exit
code for "stopped by user" and "ran successfully."

### 5. How is identity established?

Reviewers asked this in: joshuarubin, mcampo84, RichyHBM, Zephan92,
kkloberdanz.

Acceptable: `peer.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0].Subject`.
Refuse if `VerifiedChains` is empty.

Unacceptable: `PeerCertificates[0]` (works without verification),
Serial Number (rotates), client-supplied identity in a request
field.

### 6. Where does authorization happen?

Reviewers asked this in: razzam21, GevorgGal.

Acceptable: an interceptor extracts identity from the cert, attaches
it to the context. A second interceptor checks the operation against
the identity. The handler reads identity from context.

Unacceptable: authorization inside the library API. Library accepts
opaque metadata; it does not know about CAs or roles.

### 7. How do you avoid leaking job existence?

Reviewers asked this in: kkloberdanz, Zephan92.

Acceptable: the same gRPC error code (commonly `NotFound`) for "no
such job" and "not authorized to view this job."

Unacceptable: `PermissionDenied` for one case and `NotFound` for the
other. An attacker enumerates valid IDs by probing.

### 8. How does the job get into the cgroup atomically? (L5)

Reviewers asked this in: kkloberdanz, Zephan92.

Acceptable: `CLONE_INTO_CGROUP` via Go's `SysProcAttr.UseCgroupFD +
CgroupFD`. Available since Linux 5.7.

Unacceptable: fork then write PID to `cgroup.procs`. There is a race
window in which the child can spawn descendants outside the cgroup.

### 9. How do you terminate all descendants? (L5)

Reviewers asked this in: kkloberdanz, GevorgGal, rohitsakala.

Acceptable: write `"1"` to `cgroup.kill`. One operation. Linux
5.14+.

Unacceptable: kill the process group (children can `setsid` out).
Walk `cgroup.procs` and kill each (race with new descendants).

### 10. What if the worker crashes mid-job? (L5)

Asked indirectly via: pdeathsig discussions (Zephan92).

Acceptable: cgroup sweep on worker startup. Walk the parent cgroup
for `job-*` subdirectories, kill any members, remove.

Unacceptable: rely on `Pdeathsig` (unsafe in Go due to thread
migration), rely on the systemd unit to clean up (out of scope for
the challenge).

### 11. What is the test story?

Reviewers asked this in: many (almost every PR).

Acceptable: unit tests run without root, integration tests require
root and run on a Linux host. Tests use `t.Context()`, no sleeps,
no `slog.SetDefault`. `-race` and goleak.

Unacceptable: tests that call RPC handlers directly (bypass
interceptors), tests that depend on `/dev/sda` or any hardware
specific to the candidate's machine.

### 12. Is the library importable?

Reviewers asked this in: joshuarubin (via codingllama), MrChristianL.

Acceptable: library at module root or under `pkg/`. External
consumers can import it.

Unacceptable: every interesting type under `internal/`.

### 13. Do you use stdlib idioms?

Reviewers asked this in: chintamanil, MrChristianL, RichyHBM.

Acceptable: the output reader is an `io.ReadCloser`. Where a callback
mirrors `io.Writer`, use `io.Writer`. Tests use `t.Context()`. Errors
wrap with `%w`.

Unacceptable: custom callback APIs that duplicate stdlib interfaces.

## Output format

Produce a markdown file:

```markdown
# Pre-submission review of <repo>

## Q1: How does a reader discover new output?
Status: ANSWERED / PARTIAL / UNANSWERED
Evidence: <file>:<line> snippet
Predicted question: "<verbatim from corpus>"

[repeat for each question]

## Summary
- N answered
- M partial (need to firm up before submission)
- K unanswered (high priority)

## Recommended response strategy
For each partial/unanswered: <one-line action>.
```

The candidate uses this report to firm up the design doc before
submission. The reviewer asks fewer questions per round; the round-
trips compress.

## Related reference

Full reviewer quotes by theme: `../../corpus/reviewer-feedback/`.
