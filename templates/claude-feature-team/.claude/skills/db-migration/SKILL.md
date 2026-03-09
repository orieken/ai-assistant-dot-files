---
name: db-migration
description: Safe zero-downtime database migrations using Expand/Contract.
triggers:
  keywords: ["migration", "schema", "ALTER TABLE", "schema change"]
  intentPatterns: ["Write a migration for *", "Add column * to *", "Rename column * to *", "Drop column/table *", "Change column type *"]
standalone: true
---

## When To Use
Triggered whenever a database schema change is requested. Every breaking change must be split into two or three phases via Expand/Contract.

## Context To Load First
1. The requested database changes
2. `ARCHITECTURE_RULES.md`

## Process
1. Ask: "What is the schema change?" (one sentence)
2. Classify: does it require Expand/Contract or is it single-phase safe?
3. If Expand/Contract: generate three migration files, clearly named Phase 1/2/3
4. For each phase: write the migration SQL, the rollback SQL, and the data integrity check
5. Write the backfill script for Phase 2 (if data migration required)
6. Specify the deployment sequence and the verification step at each phase
7. Check: does this migration need a feature flag to safely deploy Phase 2?

## Output Format
`.claude/feature-workspace/db-migration-[description].md`

```markdown
# DB Migration: [Description]

## Change Summary
[What is changing and why]

## Pattern
Expand/Contract (3 phases) | Single-phase safe

## Phase 1 — Expand
```sql
-- Forward
ALTER TABLE users ADD COLUMN locked_until TIMESTAMPTZ NULL;

-- Rollback
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
```
Deploy when: immediately
Verify: `SELECT COUNT(*) FROM users WHERE locked_until IS NOT NULL;` → expect 0

## Phase 2 — Migrate (if required)
```sql
-- Application writes to both old and new columns
UPDATE users SET locked_until = lockout_expires_at WHERE lockout_expires_at IS NOT NULL;
```
Deploy when: Phase 1 is live and stable
Verify: `SELECT COUNT(*) FROM users WHERE lockout_expires_at != locked_until;` → expect 0

## Phase 3 — Contract
```sql
-- Forward
ALTER TABLE users DROP COLUMN lockout_expires_at;
```
Deploy when: Phase 2 backfill verified, no traffic using old column
Verify: `\d users` → `lockout_expires_at` column absent

## Risk Assessment
- Downtime risk: None / Low / [explain]
- Rollback possible at each phase: Phase 1 ✅ / Phase 2 ✅ / Phase 3 ⚠️ [explain]

## Checklist
- [ ] Phase 1 SQL reviewed
- [ ] Phase 2 backfill script tested
- [ ] Phase 3 only runs after Phase 2 verified
- [ ] No destructive operations in Phase 1 or Phase 2
- [ ] Rollback SQL tested
```

## Guardrails
- `DROP COLUMN` or `DROP TABLE` NEVER appears in Phase 1 — only in Phase 3
- `NOT NULL` without `DEFAULT` NEVER appears in a single-phase migration
- `RENAME COLUMN` is NEVER used directly — always Expand/Contract instead
- Phase 3 requires explicit user confirmation: "confirm contract phase"
- Never write a migration that cannot be rolled back in Phase 1 or Phase 2

## Standalone Mode
Pure SQL generation. No external tools required.
