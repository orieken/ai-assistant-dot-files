> **SUPERSEDED 2026-08-29** — merged into [`BUILD-ROADMAP.md`](BUILD-ROADMAP.md), the single
> authoritative build plan. Kept for provenance; see BUILD-ROADMAP Appendix A for where each item
> landed and Appendix B for corrections. Do not plan work from this file.

# Agentic Maturity Roadmap: `loom` Framework
## Phase 1: Shoring up Level 2 (Coordinated Multi-Agent Systems)

- [ ] **Decouple Tool Execution & Validation from Host IDEs**
  - **The Problem:** `loom` lacks a native tool execution engine. It generates prompt configurations and relies entirely on host platforms (Claude Code, Cursor) to execute subagents and tools, tying framework resilience and validation to proprietary LLM implementations.
  - **The Architectural Fix:** Implement a native tool orchestration kernel using standard tool-calling APIs (e.g., an MCP client runtime) wrapped in a resilience layer. Force typed Pydantic/Zod schemas and handle execution errors independently of the LLM prompt loop.
  - **Target Files:** `cmd/loom/internal/platform/*.go`, `shared/mcp/internal/server/handler.go`

- [ ] **Migrate from Markdown Artifacts to Typed State Graphs**
  - **The Problem:** Agent state is passed as bloated, unstructured markdown files (e.g., `produces: analysis.md`) tracked in `.claude/feature-workspace/pipeline-state.json`. This forces agents to re-parse massive text files on every transition, causing context decay.
  - **The Architectural Fix:** Migrate to a strict, typed internal state graph (e.g., LangGraph or a custom Finite State Machine). Agents must read/write narrowly scoped JSON state payloads instead of raw markdown files.
  - **Target Files:** `shared/orchestration/pipeline-schema.md`, `shared/skills/deliver-feature/SKILL.md`

- [ ] **Implement Native OS-Level Human-in-the-Loop (HITL) Interrupts**
  - **The Problem:** Approval gates rely entirely on brittle, prompt-level natural language instructions ("user must say 'ship'"). LLMs can easily hallucinate past these prompt-level guards to execute irreversible actions.
  - **The Architectural Fix:** Implement true execution interrupts in the native kernel that physically halt the runtime process. Yield state back to the user via a CLI prompt or webhook before executing high-risk tool calls.
  - **Target Files:** `shared/rules/approval-gates.md`, `shared/orchestration/policy-evaluator.md`

## Phase 2: Achieving Level 3 (Autonomous Orchestration Layer)

- [ ] **Replace Hardcoded DAGs with Dynamic Routing**
  - **The Problem:** Workflows use a static, hardcoded YAML DAG with primitive boolean evaluation (`condition: "feature.hasUI == true"`). The framework cannot dynamically adapt to novel intents or re-route course mid-execution.
  - **The Architectural Fix:** Build a native "Router" or "Planner" agent node equipped with tool access. This node must dynamically evaluate the state graph and emit the next required sub-workflow or agent ID, fully decoupling execution from static YAML.
  - **Target Files:** `shared/orchestration/pipeline-schema.md`, `shared/skills/orchestrate/SKILL.md`

- [ ] **Abstract Semantic and Episodic Memory**
  - **The Problem:** Memory relies on lexical file-system reads of Markdown `Knowledge Items` loaded fully into the prompt context via `search_ki`. There is no abstract semantic or episodic memory layer to manage historical context safely.
  - **The Architectural Fix:** Integrate a native vector database or graph store (e.g., sqlite-vss) to abstract episodic memory (pipeline traces) and semantic memory. Enable agents to retrieve only highly relevant vectors without dumping entire files into the prompt.
  - **Target Files:** `shared/mcp/internal/tools/search_ki.go`, `shared/memory-registry.json`

- [ ] **Decouple Governance into Async CI/CD Gates**
  - **The Problem:** Auditing agents (e.g., `context-auditor`) run synchronously inside the runtime pipeline. This burns tokens, delays execution, and bloats the synchronous LLM prompt loop.
  - **The Architectural Fix:** Move all auditing and governance checks into an asynchronous CI/CD gating pattern (e.g., GitHub Actions hooks). Evaluate pipeline artifacts out-of-band and block deployment without tying up the runtime LLM loop.
  - **Target Files:** `shared/orchestration/pipeline-schema.md`, `shared/telemetry/event-recorder.md`

## Phase 3: Scaling to Level 4 (Self-Learning Agentic Ecosystems)

- [ ] **Implement Isolated "Reflexion" Error Recovery Cycles**
  - **The Problem:** The pipeline uses a primitive `maxRetries: 3` counter that merely repeats the failed tool call in the same context window, practically guaranteeing a hallucination spiral as the context degrades.
  - **The Architectural Fix:** Implement a native "Reflexion" (self-critique) cycle: when a tool fails, spawn a fresh, isolated agent session tasked exclusively with analyzing the error and updating the state graph. This completely prevents context contamination.
  - **Target Files:** `shared/orchestration/pipeline-schema.md` (audit behavior logic), `cmd/loom/internal/platform/`

- [ ] **Automate Prompt & Tool Meta-Evolution**
  - **The Problem:** Optimization of prompts and rules requires manual human promotion via the `learning-engine` or `promote-memory` skills based on past retrospectives.
  - **The Architectural Fix:** Implement a meta-learning agent that operates in the background. It must automatically evaluate `events.jsonl` telemetry and rewrite agent `.md` prompt templates or tool schemas programmatically based on historic failure rates and latency.
  - **Target Files:** `shared/skills/learning-engine/SKILL.md`, `shared/telemetry/events.jsonl`

- [ ] **Build a Native MCP Client Runtime**
  - **The Problem:** While `loom` provides a scaffold to *be* an MCP server, it lacks the native capability to act as an MCP *client* to discover and utilize external ecosystem agents dynamically.
  - **The Architectural Fix:** Expand the core framework to include a native MCP Client runtime layer. Allow the Router agent to dynamically discover, handshake with, and delegate sub-tasks to third-party MCP servers, rather than relying exclusively on the host IDE.
  - **Target Files:** `shared/mcp/README.md`, `cmd/loom/cmd/`
