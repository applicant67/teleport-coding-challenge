# Support track

The current support challenge is a **live debug session** conducted
via Teleport itself during the interview, against a simulated customer
environment. The candidate produces a short writeup per issue
diagnosed.

## Status of this directory

**Structurally non-public.** Because the support track does not
produce a candidate-owned GitHub artifact, there is no public corpus
to triage. Fragment-search across distinctive challenge phrases
confirmed: no public Teleport support submissions exist on GitHub.

## What the challenge typically covers

Fault categories worth preparing for, in rough order of frequency
seen in public Teleport documentation and incident postmortems:

1. **Kernel / systemd / cgroups.** Process limits, OOM, service-unit
   misconfiguration, cgroup v1 vs v2 confusion.
2. **Filesystem / volumes.** Mount namespaces, full disks, permission
   inversions, symlink edge cases.
3. **Networking.** DNS resolution, firewall rules, TLS termination,
   reverse-proxy timeouts.
4. **TLS.** Cert chains, mTLS misconfiguration, clock skew.

## Writeup template that works

Half a page per problem, three sections:

1. **What was wrong** (root cause in one sentence, then evidence).
2. **How resolved** (the specific change made, with the command).
3. **Future action items** (what monitoring, alerting, or process
   change would catch this earlier).

## What to do before the live session

1. Read the official prompt carefully. Confirm the artifact format
   with the recruiter ahead of time.
2. Apply [`../../docs/HUMANIZATION.md`](../../docs/HUMANIZATION.md)
   to any writeup if you used AI assistance - generic AI-flavored
   "let's verify the logs and check the configuration" prose is
   especially detectable in support writing.
3. Apply [`../../docs/PROCESS_NOTES.md`](../../docs/PROCESS_NOTES.md)
   for recruiter / cadence expectations.

## Cross-track principles that apply

- Cut scope to the prompt. Don't add invented context.
- Acknowledge unknowns explicitly. "I would ask the customer X" is
  stronger than guessing.
- Specific beats general.
