#!/usr/bin/env bash
# Recover events for hidden repos by streaming raw .json.gz hour files from
# gharchive.org and filtering for the 29 hidden repos. Free (no BigQuery).
#
# Reads hidden.tsv (list of repos) and event_months.tsv (YYYYMM activity
# histogram from the cheap BigQuery metadata pass). Only those months are
# downloaded.
#
# Per-repo output lands in:
#   tracks/systems/corpus/raw/<owner>__<repo>/gh-archive-events.ndjson
#
# Usage:
#   ./recover-from-gharchive.sh           # uses event_months.tsv
#   PARALLEL=8 ./recover-from-gharchive.sh
#
# Requires: curl, gunzip, jq.
set -uo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
HIDDEN="${ROOT}/tools/preserve/hidden/hidden.tsv"
MONTHS_FILE="${ROOT}/tools/preserve/hidden/event_months.tsv"
WORK="${ROOT}/tools/preserve/hidden/gharchive-work"
PARALLEL="${PARALLEL:-8}"

[ -f "${HIDDEN}" ] || { echo "Missing ${HIDDEN}"; exit 1; }
[ -f "${MONTHS_FILE}" ] || { echo "Missing ${MONTHS_FILE} (run metadata BigQuery pass first)"; exit 1; }

mkdir -p "${WORK}"
LOOKUP="${WORK}/repo_lookup.json"

# Build a jq object {repo_name: true, ...} for O(1) hit-testing in the stream.
awk -F'\t' '{print $1}' "${HIDDEN}" \
  | jq -R . \
  | jq -s 'map({(.):true}) | add' \
  > "${LOOKUP}"
n_repos=$(jq 'length' "${LOOKUP}")
echo "==> Filter built for ${n_repos} hidden repos"

# Unique months to download.
months=$(awk -F'\t' '{print $1}' "${MONTHS_FILE}" | sort -u)
n_months=$(echo "${months}" | wc -l | tr -d ' ')
echo "==> ${n_months} months to scan: $(echo $months | tr '\n' ' ')"

# Days in a month (Bash arithmetic, no `cal` portability issues).
days_in() {
  local y="$1" m="$2"
  case "$m" in
    01|03|05|07|08|10|12) echo 31 ;;
    04|06|09|11) echo 30 ;;
    02)
      if (( (y % 4 == 0 && y % 100 != 0) || y % 400 == 0 )); then echo 29; else echo 28; fi
      ;;
  esac
}

# Generate the full URL list.
URLS="${WORK}/urls.txt"
: > "${URLS}"
for ym in ${months}; do
  y="${ym:0:4}"; m="${ym:4:2}"
  d_max=$(days_in "$y" "$m")
  for d in $(seq 1 "$d_max"); do
    dd=$(printf '%02d' "$d")
    for h in $(seq 0 23); do
      echo "https://data.gharchive.org/${y}-${m}-${dd}-${h}.json.gz" >> "${URLS}"
    done
  done
done
n_urls=$(wc -l < "${URLS}" | tr -d ' ')
echo "==> ${n_urls} hour files to fetch (concurrency=${PARALLEL})"

# Allow today's incomplete month — drop URLs at/after current UTC hour.
now_y=$(date -u +%Y); now_m=$(date -u +%m); now_d=$(date -u +%d); now_h=$(date -u +%H)
now_h="${now_h#0}"  # strip leading zero so arithmetic works
cutoff="${now_y}-${now_m}-${now_d}-${now_h}"
grep -v "/${cutoff}\.\|/${now_y}-${now_m}-${now_d}-2[0-3]\.\|/${now_y}-${now_m}-${now_d}-1[$(((now_h>=10 ? now_h%10 : 0)))-9]\." "${URLS}" > "${URLS}.filt" 2>/dev/null || cp "${URLS}" "${URLS}.filt"
mv "${URLS}.filt" "${URLS}"
n_urls=$(wc -l < "${URLS}" | tr -d ' ')
echo "    ${n_urls} after dropping future hours"

