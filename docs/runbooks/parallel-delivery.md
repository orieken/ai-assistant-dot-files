# Runbook: Running Multiple Features in Parallel

How to work on two or more features concurrently without the pipelines colliding.

**Implemented**: Epic 63 (`v3.3.8`). Design rationale: `docs/aos/parallel-delivery-isolation-design.md`.

---

## The Problem (Before Epic 63)

The `deliver-feature` pipeline wrote all artifacts to a single directory:

```
.claude/feature-workspace/
  pipeline-state.json
  analysis.md
  implementation-notes.md
  ...
```

Starting a second feature in the same project overwrote the first feature's state. There was no safe way to have two deliveries in-flight at once.

---

## The Solution: Named Workspaces

Each feature now gets its own subdirectory:

```
.claude/feature-workspace/
  user-registration/
    pipeline-state.json
    analysis.md
    implementation-notes.md
    ...
  payment-flow/
    pipeline-state.json
    analysis.md
    ...
```

The `<feature-name>` slug is the same kebab-case name you pass to `deliver-feature` and that appears in `docs/features/<feature-name>/`.

---

## How to Run Two Features in Parallel

**The simplest approach today: use separate Claude Code sessions.**

Each session drives its own pipeline. Since each pipeline writes to `.claude/feature-workspace/<feature-name>/`, there is no artifact collision between sessions.

```
Terminal A                          Terminal B
──────────────────────────────────  ──────────────────────────────────
/deliver-feature user-registration  /deliver-feature payment-flow
→ writes to feature-workspace/      → writes to feature-workspace/
  user-registration/                  payment-flow/
```

You can also drive both from one session sequentially — `deliver-feature` will detect an in-progress workspace and hand off to `resume-pipeline` rather than overwriting it.

---

## Constraints

### 1. Git index is shared

Both deliveries operate on the same git working tree and the same git index. Two concurrent deliveries MUST NOT attempt Gate #2 (git commit) at the same time.

**Rule**: coordinate your commit gates sequentially. When both pipelines reach Gate #2, confirm one commit, then the other. A human reviews `git diff` at each gate anyway — that review naturally serializes the commits.

### 2. Test runner ports and `/tmp`

If both deliveries run tests simultaneously and the test suite picks a fixed port or writes to a shared `/tmp` path, they can collide. This is user-space and outside the framework's control — use different ports per run if your test config allows it.

### 3. `events.jsonl` concurrent writes

There is no shared telemetry file to contend over: `.claude/telemetry/events.jsonl` was retired in roadmap L3.9. Each `loom run` writes its own `run-events.jsonl` inside its own feature workspace, so two parallel deliveries share no append target. The interleaving risk this section used to document no longer exists.

---

## Legacy Migration (Automatic)

If you have an existing in-flight delivery from before Epic 63 (flat workspace layout with no `<feature-name>/` subdirectory), `deliver-feature` detects and migrates it automatically:

1. Reads `.claude/feature-workspace/pipeline-state.json` for the `feature` field.
2. Moves all workspace files into `.claude/feature-workspace/<feature>/`.
3. Emits a `workspace.migrated` telemetry event.
4. Continues from the last checkpoint.

If `feature` is absent from the state file, the name `default` is used and the user is informed. The migration is non-destructive (move, not copy-then-delete) and idempotent.

---

## Resuming a Suspended Delivery

If a delivery was interrupted mid-run, `deliver-feature` detects the named workspace and hands off to `resume-pipeline` automatically (Phase 0, step 3). You can also invoke it directly:

```
/resume-pipeline <feature-name>
```

Or jump to a specific phase:

```
/resume-pipeline <feature-name> --from-phase 5
```

---

## Fitness Function

After any change to pipeline skill or agent files, run:

```bash
grep -r '\.claude/feature-workspace/[^<{]' \
  shared/skills/deliver-feature/ \
  shared/skills/deliver-atdd/ \
  shared/skills/resume-pipeline/ \
  shared/skills/pipeline-trace/ \
  shared/skills/context-engineer/ \
  shared/skills/validate-artifact/ \
  shared/skills/orchestrate/ \
  shared/agents/analyst.md \
  shared/agents/architect.md \
  shared/agents/developer.md \
  shared/agents/code-reviewer.md \
  shared/agents/security-reviewer.md \
  shared/agents/performance-engineer.md \
  shared/agents/data-engineer.md \
  shared/agents/accessibility-engineer.md \
  shared/agents/sre-engineer.md \
  shared/agents/qa-engineer.md \
  shared/agents/visual-qa-engineer.md \
  shared/agents/tech-writer.md \
  shared/agents/devops-engineer.md \
  shared/agents/privacy-auditor.md \
  shared/agents/context-auditor.md \
  shared/agents/context-engineer.md
```

Expected result: **zero matches** (two exceptions allowed — the legacy-detection and migration-step lines in `deliver-feature/SKILL.md` that intentionally reference the old flat path).

---

## What Was Not Changed

Standalone skills (`five-whys`, `adr`, `threat-model`, etc.) and non-pipeline agents (`chaos-engineer`, `release-manager`, etc.) still write to `.claude/feature-workspace/` as a scratch area. They are invoked without a feature-name context and do not participate in `deliver-feature`'s concurrency model.

---

*Design decisions, isolation model ruling, and blast-radius inventory: `docs/aos/parallel-delivery-isolation-design.md`.*
