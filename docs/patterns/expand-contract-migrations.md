# Expand/Contract Migrations

Already the single most heavily-gated pattern in this framework — a hard architectural guardrail
(`shared/rules/architecture-guardrails.md` #2), a dedicated skill (`db-migration`), and its own named
approval gate (`shared/rules/approval-gates.md` #4, specifically for the destructive phase). This file
gives it the same Context/Structure/Example/Trade-offs treatment `stability-patterns.md` gave Circuit
Breaker — upgrading it from "a rule that's referenced" to a documented pattern in its own right.

## The Pattern

**Context**: A schema change that would break running application code if applied in one step (renaming
a column, dropping a column, changing a type) needs to be split into phases so old and new code can both
run correctly during the deploy window — no moment where the schema and the deployed code disagree.

**Structure**: Up to three phases:
1. **Expand** — additive only. Add the new column/table without touching the old one. Always safe,
   deployable immediately, rollback is trivial (drop the column that was just added).
2. **Migrate** — backfill data from old to new, with application code writing to both during the
   transition. Deployed once Expand is live and stable.
3. **Contract** — remove the old column/table once nothing reads from it anymore. Destructive,
   deployed only once Migrate is verified and no traffic touches the old column.

`shared/rules/architecture-guardrails.md` #2 states the two hard rules this structure exists to enforce:
never `DROP COLUMN`/`RENAME COLUMN`/`DROP TABLE` in a single-phase migration, never add a `NOT NULL`
column without a `DEFAULT` value (which would break existing inserts the instant the migration lands).

**Example** (from `db-migration/SKILL.md`'s own worked example — renaming `lockout_expires_at` to
`locked_until`):
- **Phase 1 (Expand)**: `ALTER TABLE users ADD COLUMN locked_until TIMESTAMPTZ NULL;` — deploy
  immediately, verify `SELECT COUNT(*) FROM users WHERE locked_until IS NOT NULL` returns 0.
- **Phase 2 (Migrate)**: `UPDATE users SET locked_until = lockout_expires_at WHERE lockout_expires_at IS
  NOT NULL;` — deploy once Phase 1 is live and stable, verify the two columns agree for every row.
- **Phase 3 (Contract)**: `ALTER TABLE users DROP COLUMN lockout_expires_at;` — deploy only once Phase 2
  is verified and no traffic reads the old column.

**Trade-offs**: Three deploys and three verification steps instead of one. That's the entire cost, paid
in exchange for the thing a single-phase destructive migration can't offer at all: a safe rollback point
at every stage before the point of no return, and zero-downtime deploys where the schema and the running
code never have to agree perfectly at the exact same instant.

## The Contract Gate

**Context**: Phase 3 specifically — the actually-destructive step — has its own named approval gate
(`shared/rules/approval-gates.md` #4), separate from the general "running database migrations" gate
(#3). Reaching Expand or Migrate doesn't require special sign-off beyond the normal migration gate; the
Contract phase does, explicitly, because it's the one phase in the sequence that can't be undone.

**Structure**: The gate's own text: "Gate: user must say 'confirm contract phase'." Any edit to the
pending migration plan resets it, same as every other gate in this framework — a plan that gets revised
after approval needs to be re-approved, not silently carried forward.

**Trade-offs**: This means a fully-verified, ready-to-ship Contract phase still waits on an explicit
human confirmation separate from whatever approved Phases 1 and 2. That's deliberate friction at exactly
the one point in the sequence where undoing a mistake is no longer possible.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
