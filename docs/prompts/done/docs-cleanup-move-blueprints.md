# Docs Cleanup — Move Blueprint Prompt Files to docs/blueprints/

Source: `docs/TODO.md` §"Repository Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

Three blueprint files live at the repo root:
- `API_FRAMEWORK_BLUEPRINT_PROMPT.md` — live dependency of the `api-ingest` skill
- `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md` — consumer status depends on legacy install path decision
- `BLUEPRINT_GENERATOR_PROMPT.md` — consumer status depends on legacy install path decision

**Prerequisite**: this prompt assumes the legacy extensionless `install`/`uninstall` path has
been evaluated (see `docs-cleanup-legacy-install-decision.md`). If the legacy path is still
supported, root filenames must be preserved via symlinks or the install script must be updated
to reference the new path. If it has been retired, E2E and BLUEPRINT files may be deleted
instead of moved.

This prompt handles the move-to-`docs/blueprints/` path. Do not proceed if the legacy install
decision has not been made — stop and report.

## Scope

**Op 0 — Verify prerequisite:**

Check whether the legacy install decision has been documented:
```bash
ls install uninstall 2>/dev/null
grep -n "docs/blueprints\|legacy.*install\|extensionless" docs/TODO.md | head -10
```

If `install` or `uninstall` (extensionless) still exist at root AND the TODO item for the
legacy path decision is not yet marked `[x]`, stop and report — do not proceed.

**Op 1 — Create docs/blueprints/ and move files:**

```bash
mkdir -p docs/blueprints
git mv API_FRAMEWORK_BLUEPRINT_PROMPT.md docs/blueprints/API_FRAMEWORK_BLUEPRINT_PROMPT.md
git mv E2E_FRAMEWORK_BLUEPRINT_PROMPT.md docs/blueprints/E2E_FRAMEWORK_BLUEPRINT_PROMPT.md
git mv BLUEPRINT_GENERATOR_PROMPT.md docs/blueprints/BLUEPRINT_GENERATOR_PROMPT.md
```

**Op 2 — Update all consumers:**

Find and fix every reference to the old root paths:

```bash
grep -rn "API_FRAMEWORK_BLUEPRINT_PROMPT\|E2E_FRAMEWORK_BLUEPRINT_PROMPT\|BLUEPRINT_GENERATOR_PROMPT" \
  . --include="*.md" --include="*.sh" --include="*.yaml" --include="*.json" --include="*.ts"
```

Known consumers to update:
- `shared/skills/api-ingest/SKILL.md` — references `API_FRAMEWORK_BLUEPRINT_PROMPT.md`
- `docs/patterns/sunday-framework-patterns.md` — may reference blueprint files
- `scripts/health-check.sh` — root-Markdown scan may reference these filenames
- `install` / `uninstall` scripts — if they reference these files, update path or add symlink

For each reference found, update the path to `docs/blueprints/<filename>`.

**Op 3 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Move `API_FRAMEWORK_BLUEPRINT_PROMPT.md`..." item as `[x]`.

## Guardrails

- Conventional commit: `docs(blueprints): move blueprint prompt files to docs/blueprints/`
- Stage files explicitly with individual `git add` calls — do not `git add -A`.
- Verify `scripts/check-parity.sh` passes after the move (run it if it exists).
- If any consumer cannot be safely updated in this pass, stop and report before committing.

## Escalation

Stop and report if:
- The legacy install decision has not been made (Op 0 check).
- Any consumer reference cannot be updated cleanly (complex conditional logic, generated files).
- `scripts/check-parity.sh` fails after the move.

## Report

On completion, confirm:
- Whether the Op 0 prerequisite was satisfied
- Which files were moved
- Which consumers were updated (list file + line)
- Whether `check-parity.sh` passed
- Commit hash
