# AOS Migration — Phase 1: Foundations (v3.0)

You are executing Phase 1 of the AOS migration. Scope: 7 operations, all purely additive. Do NOT expand beyond it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits go here directly). Do NOT push; that's the human's step after review.

## Source of truth

Read `docs/aos/migration-plan.md` in full before starting. Phase 1 section defines exactly what to build. The design pack at `docs/aos/AOS_Governance_Design_Pack/` is the underlying vision (00-Vision, 01-Governance-Checks-and-Balances, 03-Memory-Governance, 05-AOS-Directory, and 09-Entropy-Manager are all relevant to Phase 1).

## Backward-compatibility guarantee (non-negotiable)

A team on any prior framework version that upgrades to v3.0 without invoking any new AOS-specific capability MUST see zero behavior change in what they use today. Every addition in Phase 1 is opt-in. If nothing invokes it, nothing changes.

Op 1.7 verifies this identity — do not skip it.

## Scope: 7 ops

Do them in this order:

### Op 1.1 — Create `shared/telemetry/`

Files to create:
- `shared/telemetry/README.md` — what telemetry captures, that it is opt-in, retention convention (append-only JSONL, project-local under `.claude/telemetry/events.jsonl`, never leaves the repo).
- `shared/telemetry/event-schema.md` — human-readable JSON schema for pipeline events. Minimum event types: `agent.invoked`, `agent.completed`, `artifact.written`, `validation.passed`, `validation.failed`. Every event has `timestamp`, `event_type`, `agent_or_skill_name`, `artifact_path` (optional), `outcome` (optional), `metadata` (free-form object).
- `shared/telemetry/event-recorder.md` — a skill (yes, a skill under `shared/telemetry/`, not `shared/skills/` — telemetry is its own top-level concern per design pack `05-AOS-Directory.md`). Skill appends one JSON line per event to `.claude/telemetry/events.jsonl`. Standalone-mode discipline: works offline, no external dependencies, creates the file + parent dir on first write.

Commit: `feat(telemetry): scaffold shared/telemetry/ layer (AOS Phase 1 Op 1.1)`

### Op 1.2 — Create `shared/evaluation/`

Files to create:
- `shared/evaluation/README.md` — model: continuous evaluations (triggered by telemetry events) vs on-demand (existing skill invocation).
- `shared/evaluation/pipeline-retrospective.md` — evaluation spec pointing at the existing `shared/skills/pipeline-retrospective/SKILL.md` (that skill is the on-demand version; the evaluation spec is the "how this gets triggered continuously later" definition — Phase 3 will actually wire the trigger). For Phase 1, the spec is documentation only.

Do NOT move or modify `shared/skills/pipeline-retrospective/SKILL.md`. The skill keeps working exactly as today.

Commit: `feat(evaluation): scaffold shared/evaluation/ layer (AOS Phase 1 Op 1.2)`

### Op 1.3 — Add first counter agent: `memory-auditor`

Create `shared/agents/memory-auditor.md` following the format of existing agents (see `shared/agents/memory-engineer.md` as the paired producer for reference on frontmatter, description style, tools list). Purpose: audits changes to `shared/knowledge/` and `.claude/knowledge/` for KI schema compliance, duplicate detection, stale-tag flags. Should list `Read, Glob, Grep` as tools (read-only). Version starts at `1.0.0`. Frontmatter matches existing conventions exactly.

The agent should:
- Read every KI under both `shared/knowledge/` and `.claude/knowledge/`
- Validate frontmatter against the KI schema documented in `shared/knowledge/README.md`
- Flag duplicates (same title, same tags, or high semantic overlap the agent can spot)
- Flag KIs with stale metadata (last-referenced date > 6 months + no recent linking anywhere in the corpus)
- Produce a report to stdout OR to `.claude/audits/memory-audit-<date>.md`

Commit: `feat(agents): add memory-auditor as first AOS counter agent (Phase 1 Op 1.3)`

### Op 1.4 — Update `shared/skills/health-check/SKILL.md`

