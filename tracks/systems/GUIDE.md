# A guide to the Teleport coding challenge

You are about to spend roughly 60 to 100 hours on a take-home coding
challenge. This guide tries to make those hours count.

## What the challenge actually is

Teleport's "Job Worker Service" challenge (variants: "linux process as a
service", "remote process execution") asks you to build a small system
that:

- Exposes a gRPC API to start, stop, query status of, and stream output
  from arbitrary Linux processes.
- Uses mTLS for authentication and a simple authorization scheme for
  identifying who can do what.
- For Level 5: applies cgroup v2 resource limits (memory, cpu, io) to
  every job.
- Implements a reusable library that does the heavy lifting, plus a thin
  server and CLI on top.

It is a small project. The trap is that "small" is read as "simple" and
candidates over-build. The candidates who get advanced are the ones who
write the smallest plausible version of each requirement, very carefully.

## How this repo is laid out

- [`REJECTION_LESSONS.md`](REJECTION_LESSONS.md) - the four reasons a
  recent submission was rejected. Read this first.
- [`REVIEWER_PREFERENCES.md`](REVIEWER_PREFERENCES.md) - the question
  set reviewers reliably ask. Treat as a checklist.
- [`TECHNICAL_PITFALLS.md`](TECHNICAL_PITFALLS.md) - concrete code-level
  traps with code snippets.
- [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md) - if you
  use AI assistance, the cross-track process for removing AI tells.
- [`../../docs/PROCESS_NOTES.md`](../../docs/PROCESS_NOTES.md) -
  cross-track timing, review cadence, recruiter dynamics.
- [`corpus/submissions/INDEX.md`](corpus/submissions/INDEX.md) -
  inventory of past public systems submissions and reviewer feedback.
- [`corpus/reviewer-feedback/`](corpus/reviewer-feedback/) -
  systems reviewer quotes reorganized by theme.
- [`prompts/MASTER_PROMPT.md`](prompts/MASTER_PROMPT.md) - a long
  prompt for AI agents that bakes in everything in this track.
- [`skills/`](skills/) - Claude Code skills for individual review
  passes (systems-specific; the cross-track humanizer lives at
  `../../skills/teleport-humanizer/`).

## A suggested order of operations

### Hour 0 to 4: read the corpus

- Read the official challenge spec from `teleport-careers`. There are two
  versions: `challenge-1.md` and `challenge-2.md` in the `challenges/`
  directory of the official repo. Read both even if you are only doing
  one.
- Read `REJECTION_LESSONS.md` and `REVIEWER_PREFERENCES.md` in this
  track directory.
- Pick three submissions from `corpus/submissions/INDEX.md` whose
  reviewer PRs are publicly visible. Read every reviewer comment on
  them.

You will end this phase with a strong intuition for what reviewers care
about. If you skip this phase you will discover those preferences one at a
time, after each PR round-trip.

### Hour 4 to 20: design doc

The design doc is the most leveraged work you will do. Spend disproportionate
time here.

Structure that works:

1. Scope. One paragraph naming exactly what is in and what is out. Be
   explicit about graceful shutdown, stdin, job listing, dynamic resource
   limits, retry, multi-user roles, namespace isolation. Each of these
   that you cut, name in writing.
2. Architecture. One diagram or one prose paragraph showing the layers:
   library, server, CLI. Note the package locations. Library at the
   module root or `pkg/`. Server under `cmd/workerd`. CLI under
   `cmd/workerctl`.
3. Library API. List the public types and methods. State the contract for
   each. Especially: how does a reader discover new bytes are available?
   What happens if a reader is slow? What happens if `Stop` is called
   twice? What is the exit code for a stopped job?
4. Cgroup design. Which controllers. Which file. CLONE_INTO_CGROUP for
   atomic placement. cgroup.kill for termination. Sweep on startup.
   Kernel version floor.
5. Output streaming. State the architecture explicitly: shared backing
   store, per-reader file descriptors. Note that the design avoids the
   slow-reader cross-contamination problem. State the close-vs-drain
   ordering.
6. Authentication and authorization. mTLS, identity from VerifiedChains,
   Subject (not Serial) as identity, simple owner-only or owner-plus-admin
   authorization, enforced at the server layer not the library.
7. Tests. What runs in CI without root. What requires root. What requires
   a real Linux kernel. How CI-without-root proves the unit logic; how a
   root-required test proves the integration story.
8. Out of scope. Explicit list of things you are not building. This is the
   section reviewers most often nod through. It tells them you have
   read the spec and made decisions.

Open the design doc as a separate PR before any code. Many reviewers will
not look at code PRs until the design is approved.

### Hour 20 to 60: implementation

Stage your PRs:

1. Design doc PR.
2. Cgroup primitives. Cgroup parent setup with controller verification.
   Per-job cgroup with limits. cgroup.kill. KillAndRemove with retry on
   EBUSY. Sweep.
3. Library: Job, Manager, output store, process spawn. One method per
   operation. No graceful shutdown.
4. gRPC server and mTLS. Start, Stop, Status, StreamOutput. Identity
   interceptor + authorization interceptor.
5. CLI client.

Each PR is independently reviewable. Stack them so the next is `--base` of
the previous on GitHub.

Tests in each PR cover the layer added in that PR. Resist the urge to
re-test the lower layers from a higher PR; that is reviewer noise.

### Hour 60 to 80: response cycles

Each reviewer round will produce a few questions and a few suggestions.
Address them in order. Push fixes as new commits, do not rewrite history
on a PR under review. The reviewer wants to see the diff between rounds.

For pushback: if you agree, push the fix and reply with one line. If you
disagree, say so once with the reason and accept the next round.

### Hour 80 to 100: humanization and final QA

Pass over the whole submission with `HUMANIZATION.md` in mind. Then run
the final QA checklist in [`prompts/MASTER_PROMPT.md`](prompts/MASTER_PROMPT.md).

## The hardest call: cutting features

Almost every candidate's instinct is to add things. Graceful shutdown
feels professional. Dynamic resource limits feel flexible. Stdin support
feels complete. Each of them is a reviewer red flag because each of them
expands the surface that has to be correct under concurrency, that has
to be tested, and that has to be defended in the design doc.

The right framing: the challenge tests whether you can build the small
thing precisely. The wrong framing: the challenge tests whether you can
build the big thing.

When in doubt, cut.

## A note on AI assistance

This guide makes no claims about whether you should use AI assistance.
The challenge instructions ask you to attest to authorship; whether and
how you do that is up to you.

If you do use AI assistance, `HUMANIZATION.md` describes a process that
in one case produced a submission that, on the final iteration, did not
prompt the reviewer to comment on AI authorship. Note the result is
sample-of-one; do not read it as a method.

If you do not use AI assistance, you can skip the humanization doc, but
the rest of this repo is still useful as a checklist of what reviewers
care about.

## What the reviewer is actually looking for

Stated bluntly:

- Did you cut scope to the spec?
- Is your design doc clear?
- Do your APIs survive misuse?
- Are your tests hermetic and fast?
- Did you understand the cgroup story or did you copy-paste it?
- Can you defend every line on a 60-minute call?

Aim for "yes" on all six. The honest fail modes in the rejected sample
are mostly "yes on five, no on one."

Good luck.
