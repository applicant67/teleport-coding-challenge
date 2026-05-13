# njayp/jobber (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404_missing).
The data here was reconstructed by streaming raw `.json.gz` hour files
from data.gharchive.org and filtering for this repo's events. The Wayback
Machine has zero snapshots, so GH Archive is the only surviving record.

## Stats

- Events recovered: 133
- Earliest: 2024-06-12T09:43:24Z
- Latest: 2024-06-14T04:47:19Z

### By event type
- CreateEvent: 4
- DeleteEvent: 2
- MemberEvent: 2
- PullRequestEvent: 4
- PullRequestReviewCommentEvent: 52
- PullRequestReviewEvent: 37
- PushEvent: 32

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

- Captured: 2026-05-12T04:46:05Z
- Source: data.gharchive.org (raw hour files)
- Method: `tools/preserve/recover-from-gharchive.sh`
