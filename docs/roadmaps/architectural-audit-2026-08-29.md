# Architectural Audit & Decoupling Roadmap — 2026-08-29

**Framework version**: v3.3.14
**Repository**: `loom` (`orieken/loom`)
**Scope**: architectural coupling, Level 2/3 readiness, audit & telemetry
**Method**: direct inspection of the tree — 566 markdown files (~52k lines), 113 Go files (~8.9k lines),
`.github/workflows/`, and both compiled Go modules. Every finding below cites a file that was read.

Companion documents: [`BUILD-ROADMAP.md`](BUILD-ROADMAP.md) — the live build plan derived from these
findings (start there for what to do next). [`agy.md`](agy.md) — earlier maturity roadmap, now
superseded and merged into BUILD-ROADMAP.

---

## 1. Executive Verdict

**`loom` is a prompt distribution system, not an agentic framework.** Of 8,896 lines of Go, ~2,000
are a Cobra file-installer (`cmd/loom/`) and ~1,264 are an MCP tool server (`shared/mcp/`). There is
no kernel, no execution loop, no scheduler, no state machine, no router, and no retry governor.
Every subsystem this audit was asked to evaluate — orchestration runtime, telemetry recorder, policy
evaluator, retriever tiers, checkpointing, parallel branches — is a markdown file instructing an LLM
to *role-play* that subsystem. `shared/orchestration/interface.md` defines a `stages:` DAG schema
that nothing parses. `shared/telemetry/event-recorder.md` is a 150-line prompt asking a model to
append one JSON line to a file. `shared/orchestration/policy-evaluator.md` delegates auto-approval of
git-commit gates to a language model interpreting YAML glob conditions, with no deterministic
evaluator anywhere in the repo. The tell is in the spec itself: `parallelStrategy:
sequential-simulation` is the *default*, defined as "execute each in sequence." Parallelism is
documented, not implemented. Every guarantee — retry ceilings, checkpoint integrity, gate
enforcement, schema conformance — is a suggestion that degrades under exactly the load, context
pressure, and adversarial input that enterprise use implies.

**The second sin: none of it has ever run.** There is not a single `events.jsonl` in this
repository. Not one `pipeline-state.json` or `pipeline-trace.json`. `shared/evaluation/agent-evals/`
contains a README and nothing else. `docs/features/` holds two deliveries, and the larger has three
artifacts out of the fifteen the pipeline mandates. `.claude/feature-workspace/analysis.md` sits flat
at the root — the exact "legacy pre-Epic-63 singleton" state `deliver-feature` ships migration code
to repair. The MCP server has **0.0% test coverage in every package** against a self-declared
non-negotiable 85% rule, and `framework-ci.yml` never runs `go build`, `go test`, `go vet`, or
`golangci-lint`. There is no `.golangci.yml`. And `health-check.sh` reports **283 passed, 0 failed**.
That number is the whole problem: the fitness functions verify that files exist and frontmatter
parses. They verify no architectural property, which is the inverse of what
`architecture-guardrails.md` #7 demands.

---

## 2. High-Severity Findings

Severity is ordered by blast radius under enterprise load, not by fix cost.

### H1 — Unbounded agent loop; no circuit breaker
`shared/skills/deliver-feature/SKILL.md:102` (step 21): on `CHANGES REQUESTED`, "repeat from step 18
**until APPROVED**." The text explicitly scopes `maxContractRetries` to the *structural* check and
calls the qualitative verdict "independent of structural check." The one loop that actually
oscillates — developer ↔ code-reviewer, two LLMs disagreeing about craftsmanship — has no ceiling,
no backoff, no no-progress detector, and no cost budget. Directly violates
`shared/rules/architecture-guardrails.md` #5, which forbids hand-rolled retry loops and mandates
`CircuitBreaker` or `ExponentialBackoffStrategy`.

### H2 — Full corpus re-index on every search query
`shared/mcp/internal/tools/search_docs_tool.go:82` calls `EnsureIndex` inside `Execute`.
`shared/mcp/internal/tools/bm25_retriever.go:70` walks the entire docs tree, `os.ReadFile`s every
`.md`, and runs **one sqlite transaction per file** — no mtime check, no dirty tracking, no content
hash. Every `search_docs` call is O(corpus) disk I/O plus O(n) transactions. The drifted copy
documents the second half of the bug at `shared/mcp-patterns/go/tools/bm25_retriever.go:40`:
"EnsureIndex is not safe to call concurrently with itself" — and MCP servers field concurrent tool
calls. Deleted docs are never evicted; the `DELETE` only fires for paths being re-inserted.

