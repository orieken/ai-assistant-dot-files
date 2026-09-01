# Epic 84 — Emit OpenTelemetry with GenAI Semantic Conventions (L3.8)

Source: session discussion, 2026-09-01, following epic 83. Implements roadmap item **L3.8**, the
last of the three items the roadmap's own summary names as "if only three things get built" — the
other two (L2.9 typed state; M0.4 + L2.13 executor owning gates) have shipped. L3.8 blocks five
items: **L3.5** (episodic memory), **L3.9** (telemetry schema), **L4.3** (budget governor),
**L4.5** (correction signal), **L4.6** (lessons reach prompts). Roughly all of Milestone 4 is
dammed behind it.

The question "what did this pipeline cost" is currently unanswerable. This epic answers it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L3.8**, and **L3.9** immediately after it. L3.9 owns the
   event-type enum and the `event-recorder` contradiction. Read it to know what this epic must
   *not* absorb.
2. `internal/orchestrator/timeline.go` — the existing append-only `run-events.jsonl`. Its header
   comment already says "deliberately NOT OpenTelemetry (roadmap L3.8)". That comment stops being
   an excuse and becomes a boundary this epic has to state properly.
3. `internal/orchestrator/executor.go` — the run loop, and the only place that knows when a stage
   starts and ends by measurement rather than by recall.
4. `internal/provider/claude/claude.go` — `runSubprocess` spawns `claude -p` with stdout captured
   straight to the artifact file. Read this closely; Phase B changes it, and it is the only place
   token counts can come from.
5. `shared/telemetry/event-schema.md`, `event-recorder.md`, `README.md` — what the framework
   currently *claims* about telemetry, in the present tense, with no emitter.
6. `shared/rules/architecture-guardrails.md` #8 — "No OTel instrumentation logic is allowed inside
   domain entities... Traces and spans must only be emitted from the adapter layer or interceptor
   layer." This epic is the framework turning that rule on itself; `internal/state/` must end this
   epic with zero telemetry imports.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| OTel is a real dependency | Add `go.opentelemetry.io/otel`, `otel/sdk`, and the OTLP trace exporter to `go.mod` | This is the largest third-party addition the module has taken. Hand-rolling spans to avoid it would produce a worse thing that no collector can read, which defeats the item. Accept the dependency deliberately |
| Network export is off by default; local export is on | The OTLP **endpoint** is opt-in via the standard `OTEL_EXPORTER_OTLP_ENDPOINT` env var. The **file** exporter is on by default, writing `traces.jsonl` into the feature workspace beside `run-state.json`; `--otel-file <path>` overrides the location, and an explicit off switch disables it | A CLI that phones somewhere by default is a bad citizen — but that concern is about network egress, and a file beside the run's own state leaves nothing. Making it opt-in would mean the cost of a run is unanswerable unless someone predicted beforehand that they'd want to know, including for runs already finished. L4.3 and L3.13 both want to read historical runs |
| File format is OTLP JSON | The file exporter writes OTLP's JSON encoding, not the SDK's `stdouttrace` format | A saved file can be replayed into a collector or read by anything that already understands OTLP. `stdouttrace` is explicitly not a stable wire format, so the downstream items would be parsing something OTel does not promise to keep |
| Spans do not replace the timeline | `run-events.jsonl` stays exactly as it is | They answer different questions. The timeline is the audit record of gates, digests, and staleness — point events, read by `loom state timeline`, and durable with no collector. Spans are the timing and cost record. Merging them would make the audit log depend on an exporter being configured |
| Span shape | One root span per run; one child per stage; one grandchild per provider invocation. Loop iterations are sibling stage spans distinguished by an iteration attribute, not nested | Nesting iterations would make a three-round review loop look like three levels of depth. They are retries of one stage, and a trace viewer should show them side by side |
| GenAI semconv | `gen_ai.operation.name`, `gen_ai.request.model`, `gen_ai.usage.input_tokens`, `gen_ai.usage.output_tokens` on the provider span. Loom-specific facts (`loom.stage.id`, `loom.gate`, `loom.route.skipped`) use a `loom.` prefix | Follow the convention where one exists so a generic GenAI dashboard works, and do not squat on the `gen_ai.` namespace for things it does not define |
| Where token counts come from | `claude -p --output-format json`, whose envelope carries `usage` and `total_cost_usd`. The provider parses the envelope and writes `.result` to the artifact | This is the only honest source. It also means **no pricing table in this repo** — the CLI reports the cost it was actually charged, and a hardcoded price list would be wrong within a quarter |
| Provider output parsing failure | If the envelope does not parse, the stage FAILS with the raw output preserved. No fallback to treating stdout as the artifact | A silent fallback would make a malformed run look like a successful one with zero tokens, which is exactly the class of unmeasured-but-reported number this item exists to remove |
| Instrumentation lives at the adapter layer | `internal/telemetry/` is imported by `internal/orchestrator/` and `internal/provider/`, never by `internal/state/` | Guardrail #8. A fitness function enforces it — an import-graph test, not a code review note |
| Cost is reported, not computed | The run's total is the sum of what the CLI reported per stage | Same reason as the pricing table |
| Event-type unification is NOT here | `event-schema.md`'s six-vs-fifteen contradiction, the `event`/`event_type` key mismatch, and deleting `event-recorder.md` are **L3.9** | L3.9 is a separate numbered item that is blocked by this one. Absorbing it silently is the failure mode epic 83 had to write a roadmap correction for. Phase D says what L3.8 leaves for it |
| The markdown pipeline gets honesty, not spans | Docs state plainly that traces come from `loom run` only | Same posture as every gate and state item before it. A host LLM cannot measure elapsed time; that is the Problem statement, not something prose can fix |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**62.2%** as of epic 83). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- No network call in any test. The OTLP exporter is exercised against an in-process collector stub
  or the file exporter, never a real endpoint.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — `internal/telemetry/` and spans around the run loop — UNBLOCKED

