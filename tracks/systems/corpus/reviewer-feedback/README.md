# Reviewer feedback by theme

The files in this directory pull verbatim reviewer comments out of the
public PRs in `../submissions/INDEX.md` and group them by what the
reviewer was actually pushing on. Quotes are attributed by reviewer
handle and source repo.

The point of organizing this way is that any given submission will only
hit a fraction of the themes, but reviewers come back to the same
themes across submissions. If you read the corpus by-submission you see
each reviewer once; if you read it by-theme you see the same concern
asked five different ways.

## Themes

- [output-streaming.md](output-streaming.md) - how new bytes are
  signaled, how multiple readers share a backing store, slow-reader
  handling, final-write-vs-done races.
- [status-lifecycle.md](status-lifecycle.md) - distinguishing stopped
  vs failed vs running, exit codes for killed jobs, async-vs-sync stop.
- [auth-identity.md](auth-identity.md) - identity from cert, Subject
  vs Serial, VerifiedChains vs PeerCertificates, info leakage through
  error responses.
- [scope-cutting.md](scope-cutting.md) - features reviewers tell
  candidates to remove. Almost always the same ones.
- [process-containment.md](process-containment.md) - cgroup membership
  races, kill semantics, child processes, pdeathsig.
- [testing.md](testing.md) - hermetic tests, t.Context, real listeners,
  no sleeps.
- [api-shape.md](api-shape.md) - library importability, io.Reader vs
  callback, protobuf types in public API, snapshot atomicity.
- [concurrency-and-locking.md](concurrency-and-locking.md) - lock
  scope, write-after-close, double-stop, RWMutex vs Mutex, goroutine
  leaks on disconnect.
- [reproducibility-and-tooling.md](reproducibility-and-tooling.md) -
  buf, pinned protoc, committed .pb.go, golangci-lint config.

## How to use

If you are about to write the design doc, skim every file once. For
each theme, decide your answer in one or two sentences before reviewers
ask. Drop your answer into the relevant section of the design doc. By
the time the PR goes up, every recurring concern has a preemptive
answer somewhere in the doc.

If you are responding to a specific comment, find the theme it touches
and read the rest of the comments under it. The reviewer's next
question is almost always also in there.
