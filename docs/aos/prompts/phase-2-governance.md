# AOS Migration — Phase 2: Governance Skeleton (v3.1)

You are executing Phase 2 of the AOS migration. Scope: 8 operations, all purely additive on top of Phase 1. Do NOT expand beyond it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits go here directly). Do NOT push; that's the human's step after review.

## Prerequisite

**Phase 1 (v3.0) must be complete before starting Phase 2.** Verify:
- `shared/telemetry/` and `shared/evaluation/` directories exist with README + schema files
- `shared/agents/memory-auditor.md` exists (v1.0.0)
- `shared/skills/health-check/SKILL.md` has an "AOS Layers" detection step
- `shared/agents/CHANGELOG.md` has a v3.0.0 entry with the backward-compat promise

If Phase 1 isn't complete, halt and execute `docs/aos/prompts/phase-1-foundations.md` first.

## Source of truth

`docs/aos/migration-plan.md` — Phase 2 section defines exactly what to build. Also read:
- `docs/aos/AOS_Governance_Design_Pack/01-Governance-Checks-and-Balances.md` — the 15 producer/counter pairs
- `docs/aos/AOS_Governance_Design_Pack/03-Memory-Governance.md` — informs memory-auditor cousins
- `shared/agents/memory-auditor.md` — the pattern exemplar every new counter agent should follow (frontmatter shape, read-only tool list, "audit-only, produces findings, never mutates" discipline)

## Backward-compatibility guarantee (non-negotiable)

A team on v3.0 that upgrades to v3.1 and doesn't invoke any new counter agent or configure any hook MUST see zero behavior change. Every addition in Phase 2 is opt-in. `validate-artifact` remains structural-only by default; counter-agent invocation is a per-project config opt-in.

## Scope: 8 ops

Do them in this order:

### Op 2.1 — Create 10 audit-relationship counter agents

Each new agent follows `shared/agents/memory-auditor.md`'s shape: frontmatter (name/description/tools/model/version), read-only tools (Read/Glob/Grep), "audit-only" disposition. Version starts at `1.0.0`. Each produces a findings report; none mutate the artifact they audit.

Files to create in `shared/agents/`:
- `context-auditor.md` — audits `.claude/feature-workspace/context-manifest.md` produced by `context-engineer` for pruning discipline (files referenced but never read, KI links that resolve to nothing)
- `knowledge-auditor.md` — audits `create-ki` skill output for KI schema compliance (frontmatter shape from `shared/schemas/ki-frontmatter.schema.json`, semantic overlap with existing KIs)
- `prompt-evaluator.md` — audits agent/skill prompt files for prompt-engineering hygiene (no fabricated URLs, no leaked secrets in examples, consistent voice)
- `agent-evaluator.md` — promotes existing `agent-eval` skill logic into an agent. Runs golden-file eval against `shared/agents/` frontmatter contracts and prompt-behavior expectations
- `rule-auditor.md` — audits `shared/rules/*.md` for internal consistency (conflicting statements across files, dead references)
- `pattern-reviewer.md` — audits `docs/patterns/*.md` for accuracy against the current codebase state
- `tool-validator.md` — audits `shared/skills/*` for the standalone-mode declaration and any hidden MCP dependencies
- `documentation-auditor.md` — audits `README.md`, `docs/architecture.md`, and other prose docs for staleness against the current agent/skill inventory
- `retrieval-evaluator.md` — audits KI + ADR corpus retrievability (per ADR-002 telemetry) — flags "queries that never matched" as either missing-KI or bad-metadata cases
- `privacy-auditor.md` — pairs with existing `security-reviewer`. Audits pipeline artifacts for accidental PII inclusion, secrets in prompts, etc.

**One commit per agent** — 10 commits. `feat(agents): add <name> counter agent (AOS Phase 2 Op 2.1)`.

Update `shared/agents/CHANGELOG.md` with one entry per agent added (batch these into one CHANGELOG commit at the end of Op 2.1 — count as its own 11th commit or fold into the last agent commit, your call).

