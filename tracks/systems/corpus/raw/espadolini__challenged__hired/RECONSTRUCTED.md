# espadolini/challenged (recovered from gharchive.org)

Repo is no longer accessible on GitHub (HTTP 404). Reconstructed by streaming
raw `.json.gz` hour files from data.gharchive.org for the 2021-09 → 2021-11
window.

## Why this entry matters

espadolini is a **current Teleport reviewer**. This is his first challenge
submission; he later created a second (slimbox) the following month. The
recovered events include 34 PR review comments from Teleport employees,
making this one of the richer corpora in the collection. This sweep also
surfaced previously unknown reviewer handles: **knisbet** and **xacrimon**.

## Stats

- Events recovered: 70
- Window: 2021-10-07 → 2021-10-16 (~9 days, classic challenge submission pattern)

### By event type
- PullRequestReviewCommentEvent: 34
- PullRequestReviewEvent: 17
- PushEvent: 6
- PullRequestEvent: 6
- CreateEvent: 5
- DeleteEvent: 2

### By actor
- espadolini: 37 (the candidate)
- knisbet: 13 (Teleport reviewer — newly surfaced)
- russjones: 10 (Teleport co-founder)
- xacrimon: 5 (Teleport reviewer — newly surfaced)
- fspmarshall: 4 (Teleport reviewer)
- r0mant: 1 (Teleport engineer)

## How to read

```sh
# Just the reviewer feedback:
jq -r 'select(.type == "PullRequestReviewCommentEvent" and .actor.login != "espadolini") |
       "[\(.created_at[:10]) \(.actor.login)] \(.payload.comment.body)"' \
  gh-archive-events.ndjson | less
```

## Provenance

- Source: data.gharchive.org (raw hour files, 2021-09-01 → 2021-11-30)
- Method: focused mini-scan (similar to `tools/preserve/recover-from-gharchive.sh`)
- Captured: 2026-05-12
