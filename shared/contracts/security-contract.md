# Contract: security-report.md

**Produced by**: security-reviewer
**Consumed by**: qa-engineer, tech-writer, orchestrator (deliver-feature)

## Typed State (`loom run`)

Under `loom run`, this artifact is **typed state**, not a markdown document: the stage returns JSON
conforming to `shared/schemas/pipeline/security.schema.json` (generated from `internal/state/` —
never hand-edit it), the executor validates it, and `security-report.md` is *rendered* from that
state as a human-readable view (roadmap L2.9). The rendered view reproduces every heading below.

The typed form makes the severity ladder a real enum (`CRITICAL`, `HIGH`, `MEDIUM`, `LOW`, `INFO`)
rather than a string in prose, and it turns the Validation Rule below into an **invariant checked at
load time**: a CRITICAL or HIGH finding with an empty `fixApplied` is a validation error that names
the finding, and the stage does not complete. Under the markdown pipeline the same rule is a
`validate-artifact` check run after the fact; under the executor a report that violates it never
becomes state at all. The "block on Critical findings" guardrail stops depending on a grep matching
the right words.

`qa-engineer` receives a projection of the findings and the QA notes, not the whole report;
`tech-writer` reads the rendered view.

For the markdown pipeline (the `deliver-feature` skill), everything below remains authoritative
exactly as written.

## Required Sections (exact heading text and level)
- `## Threat Model Summary`
- `## Dependency Audit`
- `## STRIDE Analysis`
- `### Spoofing`
- `### Tampering`
- `### Repudiation`
- `### Information Disclosure`
- `### Denial of Service`
- `### Elevation of Privilege`
- `## Findings`
- `## Files Modified`
- `## Security Checklist`
- `## Notes for QA`
- `## Notes for Tech Writer`

## Validation Rule
`validate-artifact` checks presence of every heading above, all six STRIDE categories included. Any `## Findings` entry marked `CRITICAL` or `HIGH` must have a non-empty `**Fix applied**` line — per the agent's own rule that Critical/High findings are fixed directly, not left as recommendations. A Critical/High finding with `Fix applied: Recommendation only` is a FAIL; the pipeline's "block on Critical findings" guardrail exists precisely to catch this.

This is a structural check only. It does not re-run the threat model — that's the security-reviewer's job, checked here only for completeness of the artifact shape.

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
