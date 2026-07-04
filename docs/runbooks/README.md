# Operational Runbooks

This directory contains runbooks for setup, troubleshooting, agent operations, and pipeline execution. Each runbook is a self-contained guide that can be followed step-by-step.

---

## Categories

### Installation and Setup

Guides for setting up the development environment, installing dependencies, and configuring agents and skills.

- Initial repository setup and dependency installation.
- Shell configuration and dotfile integration.
- Agent and skill registration in `.claude/` directories.

### Troubleshooting

Diagnostic procedures for common issues encountered during development and pipeline execution.

- Environment and PATH debugging.
- Agent configuration resolution failures.
- Pipeline stage errors and recovery steps.
- Dependency conflicts and version mismatches.

### Agent and Skill Reference

Operational guides for the agent ecosystem.

- Agent invocation patterns and context loading.
- Skill configuration and registration.
- Pipeline orchestration and stage sequencing.
- Approval gate procedures (commit, deploy, external API calls).

### Pipeline Execution

Step-by-step guides for running the delivery pipeline and its individual stages.

- Full pipeline execution from analysis through delivery summary.
- Running individual stages (code review, security scan, QA).
- Persisting artifacts to `/docs/features/<feature-name>/`.
- Reviewing and shipping delivery summaries.

---

## Contributing a Runbook

Each runbook follows this structure:

1. **Title** -- What this runbook covers.
2. **When to use** -- The scenario or symptom that triggers this runbook.
3. **Prerequisites** -- Tools, access, or configuration required.
4. **Steps** -- Numbered, sequential actions. Each step is atomic and verifiable.
5. **Verification** -- How to confirm the runbook completed successfully.
6. **Escalation** -- What to do if the runbook does not resolve the issue.

---

## Existing Runbooks

- [context-engineering.md](context-engineering.md) — context taxonomy and principles, plus the full map of every Memory (KIs, ADRs, TEAM_TOPOLOGY.md, feature archive) and Learning (retrospective, agent-scorecard, agent-eval, extract-lessons) mechanism and how they differ
- [adding-a-new-platform.md](adding-a-new-platform.md) — wiring a new AI tool into the `shared/` -> generated-config pipeline
- [editing-agent-prompts.md](editing-agent-prompts.md) — versioning, changelog, and testing workflow for `shared/agents/` edits
- [scaling-cross-feature-learning.md](scaling-cross-feature-learning.md) — when and how to replace `context-engineer`'s grep-based same-bounded-context lookup with a real index, once the feature archive has grown enough to need one
- [self-audit-prompt.md](self-audit-prompt.md) — a reusable prompt for having a *different* agent independently audit this framework's internal consistency (twin-file drift, missing contracts, stale docs, dead references)

General runbook summaries are also available in `/docs/RUNBOOKS.md`. As new individual runbooks are created, add them to this directory and link them here.
