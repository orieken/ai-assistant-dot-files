# Contract: security-report.md

**Produced by**: security-reviewer
**Consumed by**: qa-engineer, tech-writer, orchestrator (deliver-feature)

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
