# Contract: architecture-notes.md

**Produced by**: architect
**Consumed by**: performance-engineer, data-engineer, developer, code-reviewer

## Required Sections (exact heading text and level)
- `## Structural Decisions`
- `## Component Placement`
- `## Bounded Context`
- `## Stability Design`
- `## Observability Design`
- `## Layer Boundary Checks`
- `## Anti-Pattern Check`
- `## Fitness Functions`
- `## Refactoring Opportunities`
- `## Developer Handoff Notes`
- `## Open Architectural Questions`

## Validation Rule
`validate-artifact` checks presence of every heading above. Each "Structural Decisions" sub-entry must contain a `**Fitness Function**` line (per `architect.md`'s own rule: a decision without one must be explicitly justified and flagged `judgment-only`) — `validate-artifact` flags a decision block missing this line as a WARNING, not a FAIL, since judgment-only is a legitimate documented exception.

If the architect wrote an RFC (`rfc-*.md`), that file is validated separately by human review at the RFC PAUSE checkpoint — it has no machine-checked contract.

This is a structural check only. It does not verify decisions are sound — that's the human RFC checkpoint and code-reviewer's job.

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
