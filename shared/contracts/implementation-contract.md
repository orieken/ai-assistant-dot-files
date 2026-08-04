# Contract: implementation-notes.md

**Produced by**: developer
**Consumed by**: code-reviewer, accessibility-engineer, security-reviewer, qa-engineer, sre-engineer, tech-writer, devops-engineer

## Required Sections (exact heading text and level)
- `## Files Created`
- `## Files Modified`
- `## Interface Design`
- `## Named Refactoring Log`
- `## Self-Review Checklist`
- `## Simple Design Verification`
- `## Key Decisions`
- `## Deviations from Analysis`
- `## Dependencies Added`
- `## Notes for QA`
- `## Notes for DevOps`

## Validation Rule
`validate-artifact` checks presence of every heading above. `## Self-Review Checklist` and `## Simple Design Verification` must each contain at least one checked item (`[x]` or an explicit `yes`/`no`) — an unfilled checklist template is a FAIL, not a pass, since it means the developer skipped self-review rather than completed it.

This is a structural check only. It does not verify the code itself is correct or clean — that's code-reviewer's job.

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
