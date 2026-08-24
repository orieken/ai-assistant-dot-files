---
name: mcp-add
description: Extend or retrofit an existing MCP (Model Context Protocol) server to follow the framework's tool / persona / workflow pattern with Clean Architecture layers. Handles three sub-flows — retrofit a non-conformant MCP, add a new tool/persona/workflow to a conformant one, or do both together. Auto-detects language (Go / TypeScript) and transport (stdio / SSE / streamable HTTP) from the existing repo. Uses shared/blueprints/mcp-server.md as the authoritative pattern source.
triggers:
  keywords: ["mcp-add", "add mcp tool", "add tool to mcp", "extend mcp", "retrofit mcp", "mcp refactor", "add persona", "add workflow", "expose mcp"]
  intentPatterns: ["Add a tool to *", "Extend the MCP *", "Retrofit * to follow our pattern", "Expose * as an MCP tool", "/mcp-add *"]
standalone: true
---

## When To Use

When the user has an **existing MCP server** and wants to either:
- **Retrofit** it to follow the framework's tool / persona / workflow pattern (it currently mixes concerns, has business logic in the entrypoint, calls SDKs directly from tools, etc.), or
- **Extend** it with a new tool / persona / workflow following the pattern, or
- **Both** — the common case: the existing MCP is partly there, needs some retrofit AND some new capabilities.

Do NOT use when:
- The user is creating a brand-new MCP server → use `/bootstrap-project` with `mcp-server` pattern.
- The user wants to add a feature that doesn't fit tool / persona / workflow (e.g., a background job scheduler) → use `/new-feature`.
- The user wants to refactor a non-MCP project to a different pattern → use `/refactor-to-pattern`.
- The MCP server is in a language other than Go or TypeScript — this skill only handles those two; flag and stop.

## Context To Load First

1. `shared/blueprints/mcp-server.md` — the authoritative pattern definition (layer structure, trinity vocabulary, testing pyramid, non-negotiables). Every plan produced by this skill measures against it.
2. `shared/rules/architecture-guardrails.md` — non-negotiable hard constraints (Clean Architecture direction, no hardcoded secrets, explicit timeouts, resilience primitives, OTel boundaries).
3. `shared/rules/design-principles.md` — Sandi Metz limits, Fowler refactorings by name.
4. `shared/rules/testing-conventions.md` — testing pyramid, per-language framework picks.
5. `shared/rules/<language>-conventions.md` — loaded once the language is detected in Phase 1.
6. The **existing MCP server's source tree** — this skill reads it to auto-detect language, transport, tool inventory, and pattern conformance before proposing changes.

## Process

### Phase 1 — Discovery Interview (one question at a time)

Follow the one-question-per-message discipline. Do not batch.

1. **Where is the existing MCP server?** — Ask for the repo path. Verify it exists. Verify it's actually an MCP server by checking for one of:
   - Go: `go.mod` importing a Model Context Protocol SDK.
   - TypeScript: `package.json` with `@modelcontextprotocol/sdk` in dependencies.
   - If neither is found, ask the user to confirm this is an MCP server and describe what SDK/protocol layer it uses. If they can't, stop and recommend `/bootstrap-project` instead.

2. **Which sub-flow?** — Present three options:
   - `retrofit` — the existing MCP has drift from the framework's pattern (business logic in the entrypoint, SDKs called from tools, missing tests, etc.). Refactor without adding new capabilities.
   - `extend` — the existing MCP already follows the pattern. Add a new tool / persona / workflow.
   - `both` — do a retrofit AND add new capabilities in one session.

3. **What is the new capability?** (only if sub-flow is `extend` or `both`) — Ask ONE at a time:
   - Is it a **tool**, a **persona**, or a **workflow**? Trinity only — no "handlers", "managers", "services" at the domain layer.
   - Name (kebab-case for the file; PascalCase for the type if Go, PascalCase-with-`.type.ts` if TypeScript per `CLAUDE.md`).
   - Plain-language purpose (what problem it solves for the LLM caller).
   - Inputs (fields + types, will become the JSON schema).
   - Outputs (fields + types, will become the JSON schema).
   - External systems it touches (each becomes an adapter behind a domain interface).

