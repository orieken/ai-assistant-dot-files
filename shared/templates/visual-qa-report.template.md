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
Template for visual-qa-report.md. Consumed by the visual-qa-engineer agent.
Contract in shared/contracts/visual-qa-report-contract.md validates required headings.
Preserve every heading exactly. Write "None" for inapplicable sections — never delete headings.
-->

# Visual QA Report: [Feature Name]

## Summary

[PASS | FAIL | UNCONFIGURED]

- Heatmap analysis: [PASS / FAIL / UNCONFIGURED — Coverage Score: N%]
- Visual regression: [PASS / FAIL / SKIPPED — N diffs detected | no baselines yet]
- Overall verdict: [one sentence]

## Visual Regression

[Describe screenshot baseline comparison results. If PASS: "All N visual scenarios passed against existing baselines."
If FAIL: name each failing scenario, the component affected, and the path to the diff file.
If SKIPPED: "No Playwright snapshot baselines found. See Recommendations for setup guidance."
If no visual tests exist: "No @visual-tagged tests found in the test suite."]

| Scenario | Status | Diff Path |
|---|---|---|
| [scenario name] | PASS / FAIL | [path or —] |

## Heatmap Coverage

[Results from @orieken/saturday-ml-analyzer. If UNCONFIGURED: "heatmap-data/ directory not found — @orieken/saturday-playwright-heatmap not configured in this project."]

- Coverage Score: [N%] (threshold: 80%)
- Hotspots identified: [N]
- Cold Spots identified: [N]
- Scenarios analyzed: [N heatmap JSON files]

## Cold Spots

[List each unexercised interactable element. Classify severity per element.
If no cold spots: "None — all interactable elements are covered."
If UNCONFIGURED: "None — heatmap instrumentation not present."]

| Element | Selector | Severity | Notes |
|---|---|---|---|
| [element type + label] | [CSS/ARIA selector] | HIGH / LOW | [primary journey / decorative] |

## Recommendations

[Actions to improve visual coverage. Must be non-empty if Summary is FAIL.
Examples: "Add @visual tag to login scenario and capture baseline", "Wrap test fixture with heatmapTest to enable coverage tracking", "Add visual test for [component] — cold spot on primary CTA."]

- [ ] [Recommendation 1]
- [ ] [Recommendation 2]

## Notes for QA

[Any observations for the qa-engineer or future visual QA work.
If nothing to add: "None."]
