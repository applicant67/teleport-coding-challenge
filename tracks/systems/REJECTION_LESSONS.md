# Why one submission was rejected, and how to avoid the same mistakes

A postmortem from a candidate who completed the challenge end-to-end and was
not advanced. The reviewers articulated four reasons. None of them were
"unclear spec" or "missed feature". Every one is a process or design call
that was avoidable.

This doc lists each reason, explains what it actually looked like in the
submission, and gives a concrete rule for the next person.

---

## 1. Tests modified global state and depended on host hardware

**What got flagged:** The `log/slog` test setup mutated the process-wide
default logger via `slog.SetDefault(...)` in `TestMain`. The `io.max`
device-discovery test asserted on the specific block devices present on the
host running the tests, so it passed on a developer laptop and failed in CI
or another reviewer's VM.

**Why this matters more than it sounds:** A reviewer with twenty submissions
in their queue will clone the repo, run `go test ./...`, and expect green.
If your test suite passes only on your machine, you have unfortunately just
told the reviewer that you do not have an internal model of test isolation.

**Rule:**

- Never modify `slog.SetDefault`, `os.Setenv`, `flag.CommandLine`, or any
  other process-wide singleton from tests. Inject a logger or use a per-test
  context.
- Any test that touches the filesystem, the network, block devices, or
  cgroups must either skip cleanly when the resource is unavailable (and you
  must verify the skip path) or stub the resource. The simplest pattern: a
  small interface in production code, a fake in tests, plus one optional
  integration test that runs only when the real resource is present.
- For device-discovery and `/proc` parsers specifically: write the
  filtering and formatting logic against an `io.Reader` so you can feed it
  hand-crafted `/proc/partitions` content. Keep a separate, gated test that
  exercises the real path.

---

## 2. Adding "graceful shutdown" expanded scope and introduced a race

**What got flagged:** The submission supported a coordinated `Shutdown(ctx)`
on the job manager that waited for in-flight jobs to terminate. Two
problems:

1. Graceful shutdown was not in the challenge requirements. It added a
   sizeable amount of state machine surface (a `StatusStopping` state, a
   shutdown gate on the start path, a wait group, a context contract for
   each callsite) without delivering anything the spec asked for.
2. The shutdown gate had a race: it was possible for a `Start` call to
   begin after `Shutdown` had been observed, because the check and the
   registration of the new job were not atomic.

**Why this matters:** Reviewers repeatedly tell candidates to cut scope.
The most cited reviewer line across past submissions is some variant of
"This is not a requirement, consider cutting it." Every feature you add
that the spec does not require is a feature you must implement correctly
under concurrency, document, and explain. It is a tax with no income.

**Rule:**

- Read the official challenge spec. List the explicit requirements. Build
  exactly that list. Anything you want to add beyond it goes in a note in
  the design doc called "Out of scope" with a one-line reason.
- Common things to cut: graceful shutdown, SIGTERM-then-SIGKILL escalation,
  job listing, dynamic resource limits (hardcode them), follow / offset
  options on streaming (always follow, always from the start), retry
  logic, rate limiting, audit logging, stdin support, namespace isolation,
  command allowlists, per-job timeouts, multi-user authorization beyond
  owner + admin.
- If you must add something not in the spec, the bar is the same as a
  spec-required feature: race-free, tested, and documented. Most candidates
  do not have time for that bar on extras.

---

## 3. The API split a single logical operation across two functions

**What got flagged:** The cgroup and output-store APIs each had a pair of
calls that the user had to make in the right order to get correct behavior.
For the output store, there was a separate `SetDone()` and `Close()`. For
the cgroup, kill and remove were independent calls a caller could get
wrong. Splitting like this creates two failure modes:

1. The caller forgets one of the calls and you get a quiet bug (output
   readers hang forever, cgroups leak).
2. The two calls race each other and you get an explicit bug.

**Why this matters:** A library API is judged by what the worst-case caller
can do. Reviewers ask "can a user call this in the wrong order, or twice,
and break something?" If the answer is yes, the API is wrong.

**Rule:**

- Default to one method per operation. If your "operation" naturally has
  multiple steps, hide that inside the implementation. `Close()` should
  mark the stream done and release the writer fd. `KillAndRemove()` should
  kill members and remove the directory.
- Methods must be idempotent or document loudly that they are not. Calling
  `Close()` twice should not panic, should not double-close a channel,
  should not leak.
- The simplest way to know your API has this problem: write the canonical
  use site. If you find yourself writing `x.Foo(); x.Bar()` and there is
  no plausible callsite that wants to call just `Foo` without `Bar`,
  collapse them.

---

## 4. The output writer violated the `io.Writer` contract

**What got flagged:** The output store implemented `io.Writer`. One branch
of `Write` could return `len(p), nil` without actually accepting the bytes
into the underlying buffer (specifically, a closed-store fast path that
returned success without writing). This is a silent data-loss bug. The
`io.Writer` contract says a non-error return must mean the full slice was
written.

**Why this matters:** Reviewers are systems engineers. `io.Writer` is the
most heavily relied-on interface in Go. Violating its contract is the kind
of bug that, once observed, casts doubt on every other piece of code in
the submission.

**Rule:**

- If you implement `io.Writer`, the only valid return shapes are
  `(len(p), nil)` after writing all of `p`, or `(n < len(p), err)` for
  short writes with a non-nil error.
- The same goes for `io.Reader`: return `(n, nil)` or `(n, err)`. Do not
  return `(0, nil)` to mean EOF. Do not invent your own sentinel.
- After you write a `Write` or `Read` method, re-read the standard library
  doc for the interface and walk every branch of your implementation
  against it. This takes five minutes and catches the bug.

---

## Meta-lesson: design doc problems compound

The rejection reasons above are code-level, but the underlying cause was a
design that was a little too ambitious and a little too coupled. If the
design doc had said "one method per operation, no graceful shutdown, no
host-dependent tests" before any code was written, none of these failures
would have happened.

Spend disproportionate time on the design doc. The reviewers do.
