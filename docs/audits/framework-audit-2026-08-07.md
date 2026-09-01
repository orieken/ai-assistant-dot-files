# Comprehensive Framework Audit & Strategic Roadmap — 2026-08-07

**Framework Version**: v3.3.13  
**Audited Repository**: `ai-assistant-dot-files` (Context Engineering Framework by Oscar Rieken)  
**Audit Date**: August 7, 2026  

---

## 1. Executive Summary & Verification Health

`ai-assistant-dot-files` is a single-source-of-truth context engineering and AI multi-agent orchestration framework. Its core design principle—**defining agent prompts, skills, and rules once in `shared/` and generating/symlinking them into 10 target platform environments** (Claude Code, Cursor, Windsurf, Copilot, Antigravity, OpenAI/Codex, Junie, Roo Code, Cline)—provides exceptional multi-tool context synchronization.

### Empirical Verification Results
- **Platform Parity**: 10/10 platform configurations verified in sync via `scripts/check-parity.sh`.
- **Core Roster**: 39 agents (`shared/agents/`) and 69 skills (`shared/skills/`) active and mapped.
- **Agent Golden-File Fixtures**: 32/38 agents fully covered by test input, pattern, and rubric fixtures; 6 specialists documented as deferred in `tests/agents/README.md`.
- **Install Verification Matrix**: 36/36 passed across `--global` and `--project` modes in clean containers (`scripts/test-install.sh`).
- **CI Container Suite (`scripts/ci-check.sh`)**: 4/4 suites passed (0 failed).
- **Bug Fix Applied During Audit**: Resolved a cross-platform compatibility bug in `scripts/health-check.sh` (line 682). On Linux, GNU `stat -f %m` printed `File: "%m"` text rather than failing silently, causing arithmetic comparison errors under `set -u` in `ci-check.sh`. Implemented OS-aware mtime resolution (`stat -f %m` on macOS/Darwin vs `stat -c %Y` on Linux).

---

## 2. Comprehensive Gap Analysis

### A. Operational & Integration Gaps
1. **Unsupported Emerging CLI/IDE AI Assistants**:
   - CLI/IDE tools such as **Aider** (`.aider.conf.yml`) and **Continue.dev** (`.continue/config.json`) do not have native target generators in `scripts/generate-configs.sh`.
2. **Missing MCP Packaging Exporter (Epic 43)**:
   - `shared/mcp/` provides a Go MCP server scaffold, but there is no automated tool to bundle framework skills (e.g. `search-ki`, `check-ubiquitous-language`, `team-topology-check`) into standalone MCP servers for Cursor, Claude Desktop, or VS Code.
3. **Absence of `shared/configs/` Ready-to-Use Linters (Epic 64)**:
   - Language conventions (`<language>-conventions.md`) enforce caps for linters (ESLint, Clippy, RuboCop, Detekt, golangci-lint), but the framework does not ship standard linter configuration presets in `shared/configs/`.

### B. Architectural & Governance Gaps
1. **Singleton Feature Workspace Isolation (Epic 63)**:
   - `.claude/feature-workspace/` is currently a singleton per repository checkout. Concurrent feature delivery runs or parallel agent threads on different git branches collide on state files.
2. **Telemetry Gaps in Gate Decisions & Token Spend (Epic 62)**:
   - `shared/telemetry/event-schema.md` captures tool actions, but human gate decisions (approvals/rejections/edits) are not logged as telemetry events. `pipeline-trace.json` tracks execution duration and iteration count, but lacks empirical token usage / cost metrics.
3. **Headless Prompt-Regression Runner (Epic 61)**:
   - Golden-file agent fixtures in `tests/agents/` are validated against committed outputs, but there is no automated headless runner in CI that executes dynamic prompt evaluations when agent definitions change.
4. **Self-Threat Model Audit (Epic 65)**:
   - Enterprise KI sync (ADR-003) pulls remote memory items into context, and hooks execute on events. While the framework provides a `threat-model` skill for user code, it has not undergone an internal STRIDE threat model audit of its own memory/hook ingestion paths.

