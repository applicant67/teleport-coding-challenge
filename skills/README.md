# Claude Code skills

This directory holds **cross-track** Claude Code skills - ones that
apply regardless of which Teleport challenge you're submitting.

Track-specific skills live under `tracks/<track>/skills/`. The
systems track currently has four:

- `tracks/systems/skills/teleport-scope-disciplinarian/`
- `tracks/systems/skills/teleport-reviewer-mindset/`
- `tracks/systems/skills/teleport-design-doc-reviewer/`
- `tracks/systems/skills/teleport-cgroup-reviewer/`

## Cross-track skills here

- **teleport-humanizer** - strips AI tells from a finished
  submission (any track, any format).

## How to use

Either add this repo as a Claude Code plugin, or copy individual
skill directories into `~/.claude/skills/`.

## Suggested order (systems track)

1. Before writing the design doc: read
   `tracks/systems/corpus/reviewer-feedback/` so you know what
   reviewers ask.
2. After writing the design doc: run `teleport-design-doc-reviewer`.
3. After implementing cgroup primitives (L5):
   `teleport-cgroup-reviewer`.
4. Before opening the first PR: `teleport-scope-disciplinarian`.
5. Before the final submission round: `teleport-reviewer-mindset`
   and `teleport-humanizer`.

## Suggested order (other tracks)

The track-specific skills are not yet built for non-systems tracks.
For now: run `teleport-humanizer` on any written artifact before
submitting, and apply the cross-track guidance in
`docs/PROCESS_NOTES.md`.
