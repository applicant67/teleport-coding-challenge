# SRE track

The SRE track asks you to build a small infrastructure / observability
artifact (Kubernetes operator, deployment automation, monitoring
pipeline - exact scope varies by year).

## Status of this directory

**Triaged twice, very thin public corpus.** A second pass in 2026-05
found exactly one repo with Teleport reviewer comments
(`Chili-Man/teleport-sre-challenge`, with only 2 reviewer comments
total). The SRE challenge was substantially refactored in late-2024
and early-2025 (PRs #117, #119, #150 in `gravitational/careers`);
most current candidate submissions appear to be private.

If you find a public reviewer comment thread for this track, please
add it to the index below.

## Public submissions

- `Chili-Man/teleport-sre-challenge` - jimbishopp + sclevine, 2
  comments. sclevine links to an internal `gravitational/careers`
  PR as a hint.
- `mrkooll/teleport-sre-challenge`
- `calvinbui/teleport-sre-challenge`
- `adalton/teleport-sre-challenge`
- `atburke/teleport-sre-challenge`
- `paulkarayan/rust-k8s-manager` - copied the challenge text into a
  README; no reviewer PRs.

Search GitHub for `teleport-sre-challenge` and `teleport sre take home`
for fresher examples.

## What to do before submitting

1. Read the official spec carefully. Note the level / variant you're
   targeting.
2. Borrow the systems-track *process*: scope-cutting discipline,
   design-doc-first, stacked PRs, hermetic tests. The track-specific
   technical questions differ, but the review style is consistent
   across tracks (rosstimothy, codingllama, and others sometimes
   review across tracks).
3. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md) if
   you used AI assistance.

## What we don't have yet

- Themed reviewer-feedback corpus.
- Track-specific scope-cutting checklist.
- Track-specific reviewer-mindset skill.

Time-to-build estimate from triage: ~19 hours, but this depends
heavily on the specific variant assigned.

## Cross-track guidance that still applies

- Cut scope. Reviewers prefer a smaller correct thing.
- Write the design doc first, as a separate PR.
- mTLS identity from `VerifiedChains[0][0]`, not `PeerCertificates`.
- Same gRPC error code for "not found" and "not authorized" to avoid
  ID enumeration.
- Hermetic tests, no sleeps, race detector enabled.
