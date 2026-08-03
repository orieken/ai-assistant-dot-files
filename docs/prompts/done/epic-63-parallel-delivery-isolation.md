# Epic 63 — Parallel-Delivery Workspace Isolation

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #3 — unblocks team-scale
use). The gap: one feature in flight per project, ever.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

`.claude/feature-workspace/` is a singleton. Everything assumes it:

- `deliver-feature` Phase 0 halts if a `pipeline-state.json` for another run exists (correct
  crash-safety behavior, but it means two features — or two humans, or a human plus a scheduled
  agent — cannot deliver concurrently).
- `resume-pipeline` (Modes 1–3), rollback via `.claude/feature-workspace/.history/`,
  `pipeline-state.json` checksums, `context-manifest.md`, `pipeline-trace.json` — all read/write
  the singleton path.
- `.gitignore` ignores `.claude/feature-workspace/` wholesale.
- Two candidate isolation models: (a) workspace-per-feature
  (`.claude/feature-workspace/<feature-name>/`), (b) git-worktree-per-feature (each delivery in
  its own worktree, workspace stays singleton *within* each worktree).

## Scope

**Phase A — Design (one commit, then PAUSE for user approval):**

Commit `docs(aos): parallel-delivery isolation design (Epic 63 Phase A)` ruling on:

1. **Isolation model**: (a) vs (b) vs both-layered. Weigh: (a) is simpler but two deliveries
   still share one git working tree (staged-file collisions, test-run interference); (b) gives
   true isolation but adds worktree lifecycle management and confuses tools that assume repo
   root. Recommend one.
2. **Backward compatibility**: an existing install mid-delivery with the OLD singleton layout
   must resume cleanly. Propose the detection + migration story (e.g., singleton layout detected
   → treat as legacy workspace named `default`).
3. **Cross-cutting inventory**: enumerate EVERY file that hardcodes the singleton path (grep
   `feature-workspace` across `shared/` — expect deliver-feature, resume-pipeline,
   pipeline-trace, context-engineer, validate-artifact, several agents' pre-read blocks) and the
   change each needs.
4. **Concurrency limits**: does anything actually break with N>1 (e.g., telemetry events.jsonl
   appends from two runs — fine; Friday reporting — check)? Name the real constraints.

**Phase B — Implementation (after approval; one commit per touched subsystem, sequenced so the
repo stays consistent at every commit):**

Expected shape (Phase A may amend): path-resolution helper first (single place that maps
feature name → workspace dir, legacy-aware), then deliver-feature, then resume-pipeline, then
pipeline-trace/context-engineer/validate-artifact, then agent pre-read blocks, then docs
(`docs/patterns/deliver-feature-workflow.md`, RUNBOOKS).

After every commit: `bash scripts/health-check.sh` green; final verification = two scratch
deliveries with distinct feature names progressing independently (drive each at least through
Phase 0 + one artifact, confirm no cross-contamination of state files or history dirs).

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- The cross-cutting inventory exceeds ~15 files needing edits — halt after Phase A regardless
  (that's the design-approval point anyway) and flag the blast radius explicitly.
- Worktree model chosen but any core skill breaks under a non-root working directory — halt with
  specifics.
- Any change would break `resume-pipeline` for a legacy singleton workspace — halt; crash
  recovery of in-flight deliveries is non-negotiable.

## Report (under 150 words)

```
Phase A commit: <sha>
Ruling: <a|b|layered> — <one-line rationale>
Singleton-path inventory: <n> files
Phase B commits: <sha> <subsystem> ...
Legacy migration: <detection + behavior>
Two-delivery verification: <what was run, result>
```

Go.
