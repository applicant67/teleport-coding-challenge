# What reviewers tell candidates to cut

The single most consistent message across the corpus: candidates add
features the challenge does not require, and reviewers ask them to
remove the features. Every category below has been called out by
multiple reviewers across multiple submissions.

## Graceful shutdown (SIGTERM then SIGKILL)

> **rosstimothy** (GevorgGal): "You may simplify and omit graceful
> termination via SIGTERM."

> **nklaassen** (MrChristianL): "it is okay to skip graceful shutdown
> and just send SIGKILL"

> **rosstimothy** (MarkDHarris): "This is _not_ a requirement for L4.
> I suggest moving it out of scope as it adds additional complexity to
> do correctly."

> **tigrato** (kkloberdanz): "For simplicity, feel free to force kill
> the process."

Expected answer: cgroup.kill is enough. The graceful path adds a
timeout, a select between context and the cleanup goroutine, and a
race window in which the job could exit naturally between SIGTERM and
SIGKILL. None of that is a requirement.

## ListJobs

> **rosstimothy** (GevorgGal): "There is no requirement to list jobs.
> I suggest cutting it from scope to reduce the amount of work you
> have for this exercise."

Expected answer: do not implement it. The four required operations
are Start, Stop, Status, StreamOutput.

## Stdin

> **greedy52** (MarkDHarris): "I don't know if SSH or interactive
> shell is relevant to the goal of the challenge. Just stdin is not
> part of the requirement. Up to you though."

Expected answer: out of scope. `cmd.Stdin = nil`. If the candidate
README or proto mentions stdin, delete it.

## Process isolation / namespaces

> **rosstimothy** (sabernabil12): "Process isolation is not a
> requirement for level 4."

> **espadolini** (Zephan92): "The challenge text does not mention
> using namespaces for isolation, I would recommend sticking to the
> requirements and leaving anything extra for the end if there's time
> left."

Expected answer: namespaces are L5+ territory and even there
optional. cgroups give you resource limits; that is what the spec
asks for. Mount/PID/network namespaces are a separate undertaking.

## Stream offset / follow flags

> **zmb3** (razzam21): "The challenge only requires that we offer the
> ability to stream output from the very beginning of process
> execution, so feel free to simplify by removing the from_offset
> field. Similarly, we can make follow the default behavior to further
> reduce scope."

> **rosstimothy** (nombiezinja): "For this challenge follow is the
> only required behavior. I'd suggest omitting this to reduce scope."

> **creack** (razzam21): "I would recommend to keep it simple for the
> exercise, there is no need for resume or gap detection. The
> requirements mention 'Output should be from start of process
> execution.'"

Expected answer: one streaming mode. Always from the start, always
follows until the job exits or the client disconnects.

## Job ID in response

> **zmb3** (razzam21): "We don't really need to include the job ID in
> the response since it was specified in the request. (Same goes for
> output streaming, where it's even more beneficial to avoid sending
> unnecessary data on the wire)"

> **eriktate** (GevorgGal): "You might consider only including the
> job ID as output. This makes it easier to do things like start a
> job and stream log output."

> **eriktate** (GevorgGal): "nit: consider removing the
> `Job started:` prefix. It makes it easier to pipe jobctl start into
> other commands."

> **Joerger** (MrChristianL): "Since we failed to start the job,
> isn't the job ID irrelevant? Especially if this is a user-facing
> error."

Expected answer: CLI prints the bare job ID. RPC responses omit
echoed input. CLI is pipeable.

## Explicit end-of-stream message

> **rosstimothy** (nombiezinja): "Why do you need to send an extra
> message to communicate that all the data was sent? Can you just
> close the stream after sending the last byte?"

> **nklaassen** (MrChristianL): "what is an EOF message?"

Expected answer: closing the stream is end-of-stream. No EOF marker
in the proto.

## Dynamic / client-set resource limits

> **tigrato** (kkloberdanz): "No need to receive dynamic resource
> limits. feel free to hardcode them server side"

Expected answer: server-side static limits. Clients do not specify
limits in the RPC. The CLI does not have flags for them.

## Multi-user ACLs / sharing / allowed-viewers

> **rosstimothy** (MarkDHarris): see above under graceful shutdown.

Expected answer: owner-only, or owner-plus-admin. No ACL list.

## Sleep / SetDefault / globals in tests

Out of scope is also a hygienic boundary. See `testing.md`.

## Command whitelisting

> **fspmarshall** (mcampo84): "Feel free to omit this. It isn't
> required by the challenge, and implementing a command whitelist
> tends to be a pretty iffy undertaking at the best of times."

Expected answer: do not whitelist. Any binary the server user can
exec is fair game.

## Resource limits must exist, but with defaults

This is the one place reviewers push *toward* something rather than
away:

> **rosstimothy** (joshuarubin): "The challenge requires that all
> jobs started have resource limits enforced. Perhaps instead of
> defaulting to no limits the flags can default to some arbitrarily
> chosen limits?"

Expected answer: server applies non-zero defaults for memory and CPU
limits. A "no limit" config is acceptable for testing but the
server's normal mode has limits enforced.
