# Epic 79 — Typed Graph State: Markdown Becomes a View, Not the Transport (L2.9)

Source: epic 78 completion, 2026-08-30. Operationalizes roadmap item **L2.9** (BUILD-ROADMAP.md,
KERNEL — State Management). With M0.4, L2.13, and L2.12 landed, the executor runs stages, halts at
gates, and owns integrity. What it still does not have is *data*: agents hand each other whole
markdown documents, and every downstream stage re-parses the full text. L2.9 introduces a typed
state graph. It blocks **L2.10** (deterministic projections), **L2.11** (semantic validation),
**L3.1** (planner/router), **L3.2** (capability registry), and **L4.4**.

**This epic is the first cut, not the whole migration.** It types one hop — `analyst → architect` —
end to end and leaves the other sixteen artifacts on markdown. That satisfies L2.9's done-when
("two consecutive stages exchange data with no markdown file on the path") and keeps a reviewable
blast radius. Typing the remaining stages is later epics, once this shape is proven.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item L2.9 (the spec), plus L2.10 and L2.11 for the boundary
   of what NOT to build. Respect Blocked-by: L2.9 needs only M0.4 (done).
2. `docs/adrs/ADR-006-loom-executes-pipelines.md` — Accepted; typed state is the executor's data
   model, and this epic is a direct consequence.
3. `internal/orchestrator/` — `provider.go` (`StageInput`/`StageOutput`, the interface this epic
   extends), `executor.go`, `state.go` (schema-versioned run state), `integrity.go` (digest
   verification — typed payloads inherit it for free), `timeline.go`.
4. `internal/provider/claude/` — how a stage prompt is built today from
   `shared/agents/<agent>.md`, and how the response is captured.
5. `shared/contracts/analysis-contract.md` and `architecture-contract.md` — the two contracts this
   epic turns into schemas. Note what they say about downstream agents grepping exact headings.
6. `shared/schemas/*.schema.json` + `shared/skills/validate-artifact/SKILL.md` — the existing
   JSON-Schema and structural-validation conventions to follow, not reinvent.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Scope of the slice | `analyst → architect` only. The analyst produces typed `AnalysisState`; the architect consumes a projection of it and produces typed `ArchitectureState`. Every other stage keeps writing markdown exactly as today | The done-when needs one hop. One hop also surfaces every modelling problem the full migration would hit, at a fraction of the review cost |
