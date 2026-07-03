---
name: resume-pipeline
description: Resumes an in-progress or interrupted deliver-feature run from its last checkpoint, jumps to an explicit phase with --from-phase N, or rolls back a specific agent's artifact and re-runs the pipeline from that point. Reads and writes .claude/feature-workspace/pipeline-state.json.
triggers:
  keywords: ["resume-pipeline", "resume delivery", "resume feature", "rollback pipeline", "restart from phase"]
  intentPatterns: ["Resume delivery on *", "Resume the pipeline for *", "/resume-pipeline *", "Roll back * to * and re-run", "Restart delivery on * from phase *"]
standalone: true
---

## When To Use
- `deliver-feature` found an existing `.claude/feature-workspace/pipeline-state.json` for the requested feature (Phase 0, step 3) and handed off here instead of starting over.
- The user explicitly asks to resume an interrupted or crashed delivery run.
- The user asks to jump to a specific phase (`--from-phase N`) — usually after manually fixing something in the workspace.
- The user asks to roll back a specific agent's artifact to its previous version and re-run downstream.

Do NOT use when there's no `pipeline-state.json` for the feature — that's a fresh run, use `deliver-feature` directly.

## Context To Load First
1. `.claude/feature-workspace/pipeline-state.json`
2. The feature file referenced in `pipeline-state.json` (`featureFile`)
3. Every artifact listed in `completedAgents` (to verify checksums before trusting them)
4. `deliver-feature/SKILL.md` (for the numbered step sequence and contract mapping)

## Process

### Mode 1: Resume (default)
1. Read `pipeline-state.json`. For each `completedAgents` entry, recompute the checksum of its artifact and compare. If any mismatch: treat that step and everything after it as not completed (the on-disk file was hand-edited or corrupted since the checkpoint) — report this to the user before proceeding.
2. Determine the true `lastCompletedStep` (the highest step whose checksum still matches, or the value in the state file if all match).
3. Continue executing `deliver-feature`'s numbered steps starting at `lastCompletedStep + 1`, using the existing `.claude/feature-workspace/` artifacts as-is for everything already completed. Do not re-run completed agents.
4. Keep updating `pipeline-state.json` as normal from this point on.

### Mode 2: Jump to phase (`--from-phase N`)
1. Read `pipeline-state.json`.
2. Find the first step number belonging to Phase N in `deliver-feature/SKILL.md`.
3. Mark every `completedAgents` entry at or after that step as `"stale": true` (do not delete — they're useful history) rather than removing them.
4. Resume execution at that step, re-invoking every agent in Phase N onward even if a stale artifact already exists for it.
5. This is a deliberate override — confirm with the user before discarding validated downstream artifacts as stale, since re-running earlier phases invalidates work that later phases already built on.

### Mode 3: Rollback a specific agent's artifact
1. Identify the target agent and its current artifact (e.g., "roll back the developer's implementation-notes.md").
2. Find the most recent file for that artifact in `.claude/feature-workspace/.history/` (filenames are `<artifact-name>.<unix-timestamp>.md`; pick the highest timestamp less than the current artifact's checkpoint time, i.e. the version immediately before the one being discarded).
3. If no history entry exists, stop and tell the user — there's nothing to roll back to; a fresh re-run of that agent is the only option.
4. Restore the history file over the current artifact.
5. In `pipeline-state.json`, mark every `completedAgents` entry for that agent and every step after it as `"stale": true`.
6. Resume execution at that agent's step (re-invoke it, not just validate-artifact — the point of a rollback is usually to re-run with different guidance, not to re-validate the same restored content).

## Output Format
Report back to the user before resuming:

```markdown
# Resume Report: [Feature Name]

## Mode
Resume | Jump to Phase N | Rollback [agent]

## Checkpoint Integrity
- [step]: [agent] — checksum OK / MISMATCH (treated as incomplete)

## Resuming From
Step [N]: [agent/action]

## Stale Entries (if any)
- [agent] (step [N]) — marked stale, will re-run
```

## Guardrails
- **Never** resume past a checksum mismatch — a mismatched artifact means the recorded checkpoint no longer describes what's on disk.
- **Never** silently discard `.history/` files — they're the only rollback mechanism; only `deliver-feature`'s own Phase 4 persistence step touches permanent storage.
- **Confirm before jumping phases or rolling back** — both actions invalidate downstream work; treat them like the CHANGES REQUESTED loop, which already requires explicit acknowledgment before re-running.
- **Do not** re-derive `pipeline-state.json` from scratch by guessing which artifacts look complete — if the state file itself is missing or unreadable, tell the user this is a fresh run, not a resume.

## Standalone Mode
Pure local file operations — reads/writes JSON and markdown files, computes checksums. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
