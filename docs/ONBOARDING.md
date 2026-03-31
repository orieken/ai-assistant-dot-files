# AI Platform Onboarding

Welcome! This repository leverages an advanced AI Feature Team methodology. Your AI assistants (Claude, Cursor, Antigravity, Copilot) are globally pre-configured to adhere to our Clean Architecture patterns and to act as specialized agents in your workflow.

## The Horizon Swarm (Agents)

We have dedicated autonomous **Agents** (defined under `.claude/agents/`) designed to help execute multi-step engineering workflows.

### 1. Test-Driven Developer
- **Agent File:** `test-driven-development-agent.md`
- **What it does:** Reverses the typical AI code-generation cycle. When given a feature and acceptance criteria, it autonomously writes the tests *first*, executes them (red), and then iteratively implements the production code until the entire suite passes (green).
- **How to use:** Invoke the agent by name and supply a new feature request with unambiguous acceptance criteria. Do not interrupt its red-green-refactor loop.

### 2. Modernization Swarm Supervisor
- **Agent File:** `modernization-swarm.md`
- **What it does:** Acts as a coordinator for mass codebase updates. It oversees pseudo-threads that asynchronously update outdated dependencies, refactor deprecated design patterns, and close gaps in code coverage.
- **How to use:** Invoke the supervisor and instruct it to begin a modernization sweep on a specific domain (or globally). Expect it to provide branch diffs and conflict flags.

### 3. Documentation Manager
- **Agent File:** `documentation-manager.md`
- **What it does:** Closes the knowledge gap. It parses your recent terminal sessions, debugging logs, or feature conversations and safely commits those learnings into persistent documents like `ARCHITECTURE.md`, `RUNBOOKS.md`, and `ONBOARDING.md`.
- **How to use:** Simply instruct it: "@[.claude/agents/documentation-manager.md] Update docs from this session".

---

## Daily AI Workflows (Skills)

For single-shot, immediate actions, we have predefined **Skills** located in `.claude/skills/`. Assistants will often auto-detect these based on intent, but they can be manually triggered.

### `list-agents`
Quickly query the repository for all available custom AI personas.
- **Trigger Keywords:** `"list agents"`, `/list-agents`, `"what agents do we have?"`

### `scaffold-docs`
Instantly bootstrap a markdown implementation guide prior to building a large feature.
- **Trigger Keywords:** `"scaffold docs"`, `/scaffold-docs`, `"create implementation guide"`