### H3 — Every tool discards its context; no timeouts, no cancellation
All six tools sign `Execute(_ context.Context, ...)`:
`analyze_complexity_tool.go:54`, `check_accessibility_tool.go:52`,
`check_ubiquitous_language_tool.go:50`, `search_docs_tool.go:64`, `search_ki_tool.go:56`,
`verify_dependencies_tool.go:40`. `analyze_complexity` accepts an arbitrary `projectPath` from
model-controlled arguments with no root confinement and no path validation, then hands it to an
unbounded, uncancellable `filepath.Walk`. `projectPath: "/"` walks the disk and cannot be aborted.
`shared/rules/go-conventions.md` says "ALWAYS set explicit timeouts"; the tool layer sets none.

### H4 — The domain layer imports the transport SDK
`shared/mcp/internal/domain/tool.go` — commented as "the framework's first-class abstraction for
every capability" — imports `github.com/mark3labs/mcp-go/mcp` and `invopop/jsonschema`, and types its
own signatures in them (`mcp.ToolInputSchema`, `mcp.CallToolRequest`, `*mcp.CallToolResult`). That is
`architecture-guardrails.md` #1 ("Entities cannot import Frameworks/Libraries") violated in the one
file that defines the tool abstraction. The framework ships a `verify_dependencies` tool that would
catch this. It has never been run against this repository.

### H5 — The agent contract is a closed Anthropic allowlist
`shared/schemas/agent-frontmatter.schema.json`: `tools` is
`"^(Read|Write|Edit|MultiEdit|Bash|Glob|Grep|Task|Skill|WebFetch|WebSearch)(,...)*$"` with
`additionalProperties: false`. **The framework's own MCP tools — `analyze_complexity`, `search_ki`,
`search_docs` — are unrepresentable in agent frontmatter.** `model` is
`"^(inherit|claude-[a-z0-9-]+)$"`, so a non-Anthropic model is a schema violation — while
`cmd/loom/internal/platform/platform.go:39` ships installers for `gemini`, `openai-codex`,
`github-copilot`, `cursor`, `windsurf`, `roo-code`, `cline`, and `jetbrains`. No per-tool argument
constraint, rate limit, or timeout exists in the contract. The tool interface is neither standard nor
modular; it is one vendor's built-in menu, hardcoded.

### H6 — Telemetry and policy layers emit mutually incompatible events
`shared/telemetry/event-recorder.md` guardrail: "**Never** invent a new `event_type` — refuse if the
caller passes one not documented." `event-schema.md` documents exactly six.
`shared/orchestration/policy-evaluator.md` emits `policy.evaluated`, `policy.conflict`,
`policy.skipped` — none documented — under the key **`event`**, while the schema requires
**`event_type`**. `shared/telemetry/event-schema.md:5` admits it: "noted undocumented
`policy.evaluated` and `contract.retry` types … (schema entry pending)". A conformant recorder
rejects every event the policy layer produces. `workflow.completed`, emitted by
`shared/skills/orchestrate/SKILL.md` step 12, appears nowhere in the schema either.

### H7 — Gate auto-approval is an LLM interpreting globs
`policy-evaluator.md` + `shared/policies/policy-schema.md` define a condition language
(`filePaths.noneMatch: "**/security/**"`, `diffLines.lessThan`, `not:`, `any:`) whose evaluator is a
prompt. A security-relevant authorization decision — *may this pipeline commit without a human* — is
resolved by natural-language reasoning over YAML, with the kill-switch, the always-human gate list,
and the conflict-resolution table all living in the same prose the model is free to misread. The
schema's own worked example contains a duplicate `diffType:` key in one map, which is invalid YAML —
proof no parser has ever seen it.

