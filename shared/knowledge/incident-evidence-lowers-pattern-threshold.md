---
name: incident-evidence-lowers-pattern-threshold
tags: [extract-lessons, incidents, pattern-extraction, signal-quality, thresholds]
domain: memory-pipeline
created: 2026-08-05
---

Pattern-extraction thresholds scale inversely with evidence quality. Production incidents are
stronger evidence than code-review findings, so the threshold for acting on incident-sourced
patterns is lower.

## Applied thresholds in extract-lessons

| Source | Occurrences to claim a pattern | Rationale |
|---|---|---|
| Code-review findings (Steps 1–2) | 3+ across deliveries | A single code-review finding is cheap signal — reviewers flag issues speculatively; 3 occurrences establishes a genuine recurring pattern |
| Retrospective signals (Steps 3–5) | 3+ | Same rationale — retrospective items are raised by feel, not confirmed by production data |
| Incident-feature pairs (Step 6) | 2+ | A production incident is expensive, user-visible evidence that something broke; 2 incidents in the same area strongly suggests a systemic gap |

## Principle

The number of occurrences required to act should reflect how costly a false negative would be.
For code-review speculation: low cost — wait for 3. For production failure: high cost — act at 2.

## Where this is enforced

`shared/skills/extract-lessons/SKILL.md` Step 6 states the 2+ threshold explicitly with this
rationale. Steps 1-2 state the 3+ threshold. These are deliberate asymmetric choices, not an
oversight or inconsistency.

## Practical implication

When adding a new extraction step to `extract-lessons`, choose a threshold based on the signal
quality of the source:

- Speculative / heuristic signal (linter findings, review comments) → 3+
- Confirmed production signal (incident records, SLO breaches, customer-reported failures) → 2+
- Unambiguous data (SLO breached at the same threshold N consecutive weeks) → 1 may be justified
