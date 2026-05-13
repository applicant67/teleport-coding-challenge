#!/usr/bin/env bash
set -u
ROOT="/Users/user/git/teleport_challenge_helper"
OUT="${ROOT}/mirrors/INDEX.md"

tracks=$(awk 'NF==2 && $1!~/^#/' "${ROOT}/tools/preserve/repos.txt" | awk '{print $1}' | sort -u)

{
cat <<'HDR'
# Submission Mirror Index

Per-track mirror clones (bare `.git` directories) of every Teleport
coding-challenge candidate repo we've discovered, plus their captured
PR/issue/comment data and (where the upstream is gone) GH Archive
event reconstructions.

The mirrors here are byte-for-byte preservations of the public state at
capture time. The accompanying `INFO.md` per repo states the track, the
level applied for (where derivable), and a count of captured discussion.

Click through to `tracks/<track>/corpus/raw/<owner>__<repo>/` for the
raw JSON or NDJSON.

---

HDR

for track in $tracks; do
  echo "## ${track}"
  echo ""
  echo "| Repo | Level | PRs | Comments | Status |"
  echo "|------|-------|-----|----------|--------|"
  for d in "${ROOT}/mirrors/${track}/"*/; do
    [ -d "$d" ] || continue
    info="${d}INFO.md"
    [ -f "$info" ] || continue
    repo=$(grep -m1 '^# ' "$info" | sed 's/^# //')
    owner="${repo%%/*}"; name="${repo##*/}"
    level=$(grep -m1 '^- Level' "$info" | sed 's/.*: //')
    prs=$(grep -m1 '^- PRs:' "$info" | sed 's/.*: //')
    rc=$(grep -m1 'PR review comments' "$info" | sed 's/.*: //')
    ic=$(grep -m1 'conversation comments' "$info" | sed 's/.*: //')
    total_comments=$(( ${rc:-0} + ${ic:-0} ))
    status=$(grep -m1 '^- Status:' "$info" | sed 's/.*: //')
    relpath="${track}/${owner}__${name}/INFO.md"
    echo "| [${repo}](${relpath}) | ${level} | ${prs} | ${total_comments} | ${status} |"
  done
  echo ""
done

# Hidden repos section
echo "---"
echo ""
echo "## Hidden repos (404 - recovered from GH Archive)"
echo ""
echo "These repos went private or were deleted after capture. Reviewer"
echo "comments and PR context were reconstructed from \`data.gharchive.org\`"
echo "raw event files. Mirrors are not possible (no upstream code) — only"
echo "the recovered NDJSON exists."
echo ""
echo "| Repo | Recovered events | RECONSTRUCTED.md |"
echo "|------|-----------------|-----------------|"
while IFS=$'\t' read -r repo status; do
  owner="${repo%%/*}"; name="${repo##*/}"
  ndjson="${ROOT}/tracks/systems/corpus/raw/${owner}__${name}/gh-archive-events.ndjson"
  if [ -f "$ndjson" ]; then
    n=$(wc -l < "$ndjson" | tr -d ' ')
  else
    n="(recovery pending)"
  fi
  echo "| ${repo} | ${n} | [link](../tracks/systems/corpus/raw/${owner}__${name}/RECONSTRUCTED.md) |"
done < "${ROOT}/tools/preserve/hidden/hidden.tsv"
echo ""

cat <<'FTR'
---

## Notes on level inference

For repos whose name or description mentions L4/L5/L6, the level is
auto-tagged. Others show `unknown` — could be filled in by:

- Reading the candidate's README for level hints
- Cross-referencing recruiter handoff (if known)
- Matching against challenge document version

## Re-generating this index

The index is built from per-repo `INFO.md` files. Regenerate with:

```sh
tools/preserve/build-index.sh
```

(Or `bash /tmp/gen_index.sh` while editing.)
FTR
} > "${OUT}"
echo "Wrote ${OUT} ($(wc -l < "${OUT}") lines)"
