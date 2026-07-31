# Eval Rubric: developer / input-analysis.md

- **StorageAdapter interface introduced**: S3 is not imported directly in the handler — all storage calls go through an interface defined in the consumer layer, consistent with the dependency inversion rule in architecture-guardrails.md.
- **Migration is additive-only (Expand/Contract)**: the `0031_add_avatar_url.sql` migration only adds a nullable column — no `NOT NULL` without a `DEFAULT`, no `DROP`, no `RENAME`. Commit should call out that this is an Expand phase.
- **Concurrent upload race is addressed explicitly**: either a database-level unique constraint on `(user_id)` in the upload table, a mutex, or an advisory lock — not glossed over as "shouldn't happen in practice."
- **Old S3 key deletion is handled**: the implementation notes describe deleting the previous avatar key from S3 when a new one is uploaded, and account for the possibility that the old key may not exist yet (first upload).
- **Atomic write semantics for S3 + DB**: the implementation notes acknowledge that S3 upload must succeed before the DB row is updated — the failure path is explicit (upload fails → no DB write, not a partial-update state).

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
