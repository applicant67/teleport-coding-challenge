# Removing AI tells from a submission

Teleport's stated position is that AI-generated submissions are not
acceptable. The challenge instructions ask candidates to confirm they wrote
the code themselves. Reviewers are openly skeptical of submissions that
read as machine-generated.

You can read that position in different ways. The pragmatic one: if your
process used AI assistance, you are responsible for transforming the output
into work you understand, can defend in a follow-up call, and that does not
read as machine-authored. This doc is about the second part.

This is not advice on whether to use AI. That is a personal decision and a
question of how you choose to represent yourself.

## What counts as an "AI tell"

The patterns reviewers and humans pattern-match on:

- **Non-ASCII typography**. Em-dashes, en-dashes, fancy quotes, ellipsis
  glyphs, the unicode greater-than-or-equal `>=` instead of `>=`. In source
  code, in comments, in markdown. Search for `[^\x00-\x7F]` across every
  file before submission.
- **Over-commented obvious code**. A comment on every line, comments that
  paraphrase the line below in marketing tone, "this function returns the
  result of X" docstrings. Reviewers can read the code; they need comments
  only for the non-obvious why.
- **Multi-paragraph docstrings on small functions**. A two-line helper does
  not need a 12-line preamble.
- **Hedging language in comments**. "It's worth noting that...",
  "Importantly,...", "This carefully handles...". Comments should state
  facts, not perform care.
- **Symmetric error messages**. AI tends to produce
  `failed to do thing: %w` for every wrap, every time. Real code has
  variety: lowercase, no period, sometimes terse, sometimes specific.
- **README inflation**. Long, well-structured READMEs with sections like
  "Features", "Installation", "Usage", "Contributing", "License" on a
  one-binary internal coding challenge. Read peer submissions before
  writing yours.
- **Test names that read as outlines**. `TestThatTheXShouldYWhenZ` style.
  Most Go codebases use shorter names.
- **Variable names that read as descriptions**. `theJobStatusAfterStop`
  instead of `status` or `final`. AI loves long variable names; codebases
  rarely do.
- **Defensive code for impossible cases**. `if obj == nil` inside a method
  on `*obj` where the caller cannot have a nil. AI adds these; Go reviewers
  call them out.
- **Misuse of the word "robust"** and similar reviewer-bait adjectives in
  comments and PR descriptions.

## The process that worked

The process below is what one candidate developed. It is not a
guarantee of anything, but it ran cleanly enough that the reviewer in the
final iteration did not call out AI authorship.

### Pass 1: read your code as if it were a stranger's

After the implementation passes tests, sit down with the diff and read
every line out loud or near-out-loud. Mark anything that:

- explains what the code does, where the code itself is clear
- uses words you do not personally use in writing
- has more structure than the function deserves
- contains a non-ASCII character

Fix these in-place. This pass alone removes most AI tells.

### Pass 2: dedicated new-comment review

After pass 1, before the next AI-assisted change, extract every comment
that was net-added in the session to a single scratch file:

```
git diff <last-blessed-rev>..HEAD | grep -E '^\+\s*//' | grep -v '^+++' > /tmp/new-comments.txt
```

Review the scratch file as a list. AI tells are easier to spot when comments
are torn out of context. Edit them down, or delete them entirely if the code
does not need them. The default for comments in this challenge is: do not
write one.

### Pass 3: diff against your last "blessed" revision

Keep numbered checkpoints of the code at each humanization pass: rc1, rc2,
rc3, etc. Whenever you complete a humanization pass, snapshot the tree so
you can `diff` the next AI-assisted change against the last known clean
state. This is the cheapest way to catch new AI tells: they are visible as
green lines in the diff. Many will be in code you did not even ask the
agent to touch.

### Pass 4: lexical audit

Before submission, run a literal search for the specific tokens:

```
grep -rP '[^\x00-\x7F]' .
grep -rE '\b(robust|seamless|comprehensive|elegant|leverage|utilize)\b' .
grep -rE '\b(furthermore|moreover|additionally|importantly)\b' .
grep -rE 'It (is worth|should be) noted' .
```

Almost every hit is something to delete or rewrite.

### Pass 5: explain your code to a peer

The most reliable filter. If you cannot explain why the kill function takes
no arguments, or why the output store does not need a separate `SetDone`,
or why you chose `peer.FromContext` over `PeerCertificates`, your design
will not survive a 60-minute call with a Teleport engineer. Find a real
human and walk through it. The questions they ask are the questions the
reviewer will ask.

## Working with an agent

Practical instructions for the agent you are using:

- Tell it not to write attribution lines, AI co-author tags, or "Generated
  by" footers anywhere.
- Tell it not to use em-dashes, en-dashes, or fancy quotes. Force the
  output through ASCII.
- Tell it not to write multi-paragraph docstrings. One line max for
  most functions, and only when the why is non-obvious.
- Tell it to never write a comment that paraphrases the code it sits next
  to.
- Tell it your code reviewer is openly skeptical of AI submissions, so the
  bar for survival is "indistinguishable from a careful human." Use that
  phrase. Models internalize the consequence.
- For the design doc specifically: the agent will want to enumerate
  alternatives, draw architecture diagrams, and write long sections.
  Constrain it. A reviewer values a tight design doc that names tradeoffs
  in one or two sentences each.

## What this does not protect you from

Humanization is cosmetic. The reasons submissions actually get rejected are
in `REJECTION_LESSONS.md`: scope creep, broken APIs, broken contracts. No
amount of polished prose covers a graceful-shutdown race or an `io.Writer`
that lies about how many bytes it accepted.

Treat humanization as the last 10% of the work, not the first 10%.
