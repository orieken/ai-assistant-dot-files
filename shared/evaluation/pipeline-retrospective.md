# Evaluation Spec: pipeline-retrospective

**Analysis logic**: `shared/skills/pipeline-retrospective/SKILL.md` — unchanged
by this spec. That skill is the single source of truth for how the retrospective
is computed; this file is the continuous-trigger contract that Phase 3 will wire
into the hook layer.

**Phase 1 status**: documentation only. The skill remains the sole way this
evaluation runs. Nothing in v3.0 fires it automatically.

## What it evaluates

Cross-delivery trends in pipeline performance: which agent is the slowest, which
is the most-retried, whether trends are improving or degrading, whether an
agent's version bump correlates with a regression. Full detail in the underlying
skill's own doc — this spec does not restate the logic.

## On-demand invocation (v3.0 and forward, unchanged)

```
/pipeline-retrospective
```

The skill reads `docs/features/*/pipeline-trace.json` (the last N deliveries,
default 10), computes the trend table, and writes
`docs/pipeline-retrospectives/retrospective-<YYYY-MM-DD>.md`. This path stays
exactly as it works today. Any team upgrading to v3.0 without touching AOS
layers continues to invoke it the same way.

## Continuous-trigger contract (Phase 3 target, not yet wired)

The Phase 3 hook layer will let the evaluation fire automatically on a
telemetry event. This spec defines the contract Phase 3 must satisfy — it does
not implement anything.

### Trigger event

- **Event type**: `stage.completed` (from the generated `shared/schemas/telemetry/run-event-types.md`; the `agent.completed` type this referenced was specified and never emitted)
- **Filter**: `agent_or_skill_name == "deliver-feature"` AND
  `outcome == "success"`
- **Cadence**: every Nth matching event (default N=5, matching the current
  auto-invocation cadence for the single-delivery `retrospective` skill inside
  `deliver-feature`). Cadence is per-project config, not hardcoded here.

### Inputs when fired continuously

Same as the on-demand skill:

- The last N `pipeline-trace.json` files under `docs/features/*/`
- The most recent `docs/agent-metrics/scorecard-*.md`, if one exists

### Outputs when fired continuously

Same as the on-demand skill: writes
`docs/pipeline-retrospectives/retrospective-<YYYY-MM-DD>.md`. Also emits a
telemetry event of its own:

- `event_type`: `agent.completed`
- `agent_or_skill_name`: `pipeline-retrospective`
- `artifact_path`: the retrospective file it wrote
- `outcome`: `success`
- `metadata`: `{ "trigger": "continuous", "cadence_n": 5, "deliveries_analyzed": <N> }`

The `trigger` field in metadata is the key distinction — evaluators consuming
this event can tell whether the retrospective came from a human invocation
(`"on_demand"`) or the hook layer (`"continuous"`).

### Opt-out

Phase 3 will provide a per-project config knob to disable continuous invocation
without affecting the on-demand path. Default is off — teams opt in explicitly.
Matches Migration Principle #1 (every AOS addition is opt-in) and Migration
Principle #6 (counter agents / continuous evaluations run only when invoked).

## What this spec is NOT

- **Not a rewrite of the skill.** The skill's analysis logic is authoritative.
  If the skill's process changes, that's a skill edit; this spec doesn't need
  to change unless the trigger contract itself changes.
- **Not a hook definition.** Hook definitions are the machine-readable YAML/JSON
  in `shared/hooks/` (Phase 2). This spec is human-readable prose describing
  what that hook should look like.
- **Not a new skill.** Do not invoke this file. Invoke
  `shared/skills/pipeline-retrospective/SKILL.md` — this file describes it, it
  does not replace it.

## Related

- `shared/skills/pipeline-retrospective/SKILL.md` — the authoritative analysis
  logic
- `shared/schemas/telemetry/run-event-types.md` — the generated event vocabulary this spec references
- `shared/evaluation/README.md` — the layer overview and the
  continuous-vs-on-demand model
- `docs/aos/migration-plan.md` — Phase 3 will implement the hook layer that
  makes the continuous-trigger contract executable

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
