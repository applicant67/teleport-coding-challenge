---
name: teleport-humanizer
description: Strip AI tells from a Teleport Job Worker Service submission. Use when the candidate has finished implementation and wants a humanization pass. Runs the lexical audit, finds non-ASCII characters, marketing words, multi-paragraph docstrings, and AI-typical patterns. Outputs a list of fixes.
tools: Read, Glob, Grep, Edit
---

# Teleport humanizer

Teleport's stated position is that AI-generated submissions are not
acceptable. Reviewers are openly skeptical of submissions that read
as machine-authored.

This skill runs a humanization pass over a candidate's repo. It does
not generate code; it finds patterns that are common in AI output
and rare in careful human code.

## Pass 1: Non-ASCII characters

Run:

```bash
grep -rP '[^\x00-\x7F]' --include='*.go' --include='*.md' \
    --include='*.proto' --include='*.yaml' --include='*.yml' \
    --include='Makefile' .
```

Every hit is a fix. Convert:

- `\u2014` (em-dash) -> `-` or ` - ` or restructure
- `\u2013` (en-dash) -> `-` or restructure
- `\u2018`, `\u2019` (curly quotes) -> `'`
- `\u201C`, `\u201D` (curly double quotes) -> `"`
- `\u2026` (ellipsis) -> `...`
- `\u2265` (greater than or equal) -> `>=`
- `\u2264` (less than or equal) -> `<=`

The em-dash is the strongest single AI tell. A repo with em-dashes
in markdown files reads as AI-authored at a glance.

## Pass 2: Marketing language

```bash
grep -rwE '\b(robust|seamless|comprehensive|sophisticated|elegant|leverage|utilize|extensive|advanced)\b' \
    --include='*.go' --include='*.md' --include='*.proto' .
```

Each hit: rewrite to plain language or delete.

```bash
grep -rwE '\b(furthermore|moreover|additionally|importantly|notably|crucially)\b' \
    --include='*.go' --include='*.md' .
```

Same.

```bash
grep -rE 'It (is|should be) worth noting' \
    --include='*.go' --include='*.md' .
grep -rE 'It is important to' \
    --include='*.go' --include='*.md' .
grep -rE 'This (carefully|properly|safely) (handles|manages)' \
    --include='*.go' --include='*.md' .
```

These phrasings appear in AI output and almost never in production
human-written Go.

## Pass 3: Multi-paragraph docstrings on small functions

```bash
grep -rB1 '^func ' --include='*.go' .
```

Manual review: for each function with more than 4 lines of
docstring, ask whether the function is complex enough to warrant it.
Most are not. Trim to one line.

The convention: doc the exported type or method in one sentence.
Doc the non-obvious why on internal helpers only when the why is
not in the code.

## Pass 4: Comment-paraphrases-code

```bash
git diff <baseline>..HEAD | grep -E '^\+\s*//' | grep -v '^+++' \
    > /tmp/new-comments.txt
```

Open `/tmp/new-comments.txt`. Read each new comment. For each:

- Does the next line of code restate what the comment says? Delete
  the comment.
- Does the comment say "this returns X" or "checks if Y"? Delete it;
  the signature already says that.
- Does the comment hedge or perform care? ("Carefully handle...",
  "Important: ...", "Note: ...") Rewrite as a flat fact or delete.

The default for comments in this challenge: do not write one.
Reviewers will mark every comment that does not earn its keep.

## Pass 5: Defensive code for impossible cases

```bash
grep -rn 'if .* == nil' --include='*.go' . | grep -E 'func \('
```

For each `if foo == nil` inside a method on `*foo`: is it possible
for the caller to pass a nil? If not, delete the check. If yes,
either document the contract or use a constructor that prevents
nils.

```bash
grep -rn 'recover()' --include='*.go' .
```

Each `recover()`: is the code path actually panicking? If not,
delete. Defensive `recover()` in code that does not panic is an AI
tell.

## Pass 6: Symmetric error wrapping

```bash
grep -rE 'fmt\.Errorf\("failed to [^"]+: %w"' --include='*.go' .
```

If every error in the codebase is `failed to X: %w`, you have an AI
tell. Vary the wraps:

- `fmt.Errorf("open %s: %w", path, err)` - prefix + arg + wrap
- `fmt.Errorf("read partitions: %w", err)` - prefix + wrap
- `fmt.Errorf("invalid limit %d for %s", n, name)` - no wrap, terse
- `fmt.Errorf("write cgroup.kill: %w", err)` - operation + wrap

About a third of wraps should be in each style. None should be
exactly `failed to X: %w`.

## Pass 7: Long-form test names

```bash
grep -rE '^func Test\w{40,}' --include='*_test.go' .
```

Each test with a name longer than ~40 characters: shorten. Go tests
are typically `TestVerb_NounWhen` or `TestVerb_Scenario`. Subtests
do the rest.

## Pass 8: Generic AI variable names

```bash
grep -rE '\b(theResult|theData|theValue|theStatus|theError|finalResult|currentValue)\b' \
    --include='*.go' .
```

Each hit: rename. Go prefers `res`, `data`, `val`, `status`, `err`.

## Pass 9: README and PR description audit

Read every `README.md`, every PR description, every commit message
in the candidate's branch.

Patterns to remove:

- Section headers like "Features", "Highlights", "Architecture
  Overview", "Key Design Decisions" on a small project README
- Bullet lists where prose would do
- The word "production-ready" or "production-grade"
- "Designed with X in mind" phrases
- "Made with love" or attribution footers
- "Generated by Claude" or "Co-authored by" lines pointing at AI
- Em-dash separators in headings

## Pass 10: AI co-author lines

```bash
git log --all --grep='Co-Authored-By: Claude' --oneline
git log --all --grep='Generated' --oneline
git log --all --grep='Anthropic' --oneline
```

Each hit: rewrite the commit message (interactive rebase or commit
amend). The reviewer reads commits.

## Output

Produce a humanization report:

```markdown
# Humanization report

## Pass 1: Non-ASCII
N hits across M files. <list of files>

## Pass 2: Marketing language
N hits. <samples>

[...]

## Recommended actions, ordered
1. Run `make humanize-pass1` to convert non-ASCII (auto-fixable).
2. Manually rewrite N marketing-word hits.
3. Manually delete N AI-paraphrase comments.
4. Rewrite N commit messages.

Estimated effort: H hours.
```

After the report, ask the candidate whether to auto-apply the
non-ASCII fixes (safe) and require manual review of the rest
(unsafe to auto-apply).

## What this skill does not do

This skill is cosmetic. It does not fix:

- Scope creep (use `teleport-scope-disciplinarian`)
- Concurrency bugs (use `teleport-reviewer-mindset`)
- Test hygiene (use `teleport-design-doc-reviewer`)

The four rejection lessons in `../../tracks/systems/REJECTION_LESSONS.md`
are structural (these are systems-specific examples; the principle
generalizes to any track). No amount of humanization covers a
graceful-shutdown race or an `io.Writer` that lies about how many
bytes it wrote.

Treat humanization as the last 10% of submission prep, not the
first.

## Related reference

Full process: `../../docs/HUMANIZATION.md`.
