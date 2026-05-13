# Historical submissions inventory

A list of public past submissions to the Teleport Job Worker Service
challenge, with what we know about each. This list is not exhaustive; it
is what one candidate found by:

1. Searching GitHub for forks of the official challenge repos.
2. Searching GitHub for repos whose names contain "job-worker",
   "teleport-challenge", "teleport-exercise", "jobworker",
   "job-worker-service".
3. Searching public PRs commented on by known Teleport-engineer GitHub
   handles.

Levels noted here are inferred from the README, design doc, or the
presence of cgroup code when not stated explicitly. The corpus is rough;
treat it as a starting point for your own research, not as ground truth.

## Submissions with public reviewer feedback

These are the most valuable. Reading the PR threads on these repos tells
you what reviewers actually ask. Sorted approximately by recency.

### MarkDHarris/teleport-systems-challenge

- URL: https://github.com/MarkDHarris/teleport-systems-challenge
- Date: 2026-03 / 2026-04
- Level: L4 (explicit in PR 1 title)
- Reviewers: rosstimothy, greedy52, rhammonds-teleport
- Approach: Go + gRPC + mTLS. `context.AfterFunc` to wake `sync.Cond` on
  cancellation. Owner + admin authorization. Tests use `testing.T.Context()`,
  no sleep-based sync.
- Pattern worth copying: replace sleep-based test sync with channels.
- Pattern reviewers cut: "allowed viewers" multi-user authorization,
  interactive shell / stdin support, graceful SIGTERM-then-SIGKILL.

### MrChristianL/Teleport-Job-Worker-Service

- URL: https://github.com/MrChristianL/Teleport-Job-Worker-Service
- Date: 2026-03
- Level: L4 (explicit in PR 1 title)
- Reviewers: nklaassen, eriktate, rosstimothy, Joerger
- Approach: Go + gRPC + mTLS. Initially callback-based output writer
  (`func([]byte) error`), refactored to `io.Writer`. `Snapshot()` method for
  atomic status queries.
- Pattern worth copying: `Snapshot()` returning status, exit code, stop
  reason in one lock acquisition.
- Pattern flagged: EOF message in proto (redundant; removed), reusing
  callback shape that mirrors `io.Writer`.

### GevorgGal/jobworker

- URL: https://github.com/GevorgGal/jobworker
- Date: 2026-03
- Level: L4
- Reviewers: rosstimothy, eriktate, Tener
- Approach: Go + gRPC + mTLS + cgroup v2. `OutputReader` type bundles
  offset tracking, notification, and cancellation behind one `Read(ctx)`
  method. 32KB cap per read. `buf` for proto generation.
- Pattern worth copying: clean `OutputReader` abstraction; cap on per-read
  size to avoid huge buffers for late-joining readers.
- Pattern flagged: clunky early API with cloned buffers; sleep-based tests
  (replaced with channels).

### chintamanil/job-worker

- URL: https://github.com/chintamanil/job-worker
- Date: 2026-03
- Level: L4
- Reviewers: tigrato, espadolini
- Approach: Go + gRPC + mTLS + cgroup v2. `MemoryBufferReader` implements
  `io.ReadCloser`. `sync.Cond` notification. Switched from `RWMutex` to
  `sync.Mutex` after profiling.
- Pattern worth copying: `io.ReadCloser` implementation for the streaming
  output.
- Pattern flagged: blocking readers on disconnected clients; pgid-based
  termination insufficient without cgroups.

### kkloberdanz/teleport-challenge

- URL: https://github.com/kkloberdanz/teleport-challenge
- Date: 2026-02
- Level: L5
- Reviewers: tigrato, smallinsky
- Approach: Go + gRPC + mTLS + cgroup v2. `cgroups.kill` for termination.
  Per-subscriber channel for output.
- Pattern flagged: protobuf types leaking into library public API; slow
  reader on buffered channel; authorization not failing closed; job ID
  enumeration via specific "permission denied" vs "not found" responses
  (fix: return the same error for both).

### benmoss/job-worker-service

- URL: https://github.com/benmoss/job-worker-service
- Date: 2026-01
- Level: L4
- Reviewers: jimbishopp, sclevine, nklaassen, dboslee
- Approach: Go + gRPC + mTLS. File-backed output log. `sync.Cond` with
  version counter to avoid lost wakeups. `ReadCloser`-based stream
  cleanup.
- Pattern worth copying: file-backed output log; version counter to make
  notification race-free.

### joshuarubin/teleport-job-worker

- URL: https://github.com/joshuarubin/teleport-job-worker
- Date: 2024-09
- Level: L4
- Reviewers: tigrato, codingllama, rosstimothy
- Approach: Go + gRPC + mTLS. Certificate Subject (not serial) for
  identity. Three-PR structure: design, library, server+CLI.
- Pattern flagged: serial-number-based identity (breaks on cert renewal);
  UTF-8 assumption in output (should be binary-safe).

