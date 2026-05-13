# Security automation track - Auth0 workflow automation

The security-automation track (an L3 take-home) asks you to design
and automate an Auth0 tenant configuration with Terraform: tenant
setup, application registration, role-scoped access, and the
operational story around reviewing and approving changes.

## Status of this directory

**Triaged, thin reviewer corpus.** A 2026-05 pass found two public
submissions with substantive `russjones` design-doc threads. A small
themed corpus is feasible to build but not yet built in this helper
repo - the submissions and quotes below are the raw material.

## Public submissions with reviewer threads

| Repo | Reviewer | Notes |
|------|----------|-------|
| `nick-hinds/teleport_challenge` | russjones | L3 Auth0 challenge; ~15-20 design-doc comments. |
| `tedmist1/challenge-auth0` | russjones | L3; Terraform + Auth0; design-doc-only review. |

## Canonical reviewer quotes

- **russjones on nick-hinds:** "This is a really old version of the
  Auth0 Terraform Provider. Any reason you're using such an old
  version?" Versioning hygiene.
- **russjones on nick-hinds:** "This feels like a lot more than
  needed. What are you using `read:logs`, `create:connections`, and
  `update:resource_servers` for? ... Why do you need access to
  `*:actions`? Auth0 actions are not needed for this challenge."
  Canonical least-privilege probe.
- **russjones on tedmist1:** "How are changes reviewed and approved?
  Keep in mind, having an internal repository like this is a common
  pattern at companies. How do you ensure random employees are not
  making and pushing access changes?" Pushes the candidate toward
  org-control thinking, not just tool config.
- **russjones on tedmist1:** "Can you provide 3-5 sentences on how
  you plan to configure the Auth0 tenant?" Design-doc style probe:
  reviewers want prose with specifics, not bullet points.
- **russjones on tedmist1:** "Do you mind writing a Bash script
  instead? We don't use PowerShell at Teleport." Concrete tool
  preference.

## Recurring themes

1. **Least-privilege Auth0 scopes.** Don't request `*:actions` or
   any scope you can't justify in writing.
2. **Provider version hygiene.** Use a current Auth0 Terraform
   provider; document why if you pin to an older one.
3. **Design-doc prose density.** "3-5 sentences" per area is the
   expected density. Avoid pure bullet lists.
4. **Tool-stack expectations.** Bash over PowerShell. Linux-friendly
   scripts.
5. **Operational story matters.** How do you review and approve
   changes? Who can push to the repo? This is part of the answer,
   not an afterthought.

## What to do before submitting

1. Read the official spec from `gravitational/careers` carefully.
2. Open the two reviewer-threaded repos above and read every
   `russjones` comment.
3. Borrow the systems-track *process*: design-doc-first as a
   separate PR.
4. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md)
   if you used AI assistance.
5. Apply [`../../docs/PROCESS_NOTES.md`](../../docs/PROCESS_NOTES.md)
   for recruiter / cadence expectations.

## Cross-track principles that apply

- Cut scope to spec. Don't add features the challenge didn't ask for.
- The design doc is the most leveraged artifact.
- Default to deny on parse errors, unknown roles, unsigned data.
- Document the human/org workflow, not just the tooling.
