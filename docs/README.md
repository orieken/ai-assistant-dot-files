# Documentation Knowledge Base

This directory is the central knowledge base for the ai-assistant-dot-files repository. It serves two audiences: AI agents that load context for automated pipelines, and humans who need onboarding material, reference guides, and operational runbooks.

---

## Directory Structure

```
docs/
  README.md              -- This file. Index and navigation guide.
  CLAUDE.md              -- Claude Code agent configuration reference.
  ONBOARDING.md          -- New contributor onboarding guide.
  RUNBOOKS.md            -- Operational runbook summaries.
  build-out-prompts.md   -- Prompt engineering templates.
  master-build-out-prompts.md -- Extended prompt catalog.
  thoughtworks-specialist.md  -- Specialist agent configuration.
  dotfiles-remediation.md     -- Remediation notes for dotfile configs.
  dotfiles-additions/         -- Supplementary dotfile configurations.
  spec-writer/                -- Spec-writer agent assets.
  features/              -- Pipeline artifacts for delivered features.
  patterns/              -- Reusable design and framework pattern docs.
  adrs/                  -- Architecture Decision Records.
  runbooks/              -- Operational runbooks and guides.
```

---

## How Agents Use These Docs

Agents load specific subdirectories as context depending on the task at hand:

- **Feature delivery agents** read `/docs/features/<feature-name>/` to understand prior decisions, architecture notes, and review reports for a given feature.
- **Pattern-aware agents** read `/docs/patterns/` to apply consistent design patterns across the codebase (Saturday Framework, Sunday Framework, Clean Architecture).
- **Architecture agents** read `/docs/adrs/` to understand past decisions and their rationale before proposing new ones.
- **Operational agents** read `/docs/runbooks/` for troubleshooting procedures and setup instructions.

Each subdirectory is self-contained. Agents can load a single directory without pulling the entire docs tree.

---

## How Humans Use These Docs

- **New contributors** start with `ONBOARDING.md` for setup instructions and `RUNBOOKS.md` for operational context.
- **Architects and leads** review `/docs/adrs/` for decision history and `/docs/patterns/` for established conventions.
- **Feature reviewers** check `/docs/features/<feature-name>/` for the full artifact trail of any delivered feature.

---

## Feature Delivery Artifact Convention

All pipeline artifacts are persisted to `/docs/features/<feature-name>/`. When a feature moves through the delivery pipeline, each stage writes its output to the feature directory. This creates a permanent, auditable record of every decision, review, and report associated with delivery.

See `/docs/features/README.md` for the full convention and artifact list.
