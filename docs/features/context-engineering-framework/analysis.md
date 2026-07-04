# Feature Analysis: Context Engineering Framework — Cross-Platform Parity

## Summary

Transform the ai-assistant-dot-files repository from a Claude-first agent configuration repo into a true cross-platform context engineering framework. One clone, one setup script, and any machine running Claude Code, Cursor, Gemini (Antigravity), GitHub Copilot, Windsurf, or OpenAI Codex gets the full agent team, skills, rules, and guardrails — with feedback loops, agent evaluation, and observability that currently don't exist.

---

## Acceptance Criteria

### AC-1: One-command setup on any machine
**Given** a fresh machine with at least one supported AI tool installed
**When** the user runs `./install.sh`
**Then** all platform-specific config files are symlinked or copied to their correct locations, and the framework is immediately active in every installed tool.

### AC-2: Cross-platform parity
**Given** the shared rules, agents, and skills in the canonical directory
**When** any platform-specific config is generated
**Then** it includes the same guardrails, design principles, and approval gates — adapted to that platform's syntax but identical in substance.

### AC-3: Agent feedback loop
**Given** a completed pipeline run with persisted artifacts
**When** a new pipeline run starts on a related feature
**Then** agents can access relevant past delivery summaries and learn from prior code-review rejections or security findings.

### AC-4: Agent evaluation
**Given** a set of golden-file test cases per agent
**When** an agent's prompt is modified
**Then** the golden-file tests verify the agent still produces expected outputs for known inputs.

### AC-5: Agent observability
**Given** a pipeline run in progress
**When** each agent completes its phase
**Then** duration, iteration count, and pass/fail are logged to a structured trace file.

### AC-6: Context budget enforcement
**Given** an agent about to be invoked
**When** the context-engineer prepares its manifest
**Then** the manifest includes a token budget estimate and flags if the combined context exceeds a configurable threshold.

---

## Non-Functional Requirements

- **Portability**: Must work on macOS, Linux, and WSL. No platform-specific dependencies beyond the AI tools themselves.
- **Idempotency**: Running `./install.sh` multiple times must be safe — no duplicate symlinks, no overwritten user customizations.
- **Performance**: Setup script completes in < 30 seconds.
- **Maintainability**: Adding a new platform requires adding one directory and one entry in the platform registry — not modifying every existing file.

---

## Out of Scope

- Building a custom RAG retrieval server (use file-based Knowledge Items for now)
- Real-time agent-to-agent communication protocol (current file-based handoff is sufficient)
- SaaS dashboard for agent metrics (file-based traces are sufficient for v1)

---

## Technical Breakdown

### Bounded Context
- **Owning context**: Agent Orchestration
- **Context crossings**: Craftsmanship Governance (rules), Feature Delivery (pipeline), Documentation Knowledge Base (feedback loop)

### Affected Components

| Area | Current State | Gap |
|---|---|---|
| `.claude/` | 20+ agents, 35+ skills, 3 rule files | Missing: context-engineer in pipeline, agent contracts, versioning |
| `.cursor/` | 1 global rule file | Missing: agents, skills, rules parity |
| `.github/copilot-instructions.md` | Single system prompt | Missing: agent awareness, skill triggers |
| `.gemini/` | Single system prompt | Missing: agent awareness, skill triggers |
| `.openai.md` | Single system prompt | Missing: agent awareness, skill triggers |
| `.windsurfrules` | Copy of `.cursorrules` | Missing: agent awareness, skill triggers |
| `install/` | Empty directory | Missing: setup scripts |
| `uninstall/` | Listed but not found | Missing: teardown scripts |
| `templates/` | Claude-only scaffold | Missing: cross-platform templates |
| `scaffold-team.sh` | Claude-only | Missing: multi-platform support |
| `docs/` | Runbooks, ADRs, features | Missing: agent changelog, evaluation results |

---

## Task List

### Epic 1: Canonical shared layer (platform-agnostic)

> **Why**: Right now, rules, design principles, and the domain dictionary live at the repo root but each platform config re-implements them with copy-paste drift. Extract the canonical source of truth so all platforms reference the same content.

- [ ] **1.1** Create `shared/` directory structure:
  ```
  shared/
    rules/
      architecture-guardrails.md
      design-principles.md
      approval-gates.md
    agents/          ← canonical agent definitions
    skills/          ← canonical skill definitions
    domain-dictionary.md
    architecture-rules.md
  ```
- [ ] **1.2** Move `.claude/rules/`, `.claude/agents/`, `.claude/skills/` content into `shared/` as the canonical source
- [ ] **1.3** Make `.claude/rules/`, `.claude/agents/`, `.claude/skills/` symlink to (or copy from) `shared/`
- [ ] **1.4** Create `shared/platform-registry.json` — lists every supported platform with:
  - Config file path pattern (e.g., `.cursor/rules/global.mdc`)
  - Instruction format (markdown, YAML frontmatter, plain text)
  - Capability tier (full agents, personas only, system prompt only)
  - Agent/skill support (boolean)

