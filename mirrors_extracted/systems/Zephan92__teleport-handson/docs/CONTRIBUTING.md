# Contributing Guide

This document outlines the development standards and architectural constraints for the Job Worker Service.

## Architecture Constraints

1.  **Dependency Direction**:
    *   `cmd/` depends on `internal/` and `api/`.
    *   `internal/server` depends on `internal/worker`.
    *   `internal/worker` (The Library) must **not** depend on `internal/server` or `cmd`.
    *   `internal/worker` should be usable as a standalone library.

## Coding Standards

1.  **Concurrency**:
    *   Avoid global state. If necessary, protect it with `sync.Mutex` or `sync.RWMutex`.
    *   Ensure all goroutines have a lifecycle (e.g., context cancellation or quit channels) to prevent leaks.
    *   Use `go test -race ./...` to detect data races.

2.  **Error Handling**:
    *   Return errors, do not `panic` (except during strict initialization in `main`).
    *   Wrap errors with context: `fmt.Errorf("failed to start job: %w", err)`.
    *   Do not ignore errors (e.g., `_ = func()`). Log them if they cannot be returned.

3.  **Security**:
    *   Use mTLS for all network communication.
    *   Validate all inputs (especially in the gRPC handlers).

## Testing

*   **Unit Tests**: Required for `internal/worker` logic (Cgroups, OutputBuffer).
*   **Integration Tests**: Required for the full `Start -> Stream -> Stop` lifecycle.

## Pre-Submission Checklist

Before submitting a Pull Request, please ensure the following automated checks pass:

1.  **Formatting**: Run `go fmt ./...` to ensure code style consistency.
2.  **Static Analysis**: Run `go vet ./...` to catch common Go errors.
3.  **Race Detection**: Run `go test -race ./...` to detect concurrency issues.
4.  **Dependencies**: Run `go mod tidy` to clean up `go.mod` and `go.sum`.
5.  **Code Generation**: If you modified `.proto` files, run `make proto` to update generated code.
