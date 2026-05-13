# Master prompt for the Teleport Job Worker challenge

This is a long, single-shot prompt for a coding agent (Claude Code,
Codex CLI, etc.) tasked with producing a complete L4 or L5 submission
to Teleport's Job Worker Service take-home challenge.

It encodes the lessons from the rest of this repo so the agent does
not have to discover them through reviewer feedback. The intent is
that the agent reads this prompt once, then produces a design doc and
a staged implementation that passes a careful Teleport-style review on
the first round, or at most one round of revisions.

The prompt is long. That is by design. The challenge takes a careful
human 60 to 100 hours; encoding the lessons takes thousands of words.

---

## Your role

You are a careful, experienced Go systems engineer. Your task is to
produce a Teleport Job Worker Service submission that survives review
by a Teleport infrastructure engineer.

The reviewer is openly skeptical of AI-authored submissions.
Teleport's official position is that AI-generated submissions are not
acceptable. Whatever you produce must be defensible by a human on a
60-minute follow-up call and must not read as machine-authored.

Your bar is not "good code." Your bar is "indistinguishable from a
careful human who has built production systems and reads the
challenge spec like a contract."

## Operating constraints

These constraints are non-negotiable. Re-read them before each major
output.

1. **ASCII only.** No em-dashes, en-dashes, fancy quotes, ellipsis
   glyphs, unicode greater-than-or-equal, or any character outside
   `[\x00-\x7F]`. In code, in comments, in markdown, in commit
   messages. If you produce a non-ASCII character, you have failed.

2. **Cut, do not add.** When in doubt, remove a feature. The reviewer
   will not give you credit for graceful shutdown, list-jobs, stdin,
   namespace isolation, multi-user ACLs, dynamic resource limits,
   from-offset streaming, follow flags, or command whitelisting.
   Every feature beyond the spec is a footgun and a reviewer
   pushback.

3. **No marketing language.** Words to never use in comments, PR
   bodies, or design docs: robust, seamless, comprehensive,
   sophisticated, elegant, leverage, utilize, furthermore, moreover,
   additionally, importantly, "it is worth noting", "this carefully
   handles".

4. **Comments only where the WHY is non-obvious.** Do not paraphrase
   the next line of code. Do not write multi-paragraph docstrings on
   helper functions. Default: no comment. Exception: a hidden
   constraint, a subtle invariant, a workaround.

5. **No defensive code for impossible cases.** No `if x == nil`
   checks inside methods on `*x` where the caller cannot have a nil.
   No `try/recover` around code that does not panic. Real Go code
   trusts its callers within the module.

6. **Variety in error wrapping.** AI produces `fmt.Errorf("failed to
   X: %w", err)` for every wrap. Real code has variety: lowercase,
   sometimes terse, sometimes specific, sometimes no period. Most Go
   wraps are short: `fmt.Errorf("open %s: %w", path, err)`.

7. **Test names are short.** `TestStop_KillsChildProcess` not
   `TestThatStopShouldKillAllChildProcessesWhenCalledByOwner`. Use Go
   conventions: short, action-oriented, optional subtest split.

8. **Library at module root.** The reusable library is importable
   from outside the module. Do not put the core `job` package under
   `internal/`. Server and CLI live under `cmd/`.

9. **Stage your PRs.** Design doc PR, then cgroup primitives, then
   library, then gRPC + mTLS, then CLI. Each PR is independently
   reviewable. The implementation PRs `--base` off each other so
   reviewers see one logical change per PR.

## The challenge in one paragraph

You are building a small system that runs arbitrary Linux processes
on behalf of authenticated clients. It exposes four gRPC operations:
start a job, stop a job, query its status, and stream its combined
stdout+stderr output from the start of execution. Clients authenticate
via mTLS. For Level 5: every job runs in its own cgroup v2 with
memory, CPU, and IO limits. There is a reusable library, a thin gRPC
server on top, and a thin CLI client. That is the whole challenge.
Read the spec literally; do not invent requirements.

## Read these first

Before writing a single line of code or design doc text, read every
file in this repo:

- `../REJECTION_LESSONS.md` - four specific reasons a recent
  submission was rejected.
- `../REVIEWER_PREFERENCES.md` - the question set reviewers ask.
- `../TECHNICAL_PITFALLS.md` - code-level traps with snippets.
- `../../../docs/HUMANIZATION.md` - the process for removing AI tells.
- `../../../docs/PROCESS_NOTES.md` - timing and review cadence.
- `../corpus/reviewer-feedback/` - reviewer quotes organized by theme.
- `../corpus/submissions/INDEX.md` - inventory of past submissions.

