# MCP Tool Build-Out Prompts: Context Engineering Framework

These prompts extend the existing Go MCP server with tools that support the Context Engineering Framework's agent pipeline, knowledge system, and observability layer.

**Prerequisites:**
- The MCP server already exists under `mcp/` and is written in Go
- The `docs/mcp/rag-example.md` prompt has been run (project knowledge tools exist)
- The `shared/` canonical layer exists with agents, skills, rules, and `platform-registry.json`

**Sequence:**
```
Prompt 1: pipeline_trace     — structured logging for agent pipeline runs
Prompt 2: friday_ship        — POST delivery summaries to the Friday dashboard
Prompt 3: search_knowledge   — query Knowledge Items by tag/domain for proactive RAG
Prompt 4: check_parity       — verify platform configs match canonical shared/ source
Prompt 5: agent_scorecard    — query agent performance metrics from past deliveries
```

Run these in order. Each prompt is self-contained but later prompts assume earlier tools exist.

---

## PROMPT 1 — `pipeline_trace` Tool

> **Goal: structured observability for every agent pipeline run.**

Read the existing MCP server implementation under `mcp/`.
Read `shared/skills/deliver-feature/SKILL.md` to understand the agent pipeline phases.
Read `docs/mcp/rag-example.md` to understand existing patterns for tool registration, response shape, and error handling.

Add a new MCP tool: `pipeline_trace`

This tool records and queries structured trace data for agent pipeline runs. It replaces ad-hoc logging with a consistent, queryable format.

### Tool operations

The tool accepts an `operation` field that selects the action:

**`start_run`** — called at the beginning of a `deliver-feature` pipeline run.
- Input: `{ "operation": "start_run", "feature_name": "string", "agents": ["analyst", "architect", ...] }`
- Creates a new trace file at `docs/features/<feature_name>/pipeline-trace.json`
- Returns: `{ "trace_id": "uuid", "started_at": "RFC3339", "status": "running" }`

**`record_phase`** — called after each agent completes its phase.
- Input: `{ "operation": "record_phase", "trace_id": "string", "agent": "string", "status": "completed|failed|skipped", "duration_seconds": number, "token_estimate": number, "artifacts_produced": ["analysis.md", ...], "iteration_count": number }`
- Appends the phase record to the trace file
- Returns: `{ "phase_index": number, "cumulative_duration": number }`

**`end_run`** — called when the pipeline completes or fails.
- Input: `{ "operation": "end_run", "trace_id": "string", "status": "completed|failed|aborted", "summary": "string" }`
- Finalizes the trace with total duration, status, and summary
- Returns: the complete trace object

**`query_traces`** — retrieves past traces for analysis.
- Input: `{ "operation": "query_traces", "limit": number, "status_filter": "completed|failed|all" }`
- Scans `docs/features/*/pipeline-trace.json` files
- Returns: `{ "traces": [...], "total": number }`

### Trace file schema

```json
{
  "trace_id": "uuid",
  "feature_name": "string",
  "started_at": "RFC3339",
  "ended_at": "RFC3339",
  "status": "completed|failed|aborted|running",
  "total_duration_seconds": number,
  "phases": [
    {
      "agent": "analyst",
      "status": "completed",
      "started_at": "RFC3339",
      "duration_seconds": number,
      "token_estimate": number,
      "artifacts_produced": ["analysis.md"],
      "iteration_count": 1
    }
  ],
  "summary": "string"
}
```

### Architecture rules

- The trace tool is a **thin adapter** — it writes and reads JSON files, nothing more.
- No external dependencies (no databases, no collectors). Files are the storage.
- The tool must handle concurrent access gracefully (file locking or append-only).
- All responses must use the existing project response patterns (structured JSON with `summary`, `items`, `references` fields).
- Unit test: trace lifecycle (start → record 3 phases → end → query).
- Unit test: query with status filter.
- Unit test: graceful handling of missing/corrupt trace files.

### Integration point

After this tool exists, `deliver-feature/SKILL.md` will be updated to call `pipeline_trace` at each phase boundary. That wiring happens outside this prompt — this prompt only builds the tool.

