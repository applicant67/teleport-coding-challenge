# Concrete technical pitfalls

Code-level traps observed across the public corpus, plus a couple from one
candidate's own rejected submission. Each entry: what the trap looks like,
why it bites, and how to avoid it.

## Cgroup membership race at process start

The naive sequence is fork the child, then write its PID into
`cgroup.procs`. The window between fork and the write is a window where
the child can spawn descendants outside the cgroup. A shell script that
does `sleep 1 &` in its first line will demonstrate this.

The fix is `CLONE_INTO_CGROUP`, available via `clone3(2)` since Linux 5.7
and surfaced in Go's `os/exec` as `SysProcAttr.UseCgroupFD` + `CgroupFD`.
You open the cgroup directory, pass the file descriptor, and the kernel
places the child into the cgroup before its first instruction runs.

```go
fd, err := syscall.Open(cgroupPath, syscall.O_RDONLY|syscall.O_DIRECTORY|syscall.O_CLOEXEC, 0)
if err != nil { return err }
defer syscall.Close(fd)

cmd.SysProcAttr = &syscall.SysProcAttr{
    UseCgroupFD: true,
    CgroupFD:    fd,
}
```

If your design uses fork-then-cgroup-write, reviewers will catch it.

## cgroup.kill availability

`cgroup.kill` is the right way to terminate the entire process tree of a
job in one syscall. It is a single-write atomic operation that SIGKILLs
every member of the cgroup, including descendants the parent never tracked.

Available since Linux 5.14. State the version floor in your design doc.

The alternative ("send SIGKILL to the leader, then enumerate `cgroup.procs`
and kill each entry, then handle the race where new entries appear between
the enumeration and the kill") is more complex and a reviewer will ask
why you did it the harder way.

## Pdeathsig is unsafe in Go

`syscall.SysProcAttr.Pdeathsig = syscall.SIGKILL` is sometimes proposed as
"kill the child if the parent dies." It does not do that in Go. Pdeathsig
considers the *thread* that spawned the child. Go's runtime can park and
migrate goroutines between OS threads, so the thread that spawned the
child can exit while the program is healthy, causing the child to receive
a SIGKILL out of nowhere.

The correct mechanism: cgroup sweep on worker startup. List the parent
cgroup directory, find every `job-*` subdirectory, kill any remaining
members, remove the directory. This handles ungraceful crashes of the
worker without relying on thread-affinity semantics.

## Output buffer per reader

The temptation: hand each stream subscriber its own bytes.Buffer or its
own bytes channel that the writer fans out into. The bug: 1 GB of job
output and five readers is now 5 GB of memory.

The correct design: one shared log (file-backed is simplest), each reader
opens its own file descriptor and reads at its own pace. Notifications of
new bytes are per-reader channels or a condition variable broadcast.

The "file-backed log" piece deserves emphasis. Writing the output to a
real file on disk means:

- The OS handles the slow-reader problem via the page cache.
- Each reader is just a normal `os.File` with its own offset.
- Replay from the beginning is trivial; the file is already there.

In-memory designs can be made to work but require careful slow-reader
handling. Disk is simpler.

## Closing the done channel while writes are still draining

The race observed and explicitly flagged by reviewers: the goroutine that
copies from the child's stdout pipe into the log finishes copying. You
mark the stream done. But the kernel pipe still has bytes that os/exec
has not yet drained, and now your readers miss those bytes.

The fix has two parts:

1. Wait for `cmd.Wait()` to return before marking the stream done.
   `cmd.Wait` waits for both the process to exit and for the stdout/stderr
   copy goroutines to finish.
2. Mark the stream done atomically in the same step that flushes the
   writer. Do not have a separate `SetDone()` call.

This is the same bug as "API split across two functions" in
`REJECTION_LESSONS.md`. The structural fix and the race fix are the same:
`Close()` marks done and releases the writer fd.

## Stdin handling

Stdin is out of scope. The official challenge text was clarified on this
point. Some older submissions wire up stdin and reviewers tell them to
remove it. Wire `cmd.Stdin = nil` and move on.

## Discovering block devices for io.max

`io.max` requires you to apply throttling per block device. The
enumeration step is what gets tested. The reliable approach:

- Read `/proc/partitions`. Each line has `major`, `minor`, `blocks`,
  `name`.
- Filter to whole disks: `name` exists as a directory under `/sys/block`.
- Skip ramdisks (major 1) and loops (major 7).

Two failure modes from the corpus:

- Hard-coding device IDs. Fails on any machine that is not yours.
- Not handling the "no top-level block devices" case. On a tmpfs- or
  overlay-only environment, the list is empty. Decide whether to error or
  proceed without io.max enforcement, and document the choice.

For testability, isolate the parsing logic from the discovery logic. The
parser takes an `io.Reader` of `/proc/partitions` content and a
`isWholeDisk func(name string) bool` callback. The integration version
opens the real file and stats `/sys/block/<name>`. Now the parser is
unit-testable without depending on the host's actual hardware.

## Cleanup retries on EBUSY

Removing a cgroup directory immediately after writing `cgroup.kill` will
sometimes return `EBUSY` because the kernel has not yet reaped the killed
processes. Retry with a short backoff.

```go
for attempt := 0; ; attempt++ {
    err := os.Remove(path)
    if err == nil || os.IsNotExist(err) { return nil }
    if !errors.Is(err, syscall.EBUSY) || attempt >= 50 { return err }
    time.Sleep(10 * time.Millisecond)
}
```

Half a second of retries is plenty. Watch the boundary: do not poll
forever, do not give up immediately.

## Subtree control delegation

Writing `+cpu +memory +io` into `parent/cgroup.subtree_control` is the
step that makes those controllers available to child cgroups. If the
parent's parent has not delegated the controller (common in containerized
environments and some systemd configurations), the write succeeds but
the controller is silently not enabled for children.

