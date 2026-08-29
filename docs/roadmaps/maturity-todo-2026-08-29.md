> **SUPERSEDED 2026-08-29** — merged into [`BUILD-ROADMAP.md`](BUILD-ROADMAP.md), the single
> authoritative build plan. Kept for provenance; see BUILD-ROADMAP Appendix A for where each item
> landed and Appendix B for corrections. Do not plan work from this file.

# Agentic Maturity TODO — L2 → L4

**Framework version**: v3.3.14 @ `59efe14`
**Date**: 2026-08-29
**Scope**: architectural work required to carry `loom` from Level 2 through Level 4 of the Agentic
AI Maturity Model.

**Baseline**: 566 markdown files (~52k lines), 113 Go files (~8.9k lines). The Go is a Cobra
installer (`cmd/loom/`) plus an MCP tool server (`shared/mcp/`). Orchestration, telemetry, policy,
hooks, and three of four retriever backends are markdown specifications with no executor. Every
phase below assumes closing that gap is the work.

**Companion documents**:
- [`architectural-audit-2026-08-29.md`](architectural-audit-2026-08-29.md) — the findings (H1–H11)
  this TODO operationalizes. Read that first for evidence; this file is the work breakdown.
- [`agy.md`](agy.md) — earlier maturity roadmap. Overlaps on dynamic routing, typed state, and
  semantic memory; superseded where the two disagree.

Every item cites files verified as read on 2026-08-29.

---

## PHASE 1 — Shoring Up Level 2 (Coordinated Multi-Agent Systems)

### Tool Execution & Validation

- [ ] **Decouple the tool port from the MCP transport SDK**
  1. **Problem**: `domain.Tool` — the file commented as "the framework's first-class abstraction for
     every capability" — imports `github.com/mark3labs/mcp-go/mcp` and `invopop/jsonschema`, and
     types its own signatures in them (`mcp.ToolInputSchema`, `mcp.CallToolRequest`,
     `*mcp.CallToolResult`). The domain layer is welded to one wire protocol, violating
     `architecture-guardrails.md` #1. Any second transport (HTTP, gRPC, in-process) requires
     rewriting every tool.
  2. **Architectural Fix**: Hexagonal port/adapter. Define transport-free
     `ToolRequest{Name string; Args map[string]any}` / `ToolResult{Content []ContentBlock; IsError bool; Err error}`
     in `domain`, and move all `mcp.*` marshalling into a `server/mcp_adapter.go` that translates at
     the edge. `domain` must import zero third-party packages.
  3. **Target Files**: `shared/mcp/internal/domain/tool.go`,
     `shared/mcp/internal/server/registration.go`, all 6 files in `shared/mcp/internal/tools/`.

- [ ] **Enforce input schemas server-side instead of trusting the model**
  1. **Problem**: `InputSchema()` exists only to describe arguments *to the LLM*. Actual enforcement
     is unchecked type assertion: `parseComplexityArgs` does `args["projectPath"].(string)` and
     silently yields `""` on any non-string, then the tool returns a generic error. A hallucinated
     argument shape produces a soft failure the model then re-attempts blindly. No required-field
     check, no enum check, no bounds.
  2. **Architectural Fix**: Validate every `Args` map against the tool's declared JSON Schema at the
     handler boundary *before* dispatch, returning a structured `ValidationError` naming the
     offending field and the expected type — a machine-actionable repair signal, not prose.
     `github.com/santhosh-tekuri/jsonschema/v6` is **already in the dependency graph** as an indirect
     dep of `mcp-go`; promote it to direct.
  3. **Target Files**: `shared/mcp/internal/server/registration.go` (wrap the `AddTool` closure),
     `shared/mcp/internal/tools/schemas.go`, `shared/mcp/go.mod`.

- [ ] **Propagate `context.Context` and set per-tool deadlines**
  1. **Problem**: All six tools sign `Execute(_ context.Context, ...)` —
     `analyze_complexity_tool.go:54`, `check_accessibility_tool.go:52`,
     `check_ubiquitous_language_tool.go:50`, `search_docs_tool.go:64`, `search_ki_tool.go:56`,
     `verify_dependencies_tool.go:40`. Cancellation is discarded 100% of the time. A client
     disconnect, timeout, or user abort cannot stop an in-flight `filepath.Walk`.
     `go-conventions.md` mandates explicit timeouts; the tool layer has none.
  2. **Architectural Fix**: Thread `ctx` through `Execute` → analyzer → walk, checking `ctx.Err()`
     inside every `filepath.WalkDir` callback. Add a middleware in the registration closure applying
     `context.WithTimeout` from a per-tool budget declared in the registry entry.
  3. **Target Files**: all 6 `shared/mcp/internal/tools/*_tool.go`,
     `shared/mcp/internal/analyzers/*.go`, `shared/mcp/internal/server/registration.go`.

- [ ] **Confine filesystem access to an explicit root**
  1. **Problem**: `analyze_complexity` accepts an arbitrary `projectPath` from model-controlled
     arguments with no validation, no root confinement, and no symlink handling, then walks it
     unbounded and uncancellably. `projectPath: "/"` walks the disk. Same pattern in
     `check_accessibility`, `check_ubiquitous_language`, `verify_dependencies`, and `search_docs`
     (`docsPath`).
  2. **Architectural Fix**: Adopt `os.Root` (available on `go 1.26.5`) for traversal-safe rooted FS
     access, with the root supplied by server config rather than tool arguments. Reject any path that
     escapes it. Add a file-count / byte ceiling to abort runaway walks.
  3. **Target Files**: `shared/mcp/internal/analyzers/walkutil.go`,
     `shared/mcp/internal/analyzers/*_analyzer.go`, `shared/mcp/internal/server/tool_provider.go`.

