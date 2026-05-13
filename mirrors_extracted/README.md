# Extracted source trees

This directory contains the working-tree snapshot of each candidate
submission's default branch at the time of archival. It mirrors the
structure of [`../mirrors/`](../mirrors/) but is the checked-out
source rather than a bare `.git` repository — useful for reading
code, grepping, and quickly cross-referencing reviewer quotes against
the actual implementation.

Each subdirectory is the equivalent of running

```
git clone mirrors/<track>/<owner>__<repo>/<repo>.git mirrors_extracted/<track>/<owner>__<repo>
```

and then removing the `.git/` directory. For full commit history,
use the bare mirror in [`../mirrors/`](../mirrors/) instead.

Compiled binaries that the original candidates committed by accident
(e.g., a couple of pre-built Go binaries) have been stripped. Yarn
bundles and demo media included by the candidates are retained as
part of the captured submission.
