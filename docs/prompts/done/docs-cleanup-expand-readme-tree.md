# Docs Cleanup — Expand docs/README.md Directory Tree

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/README.md` contains a directory tree describing the contents of `docs/`. The tree is
incomplete — it is missing several files and subdirectories that are now present. The TODO
identifies the following missing entries:
- `THREAT_MODEL.md`
- `human-tasks.md`
- `roadmap-2026-08-07.md`
- `TODO.md` (this cleanup ledger)
- `aos/` subdirectory
- `prompts/` subdirectory
- `runbooks/` subdirectory

## Scope

**Op 1 — Establish current directory contents:**

```bash
ls docs/
ls -d docs/*/
```

**Op 2 — Update the directory tree in docs/README.md:**

Find the directory tree section in `docs/README.md` and add entries for all files and
subdirectories currently present that are not listed. Match the existing style (file entries
with a one-line description, directory entries with a trailing `/`).

Use these descriptions for the known missing items:
- `THREAT_MODEL.md` — STRIDE threat model; implementation-status annotations (Epic 65)
- `human-tasks.md` — decisions and tasks that require human action, not agent execution
- `roadmap-2026-08-07.md` — 2026-08-07 framework roadmap with Epics 69–73
- `TODO.md` — folder-by-folder documentation audit ledger
- `aos/` — AOS (Agent Orchestration System) design docs and phase prompts
- `prompts/` — self-contained agent handoff prompts for framework improvements
- `runbooks/` — operational runbooks (debug, environment, incident response)

For any other directories found by `ls -d docs/*/` not listed above, inspect them briefly
and write a one-line description in the same style.

**Op 3 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Expand `docs/README.md`" item as `[x]`.

## Guardrails

- Conventional commit: `docs(readme): expand directory tree to reflect current docs/ structure`
- Stage only `docs/README.md` and `docs/TODO.md`.
- Do not rewrite prose sections of `docs/README.md` — only update the tree.
- Do not add entries for files that don't exist.

## Escalation

Stop and report if the directory tree section cannot be located in `docs/README.md` — the
file structure may have changed and the tree may need to be created from scratch (requires
human decision on format).

## Report

On completion, confirm:
- Which entries were added to the tree
- Whether any directories were found that weren't in the TODO list (and how they were described)
- Commit hash
