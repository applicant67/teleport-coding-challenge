#!/usr/bin/env bash
set -u
ROOT="/Users/user/git/teleport_challenge_helper"

# Heuristic: infer level applied for from repo name / description.
infer_level() {
  local repo="$1" outdir="$2"
  local hay="$(echo "$repo" | tr '[:upper:]' '[:lower:]')"
  local desc=""
  if [ -f "${outdir}/repo.json" ]; then
    desc="$(jq -r '.description // ""' "${outdir}/repo.json" 2>/dev/null | tr '[:upper:]' '[:lower:]')"
  fi
  local both="${hay} ${desc}"
  if   echo "$both" | grep -qE '(\bl4\b|level[ -]4|systems[ -]?l4)'; then echo "L4"
  elif echo "$both" | grep -qE '(\bl5\b|level[ -]5|systems[ -]?l5)'; then echo "L5"
  elif echo "$both" | grep -qE '(\bl6\b|level[ -]6|systems[ -]?l6)'; then echo "L6"
  else echo "unknown"
  fi
}

while read -r track repo; do
  owner="${repo%%/*}"; name="${repo##*/}"
  mirror_dir="${ROOT}/mirrors/${track}/${owner}__${name}"
  corpus_dir="${ROOT}/tracks/${track}/corpus/raw/${owner}__${name}"
  info="${mirror_dir}/INFO.md"

  level=$(infer_level "$repo" "$corpus_dir")

  # Comment counts. Treat absent / unreachable as 0.
  pr_count=0; issue_count=0; review_comments=0; issue_comments=0; reviews_count=0
  comments_total=0
  has_comments="no"
  corpus_link="(not captured)"
  status="public"
  if [ -d "$corpus_dir" ]; then
    if [ -f "${corpus_dir}/UNREACHABLE.md" ]; then
      status="unreachable_at_capture"
    else
      pr_count=$(jq 'length' "${corpus_dir}/pulls.json" 2>/dev/null || echo 0)
      issue_count=$(jq 'length' "${corpus_dir}/issues.json" 2>/dev/null || echo 0)
      review_comments=$(jq 'length' "${corpus_dir}/pr-review-comments.json" 2>/dev/null || echo 0)
      issue_comments=$(jq 'length' "${corpus_dir}/issue-comments.json" 2>/dev/null || echo 0)
      reviews_count=$(jq '[.[].reviews[]?] | length' "${corpus_dir}/reviews.json" 2>/dev/null || echo 0)
      comments_total=$((review_comments + issue_comments))
      if [ "$comments_total" -gt 0 ]; then
        has_comments="yes"
        corpus_link="../../../tracks/${track}/corpus/raw/${owner}__${name}/"
      else
        corpus_link="../../../tracks/${track}/corpus/raw/${owner}__${name}/  (no comments)"
      fi
    fi
  fi

  cat > "$info" <<MD
# ${repo}

- Track: **${track}**
- Level applied for: ${level}
- Upstream: https://github.com/${repo}
- Mirror: \`${name}.git/\` (bare clone; \`git clone ./${name}.git\` to inspect)
- Status: ${status}

## Captured discussion

- PRs: ${pr_count}
- Issues: ${issue_count}
- PR review comments (line-anchored): ${review_comments}
- Issue/PR conversation comments: ${issue_comments}
- Review summaries: ${reviews_count}
- Has comments: **${has_comments}**

Source JSON: ${corpus_link}
MD
done < <(awk 'NF==2 && $1!~/^#/' "${ROOT}/tools/preserve/repos.txt" | sort -u)

echo "Generated $(find "${ROOT}/mirrors" -name INFO.md | wc -l) INFO.md files"
