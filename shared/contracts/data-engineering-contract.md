# Contract: data-engineering-notes.md

**Produced by**: data-engineer
**Consumed by**: developer, validate-migrations skill (if migration files were produced)

## Required Sections (exact heading text and level)
- `## Schema Changes`
- `## Migration Strategy`
- `## Files Modified/Created`
- `## Developer Handoff Notes`

## Validation Rule
`## Migration Strategy` must contain a `**Phase**` line with value `Expand`, `Contract`, or `Safe Addition` —
matching `shared/rules/architecture-guardrails.md` rule #2 (Expand/Contract is non-negotiable). A migration
strategy with no declared phase is exactly the kind of gap that rule exists to prevent.

This is a structural check only. It does not verify the migration is actually safe or non-destructive —
that's `validate-migrations`' job (invoked directly by this agent's own process) and the human PAUSE
checkpoint before deploy.
