# Auth storage: where the session lives

The fullstack challenge is a login-protected directory browser, so
where the session token lives is the first design question every
reviewer probes. Candidates default to JWT-in-localStorage; reviewers
push toward HTTP-only cookies and a server-side session.

## localStorage vs cookies for the access token

> **r0mant** (zship): "Are there any downsides to using local storage
> for the access token vs cookies?"

> **alex-kovoy** (zship): "I think that JWT makes it a bit more
> complicated as it will be harder to implement a proper logout and
> enforcing 1h session TTL seems as a bad UX. I would revisit this
> design choice as part of the challenge is to implement a logout
> mechanism."

Expected answer: HTTP-only `Secure` `SameSite=Strict` cookie for the
session ID, with the actual session state stored server-side. JWT in
localStorage cannot be invalidated server-side, which makes logout
fake and a stolen token usable until expiry. The challenge explicitly
requires a logout endpoint, so the storage choice has to support real
revocation.

## Idle timeout

> **alex-kovoy** (ibeckermayer): "Lets ignore the idle timeout. I do
> not think that it's required. It would be strange to logout a user
> after any period of inactivity for this app."

Expected answer: drop idle timeout. A fixed-duration session is fine.
This is the fullstack equivalent of the systems track's "drop
graceful shutdown" canonical scope-cut — reviewers are asking
candidates to remove features the brief doesn't require.

## CSRF mitigation when using cookies

Reviewers expect that if the session is cookie-based, the candidate
defends against CSRF. The `ibeckermayer` and `tobocop` threads both
include CSRF-token review rounds. A double-submit cookie or a
`SameSite=Strict` cookie with an explicit Origin check is sufficient;
a bare cookie with no CSRF defense is not.
