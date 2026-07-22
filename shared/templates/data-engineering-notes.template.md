<!--
Template for data-engineering-notes.md. Consumed by the data-engineer agent.
Structure defined here; contract in shared/contracts/data-engineering-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Data Engineering Notes: [Feature Name]

## Schema Changes
- [Table Name]
  - `ADDED column_name (type)`
  - `MODIFIED column_name` (using Expand/Contract)

## Migration Strategy
- **Phase**: [Expand / Contract / Safe Addition]
- **Details**: [Explain the migration steps, how dual-writes will be handled, or why it's a safe addition]

## Files Modified/Created
- `path/to/migration_file.sql` — [What it does]

## Developer Handoff Notes
- [Instructions for the developer on how to implement the application-side of the Expand/Contract pattern (e.g., dual writes)]
- [Query performance guidelines (e.g., "Make sure to DataLoader the new relation")]
