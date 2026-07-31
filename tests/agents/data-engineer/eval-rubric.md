# Eval Rubric: data-engineer / input-analysis.md

- **NOT NULL without DEFAULT is flagged as a hard block**: the proposed migration adds `NOT NULL` to an 8M-row table with no `DEFAULT` — this will fail immediately. The agent flags this as a BLOCKER and rewrites the migration as Expand (add nullable), Migrate (backfill), Contract (apply NOT NULL after backfill).
- **Table lock risk is quantified and mitigated**: the agent does not just note "45-second lock" — it recommends `ALTER TABLE ... ADD COLUMN ... DEFAULT NULL` (which is instant in PostgreSQL 11+) or a `pg_repack`/online DDL approach, and explains why plain `ALTER TABLE` locks.
- **Backfill must happen BEFORE code deploy in this case**: the agent catches that the proposed order (migration then code) leaves existing rows with `tenant_id = NULL` even after the code is live, which breaks the NOT NULL constraint — the correct order is Expand → Backfill → Contract, with code deploy happening in the Migrate phase.
- **Index on tenant_id recommended**: queries filtering `WHERE tenant_id = $1` without an index will cause full-table scans on 8M rows. The agent recommends a concurrent index build (`CREATE INDEX CONCURRENTLY`).
- **Rollback plan is explicit**: the Expand migration (nullable column, no constraint) is trivially reversible by dropping the column. The agent calls this out as the safe rollback point.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
