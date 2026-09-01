# Architectural Audit: Loom Framework

## 1. Executive Verdict

The `loom` framework, despite its stated mission to elevate engineering teams to Level 2 and Level 3 autonomous multi-agent orchestration, is fundamentally trapped in a Level 1 architecture. At its core, the Go codebase (`cmd/loom`) acts as a sophisticated dotfiles manager and prompt synchronizer, pushing markdown templates into various AI IDEs (Claude Code, Cursor, Copilot). It relies on the host platform's LLM context to execute pipelines using a "sequential-simulation" pattern. Because it lacks a native execution kernel, it delegates true orchestration and tool execution entirely to the IDE, resulting in a bloated, monolithic execution model masquerading as a decoupled pipeline.

The biggest architectural sin is the "God-Framework" trap. `loom` tightly couples orchestration routing (`pipeline-schema.md`) and CI/CD quality gates (auditing agents) directly into synchronous LLM prompt execution. By forcing opinions on DevOps and infrastructure inside a simulated LLM loop, it bloats the context window, guarantees hallucination feedback loops at scale, and reinvents asynchronous CI/CD gating inside a synchronous chat session. 

---

## 2. High-Severity Vulnerabilities

*   **Synchronous Auditing Coupling (Infinite Loop Risk):** In `shared/orchestration/pipeline-schema.md`, auditing agents (e.g., `context-auditor`) are executed synchronously within the same LLM loop using `audit: { onFail: retry, maxRetries: 3 }`. Because the LLM frequently fakes concurrency (`sequential-simulation`), failures repeatedly stuff the same context window with retry prompts, heavily compounding hallucination risks and rapidly degrading context fidelity.
*   **Bloated State Payload:** State management relies entirely on Markdown artifacts written to `.claude/feature-workspace/` and pushed back into the context. There is no strict, typed state graph payload (like LangGraph). Agents must parse raw, sprawling markdown files for every transition, which is highly inefficient and prone to information loss.
*   **Host-Bound Execution Engine:** `loom` relies entirely on proprietary platform behaviors (Claude Code subagents, Roo Code modes) to execute tool calls. While a reference MCP server exists (`shared/mcp/`), the framework does not own its core execution loop. It is tightly bound to whatever tool-calling capabilities the host IDE happens to support.
*   **Fake Parallelism & Context Bleeding:** The pipeline relies on `parallelStrategy: sequential-simulation` by default. Stages marked as parallel are fed sequentially into the same LLM session without actual process isolation, guaranteeing that context bleeds between supposedly independent agents.

---

## 3. Level 2 & Level 3 Gap Analysis

### Level 2 Readiness (Tool Utilization & Orchestration)
To achieve Level 2 maturity, a framework requires modular, provider-agnostic tool execution and strict state management. `loom` currently delegates tool execution to the host platform (e.g., Cursor's limited tool execution vs Claude's native subagents). State is loosely managed via markdown file reads rather than a deterministic state graph, preventing reliable, reproducible agent workflows in headless enterprise environments.

### Level 3 Readiness (Autonomous Planning & Collaboration)
Level 3 requires dynamic routing, where agents can evaluate their own confidence and decide the optimal next node. `loom` uses statically declared YAML pipelines with basic boolean string evaluation (`condition: "feature.hasUI == true"`). There is no mechanism for an agent to dynamically spawn a required sub-workflow or correct its routing path mid-execution outside the hardcoded DAG.

### Memory Persistence
Memory persistence is heavily tangled with prompt generation. The framework's memory system relies on lexical file system reads (Knowledge Items in `shared/knowledge/`), loading markdown directly into the LLM context. There is no decoupled vector or graph database to handle long-term semantic memory retrieval efficiently outside the prompt window.

---

## 4. Refactoring Roadmap

1.  **Decouple the Execution Kernel:** Strip the `sequential-simulation` LLM loop from the framework. Build a native, decoupled execution engine (using a standard state graph architecture) that natively orchestrates agents as isolated processes rather than simulating them inside one massive prompt session.
2.  **Move Audits to Async CI/CD Gates:** Remove synchronous audit loops (`audit: { agent: context-auditor }`) from the runtime execution DAG. Quality gates and security reviews must be pushed to asynchronous CI/CD pipelines (GitHub Actions, GitLab CI) to prevent the LLM from getting trapped in retry loops and burning tokens.
3.  **Implement Dynamic State Graphs:** Replace the artifact-based state passing (`produces: analysis.md`) with a strict, typed internal state payload that is selectively passed between agents. This limits context exposure to only what the receiving agent explicitly needs.
4.  **Abstract Tool Calling via MCP Exclusively:** Transition entirely to the Model Context Protocol (MCP) for *all* tool utilization. Stop generating platform-specific IDE rules (Roo Code modes, Copilot instructions) and make the framework fully agnostic to the underlying LLM provider, achieving true Level 2 modularity.
