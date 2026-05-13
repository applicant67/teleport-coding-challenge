# clydotron/job_worker_service

- Captured: 2026-05-11T22:24:56Z
- Track: systems
- Upstream: https://github.com/clydotron/job_worker_service

## Files

- `repo.json` - repository metadata at capture time
- `pulls.json` - all pull requests (open + closed), 1 entries
- `issues.json` - all issues (open + closed), 1 entries
- `pr-review-comments.json` - line-anchored PR review comments, 39 entries
- `issue-comments.json` - issue + PR conversation comments, 0 entries
- `reviews.json` - per-PR review summaries (APPROVED / CHANGES_REQUESTED / COMMENTED)

## Why this exists

Forks preserve code but not the PR conversation. If this upstream
repo goes private (e.g. the candidate is hired and asked to take it
down), the reviewer feedback in these JSON files becomes the only
publicly-accessible record of those interactions.
