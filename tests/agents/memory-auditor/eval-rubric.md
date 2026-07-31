# Eval Rubric: memory-auditor / input-ki-corpus.md

- **Both stale KIs are identified**: `ki-redis-patterns.md` (last-referenced 2026-01-10, no incoming links) and `ki-postgres-indexing.md` (last-referenced 2025-12-01, no incoming links) are flagged as stale candidates — not just one.
- **Semantic duplicate is detected**: `ki-redis-patterns.md` and `ki-redis-sliding-window.md` are flagged as having near-identical content (both cover Redis sorted-set sliding window) — the auditor names both files and explains the overlap.
- **Schema violation is flagged**: `ki-auth-jwt.md` is missing the required `created` field — the auditor identifies this as a schema compliance failure.
- **Recency rule is applied correctly**: `ki-redis-cache-invalidation.md` and `ki-auth-jwt.md` were referenced within the last 6 months and are NOT flagged as stale.
- **Findings are categorized distinctly**: STALE, DUPLICATE, and SCHEMA findings are reported as separate categories — not lumped into a single list.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