---

## 3. High-Value Recommended Features

1. **Native Framework MCP Server (`ai-assistant-mcp`)**: Build a unified Go binary or Node.js server that exposes framework capabilities (e.g., Knowledge Item searching, Clean Architecture validation, Ubiquitous Language checking, and Team Topology auditing) as standard Model Context Protocol (MCP) tools.
2. **Auto-Fix / Auto-Parity Flag (`health-check --fix` & `install.sh --auto-sync`)**: Provide an automatic remediation flag in `health-check.sh` that detects parity drift, stale `CODEMAP.md` files, or un-generated platform configs and fixes them non-interactively.
3. **Pipeline TUI / Visual Delivery Dashboard**: Create an interactive Terminal UI (using Go `bubbletea` or Python `rich`) that visualizes active delivery pipelines (`deliver-feature`, `deliver-bugfix`), showing stage progression, active agent, live duration, estimated token spend, and pending human gate approvals.
4. **Dynamic Provider Fallback & Multi-Model Orchestration**: Extend `shared/model-defaults.yaml` with provider failover policies (e.g., fall back from Anthropic Claude 3.5 Sonnet to Google Gemini 1.5 Pro or OpenAI GPT-4o if rate limits or API outages occur during long pipeline runs).
5. **AOS Policy Engine (`shared/orchestration/policy-evaluator.md`)**: Fulfill AOS Phase 4 by delivering the policy evaluation engine that reads `.claude/policies/*.policy.yaml` and safely auto-approves policy-eligible gates (Tier A: `git-commit`, `out-of-boundary-write`, `fitness-function-wiring`) under strict conditions.

---

## 4. Honest Feedback & Evaluation

### Strengths
* **Unrivaled Multi-Tool Synchronization**: The single-source-of-truth model in `shared/` combined with `check-parity.sh` solves a massive real-world pain point: developer context fragmenting across multiple AI tools.
* **Deep Software Engineering Craftsmanship**: Grounding AI prompt guidelines in industry standards—Clean Architecture, Kent Beck's 4 Rules of Simple Design, Martin Fowler's Refactorings, Sandi Metz's Code Limits, and Skelton & Pais's Team Topologies—ensures high code quality rather than generic AI output.
* **Closed-Loop Memory & Learning (CML)**: The framework learns from its own deliveries. KIs, ADRs, retrospectives, `promote-memory`, and `extract-lessons` form a genuine continuous improvement engine.
* **Governed Human Approval Gates**: 8 explicit approval gates guarantee human control over irreversible actions (commits, migrations, deployments, third-party API mutations).

### Areas of Friction
* **High Cognitive Load**: With 39 agents, 69 skills, and 18 subdirectories, new users face a steep learning curve. Understanding when to use `spec-writer` vs `architect` vs `developer` vs slash commands can be overwhelming.
* **Platform Experience Disparity**: Claude Code & Gemini Antigravity experience rich subagent isolation and slash commands, whereas simpler IDE plugins get concatenated text rules or flat instructions.
* **Singleton State Limitations**: Workspaces like `.claude/feature-workspace/` limit concurrent usage, requiring developers to work sequentially on single features per repository.

---

## 5. Prioritized Action Plan

| Rank | Action Item | Category | Priority |
|---|---|---|---|
| **1** | [DONE] Fix Linux `stat` mtime resolution in `scripts/health-check.sh` | Bug Fix | **P0** |
| **2** | **Epic 61**: Headless prompt-regression eval runner | Evals & QA | **P1** |
| **3** | **Epic 62**: Gate-decision & token spend telemetry | Observability | **P1** |
| **4** | **Epic 43**: Native Framework MCP Server (`shared/mcp/`) | Integration | **P1** |
| **5** | **Epic 63**: Parallel-delivery workspace isolation | Core Engine | **P2** |
