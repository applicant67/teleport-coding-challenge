# Teleport Coding Challenge Helper

Community notes for the Teleport take-home coding challenges. The goal
is to leave better signposts than the ones available today and save
time for the folks who come next.

## What this repo gives you

The core deliverable is a compiled, indexed record of past public
submissions, paired with the reviewer feedback those submissions
attracted. Every captured repo is listed in
[`SUBMISSIONS.md`](SUBMISSIONS.md) (96 submissions across 9 tracks)
with a two-grader AI quality score per row — see
[`SCORING_RUBRIC.md`](SCORING_RUBRIC.md) for what the scores mean.

A small but important subset of the corpus are **known-strong
submissions** — challenge repos authored by people who were later
hired and who are now Teleport reviewers themselves. These repos were
deleted or made private after their authors joined Teleport. Authors
include espadolini, GavinFrazar, lxea, rosstimothy, and tigrato.

For five of the five we have at least the gh-archive event record
(commit messages, PR review threads, file-change lists) reconstructed
by streaming raw hourly archives from `data.gharchive.org`. The PR
review threads from when these candidates were themselves under review
are the closest thing to a "submissions that worked" signal that
exists publicly.

For one — `GavinFrazar/teleport-challenge` — we additionally recovered
the **full source code** from the Software Heritage Archive (which
preserved a snapshot before deletion). The other four were not in
SWH, so we have only the gh-archive event record. Look for the ✓
hired marker in `SUBMISSIONS.md`; the reconstructed event records
live under `tracks/systems/corpus/raw/<author>__<repo>__hired/`, and
the recovered GavinFrazar source lives at
`mirrors_extracted/systems/GavinFrazar__teleport-challenge/`.

The rest of the repo is supporting material: per-track reviewer
feedback distilled by theme, capture/recovery tooling, humanization
guidance for AI-assisted submissions, and a master prompt for
end-to-end systems-track submissions.

## Typical usage

1. Point your AI agent at this repo and tell it which challenge you
   are taking on. The official spec is in
   [`challenges/<track>/`](challenges/); per-track guidance is in
   [`tracks/<track>/`](tracks/).

2. Ask the agent to read the reviewer-feedback corpus and the
   highest-scoring captured submissions for your track, summarize the
   strengths and weaknesses, and propose an implementation plan
   before writing code.

3. While implementing, keep the scope-cutting lesson front of mind
   (see below) and apply
   [`docs/HUMANIZATION.md`](docs/HUMANIZATION.md) to strip AI tells
   from code, comments, and design docs.

4. Before submitting, have the agent review the finished work against
   the per-track reviewer-feedback themes — not just generic code
   quality — and fix anything that pattern-matches to a known
   reviewer objection.

## The single most important lesson from the corpus

**Do not be creative. Do not add features.** Especially on the
systems track, reviewers want one solution path with the minimum
feature set the brief asks for, and nothing else. The captured PR
threads make this unambiguous: candidates who added scope — graceful
shutdown, list-jobs endpoints, interactive shells, refresh tokens,
multi-user ACLs, idle timeouts, elaborate frontend state — got more
negative reviewer comments and worse outcomes. Candidates who cut
scope aggressively got hired more often. The signal is clearest on
systems, but the same pattern shows up in fullstack and
security-automation. See
[`tracks/systems/corpus/reviewer-feedback/scope-cutting.md`](tracks/systems/corpus/reviewer-feedback/scope-cutting.md)
for the canonical reviewer quotes.

## What to expect from the process

Two things from the captured data and candidate accounts are worth
flagging up front:

- **The challenge takes the time Teleport says it does.** Their own
  guidance is honest about this — plan the calendar before you
  commit. The systems-track submissions in the corpus typically span
  one to three weeks of evening-and-weekend work.

- **Reviewer communication is sparser than the recruiter pitch
  implies — treat this as a culture-fit signal.** The candidate
  Slack channel is quiet in practice; reviewer responses on PRs are
  sporadic and terse; the kickoff call is short and focused on
  ground rules rather than on the company, the team, or the work.
  If you value strong, energetic communication — fast feedback,
  peers who reply same-day, conversations about culture and product
  alongside the technical work — Teleport may not be the right
  environment for you, and the challenge process is the first place
  you will see this. The process is professional, not hostile; it is
  also low-touch. Plan to drive your own iteration without
  back-and-forth between rounds.

