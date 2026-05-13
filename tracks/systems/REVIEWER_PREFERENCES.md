# What reviewers consistently ask about

Synthesized from public reviewer comments across roughly 15 historical
submissions. The reviewer pool rotates but their questions do not. If your
design doc and code preempt these questions, you save a review round-trip.

The reviewer handles you will most often see are `rosstimothy`, `tigrato`,
`nklaassen`, `eriktate`, `codingllama`, `espadolini`, `greedy52`, `zmb3`,
`jimbishopp`, `sclevine`, `Joerger`, `r0mant`, `fspmarshall`. Comment style
varies, but the themes are remarkably stable.

## The repeatable questions

### Output streaming

This is the single most-asked topic.

- "How does the reader find out new data is available?" Answer: a per-reader
  notification channel, fed by the writer, drained by the reader. Or a
  condition variable broadcast. Either is fine. Polling is not.
- "How will Read be unblocked when more output is available?" Same question
  framed differently.
- "Why does each reader get a full copy of the output buffer?" If your design
  is one-buffer-per-reader, change it. The expected design is one shared
  log (file-backed or in-memory) plus per-reader cursors.
- "If a job outputs 1 GB and there are 5 readers, do we use 5 GB of RAM?"
  Same point.
- "Can a slow reader block other readers?" Answer in the design doc.
  Independent file descriptors per reader is the simplest way to say no.
- "Output is not guaranteed to be UTF-8 or newline-terminated. Why is your
  proto field a string?" Make it `bytes`.
- "Is there a race between the final write and closing the done channel?"
  This is the one rosstimothy specifically asks. Make sure your design
  drains the writer completely before signaling readers that the stream is
  done. The most common bug: the goroutine that reads from a pipe finishes
  copying, you close the done channel, but the kernel still has bytes in
  the pipe that os/exec has not flushed.
- "What happens if a client disconnects while streaming?" The stream
  resources for that one reader must release promptly. Other readers
  unaffected.

### Job status / lifecycle

- "Can your status distinguish 'killed by user' from 'failed with non-zero
  exit'?" Yes. Two states. Often `STOPPED` and `FAILED`.
- "What is the exit code of a stopped job?" Either `-1`, or some sentinel
  with a clearly documented meaning. Many reviewers prefer `128 + signal`
  for signal kills.
- "Is StopJob synchronous?" The expected answer is no. Stop initiates a
  SIGKILL via cgroup.kill and returns; the caller polls or watches for the
  terminal status. Reviewers will accept synchronous if you justify it,
  but async is the lower-friction answer.
- "What if Stop is called multiple times concurrently?" Idempotent.
- "What if Stop is called on a job that already exited?" Idempotent.

### Process containment

- "How do you guarantee all descendants of the job process are killed?"
  Answer: cgroup v2 with `cgroup.kill`. Mention the kernel version floor
  (5.14+ for cgroup.kill, 5.7+ for CLONE_INTO_CGROUP).
- "Is there a race between fork and writing the PID to cgroup.procs?" Yes,
  and you fix it by spawning the child directly into the cgroup using
  `CLONE_INTO_CGROUP` via `syscall.SysProcAttr{UseCgroupFD: true, CgroupFD:
  fd}`. The fork-then-write pattern leaves a window where the child can
  spawn descendants outside the cgroup.
- "Why not pdeathsig?" Because Go runtime can migrate goroutines across
  OS threads; the thread that spawned the child can die while the job is
  healthy, sending an inadvertent SIGKILL. Mention this in the design.
- "What cleans up cgroups if the worker process dies?" A sweep on startup
  that removes leftover per-job cgroup directories under the parent.

### Resource limits

- "Does the server require limits, or are they optional?" The challenge
  requires limits. Hardcode reasonable defaults server-side, do not let
  clients pass them, and do not allow zero / unlimited.
- "Is io.max in scope?" Yes, per rosstimothy explicitly.
- "What block devices do you apply io.max to?" Top-level block devices
  enumerated from `/proc/partitions`, filtered to whole disks (those with a
  `/sys/block/<name>` entry). Typical exclusions: major 1 (ramdisks),
  major 7 (loop). Some reviewers prefer including these; ask in your design
  doc or just include them and document why.

### Authentication and identity

- "What identifies a client: serial number or subject?" Subject (or CN, or
  SAN). The cert serial number rotates on re-issue and is the wrong
  identity stable-point.
- "Where is the cert checked in the gRPC code path?" Use `peer.FromContext`
  on the verified chain (`ConnectionState.VerifiedChains[0][0]`), not
  `PeerCertificates`. The latter is set even on misconfigurations that let
  unverified clients through.
- "Can you separate authentication from authorization?" Yes. One interceptor
  pulls identity off the cert, another decides whether that identity may
  call this RPC.
- "Does this leak existence of jobs to unauthorized callers?" Watch for
  the case where an unauthorized caller can distinguish "job does not
  exist" from "job exists but you cannot access it". Return the same
  error for both.

### API / library shape

- "Is the library importable?" If your reusable library lives under
  `internal/`, no. The challenge expects an importable package, typically
  at the module root or under `pkg/`.
- "Are transport types leaking into the library API?" The library should
  expose its own types (`job.Status`, not `pb.JobStatus`).
- "Should the owner be part of the library API or the transport layer?"
  The transport layer enforces authorization; the library itself does not
  need to know about the caller's identity. Pass an owner in to start a
  job, store it, but do not let the library make authorization decisions.

### Testing

- "Why does this test sleep?" Sleeps in tests are a smell. Use a channel,
  a `t.Context()`, or a polled condition. The common case where a sleep
  feels necessary: waiting for the child process to enter the cgroup. Use
  a tight poll loop with a tight deadline.
- "Why are you using `context.Background()` in a test?" Use `t.Context()`.
  It is cancelled when the test ends, which prevents test-leak goroutines.
- "Why are the RPC handlers tested by direct method calls?" Reviewers want
  a real listener with the full interceptor chain exercised. Spin up a
  real gRPC server in tests on a chosen port (or `bufconn`) and use a real
  client.
- "Why does this test depend on the host's block devices / cgroup state /
  network config?" It should not. See `REJECTION_LESSONS.md` item 1.

### Code organization and tooling

- "Why is `.pb.go` not committed?" Commit it. Reviewers expect to clone
  and run without extra steps.
- "Why are you using fmt-style logs?" Use `log/slog` with structured
  fields.
- "Why is this concern split across two functions?" See
  `REJECTION_LESSONS.md` item 3.

### PR / submission shape

- Reviewers like staged PRs: design doc first, then the library, then the
  server, then the CLI. They explicitly approve of this structure.
- A first design-doc PR is expected, not optional. Many reviewers will not
  look at code PRs until the design doc is approved.
- Each PR should be reviewable on its own. Reviewers tolerate stacking
  (the design doc PR has to land before the lib PR can be merged) as long
  as the diffs are individually digestible.

## How to use this list

Before you open a PR, search this doc for the topic the PR touches and
make sure your design doc preempts every relevant question with a one-line
answer. Reviewers do not need essays. They need the question answered.
