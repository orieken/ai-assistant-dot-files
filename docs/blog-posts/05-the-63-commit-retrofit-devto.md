---
title: "The 63-Commit Retrofit: Bringing Clean Architecture to a Legacy MCP Server"
published: false
description: "How we transformed an 880-LOC God Object MCP server into a composite architecture with an 89% LOC reduction across 63 atomic commits."
tags: ai, architecture, golang, devtools, mcp
canonical_url:
cover_image:
---

Model Context Protocol (MCP) servers often start out simple.

You register a couple of tools, wire up an inline handler map, parse some arguments, and return text. But as your server grows from 3 tools to 17, that initial simplicity turns into a liability. 

In `saturday-mcp` (the MCP server powering parts of our multi-agent architecture ecosystem), our main HTTP/stdio handler package had grown into an **880-line God Object**. Seventeen distinct tools were registered inline as raw `map[string]interface{}` schemas inside a single `RegisterTools` function. Business logic, external process invocations (`os/exec`), and direct filesystem calls were tightly coupled to the MCP SDK's request-response types.

We didn't want to burn the server down and rewrite it from scratch. We wanted to refactor it safely while preserving full test suite parity.

Here is how we retrofitted `saturday-mcp` across **63 atomic commits**, reducing our core handler from 880 LOC down to 93 LOC—an 89% reduction—while pushing test coverage in our tools package to 86.7%.

---

## The Initial Smell: The 880-LOC God Object

Before refactoring, `internal/server/handler.go` attempted to do everything:

1. Define JSON input schemas inline using nested `map[string]interface{}` builders.
2. Unmarshal parameters directly inside raw handler functions (`handleGenerateSite`, `handleAnalyzePerformance`, etc.).
3. Instantiate infrastructure dependencies directly (like instantiating `observability.NewFileMetricsProvider` inline).
4. Invoke un-timeout'd sub-processes (`os/exec`) directly without domain interface boundaries.
5. Format and return `mcp.NewToolResultText` strings without output schema validation.

Testing any single tool required mocking or exercising the entire server handler context. Adding a 18th tool meant opening `handler.go`, adding another 50-line block, and praying that imports and schema keys didn't clash.

---

## The Target Pattern: The Trinity (Tools, Personas, Workflows)

We established a clean architectural specification for MCP servers:

- **Domain Layer (`internal/domain/`)**: Pure Go interfaces (`Tool`, `Persona`, `Workflow`) with zero SDK imports (`github.com/mark3labs/mcp-go`).
- **Tools (`internal/tools/`)**: Isolated structs implementing `domain.Tool` for single-step capabilities.
- **Workflows (`internal/workflows/`)**: Multi-step orchestrations composed from domain interfaces.
- **Adapters (`internal/adapters/`)**: Concrete implementations for process execution, filesystem access, and OpenTelemetry tracing.
- **Server Entrypoint (`internal/server/`)**: Pure dependency injection and route registration.

---

## The 63-Commit Refactoring Mechanics

Refactoring a live server in one massive pull request is a recipe for silent regressions. Instead, we executed the migration incrementally using strict commit discipline and subagent orchestration:

### 1. Eliminating Dead Code and Establishing Interfaces (Commits 1–5)
We deleted dead entrypoint stubs and established pure domain interfaces for `Tool`, `Persona`, and `Workflow` in `internal/domain/`:

```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() ToolInputSchema
    OutputSchema() JSONSchema
    Execute(ctx context.Context, args map[string]interface{}) (Result, error)
}
```

### 2. Extracting Tool Classes One by One (Commits 6–20)
For each of the 15 tools, we extracted a dedicated struct in `internal/tools/` (e.g., `generate_site_tool.go`, `analyze_performance_tool.go`). Each tool received its required domain interfaces via constructor injection rather than pulling them from global server state.

### 3. Isolating Adapters and Process Execution (Commits 21–35)
Direct sub-process execution (`exec.Command(...)`) was wrapped behind a domain `TestRunner` interface in `internal/adapters/testrunner/` with explicit `context.Context` timeouts. Direct `os.WriteFile` calls were routed through a domain `FileSystem` adapter in `internal/adapters/filesystem/`.

### 4. Schema-First Outputs and OTel Tracing (Commits 36–50)
Every tool gained a strongly typed `OutputSchema()` method. We introduced OpenTelemetry span wrapping via a functional option middleware (`WithTracer(...)`), recording tool execution latency, success status, and error classes without polluting business logic.

### 5. Test Backfill and Parity Verification (Commits 51–63)
We backfilled unit tests for every newly extracted tool and adapter package.

---

## The Results: Real Numbers

| Metric | Before Retrofit | After Retrofit | Change |
|---|---|---|---|
| `internal/server/handler.go` LOC | 880 lines | 93 lines | **-89% reduction** |
| `internal/tools/` Test Coverage | 0% (inlined) | **86.7%** | +86.7% |
| `internal/adapters/filesystem/` | 0% | **100.0%** | +100% |
| `internal/adapters/testrunner/` | 86.1% | **97.2%** | +11.1% |
| OpenTelemetry Instrumentation | None | 94.9% adapter coverage | Fully Instrumented |

Adding a new tool to `saturday-mcp` is now remarkably simple: drop a file into `internal/tools/`, implement `domain.Tool`, and append it to the tool provider slice. The server package doesn't even need to be touched.

---

## What We Learned

1. **Clean Architecture applies to MCP servers just like web APIs.** Treating MCP tools as thin adapters around domain use-cases prevents SDK lock-in and makes testing trivial.
2. **Atomic commits preserve velocity.** By breaking the refactoring into 63 distinct commits, every single step passed our e2e regression suite before moving to the next.
3. **Structured output schemas are non-negotiable.** Moving from unstructured text responses to typed JSON output schemas makes LLM tool consumers vastly more reliable.

---

## Image Prompt

> Hero image prompt: A sleek, dark-mode conceptual illustration showing a tangled monolithic server circuit board splitting into 17 organized, color-coded, modular micro-chips labeled with Clean Architecture domain boundaries. High contrast, neon teal and violet accents, modern tech aesthetic.
