#!/usr/bin/env bash
# Find Teleport-challenge candidate repos that were once public but
# may no longer be accessible (private or deleted post-hire).
#
# How it works:
#   1. Queries GH Archive (BigQuery) for every public PR/issue comment
#      left by a known Teleport reviewer in the given year range.
#   2. Collects the unique set of non-Teleport-org repos commented on.
#   3. Checks current accessibility of each via the GitHub API.
#   4. For inaccessible ones, queries Wayback CDX for any snapshots.
#
# Output is written to tools/preserve/hidden/:
#   - candidates.tsv    one row per repo: owner/name <tab> current_status
#   - hidden.tsv        the subset whose status is 404 or 403
#   - wayback.tsv       per hidden repo: snapshot count + earliest/latest
#
# Usage:
#   ./find-hidden.sh                       # default 2-year lookback
#   ./find-hidden.sh 2024 2026             # explicit inclusive year range
#   ./find-hidden.sh 2024 2026 --dry-run   # report cost, don't run
#
# Requires: bq (gcloud) authenticated, gh authenticated, jq, curl.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUTDIR="${ROOT}/tools/preserve/hidden"
mkdir -p "${OUTDIR}"

YEAR_FROM="${1:-$(($(date +%Y) - 2))}"
YEAR_TO="${2:-$(date +%Y)}"
DRY_RUN_FLAG="${3:-}"

PROJECT="${GCP_PROJECT:?set GCP_PROJECT to your billing-enabled GCP project ID}"

# Reviewer handles drawn from the corpus. Add new ones to the
# reviewer-mindset skill and here together.
REVIEWERS=(
  rosstimothy tigrato nklaassen eriktate espadolini codingllama
  jimbishopp sclevine zmb3 GavinFrazar smallinsky greedy52 awly
  russjones r0mant alex-kovoy rhammonds-teleport Joerger Tener
  jakule fspmarshall capnspacehook kontsevoy klizhentas lxea
  rcanderson23
)

reviewer_in_clause=""
for r in "${REVIEWERS[@]}"; do
  reviewer_in_clause+="'${r}',"
done
reviewer_in_clause="${reviewer_in_clause%,}"

SUFFIX_FROM="${YEAR_FROM}01"
SUFFIX_TO="${YEAR_TO}12"

QUERY="
SELECT DISTINCT repo.name AS repo
FROM \`githubarchive.month.*\`
WHERE _TABLE_SUFFIX BETWEEN '${SUFFIX_FROM}' AND '${SUFFIX_TO}'
  AND actor.login IN (${reviewer_in_clause})
  AND type IN ('PullRequestReviewCommentEvent', 'IssueCommentEvent', 'PullRequestReviewEvent')
  AND NOT STARTS_WITH(repo.name, 'gravitational/')
  AND NOT STARTS_WITH(repo.name, 'goteleport/')
  AND NOT STARTS_WITH(repo.name, 'teleport-/')
ORDER BY repo
"

echo "==> Cost estimate (dry run)"
bq --project_id="${PROJECT}" query --use_legacy_sql=false --dry_run "${QUERY}" 2>&1 | tail -1

if [ "${DRY_RUN_FLAG}" = "--dry-run" ]; then
  echo "Dry run only. Exiting."
  exit 0
fi

echo ""
echo "==> Running query for ${YEAR_FROM}-${YEAR_TO}"
bq --project_id="${PROJECT}" query --use_legacy_sql=false --max_rows=10000 \
  --format=csv --quiet "${QUERY}" \
  | tail -n +2 \
  > "${OUTDIR}/candidates_raw.tsv"

n_candidates=$(wc -l < "${OUTDIR}/candidates_raw.tsv" | tr -d ' ')
echo "    ${n_candidates} unique non-Teleport-org repos commented on"

echo ""
echo "==> Checking current accessibility via gh api"
# candidates.tsv columns: bigquery_name  status  canonical_name
# canonical_name resolves GitHub renames so we can dedupe.
: > "${OUTDIR}/candidates.tsv"
: > "${OUTDIR}/hidden.tsv"

while IFS= read -r repo; do
  [ -z "$repo" ] && continue
  status=""
  canonical="${repo}"
  if gh api "repos/${repo}" > /tmp/_status.json 2>/dev/null; then
    status="public"
    canonical=$(jq -r '.full_name // empty' /tmp/_status.json)
    [ -z "${canonical}" ] && canonical="${repo}"
  else
    # gh returns nonzero on 404 / 403 / network. Re-call with -i to see code.
    code=$(curl -sS -o /dev/null -w '%{http_code}' \
      -H "Authorization: Bearer $(gh auth token)" \
      "https://api.github.com/repos/${repo}" || echo "000")
    case "${code}" in
      404) status="404_missing" ;;
      403) status="403_forbidden" ;;
      451) status="451_takedown" ;;
      000) status="network_error" ;;
        *) status="http_${code}" ;;
    esac
  fi
  printf '%s\t%s\t%s\n' "${repo}" "${status}" "${canonical}" >> "${OUTDIR}/candidates.tsv"
  if [ "${status}" != "public" ] && [ "${status}" != "network_error" ]; then
    printf '%s\t%s\n' "${repo}" "${status}" >> "${OUTDIR}/hidden.tsv"
  fi
done < "${OUTDIR}/candidates_raw.tsv"

# Dedupe public rows on canonical name (renames produce duplicate BigQuery rows
# that all resolve to the same canonical repo).
awk -F'\t' '$2=="public" {if (!seen[$3]++) print}' "${OUTDIR}/candidates.tsv" \
  > "${OUTDIR}/candidates_public_dedup.tsv"
n_dedup=$(wc -l < "${OUTDIR}/candidates_public_dedup.tsv" | tr -d ' ')
echo "    ${n_dedup} unique public repos after rename dedup"

n_hidden=$(wc -l < "${OUTDIR}/hidden.tsv" | tr -d ' ')
echo "    ${n_hidden} repos are no longer publicly accessible"

if [ "${n_hidden}" -eq 0 ]; then
  echo ""
  echo "==> No hidden repos found. Done."
  exit 0
fi

echo ""
echo "==> Querying Wayback CDX for each hidden repo"
: > "${OUTDIR}/wayback.tsv"
printf 'repo\tstatus\tsnapshot_count\tearliest\tlatest\n' >> "${OUTDIR}/wayback.tsv"

while IFS=$'\t' read -r repo status; do
  url="https://web.archive.org/cdx/search/cdx?url=github.com/${repo}/*&output=json&limit=1000&fl=timestamp"
  resp=$(curl -fsS --max-time 30 "${url}" 2>/dev/null || echo '[]')
  count=$(echo "${resp}" | jq 'length // 0')
  # First row is the field-name header; subtract it.
  if [ "${count}" -gt 0 ]; then
    count=$((count - 1))
  fi
  if [ "${count}" -gt 0 ]; then
    earliest=$(echo "${resp}" | jq -r '.[1][0] // "-"')
    latest=$(echo "${resp}" | jq -r '.[-1][0] // "-"')
  else
    earliest="-"
    latest="-"
  fi
  printf '%s\t%s\t%s\t%s\t%s\n' "${repo}" "${status}" "${count}" "${earliest}" "${latest}" \
    >> "${OUTDIR}/wayback.tsv"
  sleep 1
done < "${OUTDIR}/hidden.tsv"

echo ""
echo "==> Done. Output in ${OUTDIR}/"
echo "    candidates.tsv  every repo + current status"
echo "    hidden.tsv      subset that's no longer public"
echo "    wayback.tsv     Wayback CDX snapshot counts per hidden repo"
