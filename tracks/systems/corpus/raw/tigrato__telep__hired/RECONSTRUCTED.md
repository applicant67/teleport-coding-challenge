# tigrato/telep (recovered from gharchive.org)

Repo is no longer accessible on GitHub (HTTP 404). Reconstructed by streaming
raw `.json.gz` hour files from data.gharchive.org for the 2022-01 → 2022-03
window.

## Why this entry matters

tigrato is a **current Teleport reviewer**. The repo name "telep" appears to be
an abbreviation of teleport — a common pattern for hidden challenge
submissions. The recovered events include 24 PR review comments from Teleport
employees.

## Stats

- Events recovered: 63
- Window: 2022-02-14 → 2022-02-24 (~10 days, classic challenge submission pattern)

### By event type
- PullRequestReviewCommentEvent: 24
- PullRequestReviewEvent: 22
- PushEvent: 11
- PullRequestEvent: 4
- CreateEvent: 2

### By actor
- tigrato: 31 (the candidate)
- nklaassen: 17 (Teleport reviewer)
- zmb3: 9 (Teleport reviewer)
- fspmarshall: 6 (Teleport reviewer)

## How to read

```sh
# Just the reviewer feedback:
jq -r 'select(.type == "PullRequestReviewCommentEvent" and .actor.login != "tigrato") |
       "[\(.created_at[:10]) \(.actor.login)] \(.payload.comment.body)"' \
  gh-archive-events.ndjson | less
```

## Provenance

- Source: data.gharchive.org (raw hour files, 2022-01-01 → 2022-03-31)
- Method: focused mini-scan (similar to `tools/preserve/recover-from-gharchive.sh`)
- Captured: 2026-05-12