### Epic 2: Cross-platform config generation

> **Why**: Cursor, Copilot, Gemini, Windsurf, and OpenAI configs are hand-maintained copies with drift. Generate them from the canonical shared layer.

- [ ] **2.1** Create `scripts/generate-configs.sh` that reads `shared/` and `platform-registry.json` and generates platform-specific configs:
  - `.claude/` — full agents + skills + rules (symlinks)
  - `.cursor/rules/global.mdc` — rules + agent awareness section
  - `.cursorrules` — same as Cursor rules (Cursor reads both)
  - `.windsurfrules` — same as Cursor rules (Windsurf reads both)
  - `.github/copilot-instructions.md` — rules + agent roster reference
  - `.gemini/antigravity/instructions.md` — rules + agent roster reference
  - `.openai.md` — rules + agent roster reference
- [ ] **2.2** For platforms that support agents/personas natively (Cursor rules, Gemini personas), generate platform-native agent configs from `shared/agents/`
- [ ] **2.3** Create a parity test: `scripts/check-parity.sh` (you have this file already — extend it to diff generated configs against shared rules and fail on drift)
- [ ] **2.4** Add a CI fitness function: "no platform config may reference a rule that doesn't exist in `shared/rules/`"

### Epic 3: Universal install/uninstall

> **Why**: The user wants to clone this repo on any machine and run one command. `scaffold-team.sh` only handles Claude.

- [ ] **3.1** Create `install.sh` at repo root:
  - Detects installed AI tools (check for `.cursor/`, `~/.config/github-copilot/`, Gemini CLI, etc.)
  - Runs `scripts/generate-configs.sh`
  - Symlinks global configs to `~/.claude/CLAUDE.md`, `~/.cursor/rules/`, etc.
  - Offers project-level install (copies to a target project) vs. global install (symlinks to home)
  - Prints a verification summary (like `scaffold-team.sh` already does)
- [ ] **3.2** Create `uninstall.sh` — removes symlinks, restores backups
- [ ] **3.3** Make `install.sh` idempotent — safe to re-run without duplication
- [ ] **3.4** Add `--dry-run` flag to show what would be installed without doing it
- [ ] **3.5** Extend `scaffold-team.sh` to call `install.sh` with `--project` mode (backwards compatible)

### Epic 4: Wire context-engineer into the pipeline

> **Why**: The context-engineer agent exists but isn't in the `/deliver-feature` pipeline. It should be Phase 0.5 — after setup, before analyst — to prune and focus each agent's context window.

- [ ] **4.1** Update `deliver-feature/SKILL.md` to invoke context-engineer between Phase 0 (setup) and Phase 1 (analyst)
- [ ] **4.2** Context-engineer produces `context-manifest.md` in `.claude/feature-workspace/`
- [ ] **4.3** Each downstream agent reads `context-manifest.md` to load only pinpointed files (not full directory scans)
- [ ] **4.4** Add context budget estimation to the manifest: estimated tokens per agent, flag if > 80% of context window

### Epic 5: Inter-agent contract schemas

> **Why**: Agents pass artifacts via markdown files but there's no enforced structure. If the analyst forgets "Architectural Flags," the architect may be invoked incorrectly.

- [ ] **5.1** Create `shared/contracts/` directory with one schema per artifact:
  - `analysis-contract.md` — required sections for `analysis.md`
  - `architecture-contract.md` — required sections for `architecture-notes.md`
  - `implementation-contract.md` — required sections for `implementation-notes.md`
  - `review-contract.md` — required sections for `code-review-report.md`
- [ ] **5.2** Each contract specifies: required headings, required fields, conditional fields (e.g., "Architectural Flags" must be present and must be "None" or a list)
- [ ] **5.3** Create `shared/skills/validate-artifact/SKILL.md` — reads a contract and an artifact, fails if required sections are missing
- [ ] **5.4** Wire `validate-artifact` into `deliver-feature` between each agent handoff

### Epic 6: Agent evaluation (golden-file tests)

> **Why**: You test the code agents produce but not the agents themselves. Prompt edits can silently regress agent quality.

- [ ] **6.1** Create `tests/agents/` directory
- [ ] **6.2** For each critical agent, create a golden-file test:
  - `tests/agents/security-reviewer/` — known-vulnerable code snippet + expected findings
  - `tests/agents/code-reviewer/` — code with known smells + expected flagged smells
  - `tests/agents/analyst/` — feature spec + expected analysis sections
- [ ] **6.3** Create `scripts/test-agents.sh` — runs each golden-file test, compares output against expected findings (fuzzy match on key phrases, not exact text)
- [ ] **6.4** Document: "when you edit an agent prompt, run `./scripts/test-agents.sh` to verify no regressions"

### Epic 7: Agent observability & feedback loop

> **Why**: You mandate OpenTelemetry for code but have no tracing on the agent pipeline itself. And agents don't learn from past runs.

- [ ] **7.1** Create `shared/skills/pipeline-trace/SKILL.md` — logs each agent invocation with:
  - Agent name, start time, end time, duration
  - Pass/fail/skipped status
  - Iteration count (for code-review loops)
  - Token estimate (from context manifest)
