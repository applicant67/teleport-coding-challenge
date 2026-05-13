# Scoring rubric (AI quality assessment)

The "AI score" column in [`SUBMISSIONS.md`](SUBMISSIONS.md) is a coarse
quality signal from two independent grader passes. Reading these scores
without reading this file is a mistake — the scores are no more reliable
than the rubric they apply, and the rubric is a community reconstruction,
not Teleport's internal rubric.

## Where this rubric comes from

Teleport published an 11-criterion rubric on their engineering blog
(see [coding-challenge guidance][1] and [hiring-process post][2]) circa
2020. That rubric scored each criterion `+1` or `-1`. We know from
candidate-side reports that Teleport later moved to `+1`/`-2` weighting
(missing or buggy criteria penalized harder than passes were rewarded),
and that current reviewers are no longer using a numeric scoring sheet at
all — recent feedback we've recovered is structured as a bullet list of
*what bothered the reviewer* with no per-criterion summing. The bullet
list still maps onto the same conceptual themes, so the public 11-point
rubric remains the best public scaffolding for what reviewers care about.

Our rubric is the union of:

1. The published 11 criteria (weighted `+1`/`-2`, matching the more
   recent reported weighting).
2. Nine recurring themes we extracted from the PR review threads in
   `tracks/systems/corpus/reviewer-feedback/`. These give finer-grained
   shading for the parts reviewers actually drill into.

[1]: https://goteleport.com/blog/coding-challenge/
[2]: https://goteleport.com/blog/our-engineering-hiring-process/

## The 11 published criteria

Each is scored `+1` (clear pass), `0` (mixed / not applicable), `-2`
(clear miss). Score range: `-22` to `+11`.

1. **Clear and modular structure** — code split into purpose-named
   packages; lifecycle, transport, and security separated.
2. **Communicated progress during the interview** — visible iteration
   in commits, design doc, PR comments; not a single drop.
3. **Reproducible build** — pinned `go.mod`/`go.sum`, `Makefile`, vendored
   protoc or `buf`, deterministic image.
4. **Clear README** — usage, certs, commands, platform caveats.
5. **Design document captures the key trade-offs** — covers lifecycle,
   auth, output retention, cancellation, scope cuts. Not a full spec.
6. **No obvious data races or deadlocks** — `go test -race` passes; no
   lock-across-network calls; explicit shutdown signaling.
7. **Tests cover the key components** — at minimum auth, lifecycle, and
   one real round-trip integration test. No sleep-based sync.
8. **Clear error handling and reporting** — domain errors map to
   structured types or gRPC status codes; failures preserve detail.
9. **Working according to the specification** — start / stop / status /
   stream output, mTLS, identity, plus job containment for the L4/L5
   variants.
10. **Handles and applies feedback** — second-round commits visibly
    incorporate reviewer feedback rather than dismissing it.
11. **Secure client-server communication** — TLS ≥ 1.3 floor, mTLS,
    identity derived from `VerifiedChains` (not `PeerCertificates` or
    serial number), no info leakage in error responses.

## The nine reviewer-feedback themes (weighted modifier, ±3)

These are not separate add-on points; they re-weight the relevant
published criterion when reviewers spent significant comment volume on
the theme. A pattern reviewers ask about across 5 submissions counts more
than one that surfaces once. Theme details are in
`tracks/systems/corpus/reviewer-feedback/`.

| Theme | Relates to published criterion | Common loss patterns |
|---|---|---|
| output-streaming | #9 spec | dropping bytes near job exit; blocking writers on slow readers; lost-wakeup on `sync.Cond` without version counter |
| status-lifecycle | #9 spec | one boolean `running` flag conflating started/stopped/failed; exit code lost; synchronous stop returning before reaper reaps |
| auth-identity | #11 security | identity from cert Subject string instead of `VerifiedChains`; serial-number-based identity; admin override leak |
| scope-cutting | #1 modular + #10 feedback | shipping interactive shell, multi-user ACLs, SIGTERM-then-SIGKILL escalation when not asked |
| process-containment | #9 spec | pgid-only termination (escapes via setpgid); cgroup-membership race at fork time; `pdeathsig` missing |
| testing | #7 tests | sleep-based sync; mocked listeners; race between `started` flag and subscriber attach |
| api-shape | #1 modular | protobuf types leaking through library API; `func([]byte) error` callback instead of `io.Writer`; non-atomic status snapshot |
| concurrency-and-locking | #6 races | holding mutex across network calls; double-stop crash; goroutine leaks on subscriber disconnect; RWMutex used where ordinary Mutex suffices |
| reproducibility-and-tooling | #3 build | unpinned `protoc`; missing `buf.lock`; committed `.pb.go` with stale ABI; no `golangci-lint` config |

## Hired vs not hired ≠ score

A high score does not predict the hire decision and vice versa. We score
the artifact, not the candidate. The reviewer corpus contains repos with
+9 scores whose candidates were not hired and +3 scores whose candidates
later became reviewers. The "hired" column in `SUBMISSIONS.md` is
provided as separate signal so readers can see the gap.

The most common gap pattern: hired candidates wrote a tighter scope and
fewer features; rejected candidates often had wider scope and richer
features but more reviewer-flagged correctness issues. If you take one
thing from the corpus, take that.

## How the two graders were run

1. Each repo is read in batches of 12 by two grader subagents (Claude
   Opus 4.7 and Claude Sonnet 4.6 — different model families to surface
   prompt-induced bias).
2. Each grader reads: the repo's README, design doc if present, the
   `pr-review-comments.json` (or the gh-archive PR comments for
   recovered repos), and a sample of commit messages.
3. Each grader emits a JSON record per repo: `{score: int, criteria:
   {[1..11]: int}, themes: {[theme]: int}, notes: str}`.
4. Scores are reconciled in `tracks/systems/corpus/submissions/SCORES.json`.
   When the two graders disagree by more than 2 points the row is
   flagged ⚠️ in the master table.

## Caveats

- The grader sees the *archived* state of the repo, not the live
  reviewer thread. Repos with private review feedback (or no captured
  review comments) are scored against the public artifact only. A
  thin-feedback repo with a clean artifact can score well even if its
  reviewers in fact found problems we never see.
- Recovered (formerly-public) repos are scored on the gh-archive event
  stream and the captured `RECONSTRUCTED.md`, not on full source. Those
  scores are noted ⚙ in the master table to flag the lower signal.
- This is intentionally a coarse signal. Treat scores ≤0 and ≥+7 as
  meaningful; treat the middle band as noise.
