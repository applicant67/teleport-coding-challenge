# rafaelvanoni/jobservice (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404_missing).
The data here was reconstructed by streaming raw `.json.gz` hour files
from data.gharchive.org and filtering for this repo's events. The Wayback
Machine has zero snapshots, so GH Archive is the only surviving record.

## Stats

- Events recovered: 63
- Earliest: 2025-05-05T23:00:20Z
- Latest: 2025-05-06T17:08:22Z

### By event type
- CreateEvent: 3
- MemberEvent: 4
- PullRequestEvent: 1
- PullRequestReviewCommentEvent: 32
- PullRequestReviewEvent: 21
- PushEvent: 2

## Files

`gh-archive-events.ndjson` — one JSON object per event, ordered by
download arrival (re-sort with `jq -s 'sort_by(.created_at)'` if needed).
Each row is the full GH Archive event including `payload` with PR/issue/
comment bodies.

## How to read

```sh
# Just the comment text:
jq -r 'select(.type | test("Comment")) | .payload.comment.body' \
  gh-archive-events.ndjson | less

# PR titles seen:
jq -r 'select(.type == "PullRequestEvent") | .payload.pull_request.title' \
  gh-archive-events.ndjson | sort -u
```

## Provenance

- Captured: 2026-05-12T04:46:06Z
- Source: data.gharchive.org (raw hour files)
- Method: `tools/preserve/recover-from-gharchive.sh`
