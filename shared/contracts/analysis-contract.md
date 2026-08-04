# Contract: analysis.md

**Produced by**: analyst
**Consumed by**: architect, performance-engineer, data-engineer, developer, qa-engineer, tech-writer, devops-engineer

## Required Sections (exact heading text and level)
- `## Summary`
- `### Acceptance Criteria`
- `### Non-Functional Requirements`
- `## Proposed Fitness Functions`
- `## Out of Scope`
- `## Technical Breakdown`
- `### Bounded Context`
- `### Domain Events (Event Storming Lite)`
- `### Affected Components`
- `### Data Model Changes`
- `### API Changes`
- `### New Dependencies`
- `## Task List`
- `### Developer Tasks`
- `### QA Tasks`
- `### Tech Writer Tasks`
- `### DevOps Tasks`
- `## Edge Cases and Risks`
- `## Definition of Done`

## Validation Rule
`validate-artifact` checks presence of every heading above, exact string and level match. Missing a heading is a FAIL — even if the content would logically live under a sibling section, downstream agents (developer, qa-engineer, tech-writer, devops-engineer) grep for these exact headings to find "their" task list.

This is a structural check only. It does not verify the content is correct, complete, or non-placeholder — that judgment belongs to the human PAUSE checkpoint after analyst and to the architect/developer who consume it.

## Retrieval Frontmatter (WARN)

Pipeline artifacts should include a YAML frontmatter block at the very top of the file. Missing or incomplete retrieval frontmatter triggers a **WARN** from `validate-artifact` — not a FAIL. Existing archived artifacts without frontmatter are unaffected.

```yaml
---
feature: "<feature-name>"             # kebab-case slug derived from the feature file name
bounded_context: "<context>"          # owning bounded context (from DOMAIN_DICTIONARY.md domain list)
domain_terms: []                      # canonical terms from DOMAIN_DICTIONARY.md used in this feature
files_touched: []                     # repo-relative paths of files created or modified
issue_refs: []                        # ticket/issue references (e.g., PROJ-123, #456)
linked_adrs: []                       # repo-root-relative paths to referenced ADRs
linked_kis: []                        # repo-root-relative paths to referenced Knowledge Items
---
```

Once frontmatter adoption is visible across a project's feature archive, this check will be promoted to FAIL in a future release.
