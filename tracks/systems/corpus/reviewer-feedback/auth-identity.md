# Authentication, authorization, and identity

What field of the cert is identity, where authorization lives, and the
information-leak traps in error responses.

## Identity source: Subject, not Serial

> **tigrato** (joshuarubin): "Is there another way where we do not
> need to use the certificate serial number? If the client certificate
> needs to be re-issued because it's about to expire, will we need to
> keep the same serial number? Are there other fields we can use
> instead of relying on a field that's expected to be unique per
> certificate?"

> **tigrato** (mcampo84): "How do you plan to identify the client? Are
> there simpler ways besides using serials? If the client renews the
> certificate, the serial will change and he won't be able to continue
> to use the application"

Expected answer: identity comes from the Subject (CN, SAN, or OU), not
the Serial. Serials rotate on cert reissue; Subject is stable.

## Read from VerifiedChains, not PeerCertificates

> **espadolini** (Zephan92): "I'd prefer using
> `ConnectionState.VerifiedChains[0][0]` rather than PeerCertificates
> so you're safeguarded against accidental misconfigurations that let
> any client cert through"

Expected answer: `peer.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0]`.
If the slice is empty, refuse the request. `PeerCertificates` is set
even when the server is misconfigured to accept unverified clients.

## Do not trust client-supplied identity fields

> **rosstimothy** (RichyHBM): "This should not be needed to convey the
> true identity of the user."

> **rosstimothy** (RichyHBM): "This is not a secure way to do
> authorization. Any user can impersonate any other user without any
> verifiable proof of their identity."

> **tigrato** (RichyHBM): "why do you need the username?"

Expected answer: the server pulls identity from the verified cert
chain via the `peer` package. The client never tells the server who
they are in a request field.

## Authorization at the server, not the library

> **zmb3** (razzam21): "This is fine, though I wonder if you
> considered any alternatives? By adding owner to the public API of
> the library, you've let your authorization scheme dictate the API of
> the library. Ideally we'd want the authorization scheme to be
> something that's applied at the server, not in the library."

> **Tener** (GevorgGal): "Separating authentication from
> authorization code would be helpful in this regard."

Expected answer: library accepts an opaque owner string in metadata
and stores it. The gRPC server extracts cert identity in an
interceptor and rejects unauthorized calls at the handler before
calling into the library. The library does not make policy decisions.

## Authorization interceptor

> **tigrato** (kkloberdanz): "Should we move IdentityFromContext into
> a grpc interceptor to ensure we always have them available?"

Expected answer: yes. An auth interceptor sets the identity on
context, an authorization interceptor checks it before each RPC. The
handler reads from context only.

## Fail closed

> **eriktate** (GevorgGal): "You might consider inverting the
> condition here and having your authorize() function fail closed
> instead of failing open."

Expected answer: if you cannot determine authorization, deny.
Default-deny is non-negotiable.

## Info leakage through error responses

> **smallinsky** (kkloberdanz): "This could potentially allow a user
> to guess another user's job ID. Do you see any way to address this?"

> **espadolini** (Zephan92): "This leaks the existence of jobs to
> unauthenticated clients, right?"

> **espadolini** (Zephan92): "This is still leaking whether or not
> the job exists to users other than the job owner, right?"

Expected answer: return the same error code (commonly `NotFound`) for
"no such job" and "not authorized to see this job." Do not let
attackers enumerate valid job IDs by comparing
`PermissionDenied`-vs-`NotFound`.

## Role from cert OU has edge cases

> **smallinsky** (kkloberdanz): "Right now, if the client presents a
> certificate with any other role than the expected first OU in the
> subject, or even a user cert with an empty OU, it may be able to
> bypass the intended client-role authorization and execute a job."

Expected answer: explicitly enumerate accepted roles. Refuse unknown
or empty role strings. The OU parser fails closed on anything not in
the allowlist.

## Server cert and EKU

> **rosstimothy** (MarkDHarris): "Can we enforce EKU for this
> exercise?"

> **rosstimothy** (MarkDHarris): "Are there any security concerns from
> using the same CA for clients and servers?"

> **greedy52** (MarkDHarris): "Can a client use its own cert to
> pretend to be a server?"

Expected answer: enforce the Extended Key Usage on each side
(`x509.ExtKeyUsageServerAuth` for the server's cert,
`x509.ExtKeyUsageClientAuth` for client certs). Separate CAs are
cleaner but not strictly required if EKU is enforced.

## Multi-user authorization is out of scope

> **rosstimothy** (MarkDHarris): "Agree with Steve, the allowed
> viewers is a nice concept but does add some additional scope. If you
> want to omit it and only support owner + admin that would still
> satisfy the challenge requirements."

> **greedy52** (MarkDHarris): "sharing is cool but i think it is
> over-engineered for a 'simple authorization scheme'"

Expected answer: owner-only or owner-plus-admin is enough. ACLs are
scope creep.