Verify:
```bash
# The tool should be registered and respond to a health check
go test ./mcp/... -run TestPipelineTrace -v
```

---

## PROMPT 2 — `friday_ship` Tool

> **Goal: ship delivery summaries to the Friday QA dashboard with approval gate enforcement.**

Read the existing MCP server implementation under `mcp/`.
Read `shared/rules/approval-gates.md` — specifically Gate #1 (Shipping to Friday).
Read `docs/mcp/rag-example.md` for existing patterns.

Add a new MCP tool: `friday_ship`

This tool sends a Cucumber JSON test summary to the Friday dashboard. It is the only tool in the system that performs an external HTTP mutation, so it enforces an explicit approval gate.

### Tool operations

**`prepare_payload`** — assembles the payload without sending.
- Input: `{ "operation": "prepare_payload", "feature_name": "string", "cucumber_json_path": "string" }`
- Reads the Cucumber JSON file, validates its structure
- Returns: `{ "payload_preview": { "scenarios": number, "passed": number, "failed": number, "pending": number }, "destination": "string", "ready": true|false, "validation_errors": [...] }`

**`ship`** — sends the payload to Friday. Requires explicit approval.
- Input: `{ "operation": "ship", "feature_name": "string", "cucumber_json_path": "string", "approval_token": "string" }`
- The `approval_token` must be the literal string the user typed to confirm (e.g., "ship" or "yes")
- POSTs the Cucumber JSON to the Friday dashboard endpoint
- Returns: `{ "status": "shipped", "response_code": number, "friday_url": "string" }`

### Configuration

The Friday dashboard URL must come from environment variables, never hardcoded:

```
FRIDAY_DASHBOARD_URL=https://friday.example.com/api/results
FRIDAY_API_KEY=${FRIDAY_API_KEY}
```

If either is missing, the tool returns a structured error explaining what to set.

### Architecture rules

- Never send without an explicit approval token matching a known gate phrase.
- Validate the Cucumber JSON structure before sending (must have `features[]` with `scenarios[]`).
- Set an explicit HTTP timeout (30 seconds).
- Use the project's existing HTTP client or create a minimal one with circuit breaker pattern.
- Log the ship event with low-cardinality message: `logger.Info({"feature": name, "scenarios": count}, "Shipped to Friday")`
- Unit test: payload assembly and validation.
- Unit test: rejection when approval token is missing or wrong.
- Integration test (optional, behind build tag): actual POST to a mock HTTP server.

### Standalone mode

If `FRIDAY_DASHBOARD_URL` is not set, `prepare_payload` still works (validates locally). `ship` returns a clear error: "Friday dashboard URL not configured. Set FRIDAY_DASHBOARD_URL environment variable."

Verify:
```bash
go test ./mcp/... -run TestFridayShip -v
```

---

## PROMPT 3 — `search_knowledge` Tool

> **Goal: proactive RAG retrieval from the Knowledge Items directory.**

Read the existing MCP server implementation under `mcp/`.
Read `docs/features/context-engineering-framework/TODO.md` — Epic 14 (Knowledge Items infrastructure).
Read `docs/mcp/rag-example.md` for existing search patterns.

Add a new MCP tool: `search_knowledge`

This tool searches Knowledge Items (KIs) — reusable patterns, bug fixes, decisions, and lessons — by tag, domain, or keyword. It is the retrieval layer for the framework's proactive RAG system.

### Knowledge Item format

KIs are markdown files with YAML frontmatter stored in two locations:
- `shared/knowledge/` — universal patterns (portable across machines)
- `.claude/knowledge/` — project-specific context (local only)

```markdown
---
name: ki-name
domain: security|architecture|testing|performance|deployment|general
tags: [sql-injection, input-validation, parameterized-queries]
created: 2026-07-01
source: delivery:feature-name | manual | auto-promoted
confidence: high|medium|low
---

## Pattern

[Description of the pattern, fix, or decision]

## When to apply

[Conditions that make this KI relevant]

## Example

[Code example or reference]
```

### Tool operations