### H8 — Audit agents are synchronous blocking stages, and four are misrouted
`shared/orchestration/audit-composition-pattern.md` makes counter-agent invocation the *default* for
every contract-bound stage, in-band and blocking: "Producer stage finishes → artifact written →
auditor invoked → result evaluated." That roughly doubles LLM invocations and wall-clock per
delivery, with no async path. The mapping table also assigns auditors to artifacts they were never
scoped for:

| Mapping in the pattern table | Auditor's actual declared scope |
|---|---|
| `analyst → analysis.md → context-auditor` | `context-manifest.md` pruning discipline only |
| `accessibility-engineer → accessibility-report.md → documentation-auditor` | `README.md` / `ARCHITECTURE.md` staleness |
| `qa-engineer → qa-report.md → tool-validator` | `shared/skills/*/SKILL.md` frontmatter |
| `architect → architecture-notes.md → rule-auditor`, **`onFail: halt`** | `shared/rules/*.md` internal consistency |

The last is the live grenade: `rule-auditor` audits *framework* rule files. Pointed at a project's
`architecture-notes.md` with `halt` semantics, it will block a customer pipeline on findings about
this repo, not theirs.

### H9 — CI does not compile, test, or lint the product
`.github/workflows/framework-ci.yml` runs five bash/python scripts and zero Go steps. The 8,896 lines
that are the only executable artifact are verified on tag-push by goreleaser and nowhere else. The
same file has **no `permissions:` block** and pins `actions/checkout@v4` — a mutable tag — across all
six jobs, both violations of `shared/rules/iac-conventions.md`. `loom-release.yml` gets both right,
so this is drift, not ignorance. Separately, `scripts/test-agents.sh` reports "20 passed, 0 failed,
32 skipped" from **one** `actual-output.md` against 33 fixture directories; SKIP exits 0 by design.
The agent suite is green by construction.

### H10 — 1,200 lines of `//go:build ignore` fork, already drifted
`shared/mcp-patterns/go/` is a copy of `shared/mcp/internal/` excluded from compilation, in no Go
module, referenced by no build. All 11 shared files have diverged (`retriever.go` by 186 diff lines,
`bm25_retriever.go` by 106). It is presented as the reference implementation teams should copy, and
it can never be compiled, tested, or kept honest.

### H11 — 532 KB of duplicated prompt payload at the repo root
`.cursorrules` and `.windsurfrules` are byte-identical (`c662f42f…`, 76 KB each). `AGENTS.md` and
`.openai.md` have drifted apart. `.roomodes` is 233 KB. That is roughly 20k tokens of static rules
prepended to every request on those platforms before any agent does any work — a fixed per-call tax
that multiplies across a 15-stage pipeline with synchronous audits (H8).

---

## 3. Level 2 / Level 3 Gap Analysis

### Level 2 — tool utilization & deterministic workflows: not reachable as built

| Gap | Evidence |
|---|---|
| No tool registry | `shared/mcp/internal/server/tool_provider.go` — `buildFrameworkTools()` is a hardcoded slice literal. Adding a tool requires editing and recompiling the handler. No discovery, manifest, versioning, or capability negotiation. |
| No provider abstraction | H5. Tool names and model IDs are Anthropic literals baked into a JSON Schema regex. Porting to Bedrock, Vertex, or an OSS model is a schema rewrite, not a config change. |
| No validated execution boundary | Input schemas exist for the LLM's benefit; server-side enforcement is `args["x"].(string)` type assertions. No argument allowlisting, output size cap, per-tool timeout, or tool-call audit log. |
| State between agents is unmanaged markdown | No typed state object. Agents pass whole files. The only bloat control is `summarize-artifact`, itself an LLM call producing a lossy ~200-word surrogate. No schema, field-level access, size ceiling, or provenance. |
| "Deterministic" means header presence | `validate-artifact` checks that required `##` headings exist — presence, not semantics. That is not determinism in any sense a compliance team will accept. |
| Parallelism does not exist | `sequential-simulation` is the documented default in `shared/orchestration/interface.md`. |

### Level 3 — autonomous planning & collaboration: no foundation at all

