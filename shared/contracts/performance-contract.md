# Contract: performance-report.md

**Produced by**: performance-engineer
**Consumed by**: developer (mandates must be implemented), data-engineer (if invoked, for query-pattern alignment)

## Required Sections (exact heading text and level)
- `## 1. Idempotency Guarantees`
- `## 2. Timeout & Circuit Breaker Mandates`
- `## 3. N+1 Query Prevention`
- `## 4. Hot Path Caching`
- `## Notes for Developer`

## Validation Rule
Each of the four numbered sections must contain a `**Status**` line — a section missing its status line means
the agent skipped judging that risk category rather than explicitly clearing it.

This is a structural check only. It does not verify the mandates are technically correct or sufficient —
that's the human PAUSE checkpoint and the developer's/code-reviewer's job.

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