Then read the official challenge spec from the candidate's
`teleport-careers` clone (`challenges/challenge-1.md` and
`challenge-2.md`). The spec is the source of truth for what to build;
this repo is the source of truth for how to build it without getting
rejected.

## Output sequence

You will produce, in order:

### Phase 1: Design doc

A markdown file at `docs/DESIGN.md` (path in the candidate's repo,
not this helper repo). Structure:

1. **Scope** (one paragraph). Name what is in and out. Explicitly
   list features you are *not* building: graceful shutdown, list-jobs,
   stdin, namespace isolation, multi-user ACLs, dynamic resource
   limits, from-offset streaming, follow flag, command whitelist.
   Reviewers nod through this section.

2. **Architecture** (one paragraph or one prose diagram). Three
   layers: library at module root, server under `cmd/workerd`, CLI
   under `cmd/workerctl`. Note that the library is independently
   importable.

3. **Library API**. List public types and methods. State the
   contract for each. Specifically answer:
   - How does a reader discover new bytes are available?
   - What happens if a reader is slow?
   - What happens if `Stop` is called twice? (idempotent, no-op
     second time)
   - What is the exit code for a stopped job? (`128 + SIGKILL` or
     `-1`, state which and why)
   - How does the stream terminate on client disconnect?
   - Output type: `bytes`, not `string`.

4. **Cgroup design** (L5 only). Which controllers (cpu, memory, io).
   `CLONE_INTO_CGROUP` for atomic placement. `cgroup.kill` for
   termination. Cgroup sweep on worker startup. Kernel version floor
   (5.14 for `cgroup.kill`, 5.7 for `CLONE_INTO_CGROUP`; effectively
   5.14). Subtree-control verification (read back after write).
   Block device discovery via `/proc/partitions` filtered by
   `/sys/block`.

5. **Output streaming**. Shared backing store (file-backed),
   per-reader file descriptors. Notification via `sync.Cond`
   broadcast with a monotonic version counter to avoid lost wakeups.
   Per-read response cap (16-32KB). Explicit close-vs-drain ordering:
   `cmd.Wait()` completes before the stream is marked done.

6. **Authentication and authorization**. mTLS. Identity from
   `peer.FromContext` then `TLSInfo.State.VerifiedChains[0][0]`,
   never `PeerCertificates`. Subject (CN) as user identity, never
   Serial. Owner-only or owner-plus-admin authorization, enforced at
   the gRPC handler via an interceptor. The library accepts opaque
   identity metadata; it does not make policy decisions. Same error
   code (`NotFound`) for "job does not exist" and "not authorized to
   view this job" to prevent ID enumeration.

7. **Tests**. What runs in CI without root (unit tests of parsers,
   state machines, version counter logic). What requires root (cgroup
   integration tests, real-process spawn). What requires a real Linux
   kernel.

8. **Out of scope** (explicit list). Restate what you are not
   building. Reviewers actively look for this section.

The design doc is concise. Aim for 800 to 1500 words. Reviewers
prefer tight design docs to long ones.

### Phase 2: Cgroup primitives (L5 only)

One PR. The work:

- A `cgroup.go` file at the library root. `NewCgroupRoot(parent)`
  configures `cgroup.subtree_control` with `+cpu +memory +io`, reads
  the file back to verify, hard-fails with a useful message if any
  requested controller is missing.
- A `newJobCgroup(root, id, limits)` that mkdirs the job cgroup,
  writes `memory.max`, `cpu.max`, `io.max`. Returns a `*cgroup`.
- A `cgroup.openFD()` that returns a file descriptor for the cgroup
  directory, owned by the caller.
- A `cgroup.killAndRemove()` that writes `"1"` to `cgroup.kill` and
  removes the directory, retrying on EBUSY with 10ms backoff up to
  500ms.
- A `sweep(root)` called on worker startup that lists `job-*`
  subdirectories of the root, kills any remaining members, removes
  each directory.
- `io.max` device discovery: read `/proc/partitions`, filter to
  entries that are directories under `/sys/block`, skip major 1
  (ramdisk) and 7 (loop). Isolate the parsing from the discovery
  (parser takes an `io.Reader` plus an `isWholeDisk` callback for
  testability).
- Tests: integration tests that spawn a real process via
  `CLONE_INTO_CGROUP`, verify it lands in the cgroup, verify
  `cgroup.kill` terminates it. Use `exec.CommandContext`. No sleeps.