### Op 2.2 — Create 4 opposing-force skill pairs (8 skills)

The design pack's Phase 2 language mentions Phase 2 lands 4 opposing-force pairs. Each pair is TWO skills at `shared/skills/<name>/SKILL.md`, following `shared/skills/SKILL_TEMPLATE.md`. Neither side runs by default — both are opt-in invocations.

Pairs to create:
- `shared/skills/memory-expansion/` + `shared/skills/memory-compression/` — expansion promotes KIs from lessons-learned; compression deduplicates + summarizes stale KIs
- `shared/skills/learning-engine/` + `shared/skills/forgetting-engine/` — learning proposes new KIs from retrospectives (draft-mode, human-approved); forgetting flags stale KIs for expiration (draft-mode, human-approved)
- `shared/skills/cost-optimizer/` + `shared/skills/quality-optimizer/` — cost recommends cheaper models/agent combos when quality-eval permits; quality flags places to trade cost for higher-fidelity models
- `shared/skills/scheduler/` — Orchestrator side handled by `deliver-feature` today; add an explicit Scheduler skill for cron/hook-driven runs (single skill, not a pair — the "Orchestrator" side is already the existing `deliver-feature`)

**One commit per skill** — 7 commits (4 expansions + 3 compressions/counterparts + 1 scheduler = 8; if you group Learning+Forgetting into one commit for the pair, adjust downward). `feat(skills): add <name> skill (AOS Phase 2 Op 2.2)`.

### Op 2.3 — Create `shared/hooks/` layer

Files:
- `shared/hooks/README.md` — hook definition format, opt-in nature, event catalog (from telemetry event schema)
- `shared/hooks/hooks-schema.md` — YAML/JSON schema for hook config; document the `on-event → skill-or-agent` shape
- 2-3 example hook definitions under `shared/hooks/examples/`:
  - `on-artifact-write.yaml` — invokes telemetry event-recorder
  - `on-validation-pass.yaml` — invokes the corresponding auditor from Op 2.1
  - `on-ki-created.yaml` — invokes knowledge-auditor

**One commit**: `feat(hooks): scaffold shared/hooks/ layer with examples (AOS Phase 2 Op 2.3)`.

### Op 2.4 — Update `validate-artifact` to optionally invoke auditor

Extend `shared/skills/validate-artifact/SKILL.md`:
- Add an "opt-in auditor invocation" mode. When enabled per-project via config (e.g., `.claude/validate-artifact.yaml`), after passing structural check the skill invokes the corresponding counter agent (per Op 2.1 producer→counter mapping) and reports its findings alongside the structural pass.
- Default behavior UNCHANGED — no config → structural check only, no auditor invocation.
- Bump the skill's version.

**One commit**: `feat(validate-artifact): support opt-in auditor invocation (AOS Phase 2 Op 2.4)`.

### Op 2.5 — Document producer/counter pairs

Create `docs/aos/governance-pairs.md` — cross-referenced table of the 15 pairs from `01-Governance-Checks-and-Balances.md` mapped to concrete agents/skills that now exist in the repo (post-Op-2.1 + Op-2.2). Include what "Producer" means for each pair (some are agents, some are humans authoring markdown), and what the counter role's invocation entrypoint is.

**One commit**: `docs(aos): document producer/counter pairs (AOS Phase 2 Op 2.5)`.

### Op 2.6 — Update health-check

Extend `shared/skills/health-check/SKILL.md`:
- Detect + report all counter agents present under `shared/agents/*-auditor.md`, `*-evaluator.md`, `*-reviewer.md`, `*-validator.md` (skip existing producer-role names like `code-reviewer`, `security-reviewer`, `pattern-reviewer` — the AOS-added counter agents follow specific naming; the health-check should distinguish)
- Validate hook config files under `shared/hooks/examples/` and `.claude/hooks/` (if exists) against the hooks-schema.md
- Never fails on absence of these — additive, opt-in check

