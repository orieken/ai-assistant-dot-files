---
name: event-recorder
description: Appends a single JSON line to .claude/telemetry/events.jsonl for one pipeline event (agent invoked, agent completed, artifact written, validation passed/failed). Fire-and-forget, opt-in, offline — nothing in the framework invokes this by default in v3.0. Producers call it explicitly when they want the event recorded.
triggers:
  keywords: ["event-recorder", "record telemetry event", "emit telemetry", "append event"]
  intentPatterns: ["Record a telemetry event", "Log this to telemetry", "/event-recorder *"]
standalone: true
---

## When To Use

Producers (agents, skills, hooks) call this skill when they want a single event
recorded to the project's telemetry log. Phase 1 has no automatic invocation —
`deliver-feature`, `validate-artifact`, and every existing agent behave exactly
as they did in v2.x. This skill is the *mechanism* the Phase 3 hook layer will
wire producers into; in v3.0 it's directly callable only when a caller explicitly
wants to emit an event.

Do NOT use this skill to:

- Read events back — that's the future evaluation layer's job (Phase 3).
- Query trends — use `pipeline-retrospective` (single-run aggregation across
  `docs/features/*/pipeline-trace.json`, a separate stream).
- Batch multiple events — call this skill once per event. It's cheap because
  it does one append and returns.

## Context To Load First

1. `shared/telemetry/event-schema.md` — the schema every recorded event must
   conform to. If the caller's event doesn't match a documented `event_type` +
   required-fields shape, halt and tell the caller — do not silently write an
   invalid line.

## Process

### 1. Validate the input

The caller supplies:

- `event_type` (required, must be one of the types documented in
  `event-schema.md`)
- `agent_or_skill_name` (required)
- `timestamp` (optional — if omitted, this skill uses "now" in UTC, ISO-8601
  with millisecond precision, `Z` suffix)
- `artifact_path` (optional)
- `outcome` (optional)
- `metadata` (optional, must be a JSON-serializable object)

If `event_type` is not one of the documented types in `event-schema.md`, refuse
and report which types are valid. Never write an undocumented event type — that
breaks the schema contract downstream evaluations depend on.

If `metadata` contains anything that looks like a secret (keys named `token`,
`password`, `api_key`, `secret`, `credential`, values matching an obvious
credential shape), refuse and tell the caller to strip it. Same rule as
everywhere else in the framework — no secrets in telemetry.

### 2. Ensure the target file exists

Target path: `.claude/telemetry/events.jsonl` under the current project's root
(not under `shared/`).

- If `.claude/` does not exist, create it.
- If `.claude/telemetry/` does not exist, create it.
- If `events.jsonl` does not exist, create it (empty file).

This is the "creates the file + parent dir on first write" discipline the
Phase 1 handoff prompt requires. First invocation in a fresh project must not
fail because the directory hasn't been made yet.

### 3. Serialize the event

Build the JSON object with fields in this order (for readability, not
correctness — JSON is order-insensitive but a stable order makes the file
easier to eyeball):

1. `timestamp`
2. `event_type`
3. `agent_or_skill_name`
4. `artifact_path` (only if provided or the type requires it)
5. `outcome` (only if provided)
6. `metadata` (only if provided)

Serialize with no embedded newlines, no trailing whitespace, no trailing
comma. The line must be exactly one line ending in `\n`.

### 4. Append

Open `.claude/telemetry/events.jsonl` in append mode. Write the serialized
line followed by `\n`. Close. Done.

Do not rewrite the file, do not read existing lines, do not sort — this is a
strict append. That's what makes the log correct under concurrent producers
(each append is atomic on POSIX for lines under `PIPE_BUF` bytes, and typical
events are far under that).

### 5. Confirm

Return the exact path written to and the byte offset the line landed at (or
just "recorded" if the caller doesn't need the offset). Do not echo the full
event body back — the caller already had it.

## Output Format

Success:

```
Recorded: .claude/telemetry/events.jsonl (event_type=<type>, agent_or_skill_name=<name>)
```

Refusal (invalid input):

```
Refused: <reason — invalid event_type, missing required field, suspected secret in metadata, etc.>
Valid event_types per shared/telemetry/event-schema.md: agent.invoked, agent.completed, artifact.written, validation.passed, validation.failed
```

## Guardrails

- **Never** invent a new `event_type` on the fly — refuse if the caller passes
  one not documented in `event-schema.md`. The schema is added to first, then
  producers emit.
- **Never** write anything that looks like a secret. Same rule as
  `shared/rules/architecture-guardrails.md` #3.
- **Never** put unbounded data in `metadata` (full artifact bodies, log dumps,
  stack traces). Reference paths, don't inline content.
- **Never** rewrite, sort, or delete lines in `events.jsonl`. Append-only.
- **Never** make a network call. This is a pure local file write.
- **Never** raise or propagate an error that would break a pipeline. Telemetry
  is fire-and-forget-optional per Migration Principle #5 — if the write fails
  (permission denied, disk full), report it to the caller and return, do not
  abort whatever pipeline the caller is running.

## Standalone Mode

Pure local file write. No external services, no network, no shell tools beyond
whatever the runtime uses to append a line. Works fully offline. Creates its
own parent directory and file on first invocation.

## Related

- `shared/telemetry/README.md` — layer overview
- `shared/telemetry/event-schema.md` — the schema this skill enforces
- `shared/evaluation/README.md` — Phase 3+ evaluations that consume this log
- `docs/aos/migration-plan.md` — the phased rollout this is part of

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
