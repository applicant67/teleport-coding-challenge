#!/usr/bin/env bash
# Capture public PR + issue + review-comment data for a Teleport-challenge repo.
#
# Why: forks preserve code but lose the PR conversation. If the upstream
# repo goes private or is deleted, the reviewer comments are unreachable.
# We grab the JSON from the GitHub REST API while it's still public.
#
# Usage:
#   ./capture.sh                            # capture every repo in repos.txt
#   ./capture.sh systems mikewurtz/taskman  # capture one repo
#
# Requires: gh CLI authenticated, jq.
# Output:   tracks/<track>/corpus/raw/<owner>__<repo>/
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
MANIFEST="${ROOT}/tools/preserve/repos.txt"

capture_one() {
  local track="$1"
  local repo="$2"
  local owner="${repo%%/*}"
  local name="${repo##*/}"
  local outdir="${ROOT}/tracks/${track}/corpus/raw/${owner}__${name}"

  echo "==> ${track} / ${repo}"
  mkdir -p "${outdir}"

  # Repo metadata. If repo is private/deleted, record that fact.
  if ! gh api "repos/${repo}" > "${outdir}/repo.json" 2>"${outdir}/repo.err"; then
    echo "    repo unreachable; marking and continuing"
    {
      echo "# ${repo}"
      echo ""
      echo "Repo was unreachable at capture time: $(date -u +%FT%TZ)"
      echo ""
      echo "Error from gh api:"
      echo ""
      cat "${outdir}/repo.err"
    } > "${outdir}/UNREACHABLE.md"
    return 0
  fi
  rm -f "${outdir}/repo.err"

  # All PRs (open + closed). Paginated.
  gh api "repos/${repo}/pulls?state=all&per_page=100" --paginate \
    > "${outdir}/pulls.json"

  # All issues. GitHub's API surfaces PRs as issues for comment threading,
  # but we keep them separate because the discussion model differs.
  gh api "repos/${repo}/issues?state=all&per_page=100" --paginate \
    > "${outdir}/issues.json"

  # All PR review comments across the repo (line-anchored diff comments).
  gh api "repos/${repo}/pulls/comments?per_page=100" --paginate \
    > "${outdir}/pr-review-comments.json"

  # All issue + PR conversation comments.
  gh api "repos/${repo}/issues/comments?per_page=100" --paginate \
    > "${outdir}/issue-comments.json"

  # Per-PR review summaries (APPROVED / CHANGES_REQUESTED / COMMENTED).
  # These aren't exposed in a bulk endpoint, so iterate PR numbers.
  local pr_numbers
  pr_numbers=$(jq -r '.[].number' "${outdir}/pulls.json")
  {
    echo "["
    local first=1
    for pr in $pr_numbers; do
      if [ "$first" -eq 0 ]; then echo ","; fi
      first=0
      echo "{\"pr\": ${pr}, \"reviews\":"
      gh api "repos/${repo}/pulls/${pr}/reviews?per_page=100" --paginate
      echo "}"
    done
    echo "]"
  } > "${outdir}/reviews.json"

  # Manifest with stats so a future reader can sanity-check the capture.
  local pr_count issue_count review_count issue_comment_count
  pr_count=$(jq 'length' "${outdir}/pulls.json")
  issue_count=$(jq 'length' "${outdir}/issues.json")
  review_count=$(jq 'length' "${outdir}/pr-review-comments.json")
  issue_comment_count=$(jq 'length' "${outdir}/issue-comments.json")

  cat > "${outdir}/MANIFEST.md" <<EOF
# ${repo}

- Captured: $(date -u +%FT%TZ)
- Track: ${track}
- Upstream: https://github.com/${repo}

## Files

- \`repo.json\` - repository metadata at capture time
- \`pulls.json\` - all pull requests (open + closed), ${pr_count} entries
- \`issues.json\` - all issues (open + closed), ${issue_count} entries
- \`pr-review-comments.json\` - line-anchored PR review comments, ${review_count} entries
- \`issue-comments.json\` - issue + PR conversation comments, ${issue_comment_count} entries
- \`reviews.json\` - per-PR review summaries (APPROVED / CHANGES_REQUESTED / COMMENTED)

## Why this exists

Forks preserve code but not the PR conversation. If this upstream
repo goes private (e.g. the candidate is hired and asked to take it
down), the reviewer feedback in these JSON files becomes the only
publicly-accessible record of those interactions.
EOF
}

if [ "$#" -eq 2 ]; then
  capture_one "$1" "$2"
  exit 0
fi

# Batch mode: walk the manifest.
while IFS= read -r line; do
  # Strip comments and blank lines.
  line="${line%%#*}"
  line="$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')"
  [ -z "$line" ] && continue
  track="$(echo "$line" | awk '{print $1}')"
  repo="$(echo "$line" | awk '{print $2}')"
  capture_one "$track" "$repo"
done < "$MANIFEST"

echo ""
echo "Capture complete."