| Gap | Evidence |
|---|---|
| No dynamic routing anywhere | Every path is a numbered step with static conditionals (`if analysis.md has Data Model Changes != "None"`). `/orchestrate` reads an ordered `stages:` list. No planner, goal decomposition, router, or agent-selects-next-agent mechanism. |
| No capability index to route *from* | 40 agents exist as markdown prose descriptions. There is no machine-readable capability catalog a router could query. |
| Error recovery absent where it matters | No circuit breaker, backoff, dead-letter, token/dollar ceiling, oscillation detector, or output-similarity check between producer attempts. The single bound (`maxContractRetries: 3`) governs the check that rarely loops and excludes the loop that does (H1). |
| Memory is a folder curated by more prompts | `shared/rag/retriever.interface.md` is a genuinely clean contract — references not content, bounded top-K, corpus isolation, no side effects. Three of its four backends are markdown. Only BM25 exists in code. No vector store, embeddings, graph, recency decay, or eviction. `shared/memory-registry.json` is hand-maintained. |
| Memory GC is four more LLM calls | `memory-engineer`, `forgetting-engine`, `memory-compression`, `memory-expansion` perform garbage collection by natural language. |
| Memory and prompt entangled by construction | KIs are read whole into context. `shared/rules/memory-trust-boundary.md` correctly identifies synced org KIs as an injection vector, then mitigates it by *asking the model in a prompt* to treat other prompt text as data. The defense occupies the same channel as the attack. |

### Audit & telemetry: the honest answer is zero

Traces, token usage, tool-call payloads, and latency cannot be extracted without rewriting the core,
because there is no core to instrument. `shared/telemetry/event-schema.md` has no `token_count`, no
`cost`, no `trace_id`, no `span_id`, no parent/child correlation — only a `pipeline_id` convention in
free-form `metadata`. `duration_ms` is nominally captured by an LLM that cannot measure elapsed time;
`pipeline-trace.json` records `durationSeconds` and `budgetUtilization` the same way. No
OpenTelemetry is emitted anywhere, despite `architecture-guardrails.md` #8 mandating OTel from the
adapter layer and `testing-conventions.md` requiring OTel on every BDD scenario. The repository
contains zero telemetry output of any kind.

The built-in audit pipelines reinvent standard testing rather than integrating it:

| Built-in | Reinvents | Already mandated by loom's own conventions |
|---|---|---|
| `analyze_complexity` | `gocyclo`, `ruff C901`, `detekt ComplexMethod`, `SwiftLint cyclomatic_complexity` | yes, all four |
| `verify_dependencies` | `go-arch-lint`, `import-linter`, ArchUnit | partially |
| `check_accessibility` | `axe-core`, `eslint-plugin-jsx-a11y` | yes |

Each reimplementation is weaker than the tool it replaces, is not wired into CI, and has 0% coverage.
The one thing that would be genuinely valuable and does not exist: an OTel GenAI-semconv exporter
emitting spans with token counts and tool-call payloads.

---

## 4. Refactoring Roadmap

Ordered by dependency, then by ratio of risk retired to effort spent.

### R1 — Split the repository; decide what `loom` actually is
**Blocking. Everything else depends on the answer.**

