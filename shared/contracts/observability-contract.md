# Contract: observability-report.md

**Produced by**: sre-engineer
**Consumed by**: devops-engineer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## 1. Service Level Indicators (SLIs)`
- `## 2. OpenTelemetry & Tracing`
- `## 3. Log Quality & Cardinality`
- `## 4. PII Data Hygiene`
- `## Notes for DevOps Engineer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## 4. PII Data Hygiene` status must be `Clean` or explicitly resolved — per the agent's own rule that PII violations must be fixed directly with `Edit`, not left as a recommendation. A status of `Violation detected` with no accompanying fix note is a FAIL.

This is a structural check only. It does not verify OTel spans actually exist in production — that's an operational concern outside this pipeline's scope.

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