1. New `internal/telemetry/`: tracer provider construction, exporter selection (OTLP from env, file
   from flag, no-op when neither is set), resource attributes identifying loom and its version
   (`cmd/loom/buildinfo.go` already has the version), and clean shutdown that flushes on exit
   including the gate-halt exit-3 path.
2. Instrument the executor: a root run span and a per-stage child. Attributes for stage ID, agent,
   sequence, and terminal status. A stage that is skipped, stale, or waiting on a gate must be
   visible as such — a routed-around stage is a real outcome, not an absence.
3. Wire `--otel-file` into `loom run` (overriding the default workspace location) plus a way to turn file export off; document that `OTEL_EXPORTER_OTLP_ENDPOINT` is honoured for network export.
4. The import-graph fitness function: a test asserting `internal/state/` imports no telemetry
   package, per guardrail #8.
5. Tests: a completed mock run produces one root span with the expected children; a run halted at a
   gate still flushes what it recorded; file export off with no endpoint set emits nothing and costs nothing; the emitted file parses as OTLP JSON.

**Done when**: a `--provider mock` run produces a single well-formed trace whose span tree matches
the stages that actually executed, and `internal/state/` provably has no telemetry import.
**Commit** (`feat(telemetry): trace the run loop with OpenTelemetry`), report, PAUSE.

## Phase B — Real token counts and cost from the claude provider — BLOCKED BY Phase A

1. Switch `runSubprocess` to `claude -p --output-format json`, parse the envelope, and write
   `.result` to the artifact path. Preserve current behavior for everything else — the typed-stage
   JSON extraction in `typed_stage.go` still runs on the result payload, not on the envelope.
2. A provider span per invocation carrying `gen_ai.*` usage attributes, the model the CLI reports,
   and the cost figure it reports. Attach the per-stage totals to the stage span so a trace viewer
   shows cost without needing to sum leaves.
3. Malformed envelope = stage failure with the raw output preserved for diagnosis. Test it.
4. Record the run's total cost where a human can see it without a collector: the run summary line
   and `loom state show`.
5. Tests: envelope parsing including the fields we do not use; usage attributes reach the span; a
   malformed envelope fails the stage rather than producing a zero-token success.

**Done when**: a completed run produces a single trace with per-stage token counts and a total cost
figure — L3.8's stated done-when, in full. **Commit** (`feat(telemetry): real token and cost
accounting from the claude provider`), report, PAUSE.

## Phase C — The MCP server surface — BLOCKED BY Phase A

1. Spans per tool call in `shared/mcp/`, with arguments and results as attributes under a **size
   cap** and with **secret redaction** — the roadmap names both explicitly and both are testable.
2. Trace context propagation: a tool call made during a `loom run` stage must land under that
   stage's span rather than starting an orphan trace.
3. Replace or wrap `shared/mcp/internal/logging/logger.go` so its output correlates with traces
   (trace and span IDs on log lines) rather than existing beside them.
4. Tests: an oversized argument is truncated with a marker rather than dropped silently; a value
   matching the secret patterns never reaches an attribute; a propagated context produces one trace
   rather than two.

**Done when**: a tool call made inside a traced run appears as a child span of the stage that made
it, with no secret and no unbounded payload in its attributes. **Commit** (`feat(mcp): trace tool
calls and correlate logs`), report, PAUSE.

## Phase D — Docs, and the L3.9 boundary — BLOCKED BY Phases B+C

1. `shared/telemetry/README.md` and `event-schema.md`: stop describing telemetry in the present
   tense where it is not emitted, and state which pipeline produces traces. Do **not** unify the
   event-type enum or delete `event-recorder.md` — say plainly that both are L3.9, and update
   L3.9's roadmap entry with what it inherits from here.
2. `cmd/loom/README.md`: how to turn telemetry on, what a trace contains, and the offline file path.
3. `README.md`, `docs/ARCHITECTURE.md`, BUILD-ROADMAP status for L3.8, `shared/DOMAIN_DICTIONARY.md`
   if a new term appears.
4. Confirm the five items L3.8 blocks are genuinely unblocked, and say so in each entry — an item
   that stays marked blocked after its blocker ships is a dropped item.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the docs describe telemetry that exists, and L3.9's inheritance is written down.
**Commit** (`docs(telemetry): describe the traces loom actually emits`), report, PAUSE — epic
complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- The event-type enum, the `event`/`event_type` key mismatch, and deleting `event-recorder.md`
  (**L3.9**)
- Metrics or logs signals — traces only. A budget governor needs metrics; that is **L4.3**
- A pricing table, or any cost figure loom computes rather than reports
- Instrumenting the markdown pipeline, or asking a host LLM to record a duration
- Sampling strategy beyond always-on; a local CLI run is not a traffic volume problem

## Report format (end of every phase)

```
## Epic 84 Phase <X> Report
- Roadmap item: L3.8 — OpenTelemetry with GenAI semantic conventions
- Blockers verified: <list, with evidence (commit SHAs / files)>
- Commits: <sha> <subject>
- Build/lint/test: go build PASS|FAIL · golangci-lint PASS|FAIL · go test PASS|FAIL (counts) · health-check PASS|FAIL|n/a
- Done-when criterion: <restate it> — MET | NOT MET (why)
- Escalations / open questions: <list or "none">
- Next phase blocked by: <what must land first>
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
