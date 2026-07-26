# Epic 57 — Context-Manifest Test Coverage & Fixtures

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 5.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

`scripts/check-context-budget.sh` validates context manifests but has no fixtures — it currently passes trivially because there's nothing to check. Real manifests get produced by `context-engineer` during `/deliver-feature` runs, land in `.claude/feature-workspace/context-manifest.md` (transient), and then get persisted to `docs/features/<name>/context-manifest.md` after delivery. This repo has no persisted delivered features besides `docs/features/context-engineering-framework/` which may or may not have a real manifest.

Result: the script provides false confidence — it can't actually catch a regression in what `context-engineer` produces.

## Scope

**One commit** (small op). Do:

1. Investigate: `ls docs/features/*/` — enumerate which delivered features have `context-manifest.md` persisted
2. If ≥ 1 real manifest exists: skip fixture creation, wire `scripts/check-context-budget.sh` to loop over `docs/features/*/context-manifest.md` and validate each. Verify at least one existing manifest validates correctly.
3. If 0 real manifests exist: create 2-3 hand-authored `tests/fixtures/context-manifests/` files representing:
   - A well-formed passing manifest (budget under limit, all sections present)
   - A budget-over-limit failing manifest (should fail validation)
   - A missing-section failing manifest (should fail validation)
   - Wire `scripts/check-context-budget.sh` to iterate these fixtures + report per-file pass/fail

4. Update the script's docstring/comments to state clearly what it validates against (real manifests vs. fixtures).

Commit: `test(fixtures): add context-manifest fixtures + wire check-context-budget (Epic 57)`.

## Discipline

- Standard — match other prompts in `docs/prompts/`.
- Fixtures should be realistic — copy structure from the real `context-engineer` output shape, don't invent a new format.
- The script must be idempotent + fast — no network, no cross-repo reads.

## Escalation

- If `check-context-budget.sh` has evolved past what the audit assumes — halt, describe current behavior.
- If real manifests exist but ALL fail the current script's validation — that's a regression the script should have been catching all along; halt and describe.

## Report (under 150 words)

```
Commit: <sha>
Real manifests found under docs/features/: <count>
Chosen strategy: <validate-real-manifests | hand-authored-fixtures>
Fixtures created (if applicable): <count>
Script now validates: <n> manifests
Post-fix state: <pass/fail counts, matches expected>
```

Go.
