# Eval Rubric: agent-evaluator / input-eval-case.md

- **Each rubric criterion is graded with a quote**: for every item in `eval-rubric.md`, the evaluator cites the specific line from `actual-output.md` that satisfies (or fails to satisfy) it — not a summary judgment like "the output was good."
- **Pattern check is mechanistic**: the evaluator reports each `expected-patterns.txt` entry as PASS or FAIL based on grep-equivalent matching — no subjective interpretation of whether the keyword "feels present."
- **Overall verdict is distinct from per-criterion verdicts**: the output has a clear top-level PASS/FAIL/REGRESSION summary separate from the detailed per-criterion breakdown.
- **No false positives on the rubric**: if a rubric criterion is partially met (e.g., edge cases are mentioned but not as three separate items), the evaluator marks it FAIL with the specific gap — not a charitable PASS.
- **Regression framing is used when applicable**: if this is a re-run after an agent prompt change, the evaluator compares to the previous eval record and explicitly calls out any criteria that regressed.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
