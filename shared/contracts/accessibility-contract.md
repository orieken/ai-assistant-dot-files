# Contract: accessibility-report.md

**Produced by**: accessibility-engineer
**Consumed by**: qa-engineer (keyboard/screen-reader test points), security-reviewer (if invoked next)

## Required Sections (exact heading text and level)
- `## Evaluation Summary`
- `## Findings & Fixes`
- `## Notes for QA`

## Validation Rule
`## Evaluation Summary` must contain all four of: `**Semantic HTML**`, `**Interactive Elements**`,
`**ARIA & Labels**`, `**Keyboard Navigation**` — each with a Pass/Fail/Notes value. A missing category means
that risk area was never actually evaluated, not that it passed silently.

This is a structural check only. It does not verify a fix is actually accessible — that requires the human
PAUSE checkpoint or real assistive-technology testing, neither of which this contract can substitute for.

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
