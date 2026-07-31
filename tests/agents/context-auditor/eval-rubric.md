# Eval Rubric: context-auditor / input-manifest.md

- **Pinned-but-never-read file is flagged**: `user.model.ts` is identified as speculative — it was pinned but never opened during implementation, bloating the budget.
- **Unrelated ADR is called out**: `ADR-003-enterprise-memory-sync.md` appears in a rate-limiting manifest and is flagged as an irrelevant pin (not a rate-limiting concern).
- **Budget discrepancy is quantified**: the 7,200-token gap (18,400 claimed vs 11,200 actual) is explicitly called out — not just noted as "off."
- **Broken KI link is identified**: `ki-jwt-validation.md` is flagged as a broken link because the file does not exist in `shared/knowledge/` or `.claude/knowledge/`.
- **Valid entries are not false-flagged**: `ki-redis-patterns.md`, `ADR-007`, and the prior delivery reference are not flagged as problems — only the four issues above are raised.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
