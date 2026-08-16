# Docs Cleanup — Reconcile docs/human-tasks.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/human-tasks.md` tracks decisions and tasks that require human action. It has drifted:
- Three prompt paths now live under `docs/prompts/done/` (shipped), not the queued table.
- The Phase 2 handoff is complete.
- The roadmap (`docs/roadmap-2026-08-07.md`) states Phase 4 is complete.
- The `saturday-mcp` external entries have not been separately verified.

## Scope

**Op 1 — Read the file in full:**

Read `docs/human-tasks.md`. Identify:
1. Any rows whose prompt files now exist under `docs/prompts/done/` (shipped).
2. Any Phase 2 or Phase 4 entries that are marked as pending but are now complete per the roadmap.
3. The `saturday-mcp` external entries — note them separately for step 3.

**Op 2 — Clean up resolved items:**

For rows in the queued table that are resolved:
- Either remove them, OR
- Move them to a clearly labeled `## Completed` or `## Historical` section at the bottom of the file.

Do not silently delete — if there is doubt about whether an item is truly resolved, leave it
and note it in the report.

**Op 3 — Isolate saturday-mcp entries:**

For any rows referencing the external `saturday-mcp` repository:
- Add a visual separator or sub-heading (`## External — saturday-mcp`) so they are clearly
  distinguished from in-repo tasks.
- Do not attempt to verify or update them — they require human judgment about an external repo.

**Op 4 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Reconcile `docs/human-tasks.md`" item as `[x]`.

## Guardrails

- Conventional commit: `docs(human-tasks): remove resolved items and isolate saturday-mcp entries`
- Stage only `docs/human-tasks.md` and `docs/TODO.md`.
- Do not modify the `saturday-mcp` entries' content — only their placement/labeling.
- When in doubt about whether an item is resolved, leave it and flag it in the report.

## Escalation

Stop and report if the file's structure is significantly different from what the TODO describes
(e.g. the "queued table" does not exist or has a different format). Flag for human review before
making structural changes.

## Report

On completion, confirm:
- How many rows were removed or archived as resolved
- Which Phase entries were cleaned up
- Whether `saturday-mcp` entries were separated (and how many there are)
- Any items left in place due to uncertainty (list them)
- Commit hash