| Schema source of truth | Go structs in `internal/state/`, with JSON Schema **generated** from them into `shared/schemas/pipeline/<stage>.schema.json` via `invopop/jsonschema` (already a direct dependency). A generator script plus a CI check that regeneration produces no diff | No new toolchain, and the executor gets compile-time types for free. CUE was considered and rejected: a language nobody else in this repo writes, needed on every machine that generates code |
| How an agent produces state | The agent emits **JSON conforming to the stage's schema**; the executor validates it, persists it as the stage's artifact, and renders a markdown *view* for humans. Markdown is an output, never the transport | This is the roadmap's "markdown becomes a rendered view". A markdown parser was considered and rejected — it keeps prose on the data path and makes extraction a permanent, lossy tax |
| Agent prompt files are NOT edited | `shared/agents/analyst.md` and `architect.md` stay exactly as they are. The **claude provider** appends the stage's schema and a JSON-only output instruction at invocation time | Those files are shared with the markdown pipeline, which must keep working unchanged. Editing them would break every non-executor host to serve the executor |
| Which pipeline | Executor-first. `loom run` gets typed state; the markdown pipeline is untouched by this epic | The same split that worked for L2.13 and L2.12. Type safety buys enforcement only where a process, not a prompt, is doing the handoff |
| Where typed state lives | One file per stage: `.claude/feature-workspace/<feature>/state/<stage>.json`. That file **is** the stage's artifact, so L2.12's digest recording and staleness cascade apply to it unchanged | Reusing the artifact path means integrity, resume, and the stale cascade all work on day one with no new mechanism |
| Rendered markdown | Written under the **contract's** filename (`analysis.md`, `architecture-notes.md`), regenerated from state, and **not** digest-tracked. This is load-bearing, not cosmetic: `developer.md` and `qa-engineer.md` are still untyped and their prompts say to read `analysis.md`, so the typed hop stays invisible to them | It is a derived view. Verifying a derived file would make hand-editing the view corrupt the run, which is exactly backwards |
| Projections | A consuming stage declares the fields it reads; the executor passes only those, as JSON, in `StageInput`. This epic ships the mechanism for the one hop | Field-level access is the point of L2.9. Applying projections across every stage and retiring `summarize-artifact` is **L2.10** — do not touch that skill here |
| Validation depth | JSON Schema conformance only: on stage output (reject and fail the stage) and on load (a state file that no longer conforms is a load-time error). Cross-field and business rules are **L2.11** | The roadmap splits these deliberately. Semantic rules without typed state to hang them on is what L2.11 exists to fix |
| Non-JSON agent responses | Accept raw JSON, or exactly one ```` ```json ```` fenced block with only whitespace around it. Anything else fails the stage with the raw response logged | Models wrap output in fences as a formatting habit; failing a run over that would report a formatting reflex as a modelling error. Scanning for the first balanced object anywhere was rejected — it would happily accept a schema example the agent quoted back |
| Schema delivery | The generated schema is inlined into the stage prompt (~6.6 KB for analysis) | Self-contained: it works whether or not the framework is installed in the target project, the agent cannot read a stale copy, and it costs no tool call with a silent missing-file failure mode |
| Invalid agent output | The stage fails, loudly, with the validation error recorded in run state and the timeline. No repair prompt, no retry | Retries and self-repair are L3.x. A silent repair loop would hide exactly the modelling failures this epic needs to surface |
| Schema versioning | Each state struct carries `schemaVersion`; a mismatched version is a load-time error with a clear message. No migration code | Consistent with run-state's policy, and for the same reason: nothing is deployed anywhere yet |
| Architect routing input | `deliver-feature/SKILL.md:91` routes on an "Architectural Flags" heading that exists nowhere in the template or contract. Phase B adds `AnalysisState.RequiresArchitect()` — a pure predicate derived from context crossings, data-model changes, new dependencies, and performance NFRs with thresholds — plus `ArchitecturalFlags` as an explicit escape hatch (derived OR explicit). The prose is corrected to name fields that exist | A routing condition asserted by the model is a self-assessment; one derived from contract-validated facts is a testable function. Crossings and schema changes alone would miss new base classes and cross-cutting concerns, which is exactly what `architect.md` says it exists for — hence four signals and a human-settable override. *Who calls the predicate* stays L3.1 |
| `Provider` interface | `StageOutput` gains a `Payload []byte` field (the raw JSON the agent emitted). The interface stays one method | Epic 77 kept the interface untouched for gates because gates are not data. State is data — this is the field it was always going to need |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**55.9%** as of epic 78; raise the floor in `framework-ci.yml`
  if measured coverage rises).
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — The state package and generated schemas — UNBLOCKED

1. `internal/state/` — a new package, no dependency on `internal/orchestrator` (the orchestrator
   consumes state, not the reverse):
   - `analysis.go`: `AnalysisState` modelling `analysis-contract.md`'s required sections as typed
     fields — acceptance criteria and NFRs as slices, bounded context and architectural flags as
     scalars, the four task lists as separate slices, data-model and API changes as structured
     entries. Model what the contract actually promises; do not invent fields nobody consumes.
   - `architecture.go`: `ArchitectureState` for `architecture-contract.md`, same treatment.
   - Every struct carries `SchemaVersion int` and a `Validate()` returning a typed error naming the
     offending field.
2. Schema generation: `cmd/loom/…` or `scripts/` — a generator producing
   `shared/schemas/pipeline/analysis.schema.json` and `architecture.schema.json` from the structs
   via `invopop/jsonschema`. Wire a CI check (extend `framework-ci.yml`) that regenerating produces
   no diff — a generated file that drifts from its source is worse than no generated file.
3. Tests: round-trip (struct → JSON → struct) preserves every field; a payload missing a required
   field is rejected with the field named; an unknown schema version is a load-time error;
   generated schemas match the committed ones.

**Done when**: both stage schemas exist as Go types and generated JSON Schema, with drift caught by
CI. **Commit** (`feat(state): typed analysis and architecture state with generated schemas`),
report, PAUSE.

## Phase B — Executor plumbing: typed handoff — BLOCKED BY Phase A

1. `internal/orchestrator/`:
   - `StageOutput.Payload []byte`; the executor validates the payload against the stage's schema,
     writes `state/<stage>.json`, and records that path as the stage's artifact (so L2.12's digest
     and staleness machinery covers it with no new code).
   - A projection: `StageInput` gains the upstream fields the stage declares. For this slice, the
     architect receives a projection of `AnalysisState` — the fields `architecture-contract.md`
     actually needs — not the whole document.
   - Stages outside the slice are unaffected: no payload, markdown artifact, same behavior.
2. Mock provider: `Script` gains a `Payload` so tests can script typed output; `loom run --provider mock` emits valid typed payloads for typed stages.
3. `AnalysisState.RequiresArchitect()` per the design table, with table-driven tests covering the
   false-negative cases (a feature with no crossing and no schema change that still needs an
   architect must be reachable through the explicit flag). Correct `deliver-feature/SKILL.md:91` to
   route on fields the template actually produces. Do not wire the predicate into execution — that
   is L3.1.
4. Tests: the architect stage receives analyst fields with **no markdown file on the path** (the
   L2.9 done-when — name the test accordingly); an invalid payload fails the stage with the
   validation error in run state and one `stage.failed` timeline event; editing a `state/*.json`
   out of band demotes the stage exactly as L2.12 does for any artifact; untyped stages still work.

**Done when**: two consecutive stages exchange typed data with no markdown on the path, and a
schema violation fails the stage. **Commit** (`feat(orchestrator): typed stage payloads and
field-level projections`), report, PAUSE.

## Phase C — The claude provider and the rendered view — BLOCKED BY Phase B

1. `internal/provider/claude/`: for a stage with a schema, append the schema and a JSON-only output
   instruction to the prompt built from `shared/agents/<agent>.md` (which is **not** edited), and
   return the response as `Payload`. A response that is not JSON fails the stage with a message
   showing what came back.
2. Markdown rendering: render `<stage>.md` from typed state as a human-readable view — headings
   matching the contract, so a human reading the workspace sees what they see today. Not
   digest-tracked; regenerated whenever state is written.
3. Tests: prompt construction includes the schema and the instruction; a fenced or prose-wrapped
   JSON response is handled or fails clearly (decide which, document it, test it); rendering is
   deterministic and round-trips through the contract's headings.

**Done when**: a real `loom run` slice produces typed state and a readable markdown view, and a
non-conforming agent response fails loudly. **Commit** (`feat(provider): schema-directed stage
output and markdown rendering`), report, PAUSE.

## Phase D — Contracts and docs — BLOCKED BY Phase C

1. `shared/contracts/analysis-contract.md` and `architecture-contract.md`: add a schema reference
   and state the split honestly — under `loom run` these artifacts are typed state validated
   against the schema; for the markdown pipeline the heading rules below remain authoritative. Do
   not delete the heading rules.
2. `cmd/loom/README.md`: document typed state, where it lives, and that markdown in the workspace
   is now a rendered view for those stages.
3. `docs/ARCHITECTURE.md` + `README.md`: L2.9 status markers — first cut shipped for one hop,
   remaining stages still markdown-transported.
4. `shared/orchestration/pipeline-schema.md`: note which parts are now generated rather than
   hand-maintained prose, without rewriting the workflow schema itself (that is L3.1's).
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when** (L2.9, first cut): the roadmap done-when holds for the typed hop and the docs are
explicit about how much of the pipeline it covers. **Commit** (`docs(state): typed graph state —
first cut covers analyst to architect`), report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Typing the other sixteen artifacts (later epics, once this shape is proven)
- Retiring `summarize-artifact` or applying projections pipeline-wide (**L2.10**)
- Cross-field, cardinality, or referential business rules; an LLM critic gate (**L2.11**)
- Conditional routing or a planner over typed state (**L3.1**), capability registry (**L3.2**)
- Editing `shared/agents/*.md` prompt files, or changing the markdown pipeline's artifacts
- Retry, repair, or reformat loops for non-conforming agent output (L3.x)
- Migrating existing archived markdown artifacts into typed state

## Report format (end of every phase)

```
## Epic 79 Phase <X> Report
- Roadmap item: L2.9 — Typed graph state (first cut)
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
