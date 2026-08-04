# Contract: data-engineering-notes.md

**Produced by**: data-engineer
**Consumed by**: developer, validate-migrations skill (if migration files were produced)

## Required Sections (exact heading text and level)
- `## Schema Changes`
- `## Migration Strategy`
- `## Files Modified/Created`
- `## Developer Handoff Notes`

## Validation Rule
`## Migration Strategy` must contain a `**Phase**` line with value `Expand`, `Contract`, or `Safe Addition` —
matching `shared/rules/architecture-guardrails.md` rule #2 (Expand/Contract is non-negotiable). A migration
strategy with no declared phase is exactly the kind of gap that rule exists to prevent.

This is a structural check only. It does not verify the migration is actually safe or non-destructive —
that's `validate-migrations`' job (invoked directly by this agent's own process) and the human PAUSE
checkpoint before deploy.

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
