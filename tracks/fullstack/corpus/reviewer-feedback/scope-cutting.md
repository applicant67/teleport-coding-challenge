# What fullstack reviewers tell candidates to cut

The fullstack track has the same "drop the extras" pattern as the
systems track. Candidates show up with idle-timeout, refresh-token
rotation, breadcrumb history, and elaborate frontend state machines
that the brief never asked for.

## Idle / inactivity timeout

> **alex-kovoy** (ibeckermayer): "Lets ignore the idle timeout. I do
> not think that it's required."

See [`auth-storage.md`](auth-storage.md) for the longer discussion.

## Refresh tokens / long-lived sessions

The brief asks for a logout endpoint and a finite session lifetime;
it does not ask for refresh tokens. Candidates who implement OAuth-
style access-plus-refresh pairs get asked why. The expected answer is
a single session cookie with a fixed TTL and a server-side store that
can invalidate on logout.

## Elaborate frontend state

The captured `zship` and `ibeckermayer` threads include reviewer
pushback on Redux-style global state for what is effectively a
two-page app (login, browse). The expected answer: keep state local;
URL is the source of truth for the current directory path; no Redux
unless the candidate is explicit about why they want it.

## Production-grade error pages and telemetry

Reviewers do not ask candidates to ship Sentry, structured logging
pipelines, or rich 4xx/5xx error pages. The brief asks for a working
demo. Candidates who add these get asked to justify the time; usually
the answer is to remove them.

---

The pattern across systems and fullstack is the same: the challenge
is scored on what the brief asks for, plus a clear design doc.
Anything else is surface area for bugs and for "why is this here?"
review comments.