- [ ] **Replace the hardcoded tool slice with a registry**
  1. **Problem**: `buildFrameworkTools()` is a slice literal returning six constructor calls. Adding
     a tool means editing and recompiling the handler. There is no discovery, no per-tool metadata
     (timeout, retry policy, permission scope, version), no enable/disable, and no way for a
     downstream project to contribute a tool without forking.
  2. **Architectural Fix**: Registry pattern — `map[string]ToolRegistration` where
     `ToolRegistration` carries the `Tool`, its timeout, retry class, and required permission scope.
     Populate via an `init()`-style `Register(name, reg)` per tool file, or an explicit
     registry-builder consuming config. `register.FrameworkTools` becomes a registry merge instead of
     a wholesale re-registration.
  3. **Target Files**: `shared/mcp/internal/server/tool_provider.go`,
     `shared/mcp/internal/server/handler.go`, `shared/mcp/register/register.go`.

- [ ] **Introduce a typed failure taxonomy — retryable vs. terminal**
  1. **Problem**: Every failure path collapses to `mcp.NewToolResultError(fmt.Sprintf(...))` with
     `err == nil` — a *successful* tool call carrying a prose error string. The caller cannot
     distinguish "bad argument, fix and retry" from "corpus missing, stop" from "transient I/O, back
     off." `search_docs` goes further and returns `Success: true, TotalHits: 0` with the failure
     reason smuggled into the `Query` field (`emptyResult`), indistinguishable from a genuine
     zero-result search.
  2. **Architectural Fix**:
     `ToolError{Kind: Validation|NotFound|Transient|Internal|Permission, Field, Message, Retryable bool}`
     serialized into a stable `error` envelope in the result. Never encode failure state into a
     success payload. Orchestrator retry policy keys off `Kind`.
  3. **Target Files**: `shared/mcp/internal/tools/responses.go`,
     `shared/mcp/internal/tools/search_docs_tool.go:101` (`emptyResult`), all `*_tool.go` error
     branches.

- [ ] **Add the resilience primitives the guardrails already mandate**
  1. **Problem**: `architecture-guardrails.md` #5 forbids hand-rolled retry loops and requires
     `CircuitBreaker` or `ExponentialBackoffStrategy`. Neither exists anywhere in the Go tree. The
     only retry logic in the framework is prose in `deliver-feature/SKILL.md` telling an LLM to count
     to three.
  2. **Architectural Fix**: `sony/gobreaker` per tool (and per downstream dependency) plus
     `cenkalti/backoff/v4` for `Transient`-classed failures, wired as registry middleware so no tool
     implements its own retry. Emit breaker state transitions as telemetry.
  3. **Target Files**: new `shared/mcp/internal/server/middleware.go`,
     `shared/mcp/internal/server/registration.go`, `shared/mcp/go.mod`.

- [ ] **Fix the per-query full-corpus re-index before any scale test** *(audit H2)*
  1. **Problem**: `search_docs_tool.go:82` calls `EnsureIndex` inside `Execute`.
     `bm25_retriever.go:70` then walks the whole docs tree, `os.ReadFile`s every `.md`, and runs
     **one sqlite transaction per file** — no mtime check, no content hash, no dirty tracking. Every
     search is O(corpus) disk I/O plus O(n) transactions. The drifted copy documents the second half:
     `shared/mcp-patterns/go/tools/bm25_retriever.go:40` states "EnsureIndex is not safe to call
     concurrently with itself" — and MCP servers field concurrent calls. Deleted docs are never
     evicted; `DELETE` only fires for paths being re-inserted.
  2. **Architectural Fix**: Move indexing out of the query path entirely. Incremental index keyed on
     `(path, mtime, size)` or content hash, in one batched transaction; a separate explicit
     `reindex_docs` tool plus optional `fsnotify` watcher; `sync.RWMutex` around writer access; a
     reconciliation sweep deleting index rows whose paths no longer exist.
  3. **Target Files**: `shared/mcp/internal/tools/bm25_retriever.go`,
     `shared/mcp/internal/tools/search_docs_tool.go`.

### State Management

- [ ] **Replace markdown-file state passing with a typed graph state**
  1. **Problem**: There is no state object. Agents hand each other whole markdown documents on disk
     (`analysis.md`, `architecture-notes.md`, `implementation-notes.md`, … 15 artifacts). The only
     real delivery in the repo has a 15 KB `analysis.md`. Every downstream agent re-parses the full
     text; there is no field-level access, no size ceiling, no provenance, and no way to pass a value
     without passing a document.
  2. **Architectural Fix**: A versioned `PipelineState` struct with per-stage typed sub-schemas (Go
     structs + generated JSON Schema, or CUE as single source of truth). Markdown becomes a
     *rendered view* of state, not the transport. Agents receive a narrowly-scoped projection of the
     fields their contract declares, not the whole graph.
  3. **Target Files**: `shared/orchestration/pipeline-schema.md` → becomes generated;
     `shared/contracts/*.md` (all 18) → become schemas;
     `shared/skills/deliver-feature/SKILL.md`; new `internal/state/`.

- [ ] **Stop using an LLM as the context-compaction mechanism**
  1. **Problem**: Context decay is handled by `summarize-artifact --persist`, invoked at
     `deliver-feature/SKILL.md` step 37a — another LLM call producing a lossy ~200-word surrogate
     that is then indexed as a retrieval target. Compaction is nondeterministic, unverifiable, costs
     a model call, and silently drops whatever the summarizer deemed unimportant.
  2. **Architectural Fix**: Deterministic projections. Each contract declares which fields downstream
     stages may read; the executor computes the projection by field selection, not summarization.
     Reserve LLM summarization for human-facing prose only, never for machine handoff.
  3. **Target Files**: `shared/skills/summarize-artifact/SKILL.md`,
     `shared/skills/deliver-feature/SKILL.md:139` (step 37a), `shared/contracts/*.md`.

