# Testing

Reviewers care a lot about test hygiene. The patterns they flag come
back constantly.

## Use t.Context, not context.Background

> **rosstimothy** (MarkDHarris): "Suggestion: Use testing.T.Context()
> instead of a background context."

> **rosstimothy** (GevorgGal): "Suggestion: Prefer using t.Context()
> in tests"

> **eriktate** (MrChristianL): "nit: prefer t.Context() over
> context.Background() in tests"

Expected answer: every test that needs a context starts with
`ctx := t.Context()`. The test framework cancels it on `t.Cleanup`.

## No sleeps for synchronization

> **rhammonds-teleport** (MarkDHarris): "Small sleeps in unit tests
> can add up over time. Can any of these tests be reworked to avoid
> them?"

> **tigrato** (RichyHBM): "Is there any viable option that doesn't
> require sleep?"

Expected answer: channels, condition variables, polling-with-deadline,
or `eventually` style helpers. A bare `time.Sleep(100ms)` is a smell;
reviewers will mention it.

## Test against a real client

> **codingllama** (adalton): "Please avoid testing RPCs by calling
> server methods directly. The recommended way is to use a proper
> client to exercise it, so you get the full server running with
> interceptors."

> **nklaassen** (Zephan92): "imo it would be nicer to just do all
> these tests with a real TLS connection and real certs. I don't
> think this case would be possible with RequireAndVerifyClientCert"

Expected answer: integration tests open a real `net.Listen` on
loopback, attach the real interceptors, and dial with a real gRPC
client. Direct method-call tests miss the interceptor stack and the
TLS / auth flow.

## Avoid context in structs

> **nklaassen** (Zephan92): "avoid storing a context in a struct and
> prefer testing the external API
> https://go.dev/blog/context-and-structs"

Expected answer: context is a parameter, not a field. The library's
operations take a context. Internal state does not.

## Don't depend on host hardware

Synthesized from the corpus (rejection-grade observation): one
rejected submission's `io.max` test hardcoded `/dev/sda`. The reviewer
ran the test on a machine without `/dev/sda` and the test failed in
review. The fix is to discover block devices via `/proc/partitions`
and pick one dynamically, or to gate the test on whether at least one
whole-disk block device exists.

## Logging in tests

> **espadolini** (Zephan92): "Please use structured logging with
> log/slog rather than string formatting."

Expected answer: tests use `slog.New(slog.DiscardHandler{})` (Go 1.25)
or a per-test handler. No `slog.SetDefault` mutation in `TestMain`,
because it leaks to parallel tests.

## What CI runs

Not asked verbatim, but the implicit expectation:

- Unit tests run without root and without a real Linux kernel.
  Cgroup paths are mocked.
- Integration tests require root and a real Linux kernel.
- The two are clearly separated. The design doc names which is which.

## Race detector and goleak

Strong recommendation across multiple submissions (and a recurring
correctness lesson the rejected sample learned the hard way):

- Run unit and integration tests with `-race`.
- Use `go.uber.org/goleak` in `TestMain` to catch goroutine leaks.

Goroutine leaks on client disconnect are one of the recurring
reviewer concerns (see `output-streaming.md`). Goleak catches them in
CI.
