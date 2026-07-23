# Design an automated deliver-feature workflow

Take the existing `/deliver-feature` skill from "orchestrator with prose-defined human PAUSE checkpoints" to something that can run substantially more autonomously — via policy-driven auto-approval gates rather than removing gates entirely. Preserve the framework's stage-3 stance (see `docs/aos/migration-plan.md`).

**Scope of this prompt is DESIGN + PROMPT-DRAFTING, not implementation.** Produce a design doc + a follow-up implementation prompt. No source code changes to the skill itself.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo.

## Prior context (read first)

- `shared/skills/deliver-feature/SKILL.md` — the current orchestrator. Chains ~10 agents (analyst → architect → developer → code-reviewer → …), pauses at human checkpoints via prose instructions
- `shared/rules/approval-gates.md` — the 8 non-negotiable irreversible-action gates (ship, commit, migration, contract-phase, external API, out-of-boundary file writes, fitness function, deploy). These are the gates that **must stay human** or must be policy-gated with an audit trail
- `docs/aos/migration-plan.md` §Phase 3 Ops 3.11-3.13 — the "trinity-native workflow refactor" that turns `deliver-feature` from a skill into a proper `FeatureDeliveryWorkflow` (stateful, resumable, with policy evaluation points)
- `docs/aos/migration-plan.md` §Phase 4 (v3.3) — the policy layer that defines auto-approval rules. Example from the plan: *"if code-reviewer approves AND security-reviewer approves AND all fitness functions pass AND diff size < N → auto-proceed past code-reviewer gate."*
- `docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md` — establishes the "graduated approach" pattern (start simple, add complexity when telemetry shows the need). Same shape applies here — graduated automation.

## The tension to navigate

The framework's stated stance: **stage-3 excellence over stage-4 fantasy** (see the "Where does the framework fit in AI adoption stages" writeup embedded in `docs/aos/AOS_Governance_Design_Pack/00-AOS-Vision.md`'s six principles). Fully removing human gates trades discipline for throughput.

Better framing: **policy-driven graduated automation**. Different pipeline stages have different risk profiles. Doc updates and test additions are lower risk than schema migrations and external mutations. A policy layer lets each team decide which stages become auto-approve for which conditions, per project. The gates that must never become policy-only (per `approval-gates.md`) stay human, always.

## What "automated" could mean — three tiers to consider

The design should articulate which tier(s) are in scope:

