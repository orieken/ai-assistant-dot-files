# Telemetry Event Schema (AOS v3.0, Phase 1)

**Schema version: v1.1.0**
- v1.0.0 — initial Phase 1 schema (agent.invoked, agent.completed, artifact.written, validation.passed, validation.failed)
- v1.1.0 — added `gate_decision` event type; extended `outcome` enum with gate-specific values; noted undocumented `policy.evaluated` and `contract.retry` types emitted by `deliver-feature` (schema entry pending)

Downstream aggregators must read this version line to detect schema evolution.

Human-readable schema for events written to `.claude/telemetry/events.jsonl`.
Format is **JSONL** — one JSON object per line, newline-delimited, append-only.
Every event is a self-contained, order-independent record.

## Wire format

Each line is a UTF-8-encoded JSON object with no embedded newlines and no
trailing comma. Example:

```json
{"timestamp":"2026-07-22T14:32:11.482Z","event_type":"agent.invoked","agent_or_skill_name":"analyst","artifact_path":null,"outcome":null,"metadata":{"pipeline_id":"deliver-feature-abc123","feature":"user-registration"}}
```

The `event-recorder` skill (see `event-recorder.md`) is responsible for
serializing to a single line before writing.

## Required fields (every event)

| Field | Type | Description |
|---|---|---|
| `timestamp` | string (ISO-8601 UTC, millisecond precision) | When the event occurred. Always UTC with `Z` suffix (`2026-07-22T14:32:11.482Z`), never a local-timezone offset. |
| `event_type` | string | One of the event types below. New types are added by updating this file first, never invented ad hoc by producers. |
| `agent_or_skill_name` | string | The name of the agent or skill emitting the event (e.g., `analyst`, `code-reviewer`, `deliver-feature`, `validate-artifact`). Matches the `name:` frontmatter of the emitter. |

## Optional fields

| Field | Type | Description |
|---|---|---|
| `artifact_path` | string \| null | Path to the artifact this event refers to, relative to the project root (`.claude/feature-workspace/analysis.md`). Omit or `null` if the event isn't tied to a specific artifact. |
| `outcome` | string \| null | One of `success`, `failure`, `changes_requested`, `skipped` (general events); or `approved`, `rejected`, `edited_then_approved` (gate_decision events only — see below). Omit or `null` for events that don't have a binary/enum outcome (e.g., `agent.invoked`). |
| `metadata` | object | Free-form JSON object for producer-specific context. Must NOT contain secrets, PII, credentials, or unbounded data (log lines, full artifact bodies). Suggested keys below. |

### Suggested `metadata` keys (conventions, not required)

- `pipeline_id` — a stable identifier for the pipeline run (`deliver-feature`
  generates one; other orchestrators may too).
- `feature` — feature slug from `docs/features/<slug>/` if applicable.
- `duration_ms` — for `agent.completed`, how long the agent ran.
- `iteration` — for retry loops (e.g., `code-reviewer` CHANGES_REQUESTED cycles).
- `violations` — for `validation.failed`, a list of the specific rule names that
  failed. Rule *names* only, not full rule bodies.

Producers may add new metadata keys as needed — the schema is deliberately
extensible on this axis. But keep the same keys stable across events so
downstream evaluations can aggregate.

## Event types (Phase 1 minimum)

### `agent.invoked`

Emitted when an agent or skill starts work.

```json
{"timestamp":"2026-07-22T14:32:11.482Z","event_type":"agent.invoked","agent_or_skill_name":"analyst","metadata":{"pipeline_id":"df-abc123","feature":"user-registration"}}
```

- `artifact_path`: usually null (nothing written yet).
- `outcome`: null.

### `agent.completed`

Emitted when an agent finishes, regardless of outcome. Pair with the earlier
`agent.invoked` via `pipeline_id` + `agent_or_skill_name`.

```json
{"timestamp":"2026-07-22T14:34:02.117Z","event_type":"agent.completed","agent_or_skill_name":"analyst","artifact_path":".claude/feature-workspace/analysis.md","outcome":"success","metadata":{"pipeline_id":"df-abc123","duration_ms":110635}}
```

- `outcome`: `success` | `failure` | `changes_requested`.
- `artifact_path`: the contract-bound artifact this agent produced, if any.

### `artifact.written`

Emitted when any contract-bound artifact is written to disk. May duplicate the
signal from `agent.completed` for agents that write exactly one artifact — that's
fine, the two events answer different questions (one about agent lifecycle, one
about artifact provenance).

```json
{"timestamp":"2026-07-22T14:34:02.140Z","event_type":"artifact.written","agent_or_skill_name":"analyst","artifact_path":".claude/feature-workspace/analysis.md","metadata":{"pipeline_id":"df-abc123","contract":"analysis-contract.md"}}
```

- `artifact_path`: required (this is what the event is *about*).
- `outcome`: null.

### `validation.passed`

Emitted by `validate-artifact` when an artifact passes its structural contract
check.

```json
{"timestamp":"2026-07-22T14:34:04.882Z","event_type":"validation.passed","agent_or_skill_name":"validate-artifact","artifact_path":".claude/feature-workspace/analysis.md","outcome":"success","metadata":{"pipeline_id":"df-abc123","contract":"analysis-contract.md"}}
```