- [ ] **Make `validate-artifact` verify semantics, not heading presence**
  1. **Problem**: Contract validation checks that required `##` sections exist.
     `agent-scorecard/SKILL.md:29` confirms the depth: the analyst's "completeness score" is
     "fraction of required sections present **and** containing real content (not leftover `[...]`
     template placeholders)." That is a template-placeholder grep. A structurally perfect,
     semantically empty artifact passes every gate in the pipeline.
  2. **Architectural Fix**: Once state is typed, validation becomes JSON Schema conformance plus
     declarative business rules (required cardinality, cross-field consistency, referential integrity
     against prior stages). Keep an LLM critic as an *additional* qualitative gate, never as the
     structural one.
  3. **Target Files**: `shared/skills/validate-artifact/SKILL.md`, `shared/contracts/*.md`,
     `shared/schemas/`.

- [ ] **Move `pipeline-state.json` ownership into the executor**
  1. **Problem**: The state file — including SHA-256 checksums used for tamper detection and
     gate-edit detection — is written *by the LLM following prose instructions*
     (`deliver-feature/SKILL.md`, "Checkpointing & Pipeline State"). A model computing and recording
     its own integrity hashes is not integrity. No `pipeline-state.json` exists anywhere in the repo,
     so this has never executed.
  2. **Architectural Fix**: The executor owns the file: atomic write (temp + `os.Rename`), a schema
     version field, real `sha256` computed in Go, and verification on resume in code rather than by
     instruction. Agents never write it.
  3. **Target Files**: `shared/skills/deliver-feature/SKILL.md` (state section),
     `shared/skills/resume-pipeline/SKILL.md`, new `internal/orchestrator/checkpoint.go`.

### Human-in-the-Loop