- **Teleport is openly anti-AI — plan your humanization around it.**
  In at least one captured interview, the interviewer stated
  explicitly that "we do not ship AI-generated code" and spoke
  mockingly of past submissions that contained `—` (em-dash)
  characters — a classic AI-prose tell. If you use AI assistance
  during the challenge, treat humanization as a required step rather
  than a final polish: read every line of code, every comment, and
  every paragraph of the design doc, and rewrite anything that reads
  like model output. Comments in particular tend to leak the
  model's voice. See [`docs/HUMANIZATION.md`](docs/HUMANIZATION.md)
  for the full checklist. And do not forget to turn off
  "Co-Authored-By: Claude" commit-message attribution before the
  first push :)

## Layout

```
SUBMISSIONS.md          master index of every captured submission,
                        scored by two independent grader passes
SCORING_RUBRIC.md       what the AI scores in SUBMISSIONS.md mean
challenges/             snapshot of the official Teleport challenge
                        briefs from gravitational/careers
docs/                   cross-track guidance (humanization, process)
skills/                 cross-track Claude Code skills
tracks/
  systems/              Job Worker Service — full corpus + skills
  systems-intern/       intern variant of systems
  fullstack/            secure remote directory browser
  sre/                  SRE / infra automation
  ai-ml/                CloudTrail anomaly detection
  mobile/               mobile client
  product/              product written exercise
  security-automation/  security workflow automation
  support/              support engineering exercise
mirrors/                bare-clone mirrors of every candidate repo,
                        organized by track. See mirrors/INDEX.md.
mirrors_extracted/      checked-out working tree of each mirror's
                        default branch — for easy reading and grep.
tools/preserve/         capture / recover / mirror tooling
```

## Track coverage at a glance

The **systems** track is the most developed — themed reviewer-feedback
corpus, master prompt, four specialized review skills. Other tracks:

- **fullstack** — themed corpus built from reviewer quotes by
  russjones, r0mant, alex-kovoy, awly.
- **security-automation** — thin but real themed corpus from russjones
  reviews on two candidates.
- **sre** — one repo with two reviewer comments; too thin for a
  corpus.
- **systems-intern** — shares the senior systems corpus; no
  intern-specific public threads found.
- **ai-ml / mobile / product** — tracks added in mid-to-late 2025; no
  public reviewer threads located.
- **support** — structurally non-public (live debug session, no
  GitHub artifact).

If you have public links or reviewer feedback for any track, please
contribute.

## Where to start

1. Open [`SUBMISSIONS.md`](SUBMISSIONS.md) and read the rubric in
   [`SCORING_RUBRIC.md`](SCORING_RUBRIC.md) before drawing
   conclusions from any single score. The top of each track's table
   is the most useful study material.
2. Find your track in `tracks/<track>/` and read its `README.md`.
   The per-theme curated quotes live under
   `tracks/<track>/corpus/reviewer-feedback/`; the frequency-ranked
   firehose of every captured reviewer comment is in
   [`docs/REVIEWER_FEEDBACK_DISTILLED.md`](docs/REVIEWER_FEEDBACK_DISTILLED.md).
3. Read [`docs/HUMANIZATION.md`](docs/HUMANIZATION.md) if you are
   using AI assistance. Teleport is openly skeptical of AI-generated
   submissions and you need a process to strip the tells.
4. Read [`docs/PROCESS_NOTES.md`](docs/PROCESS_NOTES.md) for recruiter
   cadence, review timing, and what to do when a reviewer goes
   silent.
5. If you're an AI agent driving an end-to-end submission for the
   systems track, the long single-shot prompt is at
   [`tracks/systems/prompts/MASTER_PROMPT.md`](tracks/systems/prompts/MASTER_PROMPT.md).

## Status

First public draft, May 2026. The systems-track corpus is rich; other
tracks are sparse. Contributions welcome.

## License

MIT. See `LICENSE`.
