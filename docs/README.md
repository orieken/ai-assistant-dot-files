# Documentation Knowledge Base

This directory is the central knowledge base for the ai-assistant-dot-files repository. It serves two audiences: AI agents that load context for automated pipelines, and humans who need onboarding material, reference guides, and operational runbooks.

---

## Directory Structure

```
docs/
  README.md              -- This file. Index and navigation guide.
  ARCHITECTURE.md         -- The shared/ canonical layer, the Capability Tier system.
  AGENT_REFERENCE.md      -- Every agent's role and counterbalance (contract, reviewer, gate, metric).
  CLAUDE.md               -- Claude Code agent configuration reference.
  CONTRIBUTING.md         -- How to add a new agent, skill, rule, or platform.
  MIGRATION.md            -- Breaking changes from pre-shared/ structure.
  ONBOARDING.md           -- New contributor onboarding guide.
  RUNBOOKS.md             -- Operational runbook summaries.
  clean-code-guidelines.docx -- Full language-example reference, cited from CLAUDE.md.
  adrs/                  -- Architecture Decision Records.
  agent-metrics/          -- Monthly agent-scorecard output, evals subdirectory.
  features/              -- Pipeline artifacts for delivered features (permanent archive).
  lessons-learned/        -- Persisted output from extract-lessons.
  patterns/               -- Reusable design and framework pattern docs (Saturday/Sunday/GoF/Clean Architecture).
  pipeline-retrospectives/ -- Cross-delivery trend analysis output.
  runbooks/               -- Operational runbooks and guides (see runbooks/README.md for the index).
  blog-posts/             -- Draft blog content about the framework (see blog-content-brief.md).
  audits/                 -- Ad-hoc external audit reports -- scratch inputs, not permanent
                             documentation; act on findings, then delete once resolved.
  incidents/              -- Production incident records in <YYYY-MM-DD>-<slug>.md format.
                             Cross-referenced by the bugfix pipeline (deliver-bugfix "Fixed by" link)
                             and by extract-lessons (Step 6 incident-feature pair mining). Permanent —
                             Status field updated to Resolved; records are never deleted.
```

`docs/aos/` contains opt-in v3/AOS migration and governance material. It is tracked on `main`, but the
runtime remains backward-compatible with v2 installs unless a team explicitly adopts the AOS layers.

---

## How Agents Use These Docs

Agents load specific subdirectories as context depending on the task at hand:

- **Feature delivery agents** read `/docs/features/<feature-name>/` to understand prior decisions, architecture notes, and review reports for a given feature.
- **Pattern-aware agents** read `/docs/patterns/` to apply consistent design patterns across the codebase (Saturday Framework, Sunday Framework, Clean Architecture).
- **Architecture agents** read `/docs/adrs/` to understand past decisions and their rationale before proposing new ones.
- **Operational agents** read `/docs/runbooks/` for troubleshooting procedures and setup instructions.
- **Memory/learning skills** (`agent-scorecard`, `extract-lessons`, `pipeline-retrospective`) read and write `/docs/agent-metrics/`, `/docs/lessons-learned/`, `/docs/pipeline-retrospectives/`.

Each subdirectory is self-contained. Agents can load a single directory without pulling the entire docs tree.

---

## How Humans Use These Docs

- **New contributors** start with `ONBOARDING.md` for setup instructions and `RUNBOOKS.md` for operational context.
- **Anyone extending the framework** reads `CONTRIBUTING.md` and `ARCHITECTURE.md` before adding an agent, skill, or rule.
- **Architects and leads** review `/docs/adrs/` for decision history.
- **Feature reviewers** check `/docs/features/<feature-name>/` for the full artifact trail of any delivered feature.

---

## Feature Delivery Artifact Convention

All pipeline artifacts are persisted to `/docs/features/<feature-name>/`. When a feature moves through the delivery pipeline, each stage writes its output to the feature directory. This creates a permanent, auditable record of every decision, review, and report associated with delivery.

See `/docs/features/README.md` for the full convention and artifact list.
