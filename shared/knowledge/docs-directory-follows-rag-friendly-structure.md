---
name: docs-directory-follows-rag-friendly-structure
tags: [documentation, docs-structure, rag, onboarding]
domain: documentation-knowledge-base
created: 2026-07-02
---

`docs/` is not a flat dumping ground — it follows a deliberate subdirectory convention (ADR-001) so agents
can load a focused subset instead of scanning everything:

- `docs/features/<name>/` — every delivered feature's full pipeline artifact set (permanent record).
- `docs/adrs/` — sequentially-numbered Architecture Decision Records (context, decision, consequences).
- `docs/runbooks/` — operational guides (installation, context engineering, pipeline execution).
- `docs/agent-metrics/`, `docs/pipeline-retrospectives/` — cross-delivery quality/trend reports (added
  after ADR-001, same convention: one purpose per subdirectory, README explains it).

When adding new documentation, put it in the subdirectory whose purpose matches, not in `docs/` root —
ad-hoc root-level files defeat the reason this structure exists (see ADR-001's "What becomes harder": a
transitional period of root-level content was accepted deliberately, not endorsed as an ongoing pattern).

See: `docs/adrs/ADR-001-adopt-rag-friendly-docs-structure.md`.
