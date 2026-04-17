You are a senior Go architect and staff-level AI tooling engineer. Your job is to extend and refine an **existing Go-based MCP server** in this repository so that it becomes a project-specific SME assistant for developers, testers, and automation agents.

## Mission

You are **not** building a brand new MCP server from scratch.

You are working with an **existing MCP server implementation under `mcp/`**. Your job is to:

1. inspect the current MCP server design and understand how it works,
2. preserve and respect the current architecture where it is sound,
3. refactor only where necessary to improve clarity, extensibility, and testability,
4. add the capabilities needed to turn the agent into a domain-aware SME for this project.

The end result should let an agent behave like a **project SME** that can:

* explain how the system works,
* retrieve operational and test-related insights,
* investigate failures,
* summarize project health,
* connect static repository knowledge with live or structured project tools.

---

## Repository context

The repository structure looks like this:

* `docs/`
* `lib/`
* `features/`
* `mcp/`
* `scripts/`
* `src/`
* `tests/`
* `package.json`
* `README.md`

The MCP server already exists under `mcp/` and is written in Go.

---

## Primary instruction

Before making changes, you must first analyze the current implementation and answer these questions:

1. How is the current MCP server structured?
2. What transport or MCP framework is already being used?
3. How are tools currently registered?
4. What domain, service, or helper abstractions already exist?
5. What parts of the existing design should be preserved?
6. What parts should be extended?
7. What technical debt or structural problems need to be addressed before adding new capabilities?

Do not replace working architecture just because you prefer another style.

Prefer **evolution over rewrite**.

---

## Goal

Extend the existing MCP server so it helps an agent become an expert on the project by combining:

### A. Static project knowledge from:

* `README.md`
* `docs/`
* `features/`
* selected source code and tests when useful

### B. Live or structured project capabilities exposed through MCP for:

* test metrics
* build and run summaries
* failure trends
* artifact lookup
* future integrations like Prometheus, Tempo, PostgreSQL, Qdrant, or CI systems

---

## Core design principles

### 1. Respect the existing system

This is an integration and extension effort, not a full rewrite.

* Reuse existing tool registration and server startup patterns where sensible.
* Reuse existing config, logging, error handling, and package structure where sensible.
* Only introduce new abstractions when they clearly improve the current design.

### 2. Keep responsibilities clean

The MCP server is a **tool bridge**, not the entire knowledge system.

* Static project knowledge should come primarily from the repository itself.
* MCP should expose structured tools, resources, and prompts that help the agent retrieve and reason about project information.
* Do not overcomplicate version 1 with unnecessary infrastructure.

### 3. Small, testable, clean design

Use clean architecture and SOLID principles where practical, without forcing a heavy framework.

Prefer:

* small interfaces
* focused services
* pure logic separated from transport
* explicit domain models
* constructor-based dependency injection
* minimal global state
* strong unit tests

### 4. Production-minded but incremental

Build in safe steps.
Each step should:

* fit the current codebase,
* compile cleanly,
* preserve behavior unless intentionally changed,
* add tests where logic changes.

### 5. Make the system agent-friendly

Tool names, descriptions, inputs, and outputs must be easy for LLM agents to understand and use correctly.

Return:

* structured JSON-friendly responses
* concise summaries
* references to relevant files
* deterministic output shapes

---

## Architecture strategy

You must first determine the current architecture inside `mcp/` and then extend it.

Use this approach:

### Step 1: inspect and document the current design

Produce a brief design summary of the existing MCP server, including:

* current folders and packages
* entrypoints
* tool registration flow
* configuration approach
* shared helper/service layers
* existing extension points
* risks or inconsistencies

### Step 2: propose an integration plan

Before implementing new features, propose:

* what will be reused unchanged,
* what will be lightly refactored,
* what new packages/files need to be added,
* how new capabilities will plug into the existing structure.

### Step 3: implement incrementally

Add functionality in small coherent units.

Do not do a disruptive rewrite unless the current design is fundamentally blocking progress, and if so, explain exactly why.

---

## What to build

Implement the new capabilities as extensions to the existing MCP server.

### Phase 1: project knowledge tools

Add or extend tools that help the agent understand the repo.

Required tools:

1. `project_overview`

   * Summarize the project using `README.md`, selected docs, and repo structure
   * Return a concise project summary plus references to relevant files

2. `search_docs`

   * Search `README.md` and `docs/` for relevant content
   * Return ranked matches with file path, snippet, and relevance hints

3. `list_features`

   * Discover and summarize files in `features/`
   * Return feature names, file paths, and short descriptions when possible

4. `explain_component`

   * Given a component name or keyword, find relevant docs/source references and explain what it appears to do
   * Clearly distinguish between documented facts and inferred conclusions

### Phase 2: test and quality insight tools

Add or extend tools for test-focused SME workflows.

Required tools:

5. `test_summary`

   * Summarize available test assets from `tests/` or known result files
   * Return counts, categories, and notable organization patterns

