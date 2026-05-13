# lxea/teleport_challenge (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404). Reconstructed
by streaming raw `.json.gz` hour files from data.gharchive.org for the
2021-09 → 2021-10 window.

## Why this entry matters

lxea is a **current Teleport reviewer**. This repo appears to be his own
challenge submission from when he applied — public for ~9 days then taken
private. The recovered events include 38 PR review comments by Teleport
employees, making this one of the richer corpora in the collection.

## Stats

- Events recovered: 131
- Window: 2021-09-27 → 2021-10-06 (~9 days, classic challenge submission pattern)

### By event type
- PushEvent: 40
- PullRequestReviewEvent: 43
- PullRequestReviewCommentEvent: 38
- PullRequestEvent: 5
- CreateEvent: 2
- DeleteEvent: 2
- PublicEvent: 1

### By actor
- lxea: 74 (the candidate)
- russjones: 28 (Teleport co-founder, canonical reviewer)
- timothyb89: 16 (likely a Teleport reviewer alt handle)
- tcsc: 11 (Teleport reviewer)
- quinqu: 2

## How to read

```sh
# Just the reviewer feedback:
jq -r 'select(.type == "PullRequestReviewCommentEvent" and .actor.login != "lxea") |
       "[\(.created_at[:10]) \(.actor.login)] \(.payload.comment.body)"' \
  gh-archive-events.ndjson | less
```

## Provenance

- Source: data.gharchive.org (raw hour files, 2021-09-01 → 2021-10-31)
- Method: focused mini-scan (similar to `tools/preserve/recover-from-gharchive.sh`)
- Captured: 2026-05-12