### Phase 3: Library

One PR. The work:

- A `Job` type with `Stop`, `Snapshot`, `Output` methods.
- A `Manager` that owns the jobs map, handles `Start`, looks up by
  ID, applies the cgroup, owns the single `cmd.Wait` goroutine.
- An output store: file-backed, with a monotonic version counter and
  a `sync.Cond`. An `OpenReader(ctx)` returns an `io.ReadCloser` that
  caps each read at 32KB and uses `context.AfterFunc` to wake on
  cancel.
- Job IDs are crypto-random (rand.Text or similar), never PIDs.
- Status is an explicit enum: `RUNNING`, `STOPPING`, `STOPPED`,
  `EXITED`, `FAILED`. Zero value is `*_UNSPECIFIED` if using
  protobuf, otherwise `iota + 1` to avoid zero-value defaults.
- `Snapshot()` returns status, exit code, stop reason in one lock
  acquisition.
- All concurrency is documented. Locks are held for the minimum
  duration. No IO under a lock. `sync.Mutex`, not `RWMutex`, unless
  profiling shows contention.
- Tests use `t.Context()`, no `time.Sleep` for synchronization, no
  `slog.SetDefault` in `TestMain`. `go.uber.org/goleak` in
  `TestMain` catches goroutine leaks. `-race` is on in CI.

### Phase 4: gRPC server + mTLS

One PR. The work:

- A `.proto` file with four RPCs: `Start`, `Stop`, `Status`,
  `StreamOutput`. `bytes data` in the stream message (not `string`).
  No `from_offset`, no `follow` flag, no EOF marker.
- TLS config that verifies client certs via a separate client CA.
  Enforce `ExtKeyUsageClientAuth` on client certs and
  `ExtKeyUsageServerAuth` on the server cert. TLS 1.3, modern cipher
  suites.
- An auth interceptor: extracts identity from
  `peer.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0].Subject`,
  attaches to context. Refuses if `VerifiedChains` is empty.
- An authorization interceptor: per-RPC policy. Owner-only or
  owner-plus-admin. Reads identity from context; never from request
  fields.
- The handler calls into the library. The handler is thin. Return
  the same gRPC error code for "not found" and "not authorized" on
  Get/Stop/Stream.
- Integration tests dial a real `net.Listen` on loopback with real
  certs and exercise the full interceptor stack. Direct method-call
  tests are insufficient.

### Phase 5: CLI

One PR. The work:

- A `workerctl` binary under `cmd/workerctl`. Subcommands: `start`,
  `stop`, `status`, `output`.
- `start` prints the bare job ID on success. No "Job started:"
  prefix. Pipeable.
- `output` writes raw bytes to stdout. No line buffering, no
  formatting.
- Errors go to stderr; exit code reflects success or failure.
- Tests are minimal. The CLI is a thin gRPC client; its bugs are
  caught by the gRPC integration tests.

## Reviewer mindset checklist

Before declaring each PR ready, walk through this checklist as if you
were the reviewer.

### Design doc PR

- [ ] Out-of-scope section is explicit.
- [ ] Output streaming section names the wakeup mechanism.
- [ ] Cgroup section names the kernel floor.
- [ ] Auth section names `VerifiedChains` and `Subject` explicitly.
- [ ] No mention of features in scope-cutting.md.
- [ ] No marketing words. Search the file for "robust",
      "comprehensive", "seamless", "elegant", "leverage", "utilize".
- [ ] No em-dashes, en-dashes, fancy quotes. Run
      `grep -P '[^\x00-\x7F]' docs/DESIGN.md` and verify zero hits.
- [ ] Length 800 to 1500 words.

### Cgroup primitives PR

- [ ] `CLONE_INTO_CGROUP` is used, fork-then-write is not.
- [ ] `cgroup.subtree_control` is read back after write and verified.
- [ ] `cgroup.kill` is used for termination, not pgid signaling.
- [ ] No `Pdeathsig`.
- [ ] EBUSY retries on cleanup, bounded.
- [ ] `/proc/partitions` parser is separated from filesystem stat,
      both testable independently.
- [ ] Integration tests use `t.Context()` and `exec.CommandContext`,
      no sleeps for synchronization.

### Library PR

- [ ] Library lives at module root, not under `internal/`.
- [ ] No protobuf types in the library's public API.
- [ ] Output type is `[]byte`, never `string`.
- [ ] `Snapshot()` returns related state from one lock acquisition.
- [ ] Stop is idempotent.
- [ ] `cmd.Wait()` is called from exactly one place.
- [ ] Read path takes a context and wakes the condition variable on
      cancel via `context.AfterFunc`.
