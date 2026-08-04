# Contract: docs-report.md

**Produced by**: tech-writer
**Consumed by**: devops-engineer (reads this for ops runbook notes, per its own process step 3)

## Required Sections (exact heading text and level)
- `## Files Updated`
- `## Files Unchanged (and why)`
- `## Notes for DevOps`

## Validation Rule
This is a structural check only — presence of all three sections. It does not verify the documentation
itself is accurate; `tech-writer`'s own rule ("Do NOT invent behavior — only document what was actually
implemented") is a judgment guardrail on the agent, not something this contract can mechanically enforce.

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
