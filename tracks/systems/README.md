# Systems track - Job Worker Service

The systems track is Teleport's "Job Worker Service" challenge. You build
a small Go service that exposes a gRPC API to start, stop, query, and
stream output from arbitrary Linux processes, authenticated over mTLS.
Level 5 adds per-job cgroup v2 limits (memory, cpu, io).

This is the track with the most public corpus. Almost every artifact in
this directory is derived from real reviewer feedback on real
submissions.

## What's here

- [`GUIDE.md`](GUIDE.md) - the top-level roadmap and time budget.
- [`REJECTION_LESSONS.md`](REJECTION_LESSONS.md) - four specific reasons
  a recent submission was rejected. Read this first.
- [`REVIEWER_PREFERENCES.md`](REVIEWER_PREFERENCES.md) - the question
  set reviewers reliably ask. Treat as a checklist.
- [`TECHNICAL_PITFALLS.md`](TECHNICAL_PITFALLS.md) - code-level traps
  with snippets.
- [`corpus/submissions/INDEX.md`](corpus/submissions/INDEX.md) -
  inventory of past public submissions and reviewer feedback.
- [`corpus/reviewer-feedback/`](corpus/reviewer-feedback/) - reviewer
  quotes reorganized by theme (output streaming, auth, cgroups, etc.).
- [`prompts/MASTER_PROMPT.md`](prompts/MASTER_PROMPT.md) - a long
  single-shot prompt for AI agents that bakes in everything here.
- [`skills/`](skills/) - Claude Code skills for individual review
  passes (scope, reviewer mindset, design doc, cgroup).

## Where the cross-track material lives

- [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md) - removing
  AI tells (applies across tracks).
- [`../../docs/PROCESS_NOTES.md`](../../docs/PROCESS_NOTES.md) -
  recruiter dynamics, review cadence (applies across tracks).
- [`../../skills/teleport-humanizer/`](../../skills/teleport-humanizer/) -
  humanization skill (cross-track).

## Suggested reading order

1. `GUIDE.md` - overall shape.
2. `REJECTION_LESSONS.md` - the four most expensive mistakes.
3. `REVIEWER_PREFERENCES.md` - the question set.
4. Three submissions from `corpus/submissions/INDEX.md` - read every
   reviewer comment.
5. `TECHNICAL_PITFALLS.md` - once you start writing code.

## Levels

L4 omits cgroups; L5 includes them. The same library / gRPC / mTLS
shape applies to both. Read the official spec in `teleport-careers`
(`challenges/challenge-1.md` and `challenge-2.md`) before committing
to a level.