Bump the skill's version.

**One commit**: `feat(health-check): detect counter agents + validate hooks (AOS Phase 2 Op 2.6)`.

### Op 2.7 — CHANGELOG v3.1 entry

Update `shared/agents/CHANGELOG.md` with a consolidated v3.1.0 entry: 10 new counter agents (with individual version rows), 7-8 new skills, hooks layer, validate-artifact bump, health-check bump. Restate the backward-compat promise.

**One commit**: `docs(changelog): v3.1.0 — AOS governance skeleton (AOS Phase 2 Op 2.7)`.

### Op 2.8 — Identity install verification

Same shape as Phase 1 Op 1.7:
- `bash scripts/health-check.sh --verbose` — expect 0 FAILs. The 2 pre-existing WARNs from Phase 1 are still acceptable.
- `bash scripts/check-parity.sh` — expect PASS.
- Confirm no hook fires by default (opt-in guarantee holds).
- Confirm `validate-artifact` still passes structural check without invoking auditors when no config exists.

If any regression, halt and escalate.

**One commit** (if all checks pass): `chore(release): tag AOS Phase 2 v3.1.0 candidate (AOS Phase 2 Op 2.8)` — also updates `docs/aos/migration-plan.md`'s Phase 2 section to mark all checklist items complete + record the Op 2.8 verification result.

## Commit discipline (non-negotiable)

- **One commit per op-part** — Phase 2 is ~20 commits total (10 for Op 2.1 counter agents + 7-8 for Op 2.2 skills + 1 hooks + 1 validate-artifact + 1 governance-pairs + 1 health-check + 1 changelog + 1 verify).
- Conventional Commits format matching the repo's existing style.
- **NEVER `git add -A` or `git add .`** — the repo has known pre-existing untracked directories (`docs/audits/`, `docs/blog-posts/`, `.gitignore M`). Stage explicit paths only.
- `git status --short` after staging, before commit — verify only intended files.
- Do NOT modify existing agents, skills, or rules unless the op explicitly says to (Op 2.4 for validate-artifact, Op 2.6 for health-check, Op 2.7 for CHANGELOG).
- Do NOT push.

## Escalation criteria — STOP and report back if:

- The 10 counter agents each need substantially different tools (some read-only, some need Bash, some need Write) — halt, describe. Read-only should be the default; deviations need justification.
- An opposing-force pair needs a mutual-invocation coordination pattern the plan doesn't specify — halt, propose.
- `validate-artifact`'s config schema conflicts with an existing convention — halt, describe.
- Any producer/counter mapping (Op 2.5) reveals a producer role that doesn't exist as an agent OR a human role today — halt, describe. May be a design gap in the plan.
- Health-check WARN count grows (currently 2 pre-existing) — investigate immediately.

## Report format (under 400 words)

```
PHASE 2 STATUS: <complete | stopped-at-op-N>

Commits landed (order):
  <sha> <message>
  ...

Op-by-op tally:
  Op 2.1 — 10 counter agents: <count actually landed>
  Op 2.2 — 7-8 opposing-force skills: <count actually landed>
  Op 2.3 — hooks layer: <landed | partial>
  Op 2.4 — validate-artifact opt-in auditor: <landed>
  Op 2.5 — governance-pairs.md: <landed>
  Op 2.6 — health-check counter+hook detection: <landed>
  Op 2.7 — CHANGELOG v3.1.0: <landed>
  Op 2.8 — Identity check: <passed | failed with details>

Total agent count after Phase 2: <n> (was 21 post-Phase-1)
Total skill count after Phase 2: <n> (was <M> post-Phase-1)

Any deviations from the plan (per-agent tool list, opposing-force interpretations, etc.):
  <list>

Recommended next step:
  <e.g., "human review + git push + tag v3.1.0, then execute docs/aos/prompts/phase-3-runtime.md">
```

Go.
