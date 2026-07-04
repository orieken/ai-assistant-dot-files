# Feature Delivery Artifacts

This directory holds the output of every feature that passes through the delivery pipeline. Each feature gets its own subdirectory named after the feature.

---

## Convention

```
docs/features/<feature-name>/
```

The feature name uses kebab-case and matches the feature branch or ticket identifier where possible (e.g., `user-authentication`, `dashboard-metrics`, `api-rate-limiting`).

---

## Pipeline Artifacts

Each stage of the delivery pipeline writes a specific artifact to the feature directory. A fully delivered feature contains these files:

| Artifact                     | Stage                  | Description                                                      |
|-------------------------------|------------------------|-------------------------------------------------------------------|
| `context-manifest.md`         | Context Engineering   | Bounded context, pinned files, KIs/ADRs, prior deliveries, token budget. |
| `analysis.md`                 | Analysis               | Domain analysis, requirements breakdown, acceptance criteria.    |
| `architecture-notes.md`       | Architecture           | Layer design, dependency decisions, integration points.          |
| `performance-report.md`       | Performance            | SLA verification, N+1/caching/timeout review (if invoked).       |
| `data-engineering-notes.md`   | Data Engineering       | Schema design, Expand/Contract migration plan (if invoked).      |
| `implementation-notes.md`     | Implementation         | Key implementation decisions, trade-offs, technical debt notes.  |
| `code-review-report.md`       | Code Review            | Review findings, SOLID violations, complexity flags.             |
| `accessibility-report.md`     | Accessibility          | Semantic HTML / WCAG findings (if UI changes were involved).     |
| `security-report.md`          | Security               | Threat model results, STRIDE analysis, mitigation actions.       |
| `qa-report.md`                | Quality Assurance      | Test coverage summary, edge cases, regression risk assessment.   |
| `observability-report.md`     | Observability          | SLIs, OTel spans, structured logging, PII hygiene review.        |
| `docs-report.md`               | Documentation          | Documentation completeness check, API doc updates, ADR links.    |
| `devops-report.md`             | DevOps                 | CI/CD changes, deployment notes, infrastructure requirements.    |
| `delivery-summary.md`         | Delivery               | Final summary with status, metrics, and sign-off.                |

Not every feature requires every artifact. Smaller changes may skip stages that do not apply. The delivery summary always indicates which stages were executed.

---

## Template: Completed Feature Directory

```
docs/features/user-authentication/
  context-manifest.md
  analysis.md
  architecture-notes.md
  performance-report.md
  data-engineering-notes.md
  implementation-notes.md
  code-review-report.md
  accessibility-report.md
  security-report.md
  qa-report.md
  observability-report.md
  docs-report.md
  devops-report.md
  delivery-summary.md
```

---

## Historical Reference

Agents use past feature directories as pattern reference. When delivering a new feature, agents can read completed features to:

- Identify recurring architectural patterns and decisions.
- Reuse analysis frameworks and review checklists.
- Understand how similar trade-offs were resolved previously.
- Maintain consistency in artifact format and depth.

This directory grows over time into a searchable knowledge base of every delivered capability.
