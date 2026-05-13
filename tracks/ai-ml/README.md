# AI/ML track - CloudTrail anomaly detection

The AI/ML track asks you to build a CLI that detects anomalies in AWS
CloudTrail logs. The canonical dataset is `flaws.cloud`. Single
declared level (L4 equivalent).

## Status of this directory

**Triaged, zero public reviewer quotes.** The AI-ML challenge was
added in July 2025 (`gravitational/careers` PR #129 by r0mant).
Challenge-text-fragment search across the flaws.cloud dataset
returned no candidate Teleport submissions with reviewer activity. Candidates have either kept their work private,
or there are too few of them yet for any to surface publicly.

If you find a public reviewer comment thread for this track, please
add it to the index below.

## Public submissions to read

- `jackjduggan/teleport-ai-challenge`
- `oscarmorberg-ops/teleport-ai-ml-challenge`
- `jeffrin27/teleport-ai-challenge`
- `barlevo/teleport-ai`
- `shruti-vadlamani/teleport-ai-challenge`
- `next-josealbertoarcos/teleport-ai-challenge`

Search GitHub for `teleport-ai-challenge`, `teleport-ai-ml-challenge`,
and `flaws.cloud cloudtrail anomaly` for fresher examples.

## What to do before submitting

1. Read the official spec carefully. The flaws.cloud dataset is the
   ground truth - make sure your detector finds the anomalies it's
   meant to find.
2. Borrow the systems-track *process*: scope-cutting discipline,
   design-doc-first, hermetic tests.
3. The model / detection approach is up to you, but the reviewer is
   not impressed by complexity - they want a defensible, simple
   detector with reproducible results.
4. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md) if
   you used AI assistance.

## What we don't have yet

- Themed reviewer-feedback corpus.
- Track-specific scope-cutting checklist (e.g., "don't train your own
  embedding model when a rule-based baseline suffices").
- Track-specific reviewer-mindset skill.

Time-to-build estimate from triage: 45-75 hours.

## Likely cut-list (informed guess)

- Real-time streaming detection (offline batch is the spec).
- A web UI (CLI is the spec).
- Custom-trained models when a baseline detector suffices.
- Multi-cloud support (AWS only).

When in doubt, cut.
