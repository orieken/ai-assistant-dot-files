# Docs Cleanup — Document or Retire the Legacy Extensionless Install/Uninstall Path

Source: `docs/TODO.md` §"Repository Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

The repo may have (or have had) extensionless `install` and `uninstall` scripts at the root,
distinct from the current `install.sh` / `uninstall.sh`. These legacy paths are referenced in:
- `CLAUDE.md`
- `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md`
- `BLUEPRINT_GENERATOR_PROMPT.md`

`API_FRAMEWORK_BLUEPRINT_PROMPT.md` is a live dependency of `api-ingest` regardless of the
install path decision. The status of `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md` and
`BLUEPRINT_GENERATOR_PROMPT.md` hinges on this decision: if the legacy path is retired and
these files have no other consumer, they may be deleted rather than moved.

This is the **prerequisite prompt** for `docs-cleanup-move-blueprints.md`.

## Scope

**Op 1 — Inventory the current state:**

```bash
ls -la install uninstall 2>/dev/null || echo "No extensionless scripts found"
ls -la install.sh uninstall.sh 2>/dev/null
grep -n "extensionless\|legacy.*install\|\binstall\b" CLAUDE.md 2>/dev/null | head -20
grep -n "install\b" E2E_FRAMEWORK_BLUEPRINT_PROMPT.md BLUEPRINT_GENERATOR_PROMPT.md 2>/dev/null | head -20
```

**Op 2 — Make the determination:**

Based on the inventory, determine whether the legacy extensionless path:

**(A) Does not exist** — the `install`/`uninstall` scripts were never present or have already
been removed. In this case, document the finding, update references in the three Markdown files
to point to `install.sh`/`uninstall.sh`, and close the TODO item.

**(B) Exists and should be kept** — add explicit documentation in `README.md` under an
"Installation" or "Script paths" section stating both paths are supported. Note the supported
path explicitly in `CONTRIBUTING.md` as well.

**(C) Exists and should be retired** — remove the extensionless scripts, update all references
in `CLAUDE.md`, `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md`, and `BLUEPRINT_GENERATOR_PROMPT.md` to
use `install.sh`/`uninstall.sh` instead. If `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md` and
`BLUEPRINT_GENERATOR_PROMPT.md` have no remaining consumers after the reference update,
mark them for deletion in the next cleanup pass (do not delete in this prompt — that is
`docs-cleanup-move-blueprints.md`'s decision).

**Op 3 — Record the decision:**

Create `docs/adrs/adr-NNN-install-script-path-convention.md` documenting:
- Which path(s) are officially supported
- Why the decision was made
- What was retired (if anything)

Use the next available ADR number:
```bash
ls docs/adrs/ | sort | tail -5
```

**Op 4 — Mark TODO items resolved:**

In `docs/TODO.md`, mark the "Decide whether the legacy extensionless install/uninstall path..."
item as `[x]`.

## Guardrails

- Conventional commit for Op 2 changes: `chore(install): retire/document legacy install path`
- Conventional commit for ADR: `docs(adr): record install script path convention decision`
- Do not delete `E2E_FRAMEWORK_BLUEPRINT_PROMPT.md` or `BLUEPRINT_GENERATOR_PROMPT.md` here —
  that follows from `docs-cleanup-move-blueprints.md` after this decision is recorded.
- Stage files explicitly — do not `git add -A`.

## Escalation

Stop and report if:
- The inventory is ambiguous and the correct path cannot be determined from the files alone.
- The extensionless scripts exist and contain logic materially different from `install.sh` —
  they may be separate tools, not just aliases.

## Report

On completion, confirm:
- Whether extensionless scripts were found
- Which determination was made (A/B/C) and why
- What files were updated
- ADR number created
- Commit hash(es)
