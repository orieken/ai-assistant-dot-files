# Contract: analysis.md

**Produced by**: analyst
**Consumed by**: architect, performance-engineer, data-engineer, developer, qa-engineer, tech-writer, devops-engineer

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
