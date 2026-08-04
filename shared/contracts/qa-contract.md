# Contract: qa-report.md

**Produced by**: qa-engineer
**Consumed by**: sre-engineer, tech-writer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## Test Files Created`
- `## Test Files Modified`
- `## Coverage Summary`
- `## Test Results`
- `## Accessibility Check`
- `## Bugs Found`
- `## Known Gaps`
- `## Notes for Tech Writer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## Test Results` must show `Failed: 0` — per the agent's own rule, tests must be green before the pipeline proceeds. A non-zero failure count is a FAIL, not a warning.

This is a structural check only. It does not re-run the test suite — that's already qa-engineer's job before writing the report.

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
