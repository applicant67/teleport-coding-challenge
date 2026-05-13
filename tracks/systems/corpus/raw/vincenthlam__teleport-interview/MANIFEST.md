# vincenthlam/teleport-interview

- Captured: 2026-05-11T22:29:12Z
- Track: systems
- Upstream: https://github.com/vincenthlam/teleport-interview

## Files

- `repo.json` - repository metadata at capture time
- `pulls.json` - all pull requests (open + closed), 1 entries
- `issues.json` - all issues (open + closed), 1 entries
- `pr-review-comments.json` - line-anchored PR review comments, 7 entries
- `issue-comments.json` - issue + PR conversation comments, 0 entries
- `reviews.json` - per-PR review summaries (APPROVED / CHANGES_REQUESTED / COMMENTED)

## Why this exists

Forks preserve code but not the PR conversation. If this upstream
repo goes private (e.g. the candidate is hired and asked to take it
down), the reviewer feedback in these JSON files becomes the only
publicly-accessible record of those interactions.
