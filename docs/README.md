# Documentation Knowledge Base

This directory is the central knowledge base for the ai-assistant-dot-files repository. It serves two audiences: AI agents that load context for automated pipelines, and humans who need onboarding material, reference guides, and operational runbooks.

---

## Directory Structure

```
docs/
  README.md              -- This file. Index and navigation guide.
  ARCHITECTURE.md         -- The shared/ canonical layer, the Capability Tier system.
  AGENT_REFERENCE.md      -- Every agent's role and counterbalance (contract, reviewer, gate, metric).
  CONTRIBUTING.md         -- How to add a new agent, skill, rule, or platform.
  MIGRATION.md            -- Breaking changes from pre-shared/ structure.
  ONBOARDING.md           -- New contributor onboarding guide.
  RUNBOOKS.md             -- Compatibility redirect to the canonical runbook index.
  THREAT_MODEL.md         -- STRIDE threat model; implementation-status annotations (Epic 65).
  TODO.md                 -- Folder-by-folder documentation audit ledger.
  human-tasks.md          -- Decisions and tasks that require human action, not agent execution.
  roadmap-2026-08-07.md   -- 2026-08-07 framework roadmap with Epics 69–73 (superseded by roadmaps/).
  clean-code-guidelines.docx -- Full language-example reference, cited from root CLAUDE.md.
  adrs/                  -- Architecture Decision Records.
  agent-metrics/          -- Monthly agent-scorecard output, evals subdirectory.
  aos/                    -- AOS (Agent Orchestration System) design docs and phase prompts.
  blueprints/             -- Focused framework blueprint prompts consumed by skills and agents.
  features/              -- Pipeline artifacts for delivered features (permanent archive).
  lessons-learned/        -- Persisted output from extract-lessons.
  patterns/               -- Reusable design and framework pattern docs (Saturday/Sunday/GoF/Clean Architecture).
  pipeline-retrospectives/ -- Cross-delivery trend analysis output.
  prompts/                -- Self-contained agent handoff prompts for framework improvements.
  roadmaps/               -- Forward-looking build plans, including BUILD-ROADMAP.md, the single
                             authoritative L2-L4 roadmap linked from the root README.
  runbooks/               -- Operational runbooks and guides (see runbooks/README.md for the index).
  blog-posts/             -- Draft blog content about the framework (see blog-content-brief.md).
  audits/                 -- Dated audit reports, kept as a permanent series: <kind>-<YYYY-MM-DD>.md.
                             Committed, not scratch -- health-check.sh warns when the newest
                             doc-audit-*.md is stale, and that signal only works if the series
                             persists and is visible to everyone. Superseded reports stay for
                             comparison rather than being deleted.
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

- **New contributors** start with `ONBOARDING.md` for setup instructions and
  [`runbooks/README.md`](runbooks/README.md) for operational context.
- **Anyone extending the framework** reads `CONTRIBUTING.md` and `ARCHITECTURE.md` before adding an agent, skill, or rule.
- **Architects and leads** review `/docs/adrs/` for decision history.
- **Feature reviewers** check `/docs/features/<feature-name>/` for the full artifact trail of any delivered feature.

---

## Feature Delivery Artifact Convention

All pipeline artifacts are persisted to `/docs/features/<feature-name>/`. When a feature moves through the delivery pipeline, each stage writes its output to the feature directory. This creates a permanent, auditable record of every decision, review, and report associated with delivery.

See `/docs/features/README.md` for the full convention and artifact list.