- [ ] Separate into three concerns: `loom-kernel` (a real Go executor), `loom-content` (agents,
      skills, rules, KIs as versioned data), `loom-cli` (the installer — already the best-factored
      code in the repo; `platform.go`'s installer registry and the `FileWriter` seam are correct).
- [ ] Answer explicitly, in an ADR: does the kernel **execute** pipelines, or only **validate and
      distribute** content that a host runtime executes? Both are defensible. Shipping a markdown
      specification of an engine while claiming Level 3 orchestration is not.
- [ ] Whichever answer, update `README.md` and `docs/ARCHITECTURE.md` to match it. The current docs
      describe the aspiration as though it were the implementation.

**Closes**: the framing problem underneath every other finding.

### R2 — Put the Go in CI and turn the framework's own rules on itself
**Highest ratio in the audit. Roughly one day.**

- [ ] Add `go build ./...`, `go test ./... -race`, `go vet`, `golangci-lint run` to
      `framework-ci.yml` for **both** modules (root and `shared/mcp/`).
- [ ] Write the `.golangci.yml` the framework tells every user to write, with `gocyclo` at 6.
- [ ] Add a coverage gate. Start at whatever `shared/mcp/` actually reaches once tests exist; ratchet
      toward 85%.
- [ ] Delete `shared/mcp-patterns/go/`, or promote it to a compiled example module. **(H10)**
- [ ] Add the `permissions:` block and SHA-pin every action in `framework-ci.yml`. **(H9)**
- [ ] Run the framework's own `verify_dependencies` and `analyze_complexity` against this repository
      in CI. The first will fail on **H4** — that is the point. A framework whose central claim is
      verifiable architecture must have fitness functions that can fail.

**Closes**: H9, H10; surfaces H4.

### R3 — Make deterministic the three things that must be deterministic
Not everything needs an engine. Three things do, and none should be a prompt.

- [ ] **State machine + retry governor.** A real executor owning `pipeline-state.json`, enforcing a
      global attempt budget, a per-edge retry ceiling, a token/cost ceiling, and an oscillation
      detector (hash successive producer outputs; two near-identical attempts halt rather than
      retry). **(H1)**
- [ ] **Policy evaluation.** Replace `policy-evaluator.md` with `cel-go` or OPA. The condition
      language already designed maps almost directly onto CEL. The always-human gate list becomes a
      compiled constant, not a paragraph. **(H7)**
- [ ] **Telemetry.** Delete `event-recorder.md`. Emit from the executor and the MCP server as OTel
      spans using GenAI semantic conventions — `gen_ai.usage.input_tokens`,
      `gen_ai.usage.output_tokens`, tool-call spans with payload attributes, real wall-clock latency,
      real `trace_id`/`span_id` parent-child structure. Keep `events.jsonl` as a local exporter for
      offline mode. **(H6, and the entirety of the audit/telemetry dimension.)**

**Closes**: H1, H6, H7. Delivers the first thing an enterprise buyer will ask for.

### R4 — Fix the tool and provider boundary

- [ ] Strip `mcp-go` types out of `domain.Tool`; define transport-free request/result types and adapt
      at the server edge. **(H4)**
- [ ] Replace the hardcoded `buildFrameworkTools` slice with a name-keyed registry populated at
      startup, so tools are additive without recompiling the handler.
- [ ] Propagate `ctx` through every `Execute` and every analyzer walk; set a server-side per-tool
      deadline; confine `projectPath` to a configured root. **(H3)**
- [ ] Widen `agent-frontmatter.schema.json`: allow MCP-qualified tool names, add per-tool permission
      scoping, and replace the `claude-*` model regex with `model_tier` plus a per-platform
      resolution table (`scripts/resolve-model-tier.py` is the seed). **(H5)**

**Closes**: H3, H4, H5. Unblocks Level 2.

### R5 — Move audits out of the synchronous graph

- [ ] Take counter agents out of the in-band position in `audit-composition-pattern.md`. Run them as
      an asynchronous delivery-pipeline gate — on the PR, in CI, against the persisted artifact set
      in `docs/features/<name>/` — where they fan out in parallel, cost nothing on the critical path,
      and produce reviewable output.
- [ ] Fix the four mismatched producer→auditor mappings; delete `architect → rule-auditor / halt`
      outright. **(H8)**
- [ ] Demote every auditor whose job is a deterministic check (frontmatter schema, broken links, dead
      paths) from LLM agent to linter.

**Closes**: H8. Roughly halves pipeline latency and LLM spend per delivery.

### Deferred (real, lower priority)

- [ ] **H2** — cache the BM25 index by mtime/content-hash; move `EnsureIndex` out of `Execute` behind
      an explicit reindex tool or a file-watcher; add a mutex for concurrent calls.
- [ ] **H11** — stop committing 532 KB of generated near-duplicate prompt payload; generate at
      install time from `shared/` rather than checking the artifacts in.
- [ ] Migrate `.claude/feature-workspace/analysis.md` out of the legacy flat layout, or delete it.

---

## 5. Closing Assessment

The specifications in this repository are well designed on paper — the retriever contract, the
Expand/Contract discipline, the counter-agent concept, the installer's platform registry, the
trust-boundary threat model. Almost none of them are implemented, and the verification layer reports
283/0 green across the gap. The specification quality is real and worth keeping. The path forward is
to pick the smallest kernel that makes three or four of those specs executable and falsifiable,
delete the markdown that pretends to be the other twenty, and stop shipping fitness functions that
cannot fail.

---

*Audit performed 2026-08-29 against `main` @ `59efe14`. Findings cite files as read on that date.*