### `validation.failed`

Emitted by `validate-artifact` when an artifact fails its structural contract
check. Include the specific rule names that failed in `metadata.violations`.

```json
{"timestamp":"2026-07-22T14:34:04.882Z","event_type":"validation.failed","agent_or_skill_name":"validate-artifact","artifact_path":".claude/feature-workspace/analysis.md","outcome":"failure","metadata":{"pipeline_id":"df-abc123","contract":"analysis-contract.md","violations":["missing-required-section:acceptance-criteria","missing-required-section:definition-of-done"]}}
```

### `gate_decision`

Emitted **after** a human responds to a gate halt (`shared/rules/approval-gates.md`, gates 1–8).
The event captures whether the artifact was approved as-is, rejected outright, or edited by the
human before approving — the last being the richest corrective signal the system receives.

Emitted only when telemetry is enabled (opt-in guarantee is non-negotiable — see
`shared/telemetry/README.md`). Never emitted silently or as a side-effect of default operation.

```json
{"timestamp":"2026-07-22T14:52:00.182Z","event_type":"gate_decision","agent_or_skill_name":"deliver-feature","artifact_path":".claude/feature-workspace/analysis.md","outcome":"edited_then_approved","metadata":{"pipeline_id":"df-abc123","feature":"user-registration","gate_id":1,"gate_name":"friday_ship","reason":null}}
```

**Outcome values (gate_decision only)**

| Outcome | Meaning |
|---|---|
| `approved` | Human confirmed with the gate's approval word; artifact checksum unchanged since the halt was presented. |
| `rejected` | Human declined or said "no"; pipeline halted or rolled back. |
| `edited_then_approved` | Human edited the artifact (checksum changed between gate-presented and gate-approved) before confirming. This is the corrective-signal case `extract-lessons` and `retrospective` mine. |

**Executor note (L2.14).** For runs executed by `loom run`, edited-then-approved is no longer a
heuristic: an approval binds to the digests of everything completed when it was given, and the
executor invalidates it in code when any of them changes, recording the gate, the stage responsible,
and the time in `run-state.json` and on the `run-events.jsonl` timeline (`gate.invalidated`). The
paragraph below describes the **markdown pipeline's** path, where the checksums are the model's own;
`loom state verify` computes the same detection in Go for that pipeline when the binary is present.
Emitting this as a `gate_decision` event is still unbuilt — the emitter is **L3.9**, and this schema
is not changed here.

**Edit detection heuristic**: compare the artifact's current checksum against the checksum recorded
for that step in `pipeline-state.json` at the moment the gate halt was presented. If they differ,
the artifact was edited — use `edited_then_approved` as the outcome. Checksum format matches
`pipeline-state.json`'s existing `sha256:<hash>` entries.

**Required `metadata` keys** (in addition to the standard suggested keys):

| Key | Type | Description |
|---|---|---|
| `gate_id` | integer (1–8) | The gate number from `approval-gates.md`. |
| `gate_name` | string | Short slug for the gate: `friday_ship`, `git_commit`, `db_expand_migrate`, `db_contract`, `external_api`, `out_of_boundary_write`, `fitness_function`, `deploy`. |
| `reason` | string \| null | Optional free-text human explanation, verbatim from whatever the human typed after approval/rejection, if any. Short only — never a full artifact body. `null` when no reason given. |

**Optional `metadata` keys** (same suggested keys as other events):
`pipeline_id`, `feature`, `iteration` (for gate halts that can repeat, e.g. Gate #2 on re-staged commits).

- `artifact_path`: the artifact that was gated (e.g. the analysis.md presented at the analyst checkpoint, or null for Gate #1 Friday ship which doesn't gate a single artifact).
- `agent_or_skill_name`: the skill that owned the gate halt (`deliver-feature`, `ship-feature`, `db-migration`, etc.).

## Adding a new event type

1. Add a new subsection to this file with the type name, required fields,
   optional fields, and one worked example.
2. Bump the file header to note the new type.
3. Only then update the producer to emit it. Emitting an undocumented event
   type is a schema violation — evaluators can't reason about types they don't
   know exist.

## Guardrails

- **Never** put secrets, tokens, API keys, credentials, PII, or customer data
  into any event field, including `metadata`. Same rule as everywhere else in
  the framework (`shared/rules/architecture-guardrails.md` #3).
- **Never** put unbounded data into `metadata` — no full artifact bodies, no
  log dumps, no stack traces. Reference paths instead (`artifact_path`) and
  keep metadata to short scalars and lists of names.
- **Never** edit or delete a historical line. The file is append-only. If a
  past event was wrong, the correction is a new event, not a rewrite.
- **Never** invent new `event_type` values without adding them here first.
- **Cardinality discipline:** event types should be a small, closed set of
  short strings. `metadata` keys should also be a small, reusable set. Avoid
  putting high-cardinality data (UUIDs, timestamps, error message text) in
  event types or metadata keys — put them in metadata *values* instead.

## Related

- `shared/telemetry/README.md` — layer overview and opt-in policy
- `shared/telemetry/event-recorder.md` — the skill that writes these events
- `docs/aos/migration-plan.md` — the phased rollout this is part of

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