# Per-worker filter function. curl streams .json.gz; gunzip decompresses;
# jq filters lines whose repo.name is in the lookup. Failures (404, partial)
# are silently ignored — gharchive occasionally drops an hour.
#
# Each worker writes to its own temp file under SHARDS/ — concatenating
# at the end. This avoids interleaved writes from parallel workers when
# a single jq line exceeds the kernel's write-atomicity threshold (4KB
# on macOS), which happens for PR review comment events (~15-20KB each).
ALL="${WORK}/all-events.ndjson"
SHARDS="${WORK}/shards"
rm -rf "${SHARDS}"
mkdir -p "${SHARDS}"
: > "${ALL}"

worker() {
  local url="$1"
  # Derive a stable shard name from the URL filename.
  local fname="${url##*/}"
  local out="${SHARDS}/${fname%.json.gz}.ndjson"
  curl -fsS --max-time 90 --retry 2 --retry-delay 1 "${url}" \
    | gunzip 2>/dev/null \
    | jq -c --slurpfile lookup "${LOOKUP}" \
        'select($lookup[0][.repo.name] // false)' 2>/dev/null \
    > "${out}"
  # Drop empty shards immediately to keep the dir lean.
  [ -s "${out}" ] || rm -f "${out}"
}
export -f worker
export LOOKUP SHARDS

echo "==> Streaming (this will take a while)"
< "${URLS}" \
  xargs -P "${PARALLEL}" -I{} bash -c 'worker "$@"' _ {} \
  2> "${WORK}/errors.log"

# Concatenate all non-empty shards into the final ndjson.
find "${SHARDS}" -name '*.ndjson' -type f -print0 \
  | xargs -0 cat \
  > "${ALL}"

n_matches=$(wc -l < "${ALL}" | tr -d ' ')
echo "==> ${n_matches} matching event rows recovered"

# Split per-repo and write RECONSTRUCTED.md.
echo "==> Splitting per-repo"
while IFS=$'\t' read -r repo status; do
  [ -z "$repo" ] && continue
  owner="${repo%%/*}"; name="${repo##*/}"
  outdir="${ROOT}/tracks/systems/corpus/raw/${owner}__${name}"
  mkdir -p "${outdir}"

  jq -c --arg r "${repo}" 'select(.repo.name == $r)' "${ALL}" \
    > "${outdir}/gh-archive-events.ndjson"

  count=$(wc -l < "${outdir}/gh-archive-events.ndjson" | tr -d ' ')
  if [ "${count}" -eq 0 ]; then
    echo "    ${repo}: 0 events (no recovery possible)"
    continue
  fi
  earliest=$(head -1 "${outdir}/gh-archive-events.ndjson" | jq -r '.created_at // "-"')
  latest=$(tail -1 "${outdir}/gh-archive-events.ndjson" | jq -r '.created_at // "-"')

  # Per-event-type counts.
  type_block=$(jq -r '.type' "${outdir}/gh-archive-events.ndjson" \
    | sort | uniq -c \
    | awk '{printf "- %s: %d\n", $2, $1}')

  cat > "${outdir}/RECONSTRUCTED.md" <<EOF
# ${repo} (recovered from gharchive.org)

This repo is no longer publicly accessible on GitHub (HTTP ${status}).
The data here was reconstructed by streaming raw \`.json.gz\` hour files
from data.gharchive.org and filtering for this repo's events. The Wayback
Machine has zero snapshots, so GH Archive is the only surviving record.

## Stats

- Events recovered: ${count}
- Earliest: ${earliest}
- Latest: ${latest}

### By event type
${type_block}

## Files

\`gh-archive-events.ndjson\` — one JSON object per event, ordered by
download arrival (re-sort with \`jq -s 'sort_by(.created_at)'\` if needed).
Each row is the full GH Archive event including \`payload\` with PR/issue/
comment bodies.

## How to read

\`\`\`sh
# Just the comment text:
jq -r 'select(.type | test("Comment")) | .payload.comment.body' \\
  gh-archive-events.ndjson | less

# PR titles seen:
jq -r 'select(.type == "PullRequestEvent") | .payload.pull_request.title' \\
  gh-archive-events.ndjson | sort -u
\`\`\`

## Provenance

- Captured: $(date -u +%FT%TZ)
- Source: data.gharchive.org (raw hour files)
- Method: \`tools/preserve/recover-from-gharchive.sh\`
EOF
  echo "    ${repo}: ${count} events"
done < "${HIDDEN}"

echo ""
echo "==> Done."
