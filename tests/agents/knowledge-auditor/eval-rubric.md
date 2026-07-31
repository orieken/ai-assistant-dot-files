# Eval Rubric: knowledge-auditor / input-ki-draft.md

- **Frontmatter schema is validated**: the auditor checks required fields against the KI schema (`title`, `slug`, `domain`, `tags`, `created`, `last-referenced`) and calls out any missing or malformed fields.
- **`RATE_LIMIT_MAX` is flagged as a magic number in the snippet**: the code example references an undeclared constant — the auditor flags this as an issue in the KI's embedded code (the KI teaches a pattern with a hidden magic number).
- **Semantic duplicate check is run**: the auditor checks for existing KIs on Redis or rate-limiting that overlap and reports whether any duplicates were found (even if the result is "none found").
- **ADR-007 cross-reference is verified**: the `See also: ADR-007` reference is checked — the auditor confirms or flags whether `docs/adrs/ADR-007-rate-limiting-strategy.md` exists.
- **No fabricated findings**: the auditor does not invent issues beyond what the input actually contains — it doesn't flag the `zadd`/`zremrangebyscore` approach as wrong without evidence.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
