---
feature: "<feature-name>"
bounded_context: "<owning-bounded-context>"
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

<!--
Template for implementation-notes.md. Consumed by the developer agent.
Structure defined here; contract in shared/contracts/implementation-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Implementation Notes: [Feature Name]

## Files Created
- `path/to/file.py` — [what it does]

## Files Modified
- `path/to/file.py` — [what changed and why]

## Interface Design
[Include the public interfaces, types, or signatures designed before implementation]

## Named Refactoring Log
- **[Operation Name]**: `path/to/file.py:45`
  - **Before**: [What the smell was]
  - **After**: [What it became]

## Self-Review Checklist
- [x/ ] Every public method has an intention-revealing name
- [x/ ] No function exceeds 30 LOC
- [x/ ] Cyclomatic complexity < 7 on all new functions
- [x/ ] No primitive obsession
- [x/ ] No feature envy
- [x/ ] No magic numbers or strings
- [x/ ] Dependency direction verified

## Simple Design Verification
1. **Passes all tests**: [yes/no]
2. **Reveals intention**: [yes/no] — [what was renamed/extracted]
3. **No duplication**: [yes/no] — [what was deduplicated]
4. **Fewest elements**: [yes/no] — [what was removed]

## Key Decisions
- [Decision made]: [reasoning] (e.g., "Used repository pattern here to match existing auth module")

## Deviations from Analysis
- [Any task from analysis that was skipped or changed]: [reason]

## Dependencies Added
- [package]: [version] — [reason], or "None"

## Notes for QA
- [Things the QA engineer should pay special attention to]
- [Known edge cases that should be tested specifically]

## Notes for DevOps
- [New env vars required]
- [New services or infrastructure needed]
- [Migration steps required]
