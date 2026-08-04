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
Template for docs-report.md. Consumed by the tech-writer agent.
Structure defined here; contract in shared/contracts/docs-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Documentation Report: [Feature Name]

## Files Updated
- `CHANGELOG.md` — Added entry for [feature]
- `README.md` — Added section on [feature]
- `docs/adr/ADR-005-[name].md` — New ADR for [decision]

## Files Unchanged (and why)
- `docs/setup.md` — No new setup steps required

## Notes for DevOps
- [New env vars that need to be documented in deployment runbooks]
- [New infrastructure that needs runbook entries]
