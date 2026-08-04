# Contract: devops-report.md

**Produced by**: devops-engineer
**Consumed by**: the orchestrator's `delivery-summary.md` synthesis (final pipeline agent — nothing downstream reads this directly)

## Required Sections (exact heading text and level)
- `## Files Created`
- `## Files Modified`
- `## New Environment Variables Required`
- `## Migration Steps`
- `## Deployment Notes`
- `## Manual Steps Required`

## Validation Rule
This is a structural check only. It does not scan the `New Environment Variables Required` table for
accidentally-real secret values — `devops-engineer`'s own rule ("Do NOT add real credentials") is a judgment
guardrail on the agent, and `security-reviewer` (earlier in the pipeline) is the check for real hardcoded
secrets in code, not this contract.

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