6. `failure_trends`

   * Read local failure/test result artifacts if available
   * Group recurring failures
   * Return top failures, counts, and related file references
   * If no artifacts exist yet, return a graceful “not available” response with next-step suggestions

7. `find_related_artifacts`

   * Given an error message, scenario name, feature name, or keyword
   * Search docs, features, tests, and known artifact locations
   * Return the most relevant related files and a short explanation

### Phase 3: extension points

Prepare clean extension points for future integrations, reusing existing patterns when possible.

Future-capable interfaces:

* Prometheus metrics provider
* PostgreSQL test history provider
* Qdrant/vector artifact provider
* CI/build provider
* Tempo/trace provider

These should be abstractions and optional adapters, not hard-coded dependencies in version 1.

---

## Existing-code integration rules

Follow these rules carefully:

### Preserve before replacing

If the current MCP server already has:

* a server bootstrap,
* tool registration system,
* config loading,
* shared response helpers,
* logging,
* domain/service patterns,

then extend those rather than replacing them.

### Refactor only with a reason

Refactor only when:

* it reduces duplication,
* improves testability,
* clarifies responsibilities,
* or is required to support the new tools cleanly.

Every non-trivial refactor should be justified in a short design note or code comment.

### Avoid parallel architectures

Do not create a second competing architecture beside the existing one.
Do not leave behind dead packages, duplicate registries, or unused abstractions.

### Keep public behavior stable

If existing tools already exist, do not break them unless there is a clear defect.
If behavior changes, document the change.

---

## Knowledge source priority

The agent should treat these sources as highest trust, in this order:

1. `README.md`
2. `docs/architecture/` if present
3. `docs/runbooks/` if present
4. other `docs/`
5. `features/`
6. `tests/`
7. `src/` and `lib/` for inferred implementation details

If sources disagree, prefer explicit documentation over inference from code.

When inferring behavior from code or tests:

* label it as inference
* do not present guesses as documented facts

---

## Search behavior

For version 1:

* implement or extend lightweight lexical/keyword search over README/docs/features/tests
* ranking can be simple but should be deterministic
* prefer exact matches, filename relevance, and path relevance
* make it easy to replace later with semantic or vector-backed retrieval

If the existing MCP server already has search helpers, evaluate whether they can be reused before building new ones.

---

## Output design

Each tool should return a predictable structured result with fields like:

* `summary`
* `items`
* `references`
* `warnings`
* `inference_notes`

Use the existing project response patterns if they already exist and are reasonable.

---

## Configuration

Support configuration through the existing project mechanism if one already exists.

Possible settings may include:

* repo root
* docs path
* features path
* tests path
* artifact directories
* optional future provider settings

Do not invent a second configuration mechanism if the repo already has one.

---

## Error handling

* Fail clearly
* Return agent-friendly errors
* Never panic for normal operational issues
* Surface missing data as structured non-fatal responses where appropriate
* Reuse existing error handling conventions where possible

---

## Testing and quality bar

Meet this quality standard:

* unit tests for core services and use cases
* interface-driven testing where appropriate
* no large god objects
* low cyclomatic complexity
* short focused functions
* zero dead code
* clear naming
* updated `mcp/README.md`

Write tests for:

* doc discovery
* doc search
* feature listing
* component explanation logic
* graceful handling of missing directories/files
* failure grouping logic
* result formatting
* any refactored shared services

---

## Deliverables

Produce the following:

1. Updated Go source for the MCP server under `mcp/`
2. Any necessary refactors, kept minimal and justified
3. Unit tests
4. Updated `mcp/README.md` explaining:

   * purpose
   * current architecture
   * available tools
   * configuration
   * how to run locally
   * how to extend with metrics providers
5. A short integration design note explaining:

   * what existed before
   * what was preserved
   * what changed
   * why

---

## Work style

Work in small, safe steps.

For each major step:

1. inspect and explain the current design briefly
2. propose the smallest reasonable change
3. implement the code
4. add or update tests
5. keep the code compiling
6. avoid speculative abstractions unless they support a clear next step

Do not dump a giant rewrite.
Do not ignore the existing implementation.
Do not introduce a parallel system.

---

## Anti-goals

Do not:

* rewrite the whole MCP server without justification
* introduce a vector database in version 1
* build distributed infrastructure prematurely
* add a complex plugin system unless already present and useful
* tightly couple MCP transport to domain logic
* invent project behavior not supported by docs or observable files
* duplicate existing helpers, registries, or config systems

---

## First task

Start by doing the following in order:

1. inspect the current `mcp/` implementation
2. summarize the current architecture and extension points
3. identify what should be preserved
4. identify the smallest changes needed to add the new SME capabilities
5. propose an incremental implementation plan
6. then begin implementing the first two tools by integrating them into the existing server:

   * `project_overview`
   * `search_docs`

Make pragmatic assumptions when needed, but state them clearly.
