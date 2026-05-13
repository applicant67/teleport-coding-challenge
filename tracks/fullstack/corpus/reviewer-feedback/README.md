# Fullstack reviewer feedback (themes)

Themed distillation of the public PR review threads on
`tobocop/go-teleport-directory-browser`, `zship/teleport-challenge`,
`ibeckermayer/teleport-interview`, and `atburke/teleport_interview`.

Primary reviewer voices: **russjones**, **r0mant**, **alex-kovoy**,
**awly**.

## Themes

- [`auth-storage.md`](auth-storage.md) — where the session token
  lives; logout; idle-timeout scope cut.
- [`tls-and-hashing.md`](tls-and-hashing.md) — bcrypt, timing-safe
  comparisons, self-signed certs.
- [`scope-cutting.md`](scope-cutting.md) — what reviewers ask
  candidates to remove.

This corpus is thinner than the systems track. Four candidate threads
is enough to identify recurring themes but not to support
statistical claims about how reviewers weight each one. Treat these
as canonical examples of what reviewers care about, not as
exhaustive.
