# Library API shape

What the library's exported surface should look like, what idioms
reviewers want, what footguns they flag.

## Library must be importable

> **codingllama** (joshuarubin): "You've placed the entire API inside
> internal/, so that makes it not importable. The intent for the
> 'library' part is that it could be imported and used outside of the
> project."

> **eriktate** (MrChristianL): "Is this PR meant to be the entire
> library, or are you planning to add to it? If so, you don't actually
> have a public API surface given that everything exists in
> internal/."

Expected answer: the core `job` package lives at the module root or
in `pkg/`, not under `internal/`. External consumers can import it.

## Don't leak protobuf types into the library

> **smallinsky** (kkloberdanz): "I'm a bit concerned about coupling
> the client library's Go domain types directly to the protobuf
> transport-layer types. Would it make sense to keep the transport
> types internal to the client, and expose explicit library-level
> types instead (for example job.Status instead of pb.JobStatus in
> this case)? That way, we keep a clean API boundary."

Expected answer: protobuf types live in the gRPC layer. The library
has its own status enum, exit code type, and so on. The gRPC server
translates between them.

## Use standard interfaces

> **tigrato** (chintamanil): "Can we satisfy io.ReadCloser so io.Copy
> works?"

> **nklaassen** (MrChristianL): "suggestion: func([]byte) error looks
> almost like io.Writer.Write, I'd consider either having this
> function accept an io.Writer or return an io.Reader to make it more
> composable"

> **rosstimothy** (RichyHBM): "I think you could reduce the complexity
> a bit if this was inverted. Instead of registering an io.Writer the
> API provided a means for callers to get a unique io.Reader that
> tracked their progress through the buffer."

Expected answer: implement `io.Reader` / `io.Closer` /
`io.ReadCloser` where applicable. Callback APIs that mirror these
interfaces are a code smell.

## API ergonomics: footgun reduction

> **rosstimothy** (GevorgGal): "This API is clunky, inefficient ...
> Can you think of any ways to better encapsulate some of this from
> end users to improve the API and remove potential footguns?"

> **rosstimothy** (MarkDHarris): "It might be worthwhile to separate
> the responsibilities of output management to another object."

Expected answer: each operation is one method call. The library does
not require callers to coordinate two-step sequences (get-buffer +
connect, mark-done + flush). The cleanup is in the same call as the
operation that needs cleanup.

## No empty-but-valid responses

> **rosstimothy** (MarkDHarris): "This permits writing _after_ the
> output is closed."

Expected answer: writes after close return an error. Reads after
close return `io.EOF`. There is no third valid return shape.

## Single atomic status snapshot

> **rosstimothy** (MrChristianL): "Does the lock need to be held
> while doing IO to uphold the guarantees mentioned in this comment?"

> **MrChristianL** (response): "I've implemented a Snapshot() method
> that captures the Status, ExitCode, and StopReason within a single
> lock cycle."

Expected answer: a `Snapshot()` method returns status, exit code, and
stop reason from a single lock acquisition. Three separate getters
risk an inconsistent view.

## No `*os.Process` in the public type

> **espadolini** (chintamanil): "With this definition of Job, how is
> the grpc server supposed to request termination for a job?"

> **chintamanil** (response): "Added *os.Process field and termination
> methods"

Expected answer: the Job type exposes `Stop`, `Status`, `Output`. It
does not expose the underlying `*os.Process` or `*exec.Cmd`.

## No mutable exported fields

Implicit from the corpus: every reviewer who saw fields like
`Job.Status` exported as a plain field rather than via an accessor
asked about concurrent access. Use accessors that return values.

## Documentation

> **dboslee** (benmoss): "could you add go doc comments for the API?
> some of the fields/funcs are not immediately clear how they are
> intended to be used"

Expected answer: every exported type and method has a Go doc comment.
Comment style is plain prose, one paragraph, no marketing tone.

## CLI shape

> **eriktate** (GevorgGal): "consider removing the `Job started:`
> prefix. It makes it easier to pipe jobctl start into other commands"

Expected answer: CLI prints raw IDs. No decorative prefixes. Output
is pipeable.

## RPC shape

> **zmb3** (razzam21): "We don't really need to include the job ID in
> the response since it was specified in the request."

> **nklaassen** (sabernabil12): "i think success/fail could be
> determined by returning an error from the RPC rather than including
> it in the response"

Expected answer: RPCs use gRPC's error returns, not status fields in
the response. Responses do not echo inputs.

## Context lifecycle separate from RPC lifecycle

> **tigrato** (joshuarubin, paraphrased in GITHUB_FEEDBACK): "Is it a
> good design to have the job depending on the parent context? If the
> parent context is the RPC context, won't it be canceled as soon as
> the job start RPC ends?"

Expected answer: the job's lifecycle context is detached from the
incoming RPC context. The Start RPC returns immediately; the job
keeps running.