- **Tier A — Auto-continue on green** (the plan's Phase 4 example): if a stage's outputs are structurally valid AND all fitness functions pass AND diff size < threshold, the pipeline auto-proceeds past that stage's approval checkpoint. Human still sees a summary; can interrupt. Non-controversial per the plan.
- **Tier B — Auto-retry on structural failure**: if `validate-artifact` rejects an artifact for missing-heading violations, auto-invoke the producer agent again with the specific rejection list — up to N retries — before surfacing to a human. Reduces the "same LLM drift, three different times, keep pasting the same feedback" pattern that currently burns time.
- **Tier C — Unattended-until-blocked**: pipeline runs end-to-end without prompts unless a policy explicitly requires human review or an escalation criterion fires. The current pipeline is the opposite: prompts by default, silent only where opted in. This is the biggest philosophical shift and needs the tightest policy story.

## Prerequisites (design should acknowledge)

Automation depends on infrastructure that hasn't fully landed yet:

- **AOS Phase 1** — telemetry (SHIPPED as v3.0.0 candidate `a274300`). Automation decisions need audit trails.
- **AOS Phase 2** — counter agents + hooks (NOT SHIPPED). Automation needs the auditor half of each producer/counter pair to run automatically.
- **AOS Phase 3** — `FeatureDeliveryWorkflow` as a first-class Workflow with state machine + checkpoint hooks (NOT SHIPPED). Skills-chaining-agents can't be automated cleanly; a real workflow with state can.
- **AOS Phase 4** — the policy layer itself (NOT SHIPPED). No automation without policy definition.

The design should identify what can be done TODAY (interim automation) vs. what must wait for Phase 3/4.

## Deliverables

Produce two files in `docs/aos/` (or wherever fits the framework's doc structure — check `docs/adrs/ADR-001-adopt-rag-friendly-docs-structure.md` for placement guidance):

### 1. Design document: `docs/aos/automated-delivery-design.md`

Sections:
- **Goal** — clear statement of which of Tier A/B/C are in scope. Recommend Tier A + Tier B for v1; defer Tier C until telemetry proves the pipeline is stable enough
- **Non-negotiable gates** — restate the 8 gates from `approval-gates.md` and mark which can EVER become policy-driven (Tier A candidates) vs. which stay human, always. Recommendation: gates #2 (git commit), #6 (out-of-boundary files), #7 (fitness function wiring) can safely become Tier A with the right policy; gates #1 (Friday ship), #3-4 (migrations), #5 (external API mutations), #8 (deploy) stay human forever
- **Policy schema sketch** — YAML shape a project would use to opt into auto-approval per pipeline stage. Example rules the design should include
- **Interim automation (pre-Phase-4)** — what can be done in the CURRENT `/deliver-feature` skill without waiting for the full AOS runtime. E.g., auto-retry Tier B could be added as a loop inside the current skill's PAUSE-and-re-run pattern
- **Rollback design** — how a policy-approved change gets undone if it lands broken. Every policy-gated stage MUST have a rollback path documented in the policy itself
- **Telemetry requirements** — every policy decision emits a telemetry event so `agent-scorecard` and `pipeline-retrospective` can measure "did automation help?" over time
- **Escalation criteria** — the design should explicitly enumerate when a running pipeline halts to a human regardless of policy (e.g., three consecutive tier-B retries fail, security-reviewer flags a `[RISK-HIGH]`, an ADR-worthy decision surfaces)

### 2. Handoff prompt: `docs/prompts/implement-automated-delivery-tier-a.md`

A separate ready-to-fire handoff prompt for the FIRST implementation slice — Tier A only, interim (no Phase 4 dependency). Follows the same shape as other docs/prompts/*.md files. Scope: extend the current `/deliver-feature` skill to check for a `.claude/delivery-policy.yaml` file, apply Tier A rules to the gates that safely support it, log all decisions to telemetry.

## Discipline

- **This subagent does NOT modify the deliver-feature skill.** Only produces the design doc + follow-up prompt.
- **Two commits, one per deliverable**:
  - `docs(aos): add automated delivery design`
  - `docs(prompts): add implement-automated-delivery-tier-a handoff`
- Conventional Commits. **NEVER `git add -A`** — this repo has known untracked directories (`docs/audits/`, etc.).
- Update `docs/prompts/README.md` to add the new implement-* handoff to the index (a third commit, `docs(prompts): index automated-delivery-tier-a handoff`).
- Do NOT push.

## Escalation criteria

Stop and report if:
- The design would require modifying `shared/rules/approval-gates.md` — that's a bigger change that needs its own ADR. Halt and propose the ADR shape instead
- The Tier A implementation prompt would need to touch the skill AND multiple agents simultaneously — that's really a Phase 3 workflow-refactor concern; halt and note the dependency
- Any of the 8 gates in `approval-gates.md` looks like it should be re-classified — halt with the specific gate + rationale; don't silently reclassify

## Report format (under 200 words)

```
STATUS: complete | stopped-at-<reason>

Commits:
  <sha> docs(aos): add automated delivery design
  <sha> docs(prompts): add implement-automated-delivery-tier-a handoff
  <sha> docs(prompts): index automated-delivery-tier-a handoff

Design doc: docs/aos/automated-delivery-design.md
  - Tiers scoped in v1: <A | A+B | A+B+C>
  - Gates classified as policy-eligible: <count> / 8
  - Interim (pre-Phase-4) automation identified: <yes | no>

Follow-up prompt: docs/prompts/implement-automated-delivery-tier-a.md
  - Ready to fire against v3.x: <yes | no — reason>
  - Prereqs not yet met: <list, or "none">

Notes for the user:
  - <any gate re-classification recommendations, ADR needs, or Phase-3/4 blockers>
```

Go.
