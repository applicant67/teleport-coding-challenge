# Preservation tools

Forks preserve a repo's code but not its PR conversation. If a candidate
takes their submission private (e.g. after being hired and asked by the
employer to hide it), the reviewer comments that lived in the upstream
PRs become unreachable. The tools here mitigate that.

Three preservation layers, roughly in order of durability:

## 1. Active capture (most durable)

`capture.sh` pulls every public PR, issue, comment, and review summary
via the GitHub REST API and writes them as JSON under
`tracks/<track>/corpus/raw/<owner>__<repo>/`. Once these JSON files
are committed here, they're independent of the upstream's continued
existence.

```sh
./capture.sh                            # batch: every repo in repos.txt
./capture.sh systems mikewurtz/taskman  # one repo
```

What gets captured per repo:

- `repo.json` - repository metadata
- `pulls.json` - all PRs (open + closed)
- `issues.json` - all issues
- `pr-review-comments.json` - line-anchored review comments
- `issue-comments.json` - PR + issue conversation comments
- `reviews.json` - APPROVED / CHANGES_REQUESTED summaries

Requires: `gh` CLI authenticated, `jq` available.

## 2. Wayback Machine snapshots (belt-and-suspenders)

`wayback.sh` requests `web.archive.org` snapshots of each captured
repo's PR pages. Captures the rendered HTML including comment threading
and rendered Markdown. JS-heavy GitHub pages don't always render
perfectly, so this is a backstop, not a primary record.

```sh
./wayback.sh
```

Snapshots take a few minutes to become visible at
`https://web.archive.org/web/*/<url>`.

## 3. GH Archive (recovery, for repos already private/deleted)

[GH Archive](https://www.gharchive.org/) records every public GitHub
event (pushes, PR opens, PR comments, reviews, forks, deletes) and
ships them to BigQuery in hourly partitions, going back to 2011.
Crucially: **once an event is in the archive, it stays there even if
the source repo later goes private or is deleted.**

This is exactly the dataset you want when a candidate has been hired
and asked to take their submission down. The PR comment text is
preserved as it was sent.

### Searching it

Requires a Google Cloud project with billing enabled (BigQuery has a
1 TB/month free tier; these queries cost well under that).

Find every PR comment by a Teleport reviewer, all-time:

```sql
SELECT
  created_at,
  repo.name,
  actor.login,
  JSON_EXTRACT_SCALAR(payload, '$.comment.body') AS body,
  JSON_EXTRACT_SCALAR(payload, '$.comment.html_url') AS url
FROM `githubarchive.month.20*`
WHERE type IN ('PullRequestReviewCommentEvent', 'IssueCommentEvent')
  AND actor.login IN (
    'rosstimothy', 'tigrato', 'nklaassen', 'eriktate', 'espadolini',
    'codingllama', 'jimbishopp', 'sclevine', 'zmb3', 'GavinFrazar',
    'smallinsky', 'greedy52', 'awly', 'russjones', 'r0mant',
    'alex-kovoy', 'rhammonds-teleport', 'Joerger', 'Tener'
  )
ORDER BY created_at DESC
```

Find every PR comment mentioning "job worker" or "cgroup" or "mTLS"
on a repo whose name suggests a Teleport submission:

```sql
SELECT
  created_at,
  repo.name,
  actor.login,
  JSON_EXTRACT_SCALAR(payload, '$.comment.body') AS body
FROM `githubarchive.month.20*`
WHERE type IN ('PullRequestReviewCommentEvent', 'IssueCommentEvent')
  AND (
    REGEXP_CONTAINS(LOWER(repo.name), r'(teleport|jobworker|job-worker|job_worker)')
  )
  AND (
    REGEXP_CONTAINS(LOWER(JSON_EXTRACT_SCALAR(payload, '$.comment.body')),
                    r'(cgroup|mtls|job worker|cleartext)')
  )
ORDER BY created_at DESC
LIMIT 1000
```

Find repos that have since gone private (no recent events but had
events in the past):

```sql
WITH recent AS (
  SELECT DISTINCT repo.name
  FROM `githubarchive.month.202604`
  WHERE REGEXP_CONTAINS(LOWER(repo.name), r'(teleport|jobworker)')
),
historical AS (
  SELECT DISTINCT repo.name
  FROM `githubarchive.year.2024`
  WHERE REGEXP_CONTAINS(LOWER(repo.name), r'(teleport|jobworker)')
)
SELECT repo.name AS name
FROM historical
WHERE name NOT IN (SELECT name FROM recent)
```

That last query is the "find candidates who got hired and went private"
search you suggested. Worth running periodically.

### Cost note

`githubarchive.month.*` is roughly 100-300 GB per month. Use `_PARTITIONTIME`
or restrict the wildcard (`month.2025*`, `month.20260*`) to keep queries
under the free tier.

## The manifest

`repos.txt` lists every repo we know about, one per line:

```
<track> <owner>/<repo>
```

To add a new repo: append a line, run `./capture.sh <track> <repo>`,
commit the resulting JSON.

## Discovery and recovery scripts

### `find-hidden.sh` (BigQuery)

Sweep GH Archive on BigQuery for repos that Teleport reviewers have
commented on but that are no longer publicly accessible. Output:

- `hidden/candidates.tsv` — every comment-on repo + current GitHub status
  (`public`, `404_missing`, `403_forbidden`, etc.). Third column is the
  canonical full name after following GitHub renames.
- `hidden/candidates_public_dedup.tsv` — public rows, deduped on canonical
  name. The set of *new* public candidates not yet in `repos.txt`.
- `hidden/hidden.tsv` — the subset that is no longer publicly accessible.
- `hidden/wayback.tsv` — Wayback CDX snapshot count per hidden repo.

```sh
GCP_PROJECT=your-project ./find-hidden.sh 2024 2026
```

Cost: ~200 GB BigQuery scan (well within free tier).

### `recover-from-gharchive.sh` (free)

For repos in `hidden.tsv`, stream the raw `.json.gz` hour files from
`data.gharchive.org` and filter for the hidden repos' events. Writes
`gh-archive-events.ndjson` + `RECONSTRUCTED.md` per hidden repo under
`tracks/systems/corpus/raw/<owner>__<repo>/`.

Requires `hidden/event_months.tsv` (YYYYMM histogram) first — generated
by a small BigQuery metadata-only query, ~370 GB scan, free-tier-safe.
See script header for the exact SQL.

```sh
PARALLEL=8 ./recover-from-gharchive.sh
```

Cost: ~260 GB streamed download from gharchive.org, ~1-3h wall clock at
home bandwidth. No BigQuery cost.

### `recover-hidden.sh` (BigQuery, paid)

Alternate path to `recover-from-gharchive.sh`: one BigQuery against
`githubarchive.month.*` pulling the full `payload` column for hidden
repos. Faster but ~8 TB scan = ~$38 in BigQuery cost over free tier.
Use only if the streaming-download route is unavailable.

```sh
GCP_PROJECT=your-project ./recover-hidden.sh 2024 2026
```

### `build-info.sh` / `build-index.sh`

Regenerate per-repo `mirrors/<track>/<repo>/INFO.md` and the top-level
`mirrors/INDEX.md` from the manifest + captured JSON. Run after any
`capture.sh` or `recover-from-gharchive.sh` invocation.

```sh
./build-info.sh && ./build-index.sh
```
