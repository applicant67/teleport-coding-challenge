# TLS, password hashing, and timing-safe comparisons

These three primitives come up on every fullstack thread because the
challenge requires HTTPS, password-based auth, and a session-token
check on every request. Candidates routinely get the cryptographic
details wrong in subtle ways that reviewers flag explicitly.

## bcrypt: hash already contains the salt

> **russjones** (tobocop): "The bcrypt hash is not just the hash, it's
> the salt + hash. When you do `bcrypt.CompareHashAndPassword` it just
> extracts the salt from the passed in bcrypt hash..."

Expected answer: do not store a salt column separately and do not HMAC
the bcrypt output with a server-side pepper unless the design doc
justifies it. `bcrypt.GenerateFromPassword` and
`bcrypt.CompareHashAndPassword` are the full API surface. Candidates
who combine bcrypt with their own HMAC step usually break the
verification path because they hash the bcrypt output before
comparison.

## Constant-time comparison for tokens

> **russjones** (tobocop): "I would use
> https://pkg.go.dev/crypto/subtle#ConstantTimeCompare here, otherwise
> you might open yourself up to timing attacks."

Expected answer: any equality check on a secret (session token, API
key, HMAC tag) goes through `subtle.ConstantTimeCompare`. `==` on a
byte slice or string leaks length and prefix-match timing.

## HTTPS is required; the cert source is up to the candidate

The brief requires the server to terminate TLS. Self-signed certs
with a documented `curl --cacert` invocation in the README are
accepted across all four fullstack threads. mkcert is also accepted.
What is *not* accepted: plain HTTP with a note that "TLS would be
added in production."
