# Fullstack track - Secure Remote Directory Browser

The fullstack track asks you to build a secure web app that browses a
remote filesystem over mTLS or HTTPS. Go (or Node) backend, React /
TypeScript frontend. Levels L1 / L2 / L4 vary in scope (auth strength,
browse depth, file operations).

## Status of this directory

**Public reviewer quotes located.** Four older submissions carry
substantive Teleport reviewer PR threads. A themed distillation is
under [`corpus/reviewer-feedback/`](corpus/reviewer-feedback/) covering
auth storage, TLS/hashing, and scope-cutting. The submissions and
additional raw quotes are listed below as source material.

## Public submissions with reviewer threads

| Repo | Reviewers | Notes |
|------|-----------|-------|
| `tobocop/go-teleport-directory-browser` | russjones, r0mant, alex-kovoy | Go backend; deep TLS / bcrypt / HMAC thread on PR #2-3. |
| `zship/teleport-challenge` | alex-kovoy, r0mant | Frontend-focused; JWT-vs-cookie debate. |
| `ibeckermayer/teleport-interview` | alex-kovoy, awly | Gorilla session storage; CSRF mitigations; idle-timeout debate. |
| `atburke/teleport_interview` | russjones, alex-kovoy | Older Go + React; error-handling style. |

## Other public candidate submissions (no visible reviewer thread)

- `nickstaggs/teleport-takehome`
- `json-nguyen/remote-directory-browser`
- `b1tsized/teleport-takehome`
- `rsteinkeXJ/teleport_interview`

Search GitHub for forks of the official `gravitational/careers` repo
and for repos named `teleport-takehome`, `teleport-interview`,
`remote-directory-browser` for fresher examples.

## Canonical reviewer quotes

The four reviewer-threaded repos contain ~30+ substantive reviewer
quotes. A few to internalize before writing your design doc:

- **russjones on tobocop:** "The bcrypt hash is not just the hash,
  it's the salt + hash. When you do `bcrypt.CompareHashAndPassword`
  it just extracts the salt from the passed-in bcrypt hash..."
  Canonical bcrypt + HMAC anti-pattern.
- **russjones on tobocop:** "I would use `crypto/subtle.ConstantTimeCompare`
  here, otherwise you might open yourself up to timing attacks."
- **r0mant on zship:** "Are there any downsides to using local
  storage for the access token vs cookies?" Classic auth-storage
  probe.
- **alex-kovoy on zship:** "JWT makes it a bit more complicated as
  it will be harder to implement a proper logout, and enforcing 1h
  session TTL seems as a bad UX. I would revisit this design choice
  as part of the challenge is to implement a logout mechanism."
  Fullstack analog of the systems-track "drop graceful termination"
  cut.
- **alex-kovoy on ibeckermayer:** "Let's ignore the idle timeout. I
  do not think that it's required. It would be strange to logout a
  user after any period of inactivity for this app." Scope-cut.

## Recurring themes (raw, not yet a curated corpus)

1. **Auth storage:** JWT vs cookies. Reviewers push toward cookies
   with proper httpOnly + secure + SameSite, not JWT in localStorage.
2. **bcrypt + HMAC:** if you do extra-hashing-before-bcrypt, you
   must HMAC, not plain SHA. And document why.
3. **Constant-time compare:** for tokens, signatures, password
   verification.
4. **Idle timeout / session TTL:** reviewers often tell candidates
   to *remove* these, not add them.
5. **CSRF:** required if using cookies, not if using JWT in a header
   (but see #1).
6. **TLS scope:** mTLS for the gRPC backend at L4; HTTPS-only for the
   browser at lower levels.
7. **Logout mechanism:** required at L2+; design for it from the
   start.

## What to do before submitting

1. Read the official spec carefully. Each `challenges/challenge-N.md`
   in `gravitational/careers` describes a different level.
2. Open the four reviewer-threaded repos above and read every reviewer
   comment - especially the design-doc PRs.
3. Adopt the systems-track *process*: scope-cutting discipline,
   design-doc-first as a separate PR, stacked PRs for the
   implementation, hermetic tests.
4. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md)
   if you used AI assistance.

## What this helper does not yet provide

- A scope-cutting checklist analogous to the systems-track one.
- A design-doc reviewer skill.

Contributions welcome.

## Likely out-of-scope (informed guess)

- File upload beyond what the spec requires.
- Multi-user ACLs beyond owner/admin.
- Background indexing or search.
- Anything that smells like a product feature rather than a security
  primitive.

When in doubt, cut.
