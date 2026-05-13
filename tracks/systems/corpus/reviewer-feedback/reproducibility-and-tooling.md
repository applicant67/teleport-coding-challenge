# Reproducibility and tooling

The reviewer wants to clone the repo, run `make` or `go test`, and
have it work. They notice when proto generation is host-dependent or
when generated files are missing.

## Pin protoc / use buf

> **rosstimothy** (GevorgGal): "Should we pin to a specific version of
> protoc?"

> **rosstimothy** (GevorgGal): "Could we switch to buf now for
> reproducible generation?"

> **GevorgGal** (response): "Sure - switched to buf with remote
> plugins for reproducible generation."

Expected answer: `buf` (with remote plugins) is the path of least
resistance. Alternatively, a Makefile that installs pinned
`protoc-gen-go` and `protoc-gen-go-grpc` versions via `go install`
with explicit `@vX.Y.Z` tags. Avoid relying on a system-installed
`protoc`.

## Commit generated files

Synthesized from `GITHUB_FEEDBACK.md`:

> **codingllama** (adalton): "It's usual to commit generated files in
> Go repos. Makes it possible to checkout and run without extra
> steps."

Expected answer: `*.pb.go` and `*_grpc.pb.go` are committed. The
reviewer can `go test ./...` on a fresh clone without running buf.

## Reproducible build

> **rosstimothy** (GevorgGal): "Should we specify a config file for
> consistent results?"

Expected answer: a `.golangci.yml` (or similar) pinning the lint
config. A `Makefile` or `tasks.sh` with named targets (`make
proto`, `make test`, `make lint`). One command per stage.

## Module hygiene

Not asked verbatim but a reliable smell-test: a clean module path, a
correct minimum Go version in `go.mod`, no transitive replace
directives, no vendored modules unless the user has a specific
reason.

## CI

No reviewer asked verbatim about CI in the corpus, but every accepted
submission had at least:

- `go test ./...` running on every push
- `go vet ./...`
- `go test -race ./...` running at least on a Linux runner
- A linter (golangci-lint) optionally pinned to a version

The presence of even a minimal CI file is a useful signal of care.
The absence is not a deal-breaker.