- [ ] **Implement gates as process interrupts, not prose**
  1. **Problem**: All eight gates in `approval-gates.md` are natural-language instructions ("user
     must say 'ship'"). The enforcement mechanism for an irreversible action — DB contract-phase
     `DROP`, deploy, external API mutation — is the model's willingness to comply with a paragraph.
     There is no code path that can physically prevent the action.
  2. **Architectural Fix**: The executor halts the process at gate boundaries, persists state, and
     yields to a real approval channel (CLI prompt, webhook, or queue). High-risk tool classes are
     declared in the tool registry and are *unreachable* without a signed approval token in the
     request. Enforcement lives below the model, not in it.
  3. **Target Files**: `shared/rules/approval-gates.md`,
     `shared/skills/deliver-feature/SKILL.md:130-133`, new `internal/orchestrator/gate.go`, tool
     registry entries.

- [ ] **Enforce the "any edit resets the gate" rule in code**
  1. **Problem**: Every gate in `approval-gates.md` declares "Reset condition: any edit to the
     pending artifact resets the gate." Nothing enforces it. The related `gate_decision` telemetry
     spec describes checksum-diffing to detect `edited_then_approved` — but the checksum is computed
     by the model, the event type is emitted by nobody, and no `events.jsonl` exists in the
     repository.
  2. **Architectural Fix**: Approval binds to an artifact digest. The executor computes the digest at
     halt, issues a scoped approval token over it, and re-verifies at execution. Digest mismatch
     invalidates the token — a code-level check, not a remembered rule.
  3. **Target Files**: `shared/rules/approval-gates.md`, `shared/telemetry/event-schema.md:112`
     (`gate_decision`), new `internal/orchestrator/gate.go`.

- [ ] **Make resume a real capability instead of a prompt that re-reads a file**
  1. **Problem**: `resume-pipeline/SKILL.md` implements three modes (resume, `--from-phase N`,
     per-agent rollback) entirely as instructions to an LLM to read state, recompute checksums, mark
     entries `"stale": true`, and jump to a numbered step in *another skill's* prose. There is no
     process to resume — the "pause" was only ever the model stopping. Steps are addressed by their
     position in a hand-numbered 43-step list in `deliver-feature/SKILL.md`, so renumbering silently
     breaks every resume path.
  2. **Architectural Fix**: Durable executor with content-addressed stage IDs (never ordinals), a
     real checkpoint store, and resume as a first-class operation that replays from persisted state.
     Rollback becomes state-graph surgery in code, not markdown-file restoration by instruction.
  3. **Target Files**: `shared/skills/resume-pipeline/SKILL.md`,
     `shared/skills/deliver-feature/SKILL.md`, `shared/orchestration/interface.md`.

- [ ] **Replace the LLM policy evaluator with a real one** *(audit H7)*
  1. **Problem**: `policy-evaluator.md` specifies a condition language
     (`filePaths.noneMatch: "**/security/**"`, `diffLines.lessThan`, `not:`, `any:`) whose evaluator
     is a prompt. An authorization decision — *may this pipeline commit without a human* — is
     resolved by natural-language reasoning over YAML, with the kill-switch, the always-human list,
     and the conflict-resolution table all in the same prose the model may misread.
     `policy-schema.md`'s own worked example has a duplicate `diffType:` key in one map, which is
     invalid YAML — no parser has ever seen it.
  2. **Architectural Fix**: `google/cel-go` or OPA/Rego. The condition schema maps near-directly onto
     CEL. The always-human gate list becomes a compiled constant. Add a policy unit-test harness and
     a `--dry-run-policies` path that actually executes.
  3. **Target Files**: `shared/orchestration/policy-evaluator.md`, `shared/policies/policy-schema.md`,
     `shared/policies/examples/*.policy.yaml`, new `internal/policy/`.

- [ ] **Ship MCP over an authenticated network transport**
  1. **Problem**: `cmd/mcp-server/main.go` calls `server.ServeStdio(s)` — stdio only. No streamable
     HTTP, no SSE, no authentication, no authorization, no tenancy, no per-caller rate limiting. The
     server is single-user, local-only, and has no notion of *who* is calling.
  2. **Architectural Fix**: Add streamable HTTP transport alongside stdio, with OAuth2/OIDC bearer
     validation, per-principal tool scoping enforced against the registry's permission field, and
     per-principal rate limits. Keep stdio for local dev.
  3. **Target Files**: `shared/mcp/cmd/mcp-server/main.go`, `shared/mcp/internal/server/handler.go`,
     new `shared/mcp/internal/server/auth.go`.

---

## PHASE 2 — Achieving Level 3 (Autonomous Orchestration Layer)

### Dynamic Routing

- [ ] **Build a Planner/Router node; delete the hardcoded step sequence**
  1. **Problem**: There is no routing anywhere in the codebase. `deliver-feature/SKILL.md` is a
     hand-numbered 43-step list. Branching is static prose conditionals evaluated by reading markdown
     (`if analysis.md has Data Model Changes != "None"`, steps 12/14/16). `/orchestrate` reads an
     ordered `stages:` array. Nothing in the system lets an agent decide who runs next.
  2. **Architectural Fix**: Split the executor into a graph runtime plus a `Planner` node that emits
     the next node ID (or a sub-graph) from typed state. Conditionals become CEL predicates over
     state fields, not prose over documents. Keep the current linear pipeline as one registered
     *default plan* so existing behavior is a special case of the router, not a parallel code path.
  3. **Target Files**: `shared/skills/deliver-feature/SKILL.md`, `shared/skills/orchestrate/SKILL.md`,
     `shared/orchestration/pipeline-schema.md`, new `internal/orchestrator/planner.go`.

- [ ] **Publish a machine-readable agent capability registry**
  1. **Problem**: A router needs something to route *to*. There are 40 agents as markdown prose in
     `shared/agents/`. Their frontmatter carries `name`, `description`, `tools`, `model_tier`,
     `version` — no declared inputs, no declared outputs, no preconditions, no postconditions, no
     cost/latency class. `agent-frontmatter.schema.json` sets `additionalProperties: false`, so none
     of that can be added without a schema change.
  2. **Architectural Fix**: Extend the frontmatter contract with `consumes: [state fields]`,
     `produces: [state fields]`, `preconditions: [CEL]`, `cost_class`. Generate `agent-registry.json`
     at build time and have the planner select over it. This also makes the capability catalog
     auditable and diffable.
  3. **Target Files**: `shared/schemas/agent-frontmatter.schema.json`,
     `shared/contracts/agent-frontmatter-contract.md`, all 40 `shared/agents/*.md`,
     `scripts/generate-configs.sh`.

- [ ] **Implement real parallelism — `sequential-simulation` is a fiction**
  1. **Problem**: `shared/orchestration/interface.md` documents
     `parallelStrategy: fork | sequential-simulation (default)` and defines the default as "execute
     each in sequence but without reading each other's in-progress output." The headline
     parallel-branch feature ships disabled by definition. `orchestrate/SKILL.md` step 6 asks the
     model to "collect adjacent `parallel: true` stages into a group" by reading YAML it cannot
     execute.
  2. **Architectural Fix**: Real fan-out/join in the executor — `errgroup` with bounded concurrency,
     per-branch isolated state scopes, deterministic merge at the join with declared conflict
     resolution. Delete `sequential-simulation` outright; a fake concurrency mode is worse than none.
  3. **Target Files**: `shared/orchestration/interface.md`, `shared/skills/orchestrate/SKILL.md`,
     `shared/workflows/*.md`, new `internal/orchestrator/parallel.go`.

### Long-Term Memory

- [ ] **Implement the retriever backends that are currently markdown**
  1. **Problem**: `shared/rag/retriever.interface.md` is a well-specified contract — references not
     content, bounded top-K, corpus isolation, no side effects. Then three of its four adapters
     (`llm-as-retriever.md`, `vector.md`, `source-retrieval.deferred.md`) are prose files. Only BM25
     exists in code. There is no vector store, no embedding pipeline, no graph.
     `shared/knowledge/ollama-local-embeddings.md` is documentation, not an implementation.
  2. **Architectural Fix**: Implement `Retriever` as a real Go interface with the BM25 adapter
     refactored to satisfy it, then add a `sqlite-vec` vector adapter with a pluggable embedding
     provider (local Ollama or hosted, behind an interface). Hybrid retrieval = reciprocal-rank
     fusion across adapters, not the round-robin interleave the doc currently prescribes.
  3. **Target Files**: `shared/rag/retriever.interface.md`, `shared/rag/adapters/vector.md`,
     `shared/mcp/internal/tools/retriever.go`, `shared/mcp/internal/tools/bm25_retriever.go`.

- [ ] **Add episodic memory — there is currently no execution history at all**
  1. **Problem**: "Organizational memory" is semantic only: markdown KIs plus ADRs. There is no
     episodic store — no record of what was attempted, what failed, what a retry changed, or what a
     human corrected. `pipeline-trace.json` is specified to hold that and **does not exist anywhere
     in the repo**. Consequently nothing can learn from execution, which blocks Phase 3 entirely.
  2. **Architectural Fix**: An append-only run store (sqlite) keyed by `run_id`/`stage_id` capturing
     inputs, outputs, tool calls, retries, gate decisions, and human edits. Index it as a retrievable
     corpus so the planner can condition on "how did this go last time."
  3. **Target Files**: `shared/skills/pipeline-trace/SKILL.md`, `shared/telemetry/event-schema.md`,
     `shared/rag/retriever.interface.md` (new `episodic` corpus), new `internal/memory/`.

- [ ] **Generate the memory registry; add eviction**
  1. **Problem**: `shared/memory-registry.json` is hand-maintained. KI curation is four separate LLM
     skills — `memory-engineer`, `memory-compression`, `memory-expansion`, `forgetting-engine` —
     performing garbage collection by natural language, each requiring a human approval pass. There
     is no index, no recency decay in code, no automatic eviction, and no measure of whether any KI
     was ever retrieved and used.
  2. **Architectural Fix**: Generate the registry from frontmatter at build time and validate it in
     CI. Track retrieval hit-counts and last-used from the retriever, and drive staleness/eviction
     from that signal rather than from an LLM's monthly judgement. Keep human approval for deletion;
     automate the *detection*.
  3. **Target Files**: `shared/memory-registry.json`, `shared/skills/memory-engineer/SKILL.md`,
     `shared/skills/forgetting-engine/SKILL.md`, `scripts/health-check.sh`.

- [ ] **Structurally separate retrieved content from instructions**
  1. **Problem**: KIs are read whole into the context window.
     `shared/rules/memory-trust-boundary.md` correctly identifies synced org KIs (`sync_source`
     frontmatter, ADR-003 pull) as an injection vector — and then mitigates it by *asking the model
     in a prompt* to treat other prompt text as data. The defense occupies the same channel as the
     attack. `sync-memory.sh` validates frontmatter schema only; body content enters agent context
     unaudited.
  2. **Architectural Fix**: Retrieved content is delivered in a structurally distinct, clearly
     delimited channel with provenance metadata attached, never concatenated into the instruction
     region. Add a deterministic pre-ingestion scanner in `sync-memory.sh` for imperative-override
     patterns, so untrusted bodies are flagged before they reach any model. Prompt-level caution
     stays as defense-in-depth, not as the primary control.
  3. **Target Files**: `shared/rules/memory-trust-boundary.md`, `scripts/sync-memory.sh`,
     `shared/mcp/internal/tools/search_ki_tool.go`, `shared/agents/analyst.md`.

### Governance & Auditing

- [ ] **Emit OpenTelemetry with GenAI semantic conventions**
  1. **Problem**: `event-schema.md` has no `token_count`, no `cost`, no `trace_id`, no `span_id`, and
     no parent/child correlation — only a `pipeline_id` convention buried in free-form `metadata`.
     `duration_ms` and `pipeline-trace.json`'s `durationSeconds`/`budgetUtilization` are nominally
     recorded by an LLM that cannot measure elapsed time. No OTel is emitted anywhere, despite
     `architecture-guardrails.md` #8 mandating it from the adapter layer and
     `testing-conventions.md` requiring it on every BDD scenario. Answering "what did this pipeline
     cost" is currently impossible.
  2. **Architectural Fix**: OTel SDK in the executor and the MCP server. Spans per stage and per tool
     call, using GenAI semconv (`gen_ai.operation.name`, `gen_ai.usage.input_tokens`,
     `gen_ai.usage.output_tokens`, `gen_ai.request.model`), tool-call payloads as span attributes
     with a size cap and secret redaction, real wall-clock latency, real trace propagation across
     agent handoffs. OTLP export; keep `events.jsonl` as a local file exporter for offline mode.
  3. **Target Files**: `shared/telemetry/event-schema.md`, `shared/telemetry/event-recorder.md`
     (delete), `shared/mcp/internal/logging/logger.go`, new `internal/telemetry/`.

- [ ] **Resolve the telemetry schema contradiction and generate the schema** *(audit H6)*
  1. **Problem**: `event-recorder.md` instructs: "**Never** invent a new `event_type` — refuse if the
     caller passes one not documented." `event-schema.md` documents exactly six types.
     `policy-evaluator.md` emits `policy.evaluated`, `policy.conflict`, `policy.skipped` — none
     documented — under the key **`event`**, while the schema requires **`event_type`**.
     `orchestrate/SKILL.md` step 12 emits `workflow.completed`, also undocumented.
     `event-schema.md:5` admits it: "noted undocumented `policy.evaluated` and `contract.retry` types
     … (schema entry pending)." A conformant recorder rejects every event the policy layer produces.
  2. **Architectural Fix**: One Go enum as source of truth; generate both the JSON Schema and the
     documentation table from it. Emission moves into the executor, so an undocumented event type is
     a compile error rather than a prose violation.
  3. **Target Files**: `shared/telemetry/event-schema.md`, `shared/orchestration/policy-evaluator.md`,
     `shared/skills/orchestrate/SKILL.md`, new `internal/telemetry/events.go`.

- [ ] **Build the hook executor — the event catalog has no emitter and no dispatcher**
  1. **Problem**: `shared/hooks/` defines a schema and four example hooks against events
     `on-artifact-write`, `on-validation-pass`, `on-ki-created`, `on-retrospective-written`. **No code
     emits any of these events and nothing dispatches them.** `on-retrospective-written.yaml`
     additionally carries `description:` and `guardrails:` keys that `hooks-schema.md` does not
     define — including "draftOnly: true is non-negotiable and must not be overridden," a
     security-relevant constraint expressed as free text in a file no parser reads.
  2. **Architectural Fix**: A real in-process event bus in the executor with a typed event catalog, a
     hook loader that validates against the schema (rejecting unknown keys), sandboxed `script`-type
     execution with timeouts, and enforcement of hook-level guardrails as code. Every hook invocation
     emits telemetry.
  3. **Target Files**: `shared/hooks/hooks-schema.md`, `shared/hooks/*.yaml`,
     `shared/hooks/examples/*.yaml`, new `internal/hooks/`.

- [ ] **Make the eval harness provider-agnostic and give it a baseline**
  1. **Problem**: `scripts/run-agent-evals.sh` is 428 lines of real work, but it shells out to
     `claude --bare -p` and hardcodes `JUDGE_MODEL="claude-haiku-4-5-20251001"`. Step 4 is
     "regression comparison against the previous eval in `shared/evaluation/agent-evals/`" — that
     directory contains **only a README**, so there is no baseline and the regression check has never
     fired. In CI it runs only on push-to-main behind `AGENT_EVAL_ENABLED == 'true'`, so agent prompt
     changes reach users ungated.
  2. **Architectural Fix**: Abstract the generation and judging calls behind a provider interface
     (env-selected: Anthropic, Bedrock, Vertex, local). Commit baseline eval records for every agent.
     Run evals on **pull requests** against the changed agents, gate merge on regression, and store
     scores as trend data. This is the async CI/CD evaluation pattern L3 requires, and the harness is
     ~80% built.
  3. **Target Files**: `scripts/run-agent-evals.sh`, `shared/evaluation/agent-evals/`,
     `.github/workflows/framework-ci.yml`, `shared/evaluation/agent-eval-harness-design.md`.

- [ ] **Put the Go in CI and fix the workflow's own rule violations** *(audit H9)*
  1. **Problem**: `framework-ci.yml` runs five bash/python scripts and **zero Go steps** — no
     `go build`, `go test`, `go vet`, or `golangci-lint`. There is no `.golangci.yml`. `shared/mcp/`
     has **0.0% coverage in every package** against the non-negotiable 85% rule. The workflow also
     has **no `permissions:` block** and pins `actions/checkout@v4` (mutable tag) across all six jobs
     — both violations of `iac-conventions.md`, which `loom-release.yml` gets right. Meanwhile
     `scripts/test-agents.sh` reports "20 passed, 0 failed, 32 skipped" from **one**
     `actual-output.md` across 33 fixture dirs, because SKIP exits 0 by design.
  2. **Architectural Fix**: Add build/test/vet/lint jobs for both modules; write the `.golangci.yml`
     mandated elsewhere with `gocyclo` at 6; add a coverage ratchet. Add `permissions:` and SHA-pin
     every action. Make the agent suite fail on missing fixtures rather than passing. Then run the
     framework's own `verify_dependencies` and `analyze_complexity` against this repo in CI — the
     first *will* fail on the `domain/tool.go` violation, which is the point.
  3. **Target Files**: `.github/workflows/framework-ci.yml`, new `.golangci.yml`,
     `scripts/test-agents.sh`, `scripts/ci-check.sh`, `Makefile`.

- [ ] **Move counter agents out of the synchronous execution graph** *(audit H8)*
  1. **Problem**: `audit-composition-pattern.md` makes counter-agent invocation the default for every
     contract-bound stage, in-band and blocking: "Producer stage finishes → artifact written →
     auditor invoked → result evaluated." That roughly doubles LLM invocations and wall-clock per
     delivery. Four of eleven mappings also point auditors at artifacts they were never scoped for —
     `analyst → context-auditor` (scoped to `context-manifest.md` only),
     `accessibility-engineer → documentation-auditor` (scoped to README/ARCHITECTURE staleness),
     `qa-engineer → tool-validator` (scoped to `SKILL.md` frontmatter), and worst,
     `architect → rule-auditor` with **`onFail: halt`**, where `rule-auditor` audits
     `shared/rules/*.md` — framework files. Pointed at a customer's `architecture-notes.md`, it will
     halt their pipeline over findings about *this* repo.
  2. **Architectural Fix**: Auditors become an asynchronous post-delivery gate running on the PR
     against persisted artifacts in `docs/features/<name>/`, fanned out in parallel, off the critical
     path. Delete the four wrong mappings. Demote every auditor whose job is deterministic
     (frontmatter schema, broken links, dead paths) from LLM agent to linter.
  3. **Target Files**: `shared/orchestration/audit-composition-pattern.md`,
     `shared/workflows/feature-delivery-workflow.md`,
     `shared/agents/{context,documentation,rule,privacy}-auditor.md`,
     `shared/agents/tool-validator.md`.

- [ ] **Derive agent quality metrics from execution, not from reading markdown**
  1. **Problem**: `agent-scorecard/SKILL.md` scores four metrics computed by an LLM reading persisted
     artifacts — e.g. the analyst's score is "fraction of required sections present and containing
     real content." Nothing measures retry count, human-edit rate, gate rejection rate, latency, or
     cost, because none of those are recorded (no `events.jsonl`, no `pipeline-trace.json` exist).
     The scorecard also self-describes one metric as "directional, not exact, until that
     dispute-tracking mechanism exists."
  2. **Architectural Fix**: Compute metrics from the episodic store and OTel spans: retries per
     stage, `edited_then_approved` rate per producer, gate rejection rate, p50/p95 latency, tokens
     and cost per stage. The LLM writes the *narrative*; the numbers come from telemetry.
  3. **Target Files**: `shared/skills/agent-scorecard/SKILL.md`,
     `shared/skills/pipeline-retrospective/SKILL.md`, `docs/agent-metrics/`.

---

## PHASE 3 — Scaling to Level 4 (Self-Learning Agentic Ecosystems)

### Self-Reflection & Error Recovery

- [ ] **Bound the developer↔reviewer loop and add an oscillation detector** *(audit H1)*
  1. **Problem**: `deliver-feature/SKILL.md:102` (step 21): on `CHANGES REQUESTED`, "repeat from step
     18 **until APPROVED**." The text explicitly scopes `maxContractRetries` to the *structural*
     check and calls the qualitative verdict "independent of structural check." The one loop that
     actually oscillates — two LLMs disagreeing about craftsmanship — has no ceiling, no backoff, no
     progress detector, and no cost budget. This is `architecture-guardrails.md` #5 violated in the
     flagship pipeline. There is no reflexion primitive anywhere in the tree; the only mention of the
     word is in `docs/roadmaps/agy.md`.
  2. **Architectural Fix**: Bounded Reflexion loop in the executor: max attempts, per-attempt output
     hashing, and a no-progress detector (two semantically near-identical attempts ⇒ escalate, never
     retry). Persist critic feedback to episodic memory as a structured `Reflection` record keyed to
     the stage, so attempt N+1 receives *accumulated* critique rather than only the latest report.
  3. **Target Files**: `shared/skills/deliver-feature/SKILL.md:102`,
     `shared/orchestration/audit-composition-pattern.md`, new `internal/orchestrator/reflexion.go`.

- [ ] **Define an error taxonomy and recovery strategies for novel failures**
  1. **Problem**: The entire error-handling policy is one line in `orchestrate/SKILL.md`: "On any
     unhandled error in a stage: checkpoint current state, surface the error, and stop." There is no
     classification, no recovery strategy selection, no fallback, no degraded mode, and no escalation
     ladder. Every novel error is a full stop requiring a human — the definition of *not* L4.
  2. **Architectural Fix**: Error taxonomy (`Transient | Contract | Capability | Resource | Novel`)
     with per-class strategies: retry-with-backoff, re-plan via the router, decompose-and-retry,
     substitute agent, escalate. Unclassified errors route to a reflection step that attempts
     classification before escalating, and the classification outcome is recorded so the taxonomy
     improves.
  3. **Target Files**: `shared/skills/orchestrate/SKILL.md` (Guardrails section),
     `shared/orchestration/interface.md`, new `internal/orchestrator/recovery.go`.

- [ ] **Add a global budget governor**
  1. **Problem**: No token ceiling, no dollar ceiling, no wall-clock ceiling, no per-run cap of any
     kind. `maxDiffLines` and `maxContractRetries` in `.claude/delivery-policy.yaml` are the only
     numeric limits and neither bounds spend. Combined with the unbounded loop above and doubled
     invocations from synchronous audits, a single stuck delivery can burn until a human notices.
  2. **Architectural Fix**: A `BudgetGovernor` seeded per run (tokens, dollars, wall-clock, tool
     calls), decremented from real OTel usage data, enforcing soft-warn then hard-halt. Budget
     exhaustion is a first-class terminal state that checkpoints cleanly and reports what was
     consumed where.
  3. **Target Files**: `.claude/delivery-policy.yaml` schema,
     `shared/orchestration/policy-evaluator.md`, new `internal/orchestrator/budget.go`.

### Prompt & Tool Evolution

- [ ] **Build a prompt registry with eval-gated promotion**
  1. **Problem**: "Self-learning" today is `learning-engine` scanning retrospectives and writing a
     markdown proposal to `.claude/feature-workspace/proposed-lessons.md` for a human to approve into
     `docs/lessons-learned/`. It never touches an agent prompt. Agent prompts are static markdown
     edited by hand, versioned by hand, and enforced by `scripts/check-agent-versions-ci.sh`
     requiring a semver bump. There is no mechanism — not even a manual one — for the system to
     propose, test, and adopt a prompt change.
  2. **Architectural Fix**: Treat prompts as versioned artifacts in a registry with multiple live
     variants per agent. A candidate variant is generated (from accumulated reflections and lessons),
     evaluated by the existing harness against committed baselines, and promoted only on a
     statistically meaningful win. Champion/challenger with automatic rollback on regression. This is
     the change that makes adaptation possible without code deployment, and the eval harness is
     already most of the machinery.
  3. **Target Files**: `shared/agents/*.md` → registry-backed, `scripts/run-agent-evals.sh`,
     `scripts/check-agent-versions-ci.sh`, `shared/skills/learning-engine/SKILL.md`, new
     `internal/prompts/registry.go`.

- [ ] **Capture the human-correction signal that the schema already describes**
  1. **Problem**: `event-schema.md:131` defines `edited_then_approved` — "the human edited the
     artifact (checksum changed between gate-presented and gate-approved) … the corrective-signal
     case `extract-lessons` and `retrospective` mine." This is the single highest-value learning
     signal in the design, it is precisely specified, and **it is emitted by nothing**. No
     `events.jsonl` exists. The framework describes its own reward signal and never collects it.
  2. **Architectural Fix**: Executor-side digest comparison at every gate (already required by the
     gate-reset work in Phase 1), emitting `edited_then_approved` with a structured diff of what the
     human changed. Store in episodic memory as labelled training signal keyed to the producing agent
     and stage. Feed it into prompt-variant generation.
  3. **Target Files**: `shared/telemetry/event-schema.md:112-144`,
     `shared/skills/deliver-feature/SKILL.md` (Gate Decision Telemetry section),
     `shared/skills/extract-lessons/SKILL.md`, `internal/orchestrator/gate.go`.

- [ ] **Close the loop: lessons must reach prompts automatically**
  1. **Problem**: `learning-engine` → `promote-memory` → `create-ki` → `memory-engineer` →
     `forgetting-engine` is a five-skill chain that terminates in a markdown file a human must read.
     Nothing consumes lessons-learned at runtime except a human deciding to edit a rule file.
     `forgetting-engine`'s expiry pass is itself an LLM invoked monthly by a hook that has no
     executor.
  2. **Architectural Fix**: Lessons become structured records with a scope
     (`agent | stage | global`) and are injected as retrieved context to the relevant agent at run
     time via the retriever, rather than being copied into prompt bodies by hand. Promotion to a
     permanent prompt change happens only through the eval-gated registry above. Retrieval hit-count
     drives expiry, replacing the LLM forgetting pass.
  3. **Target Files**: `shared/skills/learning-engine/SKILL.md`,
     `shared/skills/promote-memory/SKILL.md`, `shared/skills/forgetting-engine/SKILL.md`,
     `shared/rag/retriever.interface.md`.

### Open Ecosystem Integration

- [ ] **Build an MCP *client* — loom can currently only be a server**
  1. **Problem**: `cmd/mcp-server/main.go` and `register/register.go` expose framework tools
     *outward*. There is no MCP client anywhere in the codebase. `loom` cannot consume a third-party
     MCP server; it delegates that entirely to the host IDE. So the framework's own agents cannot use
     external tools except through whatever the host happens to provide — the framework has no tool
     ecosystem of its own, only a tool export.
  2. **Architectural Fix**: An MCP client runtime in the executor: server discovery from config,
     capability negotiation, tool namespacing to avoid collisions, per-server auth, timeouts, and
     circuit breakers. External MCP tools then populate the same registry as native ones, so agents
     cannot tell the difference and `tools:` declarations become portable.
  3. **Target Files**: new `internal/mcp/client.go`, `shared/mcp/internal/server/tool_provider.go`
     (registry merge), `shared/schemas/agent-frontmatter.schema.json`.

- [ ] **Break the Anthropic lock in the agent contract** *(audit H5)*
  1. **Problem**: `agent-frontmatter.schema.json` declares `tools` as a closed regex over Claude
     Code's built-in names — **the framework's own MCP tools `analyze_complexity`, `search_ki`, and
     `search_docs` are literally unrepresentable in agent frontmatter.** `model` is
     `^(inherit|claude-[a-z0-9-]+)$`, so a non-Anthropic model is a schema violation.
     `shared/model-defaults.yaml` confirms the practical consequence: `claude_code` has real tier
     mappings; **all six other platforms are `null` across all three tiers**, with `roo_code` and
     `cline` marked "TBD — depends on Epic 42 landing." The framework ships installers for nine
     platforms and portable model selection works on one.
  2. **Architectural Fix**: Replace the tool regex with a namespaced identifier pattern
     (`builtin:Read`, `mcp:<server>/<tool>`) validated against the generated tool registry. Replace
     the `model` regex with `model_tier` plus a provider-resolution table, and add a provider adapter
     layer so `light|default|heavy` resolves on Bedrock/Vertex/OpenAI/local. Fill in or explicitly
     deprecate the `null` rows.
  3. **Target Files**: `shared/schemas/agent-frontmatter.schema.json`, `shared/model-defaults.yaml`,
     `shared/platform-registry.json`, `scripts/resolve-model-tier.py`, all 40 `shared/agents/*.md`.

- [ ] **Publish agent cards so external agents can invoke loom agents**
  1. **Problem**: There is no agent-to-agent interoperability of any kind — no A2A, no agent card, no
     invocation endpoint, no capability advertisement. `register.FrameworkTools` is the only
     integration path and it requires a foreign system to be written in Go and to import the module.
     That is language-level embedding, not ecosystem participation. Interop is currently: "be a Go
     program."
  2. **Architectural Fix**: Generate agent cards from the capability registry (Phase 2) — identity,
     skills, input/output schemas, auth requirements — and expose loom agents over a networked
     invocation endpoint alongside the MCP HTTP transport. Track the A2A protocol as the likely
     standard; the capability registry work is the prerequisite either way and is not wasted if the
     standard shifts.
  3. **Target Files**: new `internal/interop/agentcard.go`,
     `shared/schemas/agent-frontmatter.schema.json`, `shared/mcp/cmd/mcp-server/main.go`.

- [ ] **Delete or compile `shared/mcp-patterns/go/`** *(audit H10)*
  1. **Problem**: ~1,200 lines of `//go:build ignore` copies of `shared/mcp/internal/`, in no Go
     module, referenced by no build, presented as the reference implementation teams should copy into
     their own projects. All 11 shared files have already diverged from the originals
     (`retriever.go` by 186 diff lines, `bm25_retriever.go` by 106). It ships bugs to downstream
     teams — including the concurrency defect its own comment documents — and can never be compiled,
     tested, or kept honest.
  2. **Architectural Fix**: Delete it. If a reference implementation is genuinely wanted, make it a
     compiled example module in the workspace with its own tests, so drift becomes a build failure.
     Copy-paste distribution of a Go library is the wrong mechanism when `register.FrameworkTools`
     already exists.
  3. **Target Files**: `shared/mcp-patterns/go/**` (delete), `shared/mcp-patterns/README.md`,
     `shared/mcp/register/register.go` (document as the supported path).

---

## Sequencing

The phases are not independent. Typed state (P1) is a hard prerequisite for dynamic routing (P2) and
for prompt evolution (P3) — a planner cannot route on markdown, and an evaluator cannot score a
document. The telemetry work in P2 is a hard prerequisite for every P3 item, because self-learning
requires a signal and there is currently none.

**If only three things get done**: typed state, an executor that owns gates and retries, and OTel
emission. Those three unblock roughly two-thirds of the remaining list.

---

*Compiled 2026-08-29 against `main` @ `59efe14`. All cited paths verified to resolve on that date.*