Add an "AOS layers" section (opt-in): detect presence of `shared/telemetry/`, `shared/evaluation/`, `shared/hooks/`, `shared/orchestration/`, `shared/rag/`, and any counter agents under `shared/agents/` matching the pattern `*-auditor.md`, `*-evaluator.md`, `*-reviewer.md`, `*-validator.md` (skip existing `security-reviewer`, `code-reviewer`, `pattern-reviewer` etc. that are producers-in-role-name, not counter agents — only flag counter agents that pair with a known producer per the design pack's 15 pairs).

Report format: keep the existing health-check output; add a new "## AOS Layers" section at the bottom that lists which layers are present. Never fails a health-check on missing AOS layers — they're optional.

Bump the skill's frontmatter version (if it has one) by a minor version.

Commit: `feat(health-check): detect and report AOS layers (Phase 1 Op 1.4)`

### Op 1.5 — Update `shared/agents/CHANGELOG.md`

Add a v3.0.0 entry summarizing:
- New: `shared/telemetry/`, `shared/evaluation/`, `memory-auditor` agent
- Updated: `health-check` skill (AOS-layer detection)
- Unchanged: everything else (deliberate — Phase 1 is purely additive)

Include the promise: "A v2.x install upgraded to v3.0 without opting into any AOS capability behaves identically to v2.x."

Commit: `docs(changelog): v3.0.0 — AOS foundations (Phase 1 Op 1.5)`

### Op 1.6 — Add `docs/aos/migration-guide.md`

Empty structural stub for v3.0 (nothing forces adoption yet). Sections:
- "Upgrading from v2.x" — one line: "Nothing to do. All AOS additions are opt-in."
- "Opting into telemetry" — one line: "Not yet — Phase 2 wires it. See migration-plan.md."
- "Opting into memory-auditor" — brief: how to invoke it standalone.
- "What's coming next" — link to migration-plan.md Phase 2.

Commit: `docs(aos): add migration-guide stub (Phase 1 Op 1.6)`

### Op 1.7 — Verify identity install

This is a verification step, not a code change. Do:
1. Note the current install script command and the exact list of files it would install/symlink from a v2.x-era invocation vs a v3.0 invocation.
2. Confirm the v3.0 install produces a strict superset of the v2.x install — no removals, no renames, no moves.
3. If any concrete v2.x pipeline (e.g., `deliver-feature <sample-spec>`) can be run in a scratch environment, run it and verify the artifacts produced match what would have been produced under v2.x (byte-identical is unrealistic; structural-identical is the bar — same artifacts in same paths with same required sections).
4. If a full run isn't feasible in this environment, document the specific checks a human should run before tagging v3.0.

Report finding in the final report. If Op 1.7 finds a regression, halt and escalate — do not commit v3.0.

Commit (if all checks pass): `chore(release): tag AOS Phase 1 v3.0.0 candidate (Phase 1 Op 1.7)` — this commit adds a note to `docs/aos/migration-plan.md`'s Phase 1 section marking checklist items complete.

## Commit discipline (non-negotiable)

- **One commit per op** — 7 commits total (or 8 if Op 1.7 adds a plan-amendment note).
- Conventional Commits format matching the repo's existing style.
- **NEVER `git add -A` or `git add .`** — this repo has untracked pre-existing directories (`docs/aos/AOS_Governance_Design_Pack/` files may be already tracked; but `docs/audits/`, `docs/blog-posts/` are untracked). Always stage explicit paths: `git add shared/telemetry/README.md ...`
- `git status --short` after staging, before commit — verify only intended files.
- Do NOT modify `docs/aos/migration-plan.md` except in Op 1.7 (to mark Phase 1 items complete). It is the source of truth being executed.
- Do NOT push. That is the human's step.

## Escalation criteria — STOP and report back if:

- Op 1.7 reveals any regression from v2.x behavior — this violates the backward-compat guarantee and blocks the phase.
- An existing agent, skill, or rule would need to be moved or renamed to complete a Phase 1 op. Phase 1 is purely additive; if you find yourself needing to modify an existing file beyond the ones listed above, that's a scope issue — halt.
- `memory-auditor` design conflicts with the existing `memory-engineer` agent's contract in a way you can't cleanly resolve.
- The telemetry event schema needs to be more elaborate than described — the design pack's `03-Memory-Governance.md` may inform. Read it first; if still unclear, halt with a proposal.

## Report format (under 300 words)

```
PHASE 1 STATUS: <complete | stopped-at-op-N>

Commits landed:
  <sha> <message>
  ...

Files added:
  shared/telemetry/README.md
  shared/telemetry/event-schema.md
  shared/telemetry/event-recorder.md
  shared/evaluation/README.md
  shared/evaluation/pipeline-retrospective.md
  shared/agents/memory-auditor.md
  docs/aos/migration-guide.md

Files updated:
  shared/skills/health-check/SKILL.md — <what changed>
  shared/agents/CHANGELOG.md — <what changed>
  docs/aos/migration-plan.md — <if Op 1.7 marked items>

Op 1.7 identity check:
  <result — passed, failed with details, or "human verification needed with these checks: ...">

Recommended next step:
  <e.g., "human review + git push + tag v3.0.0, then start Phase 2 handoff prompt">
```

Go.
