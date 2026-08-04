# Contract: code-review-report.md

**Produced by**: code-reviewer
**Consumed by**: security-reviewer, qa-engineer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## Overall Status`
- `## Design Narrative`
- `## Design Score`
- `## Security Surface`
- `## Performance Surface`
- `## Test Design Review`
- `## Verification of Developer Self-Review`
- `## Feedback for the Developer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## Overall Status` must contain exactly one of `APPROVED` or `CHANGES REQUESTED` (bolded, per the agent's own template) — anything else is a FAIL, since the orchestrator's CHANGES REQUESTED loop parses this literal string.
- `## Design Score` must contain all four dimensions (Clarity, Cohesion, Coupling, Craft) with a numeric 1-5 rating each — a missing dimension is a FAIL.

This is a structural check only. It does not re-judge whether APPROVED was the right call — that's the security-reviewer and qa-engineer's job to catch downstream, and the human's job at the CHANGES REQUESTED checkpoint.

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