### adalton/teleport-exercise

- URL: https://github.com/adalton/teleport-exercise
- Date: 2021-12 (older but same challenge family)
- Level: L5
- Reviewers: codingllama, rosstimothy
- Approach: Go + gRPC + mTLS + Linux namespace isolation (PID, mount,
  network).
- Pattern flagged: namespace isolation adds significant scope; tested
  RPCs via direct method calls instead of through a real client.

### RichyHBM/teleport-challenge

- URL: https://github.com/RichyHBM/teleport-challenge
- Date: 2025-04 / 2025-05
- Level: L4
- Reviewers: rosstimothy, tigrato, eriktate
- Approach: Go + gRPC + mTLS. Identity from `peer.FromContext`.
  Get-buffer + connect-subscriber streaming pattern.
- Pattern flagged: data race between get-buffer and connect; impersonation
  vulnerability when identity is not pulled from the verified chain.

### razzam21/job-worker-service

- URL: https://github.com/razzam21/job-worker-service
- Date: 2025-08
- Level: L4
- Reviewers: zmb3, creack, timothyb89
- Approach: Go + gRPC. Heavy RFD-first design phase. Simplified streaming
  to follow-from-start-only (no offset / no follow flag).
- Pattern worth copying: explicitly cutting `from_offset` and `follow`
  options to reduce scope.
- Pattern flagged: racy slice append (mutates memory behind pointer);
  authorization in library API instead of server.

### ilyazz/jobs

- URL: https://github.com/ilyazz/jobs
- Date: 2022-10
- Level: L5
- Reviewers: sclevine, jimbishopp, andreiko
- Approach: Go + gRPC + mTLS + cgroup. Self-exec shim mode (`/proc/self/exe`
  with mode flag) to assign child to cgroup. Read-only superuser role.
- Pattern flagged: complex shim approach; multiple cleanup paths.

## Submissions where we did not find public reviewer feedback

Source code is visible. PR threads are either private, on a non-forked
repo, or never existed. Useful as reference for patterns and against
review, but you cannot see what reviewers asked.

### bucknercd/jobworker

- URL: https://github.com/bucknercd/jobworker
- Date: 2026-02
- Level: L5 (explicit)
- Approach: Go + gRPC. cgroup v2, semi-persistent storage. Chunked
  streaming.

### rajansandeep/teleport-job-worker-svc

- URL: https://github.com/rajansandeep/teleport-job-worker-svc
- Date: 2026-04
- Level: L4
- Approach: Go + gRPC + mTLS. Admin/viewer role from cert CN. Binary-safe
  output.

### mikewurtz/taskman

- URL: https://github.com/mikewurtz/taskman
- Date: 2025-04
- Level: L4 to L5 (cgroup integration present)
- Approach: Go + gRPC + cgroup v2. Hardcodes `/dev/sda` for io.max testing
  (cautionary: this is the host-dependent test problem). User-ID based
  authorization rather than mTLS.

### kurczynski/teleport-job-worker

- URL: https://github.com/kurczynski/teleport-job-worker
- Date: 2024-06
- Level: L4
- Approach: Sparse README; assume standard gRPC + mTLS.

### gstelang/job-worker-service

- URL: https://github.com/gstelang/job-worker-service
- Date: 2024-09
- Level: unknown
- Approach: minimal README.

### renatoaguimaraes/golang-job-scheduler

- URL: https://github.com/renatoaguimaraes/golang-job-scheduler
- Date: 2021-05
- Level: L3 or L4
- Approach: Go + gRPC + mTLS. Disk-based log persistence. Channel-based
  streaming (string-typed output is not binary-safe).

### jaylane/job-scheduler

- URL: https://github.com/jaylane/job-scheduler
- Date: 2024-08
- Level: unknown
- Approach: minimal documentation.

### diptadas/job-worker

- URL: https://github.com/diptadas/job-worker
- Date: 2021-06
- Level: L2 or L3
- Approach: Go + HTTPS + mTLS. REST API. Older era of the challenge.

### mikihau/rest-job-worker

- URL: https://github.com/mikihau/rest-job-worker
- Date: 2019-12
- Level: L1 or L2
- Approach: Go + REST + mocked auth. Pre-gRPC era.

## Patterns seen across many submissions

Things multiple candidates landed on independently:

- File-backed output log with per-reader file descriptors (cleanest
  multi-reader story).
- `sync.Cond` with a monotonic version counter to avoid lost wakeups.
- 16KB or 32KB cap on per-read response chunk size to avoid huge
  allocations for late-joining readers.
- Identity from `peer.FromContext` + verified cert chain.
- Certificate Subject (not Serial) as user identity.
- cgroup v2 with `cgroup.kill` for termination.
- Cgroup sweep on worker startup for crash cleanup.
- Hardcoded server-side resource limits (no client-controlled limits).
- `context.AfterFunc` to wake condition variables on cancellation.
- `t.Context()` in tests, not `context.Background()`.
- buf or vendored `protoc-gen-go` for reproducible proto generation.

