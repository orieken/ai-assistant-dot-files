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
