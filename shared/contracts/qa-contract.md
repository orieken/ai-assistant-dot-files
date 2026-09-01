# Contract: qa-report.md

**Produced by**: qa-engineer
**Consumed by**: sre-engineer, tech-writer, orchestrator (deliver-feature)

## Typed State (`loom run`)

Under `loom run`, this artifact is **typed state**, not a markdown document: the stage returns JSON
conforming to `shared/schemas/pipeline/qa.schema.json` (generated from `internal/state/` — never
hand-edit it), the executor validates it, and `qa-report.md` is *rendered* from that state as a
human-readable view (roadmap L2.9). The rendered view reproduces every heading below.

The typed form makes the Validation Rule below an **invariant checked at load time**: test results
are numeric fields, and a non-zero `failed` count is a validation error, so a red suite cannot become
completed state. Under the markdown pipeline the same rule is a `validate-artifact` check that greps
`## Test Results` for the literal `Failed: 0`; under the executor there is no string to match.

`qa-engineer` reads three upstreams here — a projection of `implementation-notes` (files touched,
QA notes, deviations), a projection of `security-report` (findings and QA notes), and a projection of
`analysis` (acceptance criteria, edge cases, QA tasks, definition of done). Each arrives as its own
labelled block, so provenance survives and same-named fields from different upstreams cannot collide.
That third projection replaced an LLM summarization call (roadmap L2.10).

For the markdown pipeline (the `deliver-feature` skill), everything below remains authoritative
exactly as written.

## Required Sections (exact heading text and level)
- `## Test Files Created`
- `## Test Files Modified`
- `## Coverage Summary`
- `## Test Results`
- `## Accessibility Check`
- `## Bugs Found`
- `## Known Gaps`
- `## Notes for Tech Writer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## Test Results` must show `Failed: 0` — per the agent's own rule, tests must be green before the pipeline proceeds. A non-zero failure count is a FAIL, not a warning.

This is a structural check only. It does not re-run the test suite — that's already qa-engineer's job before writing the report.

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