## Anti-patterns called out repeatedly

Things multiple candidates were asked to change:

- Output type as `string` in protobuf (should be `bytes`).
- Each reader getting its own full output buffer (5 readers, 5x RAM).
- Cert Serial as identity.
- `PeerCertificates[0]` for identity (use `VerifiedChains[0][0]`).
- Authorization decisions in the library API (move them to the server).
- Library entirely under `internal/` (not importable).
- `from_offset` and `follow` options on streaming (cut them).
- Graceful SIGTERM-then-SIGKILL (just SIGKILL via cgroup.kill).
- Stdin support (out of scope).
- Job listing API (out of scope).
- Namespace isolation on L4 (out of scope).
- Direct RPC handler calls in tests (use a real listener).
- Sleep-based synchronization in tests.
- `slog.SetDefault` mutation in `TestMain` (global state).
- io.max device list hardcoded (host-dependent).

## Additional submissions (compact listing)

These repos have public reviewer PR threads but are listed compactly
because we have not yet written full per-repo entries above. For full
context, open the repo and read the PR threads directly.

- `mikewurtz/taskman` - rosstimothy, tigrato, zmb3 (~80 comments). Heavy
  thread on scope, auth, streaming, errors, TLS, design-doc style.
- `tjper/teleport` - codingllama, rosstimothy, zmb3, sclevine, jimbishopp
  (~91 review-comments across 3 PRs). L4-ish.
- `adalton/teleport-exercise` - codingllama, rosstimothy (~104 comments).
  Heavy iteration.
- `Zephan92/teleport-handson` - codingllama, espadolini, nklaassen,
  rosstimothy (7 PRs reviewed). Multiple rounds of cgroup feedback.
- `mcampo84/teleport_challenge` - tigrato. Authorization-by-PID
  anti-pattern thread.
- `samschurter/teleport-challenge` - smallinsky. Scope-cutting feedback.
- `shawon-crosen/teleport-challenge` - rosstimothy. Process-runner
  naming; design clarity.
- `rsteinkeXJ/teleport_interview` - tigrato. Recent (2026-04);
  stderr/stdout interleave guidance.
- `rexposadas/teleport` - rosstimothy. TLS/cipher and authorization
  scoping.
- `moalf/teleport` - russjones. "Drop to reduce scope" theme.
- `neildo/tjob` - rosstimothy. Streaming-default rule explicit.
- `kovyrin/teleport-exec` - zmb3. Don't-delete-output-on-status thread.
- `kiakeshmiri/process-runner` - tigrato. UTF-8 / slow-consumer canonical
  questions.
- `sabernabil12/teleport-job-worker` - rosstimothy. "Process isolation
  not required at L4."
- `xSudoNymx/job-manager` - tigrato. sigterm-then-sigkill not required.
- `kurczynski/teleport-job-worker` - GavinFrazar. Explicit "list/get jobs
  out of scope" statement.
- `atburke/teleport_interview` - russjones, alex-kovoy. Older, but
  russjones quotes are unique to this thread.

## Reviewer-quote highlights

A few canonical quotes worth knowing before you write your design doc:

- **rosstimothy on neildo/tjob:** "streaming should be the default and
  only option, and ... you could simplify by ... drop the `--tail` and
  `--follow` flags."
- **rosstimothy on mikewurtz/taskman:** "Shouldn't the admin be able to
  retrieve the status for any task? ... terminate ... stream the task
  output?" - admin-role flows expected, not just owner-only.
- **zmb3 on kovyrin/teleport-exec:** "A client should be able to stream
  the output as many times as they'd like, so let's not delete the file
  when it may still be needed." - replayable output.
- **zmb3 on kovyrin/teleport-exec:** "Can you think of a way to do
  authorization without needing access tokens?" - mTLS-cert-CN, not JWT.
- **tigrato on mcampo84:** "Is PID a good external identifier for jobs?
  Is it possible that PIDs get reused...?" - PID-reuse anti-pattern.
- **tigrato on kiakeshmiri:** "What will happen if a consumer is slow to
  consume the logs? Are logs guaranteed to be utf8 compatible?"
- **GavinFrazar on kurczynski:** "fetching/listing the jobs owned by a
  user is out of scope ... The API only needs to allow a user to start,
  stop, stream output, and get the status of an individual job."

## How to use this list

Pick three submissions whose reviewers overlap with the current Teleport
review pool (rosstimothy, tigrato, nklaassen, eriktate, codingllama,
espadolini, zmb3 are common in recent ones; russjones, alex-kovoy,
r0mant appear on older systems and fullstack threads). Read every
comment on their design doc PR and library PR. Most of what reviewers
will ask you, they asked someone else first.
