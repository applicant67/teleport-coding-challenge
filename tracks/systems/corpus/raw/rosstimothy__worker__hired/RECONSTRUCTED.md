# rosstimothy/worker (recovered from gharchive.org)

Repo is no longer accessible on GitHub (HTTP 404). Reconstructed by streaming
raw `.json.gz` hour files from data.gharchive.org for the 2021-05 → 2021-07
window.

## Why this entry matters

rosstimothy is a **current Teleport reviewer**. This appears to be his own
challenge submission ("worker" is the canonical name for the Job Worker
challenge task). The recovered events include 62 PR review comments from
Teleport employees — one of the richer corpora in the collection.

## Stats

- Events recovered: 154
- Window: 2021-06-17 → 2021-06-29 (~12 days, classic challenge submission pattern)

### By event type
- PullRequestReviewCommentEvent: 62
- PullRequestReviewEvent: 59
- PushEvent: 13
- CreateEvent: 8
- PullRequestEvent: 6
- MemberEvent: 2
- IssueCommentEvent: 2
- DeleteEvent: 2

### By actor
- rosstimothy: 88 (the candidate)
- nklaassen: 29 (Teleport reviewer)
- russjones: 22 (Teleport co-founder, canonical reviewer)
- r0mant: 15 (Teleport engineer)

## How to read

```sh
# Just the reviewer feedback:
jq -r 'select(.type == "PullRequestReviewCommentEvent" and .actor.login != "rosstimothy") |
       "[\(.created_at[:10]) \(.actor.login)] \(.payload.comment.body)"' \
  gh-archive-events.ndjson | less
```

## Provenance

- Source: data.gharchive.org (raw hour files, 2021-05-01 → 2021-07-31)
- Method: focused mini-scan (similar to `tools/preserve/recover-from-gharchive.sh`)
- Captured: 2026-05-12