4. **What does "done" look like?** — Minimum acceptance criteria for the change to be considered shipped. Feeds the qa-engineer downstream.

5. **Any constraints or non-negotiables?** — Auth model, rate limits, secrets handling, backward compatibility with existing tools, deployment window.

### Phase 2 — Auto-Detect + Gap Analysis (show your thinking)

Read the existing MCP source tree. Reason out loud before producing artifacts.

- **2a. Language detection.** Confirmed by which of `go.mod` / `package.json` is present. State the version.
- **2b. Transport detection.** Look for `StdioServerTransport`, `SSEServerTransport`, `StreamableHTTPServerTransport` (or their Go equivalents). Report which is in use.
- **2c. Tool / persona / workflow inventory.** Enumerate every tool/persona/workflow currently exposed. Report as a table: `Name | Type | File | External deps used | Has tests`.
- **2d. Pattern conformance scorecard.** For each non-negotiable in `shared/blueprints/mcp-server.md`'s Key Abstractions section, mark ✓ / ✗ / partial:
  - Domain layer has zero SDK imports?
  - External APIs behind adapter-layer interfaces?
  - Entrypoint (main.go / index.ts) is wiring-only (no business logic)?
  - Every tool has a JSON schema for inputs and outputs?
  - Every tool has unit tests with mocked adapters?
  - Every tool has an integration test against test containers or staging?
  - OTel spans emitted only from adapters, never from domain?
  - No hardcoded secrets?
  - Explicit timeouts on all outbound network calls?
  - Retries use `CircuitBreaker` or `ExponentialBackoff`, never hand-rolled loops?

- **2e. Change plan.** Based on 2d + Phase 1 inputs, produce:
  - **Retrofit list**: specific files that need refactoring, named Fowler operation for each (Extract Function, Move Method, Introduce Adapter, etc.).
  - **Additions list** (if `extend` or `both`): new files to create, matching the directory tree in `shared/blueprints/mcp-server.md`.
  - **Risk tags**: `[AMBIGUOUS]` / `[RISK-HIGH]` / `[ASSUMPTION]` for anything uncertain.

### Phase 3 — Plan Artifact + Approval Gate

Write `mcp-add-plan.md` at the repo root of the existing MCP server. This is the single approval surface — no source code is touched until the user approves this file.

```markdown
# MCP Add Plan: <server name>

**Server path**: <absolute path>
**Language**: <Go | TypeScript> <version>
**Transport**: <stdio | SSE | streamable HTTP>
**Sub-flow**: <retrofit | extend | both>
**Generated by**: mcp-add skill — <YYYY-MM-DD>

## Current Tool Inventory
<Table from 2c>

## Pattern Conformance
<Scorecard from 2d>

## Retrofit Plan
<Ordered list of Fowler-named refactorings, with file paths and one-line rationale each. Empty if sub-flow is `extend`.>

## Additions Plan
<New files to create with their layer assignment, matching the mcp-server.md directory tree. Empty if sub-flow is `retrofit`.>

## Testing Plan
<For every new/refactored tool: which pyramid levels (unit / integration / API contract), which agent writes them.>

## OTel Instrumentation Plan
<Spans to add, all at adapter/interceptor layer.>

## Risks
<Everything tagged [AMBIGUOUS] / [RISK-HIGH] / [ASSUMPTION] from 2e.>

## Downstream Agents
<Ordered slash commands to invoke after approval — see Phase 4.>
```

Then present:
```
MCP-ADD PLAN READY
==================
Server:        <name>
Sub-flow:      <retrofit | extend | both>
Retrofits:     <n files>
Additions:     <n files>
Tests to add:  <n>

Plan written to: <path>/mcp-add-plan.md

Type 'yes' to proceed with the plan (source code will be modified).
Type 'edit' to revise the plan first.
Any edit to the plan resets this gate.
```

### Phase 4 — Execution

Only run after explicit "yes." Executes in this order:

