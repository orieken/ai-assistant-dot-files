# Telemetry (AOS v3.0, Phase 1)

Telemetry is the AOS layer that captures pipeline events so continuous evaluations
(Phase 3+) and policy engines (Phase 4) have data to reason about. The event layer
described in this directory is still **scaffolding only** — nothing writes
`.claude/telemetry/events.jsonl` by default, and no downstream consumer reads it.

## Two telemetry systems, and which one exists

Do not read this directory as a description of what loom emits today. There are
two systems and only one of them is running.

| | **AOS events** (this directory) | **OTel traces** (`loom run`) |
|---|---|---|
| Status | Specified, **no emitter** | **Shipped** — roadmap L3.8, epic 84 |
| Written by | Nothing. The `event-recorder` skill is a prose spec | The Go executor and `loom mcp serve`, in `internal/telemetry/` |
| Output | `.claude/telemetry/events.jsonl` (does not exist) | OTLP/JSON `traces.jsonl` per run, plus OTLP/HTTP when an endpoint is set |
| Shape | `event_type` + free-form `metadata` | W3C trace/span tree with GenAI semantic conventions |
| Answers | "what happened", eventually | "how long did it take, and what did it cost" |
| Measured by | A model, which cannot observe elapsed time | The process doing the work, by subtraction |

The traces are the reason `duration_ms` in `event-schema.md` was never trustworthy:
a language model reporting its own elapsed time is recalling a number, not measuring
one. Under `loom run`, timing and cost now come from the process that did the work,
and token counts and dollars come from the provider's own reported accounting rather
than from a price table in this repository.

Reconciling the two — one event-type enum, generated schema, emission from the
executor — is roadmap **L3.9**, which is unblocked and not yet built. Until then this
directory describes the markdown pipeline's intended event layer, and
`internal/telemetry/` describes what actually runs.

A third file is easy to confuse with both: `run-events.jsonl`, beside `run-state.json`,
is the executor's own append-only audit log of gates, digests and staleness. It is
deliberately not OTel and deliberately not this event layer — it stays readable with
no collector configured and no exporter running, which is the property an audit record
needs and a trace does not.

## Location convention

- **This directory (`shared/telemetry/`)** — the *layer* definition: what telemetry
  is, the event schema, and the `event-recorder` skill spec. Portable across every
  install.
- **`.claude/telemetry/events.jsonl`** — the *runtime output*: append-only JSONL,
  project-local. Written to the project's own `.claude/` directory, never leaves
  the repo, never uploaded anywhere. Follows the same "project-local, not shared"
  discipline as `.claude/feature-workspace/` and `.claude/knowledge/`.

Telemetry sits alongside `shared/agents/`, `shared/skills/`, `shared/rules/`,
`shared/knowledge/` as a first-class top-level concern per the AOS design pack
(see `docs/aos/AOS_Governance_Design_Pack/05-AOS-Directory.md`). It is deliberately
not a subdirectory of `shared/skills/` — telemetry is infrastructure the skills
layer emits into, not another skill.

## Opt-in nature (non-negotiable, backward-compat guarantee)

Every telemetry emission is **opt-in and fire-and-forget-optional** — no pipeline
correctness depends on it. A team that never invokes `event-recorder` sees zero
behavior change from a pre-v3.0 install. This is Migration Principle #5 in
`docs/aos/migration-plan.md`.

Phase 1 lands the *ability* to record events. Phase 3's hook layer wires
automatic emission from `deliver-feature` and `validate-artifact` into it —
still opt-in per project, still skippable.

## What telemetry captures

Events represent things that happened during a pipeline run:

- `agent.invoked` — an agent started work
- `agent.completed` — an agent finished (success or failure)
- `artifact.written` — a contract-bound artifact was produced
  (`analysis.md`, `architecture-notes.md`, etc.)
- `validation.passed` — `validate-artifact` accepted an artifact
- `validation.failed` — `validate-artifact` rejected an artifact

The full schema, including required and optional fields, is in
`event-schema.md`. New event types get added to that file as new
producers start emitting — the recorder itself is schema-agnostic.

## Retention convention

- **Append-only JSONL.** Never rewrite in place; every event is a new line.
  Existing lines are historical record and are not edited.
- **Project-local.** `.claude/telemetry/events.jsonl` stays inside the project's
  own working tree. It never gets uploaded, synced, or aggregated across projects
  by the framework itself.
- **Rotation is the operator's call.** The framework does not auto-rotate,
  compress, or trim the file. If a project wants to cap size (e.g., rotate
  monthly, or drop lines older than N days), that's a project-level decision
  outside this layer.
- **Never contains secrets.** Producers must not put API keys, tokens,
  credentials, or PII into event `metadata`. Same rule as everywhere else in
  this framework (`shared/rules/architecture-guardrails.md` #3).

## Standalone discipline

The `event-recorder` skill is pure local file writes. No network calls, no
external services, no external dependencies. Works fully offline. Creates the
`.claude/telemetry/` directory and `events.jsonl` file on first write if
either is missing.

## Concurrent-delivery behavior (Epic 63)

As of Epic 63, `deliver-feature` supports multiple named workspaces
(`.claude/feature-workspace/<feature-name>/`) running concurrently. Two
deliveries in-flight both append to the shared `events.jsonl`. Concurrent
appends are safe at the line level (each event is one JSON line) and remain
distinguishable because every event carries a `pipeline_id` field. In the
extremely unlikely case that two simultaneous sub-second writes produce a
corrupted line, `pipeline-retrospective` and consumers that parse JSONL are
expected to skip malformed lines rather than halting — the telemetry layer
is fire-and-forget and never a source of pipeline correctness. No format or
schema change is needed for concurrent use.

## Related

- `shared/telemetry/event-schema.md` — the event format
- `shared/telemetry/event-recorder.md` — the skill spec
- `shared/evaluation/README.md` — how telemetry feeds evaluations (Phase 3)
- `docs/aos/migration-plan.md` — the phased rollout this is part of
- `docs/aos/AOS_Governance_Design_Pack/05-AOS-Directory.md` — why telemetry
  is a top-level layer rather than a skill subdirectory

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