- [ ] **7.2** Persist trace to `docs/features/<name>/pipeline-trace.json`
- [ ] **7.3** Create `shared/skills/pipeline-retrospective/SKILL.md` — reads trace files from past N deliveries and surfaces:
  - Which agents take longest
  - Which agents cause the most review loops
  - Common security findings that could become rules
- [ ] **7.4** Update analyst agent to read the 3 most recent delivery summaries from `docs/features/` when starting a new analysis (feedback loop)
- [ ] **7.5** Update `deliver-feature` to invoke `/retrospective` automatically after the 5th delivery (not every time — noise reduction)

### Epic 8: Agent versioning & changelog

> **Why**: Agent prompt changes directly affect output quality but are only tracked via git commits with no structured changelog.

- [ ] **8.1** Add `version:` field to every agent's frontmatter (semver: `1.0.0`)
- [ ] **8.2** Create `shared/agents/CHANGELOG.md` — one entry per agent version bump:
  ```
  ## security-reviewer v1.1.0 (2026-07-01)
  - Added OWASP API Security Top 10 checks
  - Reason: missed API-specific vulnerabilities in go-sunday review
  ```
- [ ] **8.3** Create a pre-commit hook (or CI check): if any file in `shared/agents/` changes, `version:` must be incremented and `CHANGELOG.md` must have a new entry
- [ ] **8.4** Delivery summaries include agent versions used

### Epic 9: Domain dictionary — persona vs. agent formalization

> **Why**: The MCP uses "persona" while your framework uses "agent." The domain dictionary should formalize both terms.

- [ ] **9.1** Update `DOMAIN_DICTIONARY.md` to add:
  - **Persona**: A context frame that shapes model identity, tone, and focus areas. A persona has no tool access and no autonomous workflow. Synonyms to avoid: "profile", "role" (too generic).
  - **Agent**: A persona + tools + autonomous process. An agent can take actions, produce artifacts, and participate in pipelines. A persona is a component of an agent.
- [ ] **9.2** Update cross-platform configs to use the correct term per platform capability:
  - Claude: "agents" (full tool access + process)
  - Cursor: "personas" in rules (no native agent orchestration)
  - Gemini: "personas" (context-shaping only)
  - Copilot: "personas" (context-shaping only)
- [ ] **9.3** Update the presentation materials to include this distinction

### Epic 10: Health check & self-test

> **Why**: After install, the user needs to know it worked. The existing `health-check` skill is a stub.

- [ ] **10.1** Implement the `health-check` skill to verify:
  - All symlinks resolve
  - All agents have valid frontmatter (name, description, tools, model)
  - All skills have valid SKILL.md with triggers
  - All rules are present in `shared/`
  - All platform configs were generated from current `shared/` (no drift)
  - Domain dictionary has no orphaned terms
- [ ] **10.2** `install.sh` runs `health-check` at the end automatically
- [ ] **10.3** Add `--verbose` flag for detailed diagnostics

---

## Edge Cases and Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Symlinks don't work on Windows (non-WSL) | Windows users can't use global install | Provide a `--copy` mode that copies instead of symlinks |
| Platform config format changes | Generated configs break after a tool update | Pin known config format versions in `platform-registry.json`, test in CI |
| Agent golden-file tests are brittle | Exact text matching breaks with minor prompt changes | Use fuzzy matching on key findings, not exact output comparison |
| Feedback loop context grows unbounded | Past delivery summaries flood the context window | Cap at 3 most recent deliveries, summarize older ones |
| `shared/` rename breaks existing scaffold-team.sh users | Backwards compatibility | Keep `scaffold-team.sh` working via fallback paths for one version |

---

## Proposed Fitness Functions

| Property | Verification | Owner |
|---|---|---|
| No platform config drift | `scripts/check-parity.sh` diffs generated vs. actual | CI |
| All agents have valid frontmatter | `health-check` skill validates YAML frontmatter | CI |
| Agent prompts don't regress | `scripts/test-agents.sh` runs golden-file tests | CI |
| No orphaned domain terms | `health-check` cross-refs DOMAIN_DICTIONARY.md against code | CI |
| Inter-agent contracts satisfied | `validate-artifact` skill checks required sections | Pipeline |
| Agent version bumped on change | Pre-commit hook checks `version:` field in changed agent files | Pre-commit |

---

## Definition of Done

- [ ] `./install.sh` works on a fresh macOS and Linux machine
- [ ] All 6 platforms generate configs from `shared/`
- [ ] `check-parity.sh` passes — zero drift between shared and generated
- [ ] Context-engineer is wired into the deliver-feature pipeline
- [ ] At least 3 agents have golden-file tests
- [ ] Pipeline trace is produced for every delivery
- [ ] Agent changelog exists with at least one entry per agent
- [ ] Domain dictionary includes persona vs. agent distinction
- [ ] Health check passes after fresh install
- [ ] README documents the full setup flow
