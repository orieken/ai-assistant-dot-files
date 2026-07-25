---
name: template-ki-do-not-use
tags: [template, knowledge-item]
domain: documentation
created: 2026-07-24
---

## Context

Explain the background, problem statement, or trigger condition that makes this Knowledge Item relevant. Describe when an agent should load or cite this item during feature analysis, context engineering, or debugging.

## Pattern

Document the established pattern, solution, best practice, or architectural guidance. Use clear prose, concrete file references, and short code snippets where appropriate.

Practical implication: explain how downstream agents (such as analyst, developer, or code-reviewer) should apply this pattern to prevent regressions or redundant re-derivation.

---

## Frontmatter reference

| Field | Required | Format / notes | Rationale |
|---|---|---|---|
| `name` | ✓ | kebab-case; must match filename base (`my-pattern.md` → `name: my-pattern`) | Referenced by `[[link]]` syntax and memory-registry index |
| `tags` | ✓ | List of kebab-case strings | Consumed by `search-ki` tag filtering and `memory-auditor` |
| `domain` | ✓ | Single kebab-case string (e.g. `testing`, `architecture`, `mcp-server-pattern`) | Consumed by `context-engineer` for bounded-context mapping |
| `created` | ✓ | ISO date (`YYYY-MM-DD`) | Immutable authoring date used for staleness calculations |

The framework enforces field presence and schema validity via `scripts/health-check.sh` and `shared/schemas/ki-frontmatter.schema.json`. See `docs/patterns/frontmatter-conventions.md` for full details.

---

*Delete this reference table + this footer when using the template. Keep the frontmatter block + the body sections above.*
