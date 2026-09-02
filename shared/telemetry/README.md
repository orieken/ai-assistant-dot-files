# Telemetry

What `loom run` records about a run, and what it does not.

This directory used to describe an event layer — `.claude/telemetry/events.jsonl`, written by an
`event-recorder` skill, carrying fifteen event types. Six of those types were documented here, nine
were specified in other files, and the recorder was instructed to refuse any type not documented.
Several skills carried instructions to emit into it. Nothing ever verified that any of it happened,
and the v3.0.0 release check recorded that a real delivery did not create the file at all.

That layer is retired (roadmap **L3.9**). What follows is what exists.

## Three records, and what each answers

| Record | Where | Answers | Written by |
|---|---|---|---|
| **Run event timeline** | `.claude/feature-workspace/<feature>/run-events.jsonl` | What happened and when — stage transitions, gate halts and approvals, staleness, human corrections | The Go executor, timestamps from the clock |
| **Run trace** | `traces.jsonl` beside it, and OTLP/HTTP when an endpoint is set | How long each stage took and what it cost | The Go executor (roadmap L3.8) |
| **Run state** | `run-state.json` beside both | Where the run stands right now, plus per-stage token counts, cost, and corrections | The Go executor |

They overlap in subject and answer different questions. The timeline stays readable with no
collector configured and no exporter running, which is the property an audit record needs and a
trace does not.

## The event vocabulary is generated

One Go enum in `internal/orchestrator/vocabulary.go` is the source of truth. Both
`shared/schemas/telemetry/run-event.schema.json` and
`shared/schemas/telemetry/run-event-types.md` are generated from it by `go run ./cmd/gen-schemas`,
and a test fails when the committed copies drift. Never hand-edit either.

Adding an event kind without documenting it fails the build: a test parses the constants out of the
source and checks each one appears in the vocabulary. The previous arrangement — two hand-maintained
lists and a recorder told to enforce one of them — had nothing to check the prose against, which is
how 60% of the specified surface ended up outside it.

## Specified, but emitted by nothing

These types are described elsewhere in the framework and have no emitter. They are deliberately
**not** in the vocabulary: an enum where entries fire from nowhere is the same trap in a new place,
and a consumer would have to learn which members are real. Each is listed with the roadmap item
that would build its emitter.

| Type | Described in | Would be emitted by |
|---|---|---|
| `policy.evaluated`, `policy.conflict`, `policy.skipped` | `shared/orchestration/policy-evaluator.md`, `shared/policies/` | **L2.16** — replace the LLM policy evaluator with a real one |
| `contract.retry` | `shared/skills/deliver-feature/SKILL.md` | **L2.18** — run contract validation under the executor |
| `audit.fail`, `audit.retry`, `audit.halt` | `shared/orchestration/audit-composition-pattern.md` | **L3.12** — move counter agents out of the synchronous graph |
| `workflow.completed` | `shared/skills/orchestrate/SKILL.md` | **L2.16** |
| `workspace.migrated` | `shared/skills/deliver-feature/SKILL.md` | no item — a one-off migration step, not a recurring signal |
| `gate_decision` | formerly `event-schema.md` | superseded: the executor emits `artifact.corrected` (**L4.5**) for the corrective half, which is the part anything consumed |

## What the markdown pipeline records

Nothing of its own. The `deliver-feature` skill run by a host platform's LLM has no process that can
measure elapsed time or observe its own decisions, which is why its emission instructions were
prompt-discipline and why they were removed rather than rewritten. Its checkpoints go through
`loom state` when the binary is present — see `cmd/loom/README.md` — and those land in the same
`run-events.jsonl` the executor writes.

**One consequence worth stating plainly.** `shared/rules/approval-gates.md` says a policy-based gate
emits `policy.evaluated` for every decision and that there are no silent auto-approvals. Nothing
records those decisions today and nothing did before: the file was never written. The requirement
itself is unchanged — a gate still requires a human unless a policy matches — but the audit trail
that would prove it is **L2.16**'s to build.

## Privacy

Everything here is project-local: written inside the project's own `.claude/` directory, never
uploaded, never aggregated across projects by the framework. The one path that leaves the machine is
OTLP export, which is off unless `OTEL_EXPORTER_OTLP_ENDPOINT` is set. Rotation is the operator's
call — the framework does not auto-rotate, compress, or trim.
