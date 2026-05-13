---
name: teleport-cgroup-reviewer
description: Audit the cgroup v2 code in a Teleport Job Worker Service submission (L5). Use after the candidate writes their cgroup primitives package. Checks for the membership race, kill semantics, sweep, EBUSY handling, subtree-control verification, and io.max device discovery.
tools: Read, Glob, Grep
---

# Teleport cgroup reviewer

The L5 challenge requires cgroup v2 with memory, cpu, and io limits.
The reviewer cares less about the absolute values of the limits than
about whether the cgroup code correctly contains the job. This skill
runs the cgroup-specific checks.

## Check 1: Atomic placement

```bash
grep -rE 'UseCgroupFD|CgroupFD' --include='*.go' .
```

There should be at least one hit. The job-spawn code path must use
`SysProcAttr.UseCgroupFD: true` with `CgroupFD: <fd>`.

Anti-pattern:

```bash
grep -rE 'WriteFile.*cgroup\.procs' --include='*.go' .
grep -rE 'Write.*cgroup\.procs' --include='*.go' .
```

If a process is written to `cgroup.procs` after fork, there is a
race window. Flag it.

## Check 2: cgroup.kill

```bash
grep -rE 'cgroup\.kill' --include='*.go' .
```

There should be at least one hit, in a function that writes `"1"`
to it.

Anti-pattern:

```bash
grep -rE 'syscall\.Kill|process\.Kill\(\)|cmd\.Process\.Kill' --include='*.go' .
```

A handful of hits is expected (cleanup paths, test teardown). But
the *primary* termination path for a job should be `cgroup.kill`,
not `process.Kill()`. If the job library's Stop method calls
`process.Kill()` and does not write to `cgroup.kill`, descendants
will survive.

## Check 3: No Pdeathsig

```bash
grep -rE 'Pdeathsig' --include='*.go' .
```

There should be zero hits. `Pdeathsig` is unsafe in Go because the
runtime migrates goroutines between OS threads; the thread that
spawned the child can exit while the program is healthy, causing the
child to SIGKILL out of nowhere.

If there is a hit, flag it as critical.

## Check 4: Subtree control verification

```bash
grep -rE 'cgroup\.subtree_control' --include='*.go' .
```

After writing to `cgroup.subtree_control`, the code must read it
back and verify the requested controllers (`+cpu +memory +io`) are
present. The kernel silently does nothing if the parent has not
delegated the controllers; the write succeeds but enables nothing.

A passing implementation looks like:

```go
if err := os.WriteFile(filepath.Join(parent, "cgroup.subtree_control"),
    []byte("+cpu +memory +io"), 0); err != nil {
    return fmt.Errorf("write subtree_control: %w", err)
}
enabled, err := os.ReadFile(filepath.Join(parent, "cgroup.subtree_control"))
if err != nil {
    return fmt.Errorf("read subtree_control: %w", err)
}
for _, want := range []string{"cpu", "memory", "io"} {
    if !strings.Contains(string(enabled), want) {
        return fmt.Errorf("controller %q not delegated by parent", want)
    }
}
```

Flag if the read-back-and-verify is missing.

## Check 5: Sweep on startup

```bash
grep -rE 'sweep|orphan|cleanup' --include='*.go' . | grep -i cgroup
```

The worker must, on startup, walk the parent cgroup directory and
clean up any `job-*` subdirectories left from a previous crash. This
handles ungraceful worker exits.

Look for a function that:
1. Lists subdirectories of the parent cgroup
2. For each `job-*` (or similar), writes "1" to its `cgroup.kill`
3. Removes the directory

Flag if missing.

## Check 6: EBUSY retries on cleanup

```bash
grep -rE 'EBUSY' --include='*.go' .
```

Cleanup code (`os.Remove` on the cgroup directory) immediately after
`cgroup.kill` will sometimes return EBUSY because the kernel has
not yet reaped the killed processes.

A passing implementation:

```go
for attempt := 0; ; attempt++ {
    err := os.Remove(path)
    if err == nil || os.IsNotExist(err) {
        return nil
    }
    if !errors.Is(err, syscall.EBUSY) || attempt >= 50 {
        return err
    }
    time.Sleep(10 * time.Millisecond)
}
```

Flag if `os.Remove` on cgroup is called without retry.

## Check 7: io.max device discovery

```bash
grep -rE 'io\.max|/proc/partitions|/sys/block' --include='*.go' .
```

`io.max` is per-block-device. The list of devices must be
discovered at runtime. The reliable approach:

1. Read `/proc/partitions`.
2. Filter to entries whose `name` exists as a directory under
   `/sys/block` (whole disks only, not partitions).
3. Skip major 1 (ramdisk) and 7 (loop).

Anti-pattern:

```bash
grep -rE '"/dev/sda"|"/dev/nvme0n1"|"/dev/vda"' --include='*.go' .
```

Hardcoded device IDs fail on any machine without that device. This
is one of the recurring rejection reasons. Flag if found.

## Check 8: Resource limits configurable

```bash
grep -rE 'memory\.max|cpu\.max|cpu\.weight|io\.max' --include='*.go' .
```

The code should write all three resource controllers. The values
come from server-side configuration, not the RPC.

Anti-pattern: RPC request fields like `memory_max` or `cpu_max`.
Resource limits are server-side static; clients do not specify
them.

## Check 9: Fail at startup if cgroup setup fails

```bash
grep -rE 'log.*[Ww]arn.*cgroup|fmt\.Printf.*WARN.*cgroup' --include='*.go' .
```

If cgroup configuration fails, the server should exit non-zero, not
log a warning and continue. Running without enforcement is a
correctness violation; the spec requires limits be enforced.

Flag any cgroup-setup error path that logs and continues.

## Check 10: Single fd ownership

After `cmd.Start()` succeeds with `CloneIntoCgroup`, the cgroup
directory fd is no longer needed by the parent. Close it.

```bash
grep -B5 -A5 'CgroupFD:' --include='*.go' .
```

Look for a `Close(fd)` or `defer Close(fd)` after the spawn. Flag
if missing; long-running parents will exhaust fds.

## Output format

```markdown
# Cgroup audit of <repo>

## Check 1: Atomic placement
Status: PASS / FAIL
Evidence: <file>:<line>
Issues: <list>

[repeat for each check]

## Summary
- Checks passed: N/10
- Critical issues: <list with file:line>
- Recommended fixes ordered by severity

## Risk assessment
- Will reviewer flag this? YES / MAYBE / NO
- Estimated fix time: <hours>
```

The output is precise: file:line for each issue, code suggestion
where useful, no narrative.

## Related reference

Full pitfall descriptions: `../../TECHNICAL_PITFALLS.md`.
Process containment quotes: `../../corpus/reviewer-feedback/process-containment.md`.
