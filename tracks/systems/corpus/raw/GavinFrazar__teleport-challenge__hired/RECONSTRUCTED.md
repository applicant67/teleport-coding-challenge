# GavinFrazar/teleport-challenge (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404).
The data here was reconstructed by streaming raw `.json.gz` hour files
from data.gharchive.org and filtering for this repo's events.

## Why this entry matters

GavinFrazar is a **current Teleport reviewer** — listed in the
canonical reviewer set used elsewhere in this corpus. The activity
window (April 2022) precedes Gavin's tenure as an active reviewer
on others' challenges, suggesting this repo is his own challenge
submission from when he applied.

All 11 captured events are Gavin's own actions (push, PR open,
self-review). No third-party reviewer comments are present — either
the review happened while the repo was private (Teleport reviewer
invited as collaborator), or it happened entirely on a separate
fork. GH Archive only captures events on the repo that was made
public.

## Stats

- Events recovered: 11
- Earliest: 2022-04-06T03:27:35Z (PublicEvent — repo went public)
- Latest: 2022-04-16T23:51:27Z
- All actor: GavinFrazar

### By event type
- PublicEvent: 1
- PushEvent: 5
- PullRequestEvent: 2
- PullRequestReviewEvent: 1
- PullRequestReviewCommentEvent: 1 (self-comment: "resolved on cli branch: c0faa382337ad7654c87b45e9d366ba9ff449633")
- IssueCommentEvent: 1

## Files

`gh-archive-events.ndjson` — one JSON object per event.

## Provenance

- Captured via `/tmp/recover-gavin-202204.sh` (focused mini-scan)
- Source: data.gharchive.org (raw hour files, 2022-04-01 → 2022-04-30)
- The main batch recovery `tools/preserve/recover-from-gharchive.sh`
  did NOT cover this repo because its `event_months.tsv` input only
  listed months where Teleport reviewers had touched the repo — and
  no reviewer activity is visible in GH Archive for this one.
