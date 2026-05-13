# parmjassal/assignment-job-worker (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404_missing).
The data here was reconstructed by streaming raw `.json.gz` hour files
from data.gharchive.org and filtering for this repo's events. The Wayback
Machine has zero snapshots, so GH Archive is the only surviving record.

## Stats

- Events recovered: 61
- Earliest: 2025-04-07T16:03:51Z
- Latest: 2025-04-04T14:01:04Z

### By event type
- CreateEvent: 1
- MemberEvent: 3
- PullRequestEvent: 1
- PullRequestReviewCommentEvent: 25
- PullRequestReviewEvent: 20
- PushEvent: 11

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
