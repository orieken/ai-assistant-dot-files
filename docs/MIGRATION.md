# Migration: v1 -> v2 (Canonical `shared/` Layer)

This document is for anyone who cloned or forked this repo **before** the `shared/` canonical layer
restructure (commit `cc841a8`, "establish canonical shared/ layer with cross-platform install"). If you
cloned recently, none of this applies to you — skip straight to [README.md](../README.md).

## What changed

### Before (v1)
```
.claude/
  agents/*.md        <- real files, the only source of truth
  skills/*/SKILL.md   <- real files
  rules/*.md          <- real files
ARCHITECTURE_RULES.md  <- real file, repo root
DOMAIN_DICTIONARY.md   <- real file, repo root
```
Every other platform's config (`.cursorrules`, `.openai.md`, etc.) was hand-maintained separately and
regularly drifted from `.claude/`'s content — there was no single source of truth and no way to detect
drift.

### After (v2)
```
shared/
  agents/*.md         <- real files, THE source of truth
  skills/*/SKILL.md    <- real files
  rules/*.md           <- real files
  contracts/           <- new: required-section contracts (Epic 5)
  knowledge/            <- new: Knowledge Items (Epic 14)
  templates/            <- new: tutorial content (Epic 19)
  ARCHITECTURE_RULES.md
  DOMAIN_DICTIONARY.md
  platform-registry.json <- new: tier/capability/format per platform

.claude/{agents,skills,rules}/  <- now symlinks to shared/ equivalents
ARCHITECTURE_RULES.md            <- now a symlink to shared/ARCHITECTURE_RULES.md
DOMAIN_DICTIONARY.md             <- now a symlink to shared/DOMAIN_DICTIONARY.md
.cursor/rules/*.mdc, .windsurfrules, .github/copilot-instructions.md,
.github/instructions/*.instructions.md, AGENTS.md, .openai.md   <- now generated from shared/, never hand-edited
```

## Is this a breaking change for me?

- **If you only ever used Claude Code** and never hand-edited `.claude/agents/*.md` directly (only read
  them): nothing breaks. The content is identical, just relocated with a symlink pointing back at it.
- **If you hand-edited `.claude/agents/*.md` or `.claude/skills/*/SKILL.md` directly**: your edits are
  preserved by the migration (it moves files, never deletes them), but going forward you must edit
  `shared/agents/`/`shared/skills/` instead — editing through the `.claude/` symlink still works (it's the
  same file), but editing a *generated* platform config (anything under `.cursor/`, `.github/`, etc.) will
  now be silently overwritten the next time `scripts/generate-configs.sh` runs.
- **If you hand-maintained `.cursorrules`, `.openai.md`, or similar separately**: those are now generated
  from `shared/` — any manual edits to them will be lost on the next `generate-configs.sh` run. Move that
  content into the relevant `shared/rules/` file instead (see
  [CONTRIBUTING.md](CONTRIBUTING.md), "Adding a new rule").

## How to migrate

```bash
# Preview what would change first
bash scripts/migrate-v1-to-v2.sh --dry-run

# Apply it — moves .claude/{agents,skills,rules} and root ARCHITECTURE_RULES.md/DOMAIN_DICTIONARY.md
# into shared/, then symlinks the old locations back to the new ones. Never deletes anything; if a
# shared/ target already exists, that item is skipped rather than overwritten.
bash scripts/migrate-v1-to-v2.sh

# Verify the result
bash scripts/health-check.sh --verbose

# Regenerate every other platform's config from the now-canonical shared/ source
bash scripts/generate-configs.sh
```

The script is idempotent — running it again after a successful migration reports everything as already
migrated (`[skip]`) rather than erroring or duplicating anything.

## Tags

- `v1.0.0` — the last commit before the restructure began (`e7d5557`).
- `v2.0.0` — the completed canonical `shared/` layer, including everything built across Epics 1-21 in
  `docs/features/context-engineering-framework/TODO.md`.