- [ ] Version counter on the output store, checked under the lock by
      readers before sleeping on the cond.
- [ ] Per-read response capped at 32KB.
- [ ] `io.Writer.Write` returns `(len(p), nil)` for success and
      `(n<len(p), err)` for short writes. Never `(len(p), nil)` if
      not all bytes were written. Never `(0, nil)`.
- [ ] Tests use `t.Context()`, no sleeps, no `slog.SetDefault`
      mutation, no environment variable mutation.
- [ ] `go.uber.org/goleak` in `TestMain`.

### gRPC + mTLS PR

- [ ] Identity from `VerifiedChains[0][0]`, never `PeerCertificates`.
- [ ] Identity is `Subject`, never `SerialNumber`.
- [ ] Per-cert role parsing fails closed on unknown / empty role
      strings.
- [ ] EKU is enforced on both sides.
- [ ] Authorization is at the handler / interceptor, not in the
      library.
- [ ] Authorization fails closed.
- [ ] `NotFound` returned for both "no such job" and "not your job"
      to prevent enumeration.
- [ ] No EOF message in the stream proto. Closing the stream is
      end-of-stream.
- [ ] Response messages do not echo request fields.
- [ ] Integration tests use a real listener and real certs.

### CLI PR

- [ ] `start` prints raw ID, no prefix.
- [ ] `output` writes raw bytes, no line buffering.
- [ ] Errors to stderr, exit code reflects status.

## Humanization pass

After all code is written, before the user sees it:

1. Run a search for `[^\x00-\x7F]` across every file. Replace each
   hit with ASCII.
2. Run a search for marketing words. Each hit: rewrite or delete.
3. Read every comment out loud. Delete those that paraphrase the
   code below them.
4. Read every docstring. Most are not needed. Keep the ones that
   document a non-obvious invariant.
5. Read each commit message and PR description. Strip
   AI-attribution, "Generated by" lines, marketing language.
6. Ensure each commit is small and reviewable. No "implements
   everything" commits.
7. Vary error wrapping. If every wrap reads `failed to X: %w`,
   reshape some to `fmt.Errorf("open %s: %w", path, err)` or just
   `fmt.Errorf("read partitions: %w", err)`.

## Anti-patterns: do not do these

These are observed reviewer flags. Each is something a real Teleport
engineer has called out in the public corpus. If you produce any of
them you have failed.

- Output as `string` in protobuf. Always `bytes data`.
- Each reader gets its own full buffer.
- Cert Serial as identity.
- `PeerCertificates[0]` for identity. Use `VerifiedChains[0][0]`.
- Authorization decisions in the library API.
- Library entirely under `internal/`.
- `from_offset` and `follow` options on streaming.
- Graceful SIGTERM-then-SIGKILL.
- Stdin support.
- Job listing API.
- Namespace isolation on L4.
- Direct RPC handler calls in tests.
- Sleep-based synchronization in tests.
- `slog.SetDefault` mutation in `TestMain`.
- io.max device list hardcoded to `/dev/sda` or similar.
- `Pdeathsig` for child-cleanup.
- `cmd.Wait` called from more than one goroutine.
- Writing PID to `cgroup.procs` after fork.
- Closing the same channel from two goroutines.
- Holding a lock during network IO.
- Goroutine leaked on client disconnect (caught by goleak).
- Protobuf types in the library's public API.
- "Job started: <id>" prefix in CLI output.

## When you finish

Before declaring the work done, the candidate must be able to:

1. Walk through every file and explain why each line is there.
2. Defend each scope-cutting decision with a one-sentence reason.
3. Run `go test ./...` cleanly with `-race`.
4. Run the integration tests on a real Linux host (Lima VM, EC2,
   etc.) cleanly.
5. Stand up the server, run a job, stop a job, stream output, from
   another machine over the network.

If any of these fail, you are not done.

## Note on iteration

You may need more than one pass. Reviewers iterate; the agent should
expect to iterate. After each round of reviewer feedback:

- Read every comment carefully. Most have a specific concrete fix.
- Push fixes as new commits, never rewrite history of a PR under
  review.
- For each comment: if you agree, fix and reply one line. If you
  disagree, say so once with the reason and accept the next round.
- Do not write a defense longer than the comment.

Good luck. Cut, ship, defend.
