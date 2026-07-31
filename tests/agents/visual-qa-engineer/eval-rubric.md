# Eval Rubric: visual-qa-engineer / input-qa-report.md

- **"Export CSV" cold spot is named and classified**: the agent identifies the bottom-right quadrant cold spot and specifically names the "Export CSV" button as the likely cause — not just reports "0 clicks in bottom-right region." Recommends either improving discoverability or removing the feature.
- **Mobile sparkline gap is flagged as a visual regression risk**: no mobile baseline exists for the three new sparkline components — the agent flags this as an uncovered risk and recommends generating baselines before shipping to production.
- **UNCONFIGURED is not used for present heatmap data**: `heatmap-data/` exists and contains data — the report must return PASS or FAIL, not UNCONFIGURED (UNCONFIGURED is only for projects that haven't instrumented Saturday ML).
- **Coverage Score is quantified**: the heatmap coverage section includes a `Coverage Score:` line (e.g., `Coverage Score: 73%`) with the method of calculation, per the visual-qa-report contract.
- **Three new snapshots are validated, not just acknowledged**: the agent either runs or recommends running `toHaveScreenshot()` against the three new sparkline baselines and reports the result — not just lists them as "added."

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