After writing, read the file back. If a requested controller is not
present, hard-fail at startup with a useful message that says "this
parent cgroup does not have cpu/memory/io delegated." Otherwise the
failure surfaces much later, in `memory.max` writes, as a confusing
`EINVAL`.

## Owner / authorization separation

Common mistake: the `job` library exposes a `StartJob(owner, ...)` and
stores the owner, then has a `CanRead(caller)` method that compares.

Better: the library exposes `StartJob(meta, ...)` where `meta` is
opaque-to-the-library identity (an owner string is fine). The server's
gRPC handler is what extracts the identity from the verified cert, passes
it as the meta, and rejects unauthorized accesses at the handler level
before calling into the library.

The library should not be making policy decisions.

## Identity from the verified chain, not PeerCertificates

In a gRPC server with mTLS:

- `peer.FromContext` gives you the peer.
- `peer.AuthInfo.(credentials.TLSInfo).State.PeerCertificates[0]` is the
  cert the peer presented, but it is set even when the server is
  misconfigured to accept unverified clients.
- `peer.AuthInfo.(credentials.TLSInfo).State.VerifiedChains[0][0]` is the
  same cert, but only set if the verification succeeded.

Always read from `VerifiedChains`. If the slice is empty, refuse the
request.

## Identity from cert Subject, not Serial Number

Cert serials rotate when certs are reissued. Subject (CN or SAN) is
stable across re-issues. Use Subject.

## Wait can be called once

`cmd.Wait()` is single-shot. If you have a separate goroutine watching the
child and another path that also calls `Wait`, you have already lost.
Centralize the `Wait` to one place.

## Output writer must obey io.Writer contract

`Write(p)` returns `(len(p), nil)` after a successful write of all bytes,
or `(n, err)` with `n < len(p)` and `err != nil` for a short write. There
is no third valid return shape. In particular:

- Do not return `(0, nil)`. That is a contract violation.
- Do not return `(len(p), nil)` from a branch that did not write all the
  bytes. That is the bug that got this candidate's submission rejected.

If the underlying store is closed and you do not want to write, return
`(0, ErrClosed)` or a similar named error. Never silently drop.

## Test isolation

Restating from `REJECTION_LESSONS.md` because reviewers care a lot:

- No `slog.SetDefault` mutation in `TestMain`.
- No environment variable mutation that leaks to other tests.
- No test that depends on hardware specifics of the host.
- Prefer `t.Context()` over `context.Background()`.
- Sleeps in tests are a smell. Use channels or polls with deadlines.
- Run with `-race` and `go.uber.org/goleak` in `TestMain`. Both catch
  real bugs that are otherwise invisible until production.