**`search`** — find relevant KIs by query.
- Input: `{ "operation": "search", "query": "string", "domain": "string (optional)", "tags": ["string"] (optional), "limit": number (default 5) }`
- Searches both `shared/knowledge/` and `.claude/knowledge/`
- Matches against: name, tags, domain, and full-text body
- Returns: `{ "results": [{ "name": "string", "domain": "string", "tags": [...], "confidence": "string", "path": "string", "snippet": "string (first 200 chars of body)" }], "total_matches": number, "searched_dirs": [...] }`

**`get`** — retrieve a specific KI by name.
- Input: `{ "operation": "get", "name": "string" }`
- Returns: the full KI content with parsed frontmatter

**`list_domains`** — list all domains with KI counts.
- Input: `{ "operation": "list_domains" }`
- Returns: `{ "domains": [{ "name": "security", "count": 5 }, ...] }`

**`record_usage`** — track that a KI was referenced by an agent during a pipeline run.
- Input: `{ "operation": "record_usage", "ki_name": "string", "agent": "string", "feature_name": "string" }`
- Appends to `shared/knowledge/.usage-log.jsonl`
- Returns: `{ "recorded": true, "total_uses": number }`

### Search ranking

For version 1, use lexical/keyword matching with this priority:
1. Exact tag match (highest)
2. Domain match
3. Name substring match
4. Body keyword match (lowest)

Make the ranking interface-driven so it can be swapped for vector/semantic search later without changing the tool API.

### Architecture rules

