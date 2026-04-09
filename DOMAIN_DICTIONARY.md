# Domain Dictionary

This file is the single source of truth for the **Ubiquitous Language** in the ai-assistant-dot-files ecosystem. All agents, skills, rules, and documentation MUST use exactly these terms. Do not invent synonyms.

**Rule:** When a new concept is introduced (via a new agent, skill, or rule), it MUST be added here before merging.

---

## 1. Core Domains (Bounded Contexts)

- **Agent Orchestration**: Coordination of specialized AI agent personas through a sequential delivery pipeline with human checkpoints.
- **Craftsmanship Governance**: Enforcement of clean code, architecture, and testing standards across all AI assistants and generated code.
- **Test Automation**: Automated verification of software quality through the Saturday Framework (E2E/UI) and Sunday Framework (API).
- **Feature Delivery**: The end-to-end lifecycle of taking a feature spec through analysis, implementation, review, and deployment.
- **Documentation Knowledge Base**: Persistent, RAG-friendly documentation that serves as context for both agents and humans.

---

## 2. Entities (Identifiable Objects)

| Term | Domain | Description | Synonyms to AVOID |
|---|---|---|---|
| **Agent** | `Agent Orchestration` | A specialized AI persona with a defined role, tools, and process. Lives in `.claude/agents/`. | `bot`, `assistant`, `model` |
| **Skill** | `Agent Orchestration` | An executable, on-demand capability triggered by keywords or slash commands. Lives in `.claude/skills/`. | `tool`, `command`, `action`, `plugin` |
| **Rule** | `Craftsmanship Governance` | A governance document that constrains agent behavior. Lives in `.claude/rules/`. | `policy`, `guideline`, `standard` |
| **Feature Spec** | `Feature Delivery` | A structured markdown document describing a work item. Lives in `features/`. | `ticket`, `story`, `issue`, `task` |
| **Pipeline** | `Feature Delivery` | The sequential chain of agents that processes a feature spec into shipped code. | `workflow`, `flow`, `process` |
| **Artifact** | `Feature Delivery` | A markdown document produced by an agent during pipeline execution. Persisted to `docs/features/<name>/`. | `output`, `result`, `report` |
| **Feature Workspace** | `Feature Delivery` | The temporary working directory (`.claude/feature-workspace/`) used during pipeline execution. | `scratch`, `temp`, `staging` |
| **Feature Archive** | `Documentation Knowledge Base` | The permanent directory (`docs/features/<name>/`) where pipeline artifacts are persisted after delivery. | `output folder`, `results` |
| **Blueprint Prompt** | `Test Automation` | A foundational prompt that establishes framework conventions for E2E or API test generation. | `template`, `boilerplate` |
| **Approval Gate** | `Craftsmanship Governance` | A mandatory human checkpoint before an irreversible action. Resets if the pending artifact is edited. | `confirmation`, `sign-off` |
| **Fitness Function** | `Craftsmanship Governance` | An automated, measurable check that verifies an architectural property in CI. | `lint rule`, `check`, `validation` |

---

## 3. Value Objects (Immutable Concepts)

| Term | Domain | Description | Example |
|---|---|---|---|
| **Delivery Summary** | `Feature Delivery` | The final synthesis artifact produced by the orchestrator after all agents complete. | `docs/features/<name>/delivery-summary.md` |
| **Readiness Critique** | `Feature Delivery` | The spec-writer's structured evaluation of whether a feature spec is ready for the pipeline. | Verdict: READY or NEEDS WORK |
| **Human Checkpoint** | `Feature Delivery` | A pause point in the pipeline where explicit user confirmation is required before proceeding. | After analyst, after architect RFC |
| **Cyclomatic Complexity Threshold** | `Craftsmanship Governance` | The maximum allowed cyclomatic complexity per function: 7. | `complexity < 7` |
| **Coverage Threshold** | `Craftsmanship Governance` | The minimum required unit test coverage: 85%. | `coverage >= 85%` |

---

## 4. Domain Events

| Event Name | Domain | Triggered When |
|---|---|---|
| **SpecReady** | `Feature Delivery` | The spec-writer's readiness critique verdict is READY. |
| **AnalysisComplete** | `Feature Delivery` | The analyst produces `analysis.md` and the user confirms scope. |
| **ArchitectureApproved** | `Feature Delivery` | The architect's structural decisions are confirmed by the user. |
| **CodeReviewApproved** | `Feature Delivery` | The code-reviewer's verdict is APPROVED (no more feedback loops). |
| **SecurityCleared** | `Feature Delivery` | The security-reviewer finds no Critical findings, or all Critical findings are resolved. |
| **PipelineComplete** | `Feature Delivery` | All agents have produced their artifacts and the delivery summary is written. |
| **ArtifactsPersisted** | `Documentation Knowledge Base` | All pipeline artifacts are copied from feature workspace to `docs/features/<name>/`. |
| **ShippedToFriday** | `Feature Delivery` | The Cucumber JSON summary is POSTed to the Friday dashboard after user approval. |

---

## 5. Operations / Actions

- **Deliver**: Execute the full agent pipeline for a feature spec. Do not use "build", "run", or "process" when referring to pipeline execution.
- **Scaffold**: Deploy template files into a target project directory. Do not use "copy", "install", or "bootstrap" for this operation.
- **Persist**: Copy pipeline artifacts from the feature workspace to the permanent feature archive in `docs/features/`. Do not use "save", "export", or "archive" as verbs for this operation.
- **Critique**: The spec-writer's evaluation of a feature spec against agent readiness criteria. Do not use "review" (reserved for code-reviewer) or "audit" (reserved for dependency-auditor).
- **Profile**: The performance-engineer's diagnostic analysis of a slow system. Do not use "benchmark" (which implies measurement only, not diagnosis).

---

## 6. Framework Terms

### Saturday Framework (E2E / UI Testing)

| Term | Description | Synonyms to AVOID |
|---|---|---|
| **BaseSite** | Root orchestrator for a web application under test. | `App`, `Application` |
| **BasePage** | Represents a single page within a BaseSite. | `PageObject`, `Screen`, `View` |
| **BaseElement** | A reusable UI component abstraction within a BasePage. | `Component`, `Widget` |
| **BaseFlow** | A multi-step user journey that spans multiple pages. | `Workflow`, `Scenario`, `Journey` |
| **SiteManager** | Manages cross-application test contexts. | `AppManager`, `ContextManager` |
| **TabManager** | Manages multi-tab browser contexts within a test. | `WindowManager`, `BrowserManager` |

### Sunday Framework (API Testing)

| Term | Description | Synonyms to AVOID |
|---|---|---|
| **BaseApiClient** | Abstract base for all domain-specific API clients. | `HttpClient`, `RestClient` |
| **IHttpAdapter** | Interface hiding HTTP implementation details from test logic. | `HttpService`, `RequestHandler` |
| **api fixture** | The custom Playwright fixture providing fluent API testing. | `client`, `request` |

---

*When designing a new agent, skill, or rule: if a new business concept is introduced, add it to this dictionary before merging.*
