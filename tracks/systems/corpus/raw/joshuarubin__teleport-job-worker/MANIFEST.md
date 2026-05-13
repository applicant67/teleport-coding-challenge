# joshuarubin/teleport-job-worker

- Captured: 2026-05-11T17:55:45Z
- Track: systems
- Upstream: https://github.com/joshuarubin/teleport-job-worker

## Files

- `repo.json` - repository metadata at capture time
- `pulls.json` - all pull requests (open + closed), 16 entries
- `issues.json` - all issues (open + closed), 16 entries
- `pr-review-comments.json` - line-anchored PR review comments, 197 entries
- `issue-comments.json` - issue + PR conversation comments, 10 entries
- `reviews.json` - per-PR review summaries (APPROVED / CHANGES_REQUESTED / COMMENTED)

## Why this exists

Forks preserve code but not the PR conversation. If this upstream
repo goes private (e.g. the candidate is hired and asked to take it
down), the reviewer feedback in these JSON files becomes the only
publicly-accessible record of those interactions.
