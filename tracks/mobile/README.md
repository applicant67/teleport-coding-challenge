# Mobile track

The mobile track asks you to build a mobile client (iOS / Android /
cross-platform - varies by year) that integrates with a Teleport
service or backend.

## Status of this directory

**Triaged, zero public reviewer quotes.** The mobile JD was only added
to `gravitational/careers` in December 2025 (PR #145), with a
design-doc-focused update in April 2026 (#154). Fragment-search and
repo-name search located no public candidate submissions for this
track.

The mobile challenge shares the filesystem-API + auth design with
fullstack, so the [`../fullstack/`](../fullstack/) reviewer quotes
(especially `alex-kovoy` on URL-encoded state, breadcrumbs, session
storage) probably carry over conceptually.

If you have public submission links or reviewer comment threads for
this track, please contribute them.

## What to do before submitting

1. Read the official spec from `gravitational/careers` carefully.
2. Borrow the systems-track *process*: scope-cutting discipline,
   design-doc-first, hermetic tests, mTLS done right.
3. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md) if
   you used AI assistance.
4. Apply [`../../docs/PROCESS_NOTES.md`](../../docs/PROCESS_NOTES.md)
   for recruiter / cadence expectations.

## Cross-track principles that apply

- Cut scope to spec. A submission that does the small thing well beats
  one that does the big thing with rough edges.
- The design doc is the most leveraged artifact. Write it first as a
  separate PR.
- mTLS identity from verified chains, not raw peer certificates.
- Same error code for "not found" and "not authorized."
- Hermetic, fast tests with no real-network or device-specific
  dependencies.
