# Pragmatic notes on the challenge process

Things one candidate noticed that the official instructions do not tell
you. None of this is mission-critical; it is timing and expectation
calibration.

## Review cadence is irregular

Reviewers are Teleport engineers reviewing on the side of their actual
job. A typical cadence in the observed sample:

- Design doc PR: first reviewer comment within 1 to 3 business days.
- Follow-up rounds: roughly one round per day, but with gaps.
- A weekend or US holiday will silently extend everything.

Plan for a multi-week back-and-forth, not a multi-day one. Do not panic if
you go 48 hours without a response. Do not bump the PR before then.

## The slack channel is essentially silent

If you are added to a candidates-and-reviewers Slack channel, do not
expect organic conversation there. In the observed sample it served as a
read-only announcement channel; questions posted by candidates went
unanswered for days or never. Direct messages or PR comments produced
faster responses.

If you have a genuine question, the PR comment thread is the place. The
recruiter is the secondary channel.

## What "level" really means

The challenge has two variants commonly referenced as Level 4 and Level 5.
The same source repo contains both. Differences in the observed sample:

- L4: no cgroups, no resource limits, no process isolation; otherwise the
  same library / gRPC / mTLS shape.
- L5: cgroup v2 with memory, cpu, io limits applied per job. The challenge
  spec for L5 was rewritten in early 2025 to make the cgroup requirement
  explicit; some older candidates' submissions predate that and target a
  spec that was vaguer about resource limits.

Reviewers will not always tell you which level you are targeting. The
challenge files in the official `teleport-careers` repo do; read them
both before committing.

If you have any choice in the matter, take the level your background fits
best. Reviewers are not impressed by an L5 submission with a sketchy
cgroup story when an L4 submission would have nailed the basics.

## Search the public corpus before you write

The most cost-effective hour you can spend is reading prior public
submissions and their reviewer comments. Public PR threads on candidates'
forks are visible if you search GitHub for forks of the official challenge
repos or for repos with names containing "job-worker", "teleport-challenge",
"teleport-exercise", "jobworker", "job-worker-service". The corpus in
`tracks/systems/corpus/submissions/INDEX.md` here is a starting point; the freshest
context will always be on GitHub.

Read at least three threads. You will find that the same five or six
questions are asked of every candidate.

## When you get a hard review comment

Engineers reviewing on the side of their day job sometimes write blunt
comments. The right response is the same as in a regular code review:

- Acknowledge.
- If you agree, push the fix with a short note.
- If you disagree, say so once, give the reason, and accept the next round
  of feedback gracefully.
- Do not write a defense longer than the comment.

In the observed sample, candidates who responded with short concrete
updates moved faster than candidates who wrote prose.

## Time investment

A realistic time budget for a careful L5 submission:

- Reading the spec, the public corpus, and the reviewer-feedback notes:
  6 to 10 hours.
- Design doc: 10 to 20 hours of writing, more of thinking.
- Implementation including tests: 30 to 50 hours.
- Humanization, packaging, response cycles: 10 to 20 hours.

Anyone telling you this is a weekend project is wrong, or is targeting a
much lower bar than the reviewers actually hold candidates to.

## Stay in scope, then ship

The submission that gets accepted is not the most polished one. It is the
one that meets the spec exactly, has no race conditions in the small surface
it does cover, and answers the predictable reviewer questions in the design
doc.

Cut, ship, then accept the offer.
