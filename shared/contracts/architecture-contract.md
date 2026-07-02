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
