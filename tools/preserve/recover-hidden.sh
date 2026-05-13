#!/usr/bin/env bash
# Recover PR/issue/comment/review event data for repos that have since gone
# 404 by querying GH Archive on BigQuery. One query for all hidden repos,
# results split into per-repo NDJSON files.
#
# Reads hidden.tsv produced by find-hidden.sh. Output lands in:
#   tracks/systems/corpus/raw/<owner>__<repo>/   (track is "systems" — all
#   known hidden repos look like systems-track job-worker submissions; the
#   track-classification step is a TODO if other tracks appear).
#
# Per-repo files written:
#   gh-archive-events.ndjson   — one JSON object per event, sorted by time
#   RECONSTRUCTED.md           — summary stats + provenance
#
# Usage:
#   ./recover-hidden.sh                       # default 2024-2026
#   ./recover-hidden.sh 2024 2026 --dry-run
#
# Requires: bq (gcloud), jq.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
HIDDEN_TSV="${ROOT}/tools/preserve/hidden/hidden.tsv"
TRACKS_ROOT="${ROOT}/tracks"

YEAR_FROM="${1:-2024}"
YEAR_TO="${2:-$(date +%Y)}"
DRY_RUN_FLAG="${3:-}"
PROJECT="${GCP_PROJECT:?set GCP_PROJECT to your billing-enabled GCP project ID}"

if [ ! -f "${HIDDEN_TSV}" ]; then
  echo "hidden.tsv not found; run find-hidden.sh first." >&2
  exit 1
fi

# Build IN clause from hidden.tsv.
repos_clause=""
while IFS=$'\t' read -r repo status; do
  [ -z "$repo" ] && continue
  repos_clause+="'${repo}',"
done < "${HIDDEN_TSV}"
repos_clause="${repos_clause%,}"
n_repos=$(echo "${repos_clause}" | tr ',' '\n' | wc -l | tr -d ' ')
echo "==> Recovering ${n_repos} hidden repos for ${YEAR_FROM}-${YEAR_TO}"

SUFFIX_FROM="${YEAR_FROM}01"
SUFFIX_TO="${YEAR_TO}12"

# Pull the full event row as JSON so we can split per-repo locally without
# losing fields. Cast payload to STRING; jq downstream parses it.
QUERY="
SELECT
  TO_JSON_STRING(STRUCT(
    created_at AS created_at,
    type AS type,
    repo.name AS repo_name,
    actor.login AS actor_login,
    payload AS payload_raw
  )) AS row_json
FROM \`githubarchive.month.*\`
WHERE _TABLE_SUFFIX BETWEEN '${SUFFIX_FROM}' AND '${SUFFIX_TO}'
  AND repo.name IN (${repos_clause})
  AND type IN (
    'PullRequestEvent',
    'PullRequestReviewEvent',
    'PullRequestReviewCommentEvent',
    'IssueCommentEvent',
    'IssuesEvent',
    'PushEvent',
    'CreateEvent',
    'DeleteEvent',
    'ForkEvent'
  )
ORDER BY repo.name, created_at
"

echo "==> Cost estimate"
bq --project_id="${PROJECT}" query --use_legacy_sql=false --dry_run "${QUERY}" 2>&1 | tail -1

if [ "${DRY_RUN_FLAG}" = "--dry-run" ]; then
  echo "Dry run only. Exiting."
  exit 0
fi

RAW="${ROOT}/tools/preserve/hidden/gh-archive-raw.ndjson"
echo "==> Running query (results to ${RAW})"
bq --project_id="${PROJECT}" query --use_legacy_sql=false --max_rows=1000000 \
  --format=json --quiet "${QUERY}" \
  | jq -c '.[].row_json | fromjson' \
  > "${RAW}"

n_rows=$(wc -l < "${RAW}" | tr -d ' ')
echo "    ${n_rows} event rows recovered"

# Split per-repo. Track classification: default systems; can be refined.
echo "==> Splitting per-repo"
while IFS=$'\t' read -r repo status; do
  [ -z "$repo" ] && continue
  owner="${repo%%/*}"; name="${repo##*/}"
  track="systems"
  outdir="${TRACKS_ROOT}/${track}/corpus/raw/${owner}__${name}"
  mkdir -p "${outdir}"

  jq -c --arg r "${repo}" 'select(.repo_name == $r)' "${RAW}" \
    > "${outdir}/gh-archive-events.ndjson"

  count=$(wc -l < "${outdir}/gh-archive-events.ndjson" | tr -d ' ')
  earliest=$(head -1 "${outdir}/gh-archive-events.ndjson" | jq -r '.created_at // "-"')
  latest=$(tail -1 "${outdir}/gh-archive-events.ndjson" | jq -r '.created_at // "-"')

  # Per-event-type counts.
  declare -A type_counts
  while read -r t; do
    [ -z "$t" ] && continue
    type_counts[$t]=$(( ${type_counts[$t]:-0} + 1 ))
  done < <(jq -r '.type' "${outdir}/gh-archive-events.ndjson")

  type_lines=""
  for k in "${!type_counts[@]}"; do
    type_lines+="- ${k}: ${type_counts[$k]}"$'\n'
  done

  cat > "${outdir}/RECONSTRUCTED.md" <<EOF
# ${repo} (recovered from GH Archive)

This repo is no longer publicly accessible on GitHub (HTTP ${status}).
The data here was reconstructed from \`githubarchive.month.*\` on BigQuery,
covering ${YEAR_FROM}-${YEAR_TO}.

Wayback Machine has zero snapshots of this repo, so GH Archive is the
only surviving record of its PRs and comments.

## Stats

- Events recovered: ${count}
- Earliest: ${earliest}
- Latest: ${latest}

### By event type
${type_lines}

## Files

- \`gh-archive-events.ndjson\` — one JSON object per event, sorted by created_at.
  Each row has: created_at, type, repo, actor, payload_raw.
  payload_raw is the full event payload as recorded by GH Archive; for
  PR/issue events it includes the PR/issue body, for comment events it
  includes the comment body, etc.

## How to read it

\`\`\`sh
# Pull just the comment text:
jq -r 'select(.type | test("Comment")) | .payload_raw | fromjson | .comment.body' \\
  gh-archive-events.ndjson | less

# Pull PR titles:
jq -r 'select(.type == "PullRequestEvent") | .payload_raw | fromjson | .pull_request.title' \\
  gh-archive-events.ndjson | sort -u
\`\`\`

## Provenance

- Captured: $(date -u +%FT%TZ)
- Source: \`githubarchive.month.*\` (BigQuery)
- Method: \`tools/preserve/recover-hidden.sh\`
EOF

  unset type_counts
  echo "    ${repo}: ${count} events"
done < "${HIDDEN_TSV}"

echo ""
echo "==> Done. Hidden repos now have gh-archive-events.ndjson + RECONSTRUCTED.md"
