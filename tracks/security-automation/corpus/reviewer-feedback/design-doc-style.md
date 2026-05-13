# Design-doc style for the security-automation challenge

The security-automation challenge is more design-doc-heavy than the
systems challenge: reviewers spend most of their comment volume on
the prose, not on the Terraform. Two patterns recur across the
captured threads.

## Prose density: 3–5 sentences per area, not bullets

> **russjones** (tedmist1): "Can you provide 3-5 sentences on how you
> plan to configure the Auth0 tenant?"

Expected answer: each design-doc section should be a short paragraph
with explicit trade-offs and a justification. A bullet list of
configuration knobs without prose is treated as incomplete. The
reviewer is looking for the candidate's reasoning, not their config
output.

## Operational controls, not just tool config

> **russjones** (tedmist1): "How are changes reviewed and approved?
> Keep in mind, having an internal repository like this is a common
> pattern at companies. How do you ensure random employees are not
> making and pushing access changes?"

Expected answer: the design doc has to cover the *process* around
the tool, not only the tool config. Required-reviewer rules on the
repo, CODEOWNERS, branch protection, and an explicit audit trail are
all expected to appear. A submission that documents only the
Terraform itself, without the organizational controls around it,
fails the bar for an enterprise security automation role.
