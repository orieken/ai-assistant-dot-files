# Docs Cleanup — Fix Broken Links in roadmap-2026-08-07.md

Source: `docs/TODO.md` §"docs/ Root" — link-check result.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/roadmap-2026-08-07.md` contains an Epics 69–73 table. Each row links to a prompt file
using a bare relative path (e.g. `[epic-69-mcp-skill-bundler.md](epic-69-mcp-skill-bundler.md)`).
Because the roadmap lives at `docs/`, those links resolve to `docs/epic-69-*.md` which does not
exist. The actual files live at `docs/prompts/epic-69-*.md`. All five links are broken.

## Scope

**Op 1 — Fix the five broken table links:**

File: `docs/roadmap-2026-08-07.md`

In the Epics 69–73 table, update each link to include the `prompts/` prefix:

| Epic | Current (broken) | Fixed |
|---|---|---|
| 69 | `(epic-69-mcp-skill-bundler.md)` | `(prompts/epic-69-mcp-skill-bundler.md)` |
| 70 | `(epic-70-health-check-autofix.md)` | `(prompts/epic-70-health-check-autofix.md)` |
| 71 | `(epic-71-pipeline-tui.md)` | `(prompts/epic-71-pipeline-tui.md)` |
| 72 | `(epic-72-multi-model-fallback.md)` | `(prompts/epic-72-multi-model-fallback.md)` |
| 73 | `(epic-73-policy-engine-completion.md)` | `(prompts/epic-73-policy-engine-completion.md)` |

Before saving, verify each target file exists:
```
ls docs/prompts/epic-69-mcp-skill-bundler.md
ls docs/prompts/epic-70-health-check-autofix.md
ls docs/prompts/epic-71-pipeline-tui.md
ls docs/prompts/epic-72-multi-model-fallback.md
ls docs/prompts/epic-73-policy-engine-completion.md
```

**Op 2 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Fix the five broken links in `docs/roadmap-2026-08-07.md`" item
as `[x]`.

## Guardrails

- One commit covering both ops.
- Conventional commit: `docs(roadmap): fix broken Epic 69-73 prompt links`
- Stage only: `docs/roadmap-2026-08-07.md` and `docs/TODO.md`
- Do not reformat or rewrite any other content in either file.

## Escalation

Stop and report if any of the five target files are missing — they may have been renamed or moved
since the TODO was written. Do not guess alternate paths.

## Report

On completion, confirm:
- Which five links were updated (before/after)
- Whether each target file was verified to exist
- Commit hash
