# Auth0 scopes and provider versioning

The security-automation challenge asks the candidate to configure an
Auth0 tenant via Terraform. Two of the canonical review threads
(`nick-hinds/teleport_challenge`, `tedmist1/challenge-auth0`) hit the
same two issues: too many Auth0 management scopes requested, and a
stale Terraform provider version pinned.

## Least-privilege Auth0 scopes

> **russjones** (nick-hinds): "This feels like a lot more than needed.
> For example, what are you using `read:logs`, `create:connections`,
> and `update:resource_servers` for? ... Why do you need access to
> `*:actions`? Maybe I am misunderstanding something but Auth0 actions
> are not needed for this challenge?"

Expected answer: enumerate exactly the Auth0 management-API scopes
needed by the Terraform run, and no more. The challenge does not
require Auth0 Actions, so anything under `*:actions` is unjustified.
A submission that requests every read+write scope on every resource
fails the least-privilege bar even if the rest of the work is
correct.

## Terraform provider version hygiene

> **russjones** (nick-hinds): "This is a really old version of the
> Auth0 Terraform Provider. Any reason you're using such an old
> version?"

Expected answer: pin the latest stable provider release at submission
time, and document the version explicitly in the design doc. The
provider is `auth0/auth0` — older `alexkappa/auth0` is community-
maintained and not the canonical source. Pinning a stale version
without justification is treated the same way as a stale `go.mod` in
the systems track.

## Tool-stack expectations

> **russjones** (tedmist1): "However, do you mind writing a Bash
> script instead? We don't use PowerShell at Teleport."

Expected answer: even when not explicitly specified, the tool stack
should match the company's preferences. Bash over PowerShell is the
canonical example. This generalizes: pick the tool the team you're
applying to already runs in production.
