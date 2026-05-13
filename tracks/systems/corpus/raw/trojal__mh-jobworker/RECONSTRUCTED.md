# trojal/mh-jobworker (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP 404_missing).
The data here was reconstructed by streaming raw `.json.gz` hour files
from data.gharchive.org and filtering for this repo's events. The Wayback
Machine has zero snapshots, so GH Archive is the only surviving record.

## Stats

- Events recovered: 108
- Earliest: 2025-05-15T20:11:32Z
- Latest: 2025-05-16T18:55:37Z

### By event type
- CreateEvent: 5
- DeleteEvent: 1
- MemberEvent: 2
- PullRequestEvent: 6
- PullRequestReviewCommentEvent: 41
- PullRequestReviewEvent: 35
- PushEvent: 18

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

- Captured: 2026-05-12T04:46:08Z
- Source: data.gharchive.org (raw hour files)
- Method: `tools/preserve/recover-from-gharchive.sh`
