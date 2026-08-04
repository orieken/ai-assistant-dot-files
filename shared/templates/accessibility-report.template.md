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
Template for accessibility-report.md. Consumed by the accessibility-engineer agent.
Structure defined here; contract in shared/contracts/accessibility-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Accessibility & UX Report: [Feature Name]

## Evaluation Summary
- **Semantic HTML**: [Pass/Fail/Notes]
- **Interactive Elements**: [Pass/Fail/Notes]
- **ARIA & Labels**: [Pass/Fail/Notes]
- **Keyboard Navigation**: [Pass/Fail/Notes]

## Findings & Fixes
- `path/to/component.tsx` — Changed `<div onClick={...}>` to `<button type="button" onClick={...}>` to ensure keyboard accessibility.
- [Finding without autofix]: [Recommendation]

## Notes for QA
- [Specific interactions QA should test via keyboard only]
- [Screen reader verification points]