- Read-only for KI content (never modify KIs through this tool — that's a separate `create-ki` skill).
- Usage logging is append-only (`.usage-log.jsonl`).
- Graceful degradation: if `shared/knowledge/` doesn't exist yet, return empty results with a helpful message suggesting `mkdir shared/knowledge/`.
- Unit test: search by tag, domain, keyword.
- Unit test: ranking order (tag match > domain match > body match).
- Unit test: graceful empty results when directory doesn't exist.
- Unit test: usage log append.

### Integration point

The context-engineer agent will call `search_knowledge` during manifest creation (Phase 0 of `deliver-feature`). The wiring happens in the context-engineer skill, not here.

Verify:
```bash
go test ./mcp/... -run TestSearchKnowledge -v
```

---

## PROMPT 4 — `check_parity` Tool

> **Goal: expose the config drift fitness function as an MCP tool for agent self-verification.**

Read the existing MCP server implementation under `mcp/`.
Read `scripts/check-parity.sh` to understand the current parity checking logic.
Read `shared/platform-registry.json` for the platform list and tiers.

Add a new MCP tool: `check_parity`

This tool runs the same checks as `scripts/check-parity.sh` but returns structured JSON instead of terminal output. Agents can call it to verify configs are in sync before or after making changes.

### Tool operations

**`check`** — run the full parity check.
- Input: `{ "operation": "check", "platform": "string (optional — filter to one platform)" }`
- Checks all generated configs against `shared/` canonical source
- Returns:
```json
{
  "status": "pass|drift",
  "checks": [
    {
      "platform": "cursor",
      "file": ".cursor/rules/architecture.mdc",
      "check": "frontmatter_valid",
      "status": "pass|fail",
      "detail": "string (only on fail)"
    }
  ],
  "summary": {
    "total_checks": 24,
    "passed": 24,
    "failed": 0,
    "missing": 0
  },
  "fix_command": "scripts/generate-configs.sh"
}
```

**`check_agent_roster`** — verify all 24 agents appear in every platform config.
- Input: `{ "operation": "check_agent_roster" }`
- Returns: per-platform agent coverage with any missing agents listed

### Checks performed

Mirror what `scripts/check-parity.sh` does:
1. Cursor `.mdc` files: exist, have valid YAML frontmatter, contain rule headings from `shared/rules/`
2. Agent roster: all agents from `shared/agents/` appear in every Tier 2/3 config
3. Core concepts: key terms (Clean Architecture, cyclomatic complexity, TDD, etc.) present in flat files
4. Claude Code symlinks: `.claude/agents`, `.claude/skills`, `.claude/rules` are symlinks pointing to `shared/`

### Architecture rules

- Read-only. Never modify configs — only report drift.
- Parse `shared/platform-registry.json` for the platform list rather than hardcoding.
- Read agent names from `shared/agents/*.md` frontmatter dynamically.
- Unit test: all-pass scenario.
- Unit test: drift detected (missing agent in roster).
- Unit test: missing config file.

Verify:
```bash
go test ./mcp/... -run TestCheckParity -v
```

---

## PROMPT 5 — `agent_scorecard` Tool

> **Goal: query agent quality metrics from past delivery artifacts.**

Read the existing MCP server implementation under `mcp/`.
Read `docs/features/context-engineering-framework/TODO.md` — Epic 13 (Agent performance metrics).
Read the `pipeline_trace` tool created in Prompt 1 (this tool reads its output).

Add a new MCP tool: `agent_scorecard`

This tool analyzes past delivery artifacts and pipeline traces to score each agent's quality over time. It answers: "Are our agents getting better or worse after prompt edits?"

### Tool operations

**`score`** — compute quality scores for all agents or a specific agent.
- Input: `{ "operation": "score", "agent": "string (optional)", "last_n_deliveries": number (default 10) }`
- Scans `docs/features/*/pipeline-trace.json` and the corresponding delivery artifacts
- Returns:
```json
{
  "period": "last 10 deliveries",
  "agents": [
    {
      "name": "security-reviewer",
      "deliveries_participated": 8,
      "metrics": {
        "completion_rate": 1.0,
        "avg_duration_seconds": 45,
        "avg_iterations": 1.2,
        "artifact_completeness": 0.95
      },
      "trend": "stable|improving|degrading",
      "trend_detail": "completion rate up 10% over last 5 vs prior 5"
    }
  ]
}
```

**`compare_versions`** — compare agent performance before and after a prompt edit.
- Input: `{ "operation": "compare_versions", "agent": "string", "before_version": "string", "after_version": "string" }`
- Reads `shared/agents/CHANGELOG.md` to find the version boundary
- Compares metrics from deliveries before vs after the version change
- Returns: `{ "before": { metrics }, "after": { metrics }, "verdict": "improved|degraded|no_change" }`

**`underperformers`** — flag agents that may need prompt tuning.
- Input: `{ "operation": "underperformers", "threshold_completion_rate": number (default 0.8) }`
- Returns agents below the threshold with specific recommendations

### Metrics computed

| Metric | Source | Meaning |
|---|---|---|
| `completion_rate` | pipeline-trace.json | % of phases this agent completed without failure |
| `avg_duration_seconds` | pipeline-trace.json | Average time per invocation |
| `avg_iterations` | pipeline-trace.json | Average feedback loops before approval |
| `artifact_completeness` | delivery artifacts | % of contract-required sections present and non-empty |

### Architecture rules

- Read-only. Never modify traces or artifacts.
- Graceful degradation: if no traces exist yet, return `{ "agents": [], "message": "No pipeline traces found. Run /deliver-feature to generate data." }`
- If `shared/agents/CHANGELOG.md` doesn't exist, `compare_versions` returns a helpful error.
- Unit test: scoring with mock trace data.
- Unit test: trend calculation (improving vs degrading).
- Unit test: empty trace directory handling.

### Future extension

When inter-agent contracts exist (Phase 3, Epic 5), `artifact_completeness` will check against the contract schema rather than heuristic section counting. Design the metric interface to support swapping the completeness checker.

Verify:
```bash
go test ./mcp/... -run TestAgentScorecard -v
```

---

## Final Verification

After all 5 prompts are complete, verify all tools are registered:

```bash
cd mcp/ && go build ./... && go test ./... -v
```

The MCP server should expose these new tools:
```
pipeline_trace    — agent pipeline observability
friday_ship       — delivery summary shipping with approval gate
search_knowledge  — Knowledge Item retrieval for proactive RAG
check_parity      — config drift detection as structured JSON
agent_scorecard   — agent quality metrics from past deliveries
```

All tools should:
- Follow existing MCP server patterns (tool registration, response shape, error handling)
- Have unit tests with >85% coverage
- Return structured JSON responses
- Handle missing data gracefully
- Use no external dependencies beyond what the MCP server already uses