1. **Retrofits first** (if any) — one Fowler operation per commit-worthy change. Delegate the actual code edits to `developer` agent, then run `code-reviewer` per operation. Never batch multiple operations into one edit.
2. **Additions next** — for each new tool / persona / workflow:
   - `developer` — implement the domain type + use-case, tests first (Red → Green → Refactor).
   - `developer` — implement the adapter(s) for external deps.
   - `developer` — register in the entrypoint (thin wiring only).
   - `code-reviewer` — verify pattern conformance against `shared/blueprints/mcp-server.md`.
   - `security-reviewer` — verify secret handling, input validation, adapter isolation.
   - `qa-engineer` — write integration tests against test containers or staging.
   - `sre-engineer` — verify OTel spans emit from adapters only, span cardinality is bounded.
3. **Regression sweep** — run the existing MCP's full test suite to verify retrofits didn't break existing tools.
4. **Documentation** — `tech-writer` updates the README's tool inventory and adds/updates ADRs for material decisions.

### Phase 5 — Confirmation

Present final summary:
```
MCP-ADD COMPLETE
================
Retrofits applied:   <n>
New tools:           <name, name, ...>
New personas:        <name, name, ...>
New workflows:       <name, name, ...>
Tests added:         <n> unit, <n> integration
Pattern conformance: <was X/10, now Y/10>

Files changed:
  <list>

Next steps:
  - Review the diff before pushing.
  - Run /deliver-feature or your CI to confirm the full pipeline is green.
  - Consider whether any of the new capabilities warrant a new ADR
    (delegate to /adr if so).
```

## Output Format

### Files Produced (retrofit or extend or both)
- `mcp-add-plan.md` at the target MCP's repo root — the approval-gated plan.
- New/modified source files per the plan.
- New test files per the testing plan.
- Optional `docs/adrs/ADR-<n>-<title>.md` if the change is architecturally material.

### Final Confirmation Message
Exactly as shown in Phase 5 above.

## Guardrails

- **No source code touched before Phase 3 approval.** The plan file is the single approval surface. Any edit to the plan resets the gate (per `shared/rules/approval-gates.md` gate #6 — writing files out of boundary).
- **Never remove or rename an existing tool without explicit approval.** External MCP clients may depend on the tool's name and schema. Retrofits preserve the tool's external contract by default.
- **Never modify the entrypoint (main.go / index.ts) except to add wiring for new tools.** Business logic never lands there. If the retrofit list would put logic in the entrypoint, that's a bug — flag and re-plan.
- **Trinity only.** New domain-layer types are `Tool`, `Persona`, or `Workflow`. Reject any Phase 1 answer that wants a "Handler", "Manager", "Service", or "Controller" at the domain layer — those live in the adapter layer or nowhere.
- **Every new tool has a JSON schema for inputs AND outputs.** No exceptions. Schemas live with the tool, not in a separate schemas/ dir at the top level (per the blueprint's directory tree).
- **Every new tool has unit tests + integration tests.** Coverage < 85% for the new code is a Phase 4 failure — send the developer back.
- **External SDKs stay in adapters.** If the retrofit finds an SDK call inside a tool, that's an `Introduce Adapter` refactoring, not an inline fix.
- **OTel spans emit only from adapters and interceptors.** Never from tool implementations directly (per `architecture-guardrails.md` #8).
- **No hardcoded secrets, ever.** New adapters read from env or the existing config layer (per `architecture-guardrails.md` #3).
- **Language boundaries.** Only Go and TypeScript are supported. If Phase 1 detects a different language, stop and tell the user.
- **Never invoke a downstream agent (`/developer`, `/code-reviewer`, etc.) that doesn't exist in `shared/agents/`.**

## Standalone Mode

Works entirely offline. The only external dependency is the target MCP server's source tree, which the user provides via path. All pattern references live in `shared/blueprints/mcp-server.md` and `shared/rules/*` — no network calls, no live SDK version lookups.

If the target MCP uses an SDK version significantly newer than what's documented in `shared/blueprints/mcp-server.md`, flag it in Phase 2a as an `[ASSUMPTION]` — the plan is still valid but a human should double-check the SDK's current transport/tool API surface.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
