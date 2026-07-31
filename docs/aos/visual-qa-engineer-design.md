# Visual QA Engineer Agent Design (Epic 46)

Drafted 2026-07-30 from investigation of saturday-monorepo visual packages.

## Saturday.ML Capability Map

| Package | What it does | What it does NOT do |
|---|---|---|
| `@orieken/saturday-playwright-heatmap` | Tracks `click`/`input`/`change` events during test runs; scans DOM for interactable elements; writes `heatmap-data/*.json` per scenario; generates HTML heatmap reports | Screenshot comparison, pixel diffing |
| `@orieken/saturday-ml-analyzer` | K-Means clustering of interaction events → Hotspots + Cold Spots; Coverage Score (% of interactable elements exercised) | Visual regression, screenshot comparison |
| C# `Saturday.ML` (via SixLabors.ImageSharp) | Screenshot diffing + heatmap overlays | Only available in C# ecosystem |

**Conclusion**: The TypeScript saturday packages provide interaction coverage analysis, not visual regression. Visual regression is handled by Playwright's native `toHaveScreenshot()` / `expect(page).toHaveScreenshot()`. The agent bridges both.

## Agent Scope

The `visual-qa-engineer` agent:
1. **Analyzes interaction heatmaps** — runs `@orieken/saturday-ml-analyzer` on `heatmap-data/` produced during the qa-engineer's test run. Surfaces cold spots and coverage score.
2. **Checks visual regression** — looks for Playwright snapshot baselines and runs visual test tags to detect pixel-level regressions.

## Escalation Decision

No escalation needed. Saturday.ML (TS) is the correct foundation for the coverage analysis half. Playwright is the correct tool for the screenshot diffing half. The agent wraps both under one pipeline artifact.

## Pipeline Positioning

**Phase 3, after step 27 (validate-artifact for qa-report)**, before sre-engineer.

Rationale:
- Needs `qa-report.md` and `heatmap-data/` (produced when qa-engineer's test suite ran)
- Does NOT block sre-engineer, tech-writer, or devops — runs first in Phase 3's tail
- Runs AFTER code-reviewer+accessibility-engineer (Phase 2) — those are static analysis; this is execution-based

**Conditional**: Only invoked when (a) feature touches UI AND (b) `heatmap-data/` exists OR visual snapshot baselines exist. Missing heatmap data → UNCONFIGURED, not FAIL. Pipeline proceeds.

## Agent vs. Skill

Agent — produces a pipeline artifact (`visual-qa-report.md`) with a validated contract, participates in the handoff chain, and has conditional invocation via the delivery policy. Not a one-off operation that could be a skill.

## Contract Sections

```
## Summary
## Visual Regression
## Heatmap Coverage
## Cold Spots
## Recommendations
## Notes for QA
```

Validation: `## Summary` must contain `PASS`, `FAIL`, or `UNCONFIGURED`. `## Heatmap Coverage` must contain `Coverage Score:`. If summary is FAIL, `## Recommendations` must be non-empty.

## Phase B Plan

- B1: `shared/agents/visual-qa-engineer.md`
- B2: `shared/contracts/visual-qa-report-contract.md`
- B3: `shared/templates/visual-qa-report.template.md`
- B4: `shared/skills/deliver-feature/SKILL.md` — insert steps 28-29 (visual-qa-engineer + validate-artifact), renumber 28→30, 29→31, 30→32, 31→33, 32→34, 33→35
- B5: `docs/patterns/deliver-feature-workflow.md` — add visual-qa-engineer to diagram
- B6: `shared/skills/validate-artifact/SKILL.md` — add visual-qa-report-contract to mapping table
