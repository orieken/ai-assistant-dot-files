# Contract: analysis.md

**Produced by**: analyst
**Consumed by**: architect, performance-engineer, data-engineer, developer, qa-engineer, tech-writer, devops-engineer

## Typed State (`loom run`)

Under `loom run`, this artifact is **typed state**, not a markdown document: the stage returns JSON
conforming to `shared/schemas/pipeline/analysis.schema.json` (generated from `internal/state/` — never hand-edit it),
the executor validates it, and `analysis.md` is *rendered* from that state as a human-readable view
(roadmap L2.9). The rendered view reproduces every heading below, so the structural check in this
contract still describes what a reader sees; the machine handoff no longer goes through it.

Two consequences worth knowing: the view is derived, so editing it changes nothing a downstream
stage reads — the state document under `state/` is what integrity tracks (L2.12). And a downstream
stage receives only the fields its projection declares, not the whole document.

For the markdown pipeline (the `deliver-feature` skill), everything below remains authoritative
exactly as written.

## Required Sections (exact heading text and level)
- `## Summary`
- `### Acceptance Criteria`
- `### Non-Functional Requirements`
- `## Proposed Fitness Functions`
- `## Out of Scope`
- `## Technical Breakdown`
- `### Bounded Context`
- `### Domain Events (Event Storming Lite)`
- `### Affected Components`
- `### Data Model Changes`
- `### API Changes`
- `### New Dependencies`
- `## Task List`
- `### Developer Tasks`
- `### QA Tasks`
- `### Tech Writer Tasks`
- `### DevOps Tasks`
- `## Edge Cases and Risks`
- `## Definition of Done`

## Validation Rule
`validate-artifact` checks presence of every heading above, exact string and level match. Missing a heading is a FAIL — even if the content would logically live under a sibling section, downstream agents (developer, qa-engineer, tech-writer, devops-engineer) grep for these exact headings to find "their" task list.

This is a structural check only. It does not verify the content is correct, complete, or non-placeholder — that judgment belongs to the human PAUSE checkpoint after analyst and to the architect/developer who consume it.

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
