#!/usr/bin/env bash
# Request Wayback Machine snapshots for every PR URL in the captured corpus.
#
# This complements capture.sh: capture.sh stores the machine-readable JSON,
# wayback.sh asks archive.org to snapshot the rendered HTML page for each PR.
# If a PR is deleted upstream, archive.org may still serve it.
#
# Usage:
#   ./wayback.sh          # request snapshots for every captured repo
#
# This script is rate-limit-conscious: it sleeps between requests.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

# Find every captured pulls.json and snapshot each PR's HTML page.
find "${ROOT}/tracks" -name pulls.json -print0 | while IFS= read -r -d '' f; do
  dir="$(dirname "$f")"
  repo_url="$(jq -r '.[0].html_url // empty' "$f" | sed 's|/pull/.*||')"
  [ -z "$repo_url" ] && continue
  echo "==> ${repo_url}"

  # Snapshot the repo root and each PR page.
  for url in "$repo_url" $(jq -r '.[].html_url' "$f"); do
    echo "    ${url}"
    # Wayback's Save Page Now endpoint. Anonymous use; respect rate limit.
    curl -fsS -o /dev/null "https://web.archive.org/save/${url}" \
      || echo "    (snapshot request failed; continuing)"
    sleep 8
  done
done

echo ""
echo "Wayback snapshot requests complete. Snapshots take a few minutes"
echo "to appear at https://web.archive.org/web/*/<url>"
