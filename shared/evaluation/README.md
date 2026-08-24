# Evaluation (AOS v3.0, Phase 1)

The evaluation layer is where **continuous** quality checks live — the automated,
event-driven counterparts to the manually-invoked skills under `shared/skills/`.

In v3.0 (Phase 1) this layer is **specifications only**. Nothing runs
continuously yet. Every existing evaluation-shaped skill (`pipeline-retrospective`,
`agent-scorecard`, `context-audit`, `retrospective`) keeps working exactly as it
does today, invoked manually on demand. The specs here document how each will
also be triggerable by telemetry events once Phase 3 wires the hook layer.

## The model: continuous vs on-demand

The AOS design pack (see
`docs/aos/AOS_Governance_Design_Pack/00-AOS-Vision.md`, Continuous Improvement
principle) treats evaluation as a first-class ongoing activity, not a periodic
ceremony. But the framework already has valuable evaluation skills that fire
on human intent. The migration keeps both paths, on purpose:

| Path | How it's triggered | Example (today) | Example (Phase 3+) |
|---|---|---|---|
| **On-demand** | Human types the skill name or a matching intent phrase | User runs `/pipeline-retrospective` after 10 deliveries | (unchanged — still available) |
| **Continuous** | A telemetry event fires a hook that invokes the same underlying analysis | Not wired yet | `on-retrospective-written` hook auto-invokes the learning engine |

Every evaluation spec in this directory has two sides:

1. **Analysis logic** — what the evaluation actually does. In v3.0 this always
   points to the existing on-demand skill under `shared/skills/`. The skill
   file is the single source of truth for the logic; the spec here is a
   pointer, not a copy.
2. **Continuous trigger contract** — what event(s) *would* fire this evaluation
   under the Phase 3 hook layer, and what its inputs/outputs look like when
   fired that way. In v3.0 this is documentation only — nothing reads it.

This split lets Phase 3 add hooks without changing any skill's behavior. The
skill keeps working exactly as it does today; the hook layer just gains the
ability to call it automatically when a matching event lands.

## Backward-compatibility guarantee

Same guarantee as the telemetry layer: an install upgraded to v3.0 that doesn't
invoke anything new sees zero change. The evaluation layer adds no runtime
behavior in Phase 1 — it's a documentation layer that Phase 3 will hook into.

Per Migration Principle #6 in `docs/aos/migration-plan.md`: counter agents and
continuous evaluations run only when invoked. The default remains "structural
check only."

## What lives here in v3.0

- `README.md` (this file) — the layer overview
- `pipeline-retrospective.md` — evaluation spec pointing at the existing
  `shared/skills/pipeline-retrospective/SKILL.md` for the analysis logic; adds
  the continuous-trigger contract for later phases.

More evaluation specs will land in Phase 2 (paired with new counter agents)
and Phase 3 (paired with the Learning/Forgetting engines).

## What does NOT live here

- The existing on-demand skills — those stay in `shared/skills/`. This layer
  never *moves* an existing skill; it points at it.
- Hook definitions — those land in `shared/hooks/` in Phase 2.
- Policy files — those land in `shared/policies/` (or equivalent) in Phase 4.

## Quality Signals (Phase 3 candidates)

When the hook layer lands (Phase 3), evaluations will gain access to telemetry signals beyond timing
and iteration count. The signal introduced in Epic 62 is **gate-rejection rate**:

- **Per-agent gate-rejection rate**: across all `gate_decision` events (see `shared/telemetry/event-schema.md`)
  tied to a given `agent_or_skill_name`, what fraction had outcome `rejected` or `edited_then_approved`
  vs. `approved`? A consistently high edit rate for a specific agent is evidence that the agent's
  output often needs human correction before a gate passes — a quality signal `agent-scorecard` should
  surface alongside its four existing scored metrics once enough baseline data exists.

In v3.0 (Phase 1), `agent-scorecard` surfaces raw gate-correction counts from `events.jsonl` when the
file exists (informational only — no scored floor yet). `extract-lessons` mines patterns across multiple
deliveries and surfaces them as candidate prompt improvements. `retrospective` shows per-delivery gate
corrections when telemetry is enabled. The scored floor and automated hook are Phase 3 work.

## Related

- `shared/telemetry/README.md` — the event source these evaluations will
  eventually subscribe to
- `shared/skills/pipeline-retrospective/SKILL.md` — the on-demand skill this
  layer's first evaluation spec points at
- `docs/aos/migration-plan.md` — the phased rollout this layer is part of
- `docs/aos/AOS_Governance_Design_Pack/00-AOS-Vision.md` — the Continuous
  Improvement principle this layer serves

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
