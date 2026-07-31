# Eval Rubric: model-tier-auditor / input-agent-frontmatter.md

- **Missing `model_tier` is flagged**: `context-engineer.md` has no `model_tier` field — the auditor flags this as a FAIL, naming the agent explicitly.
- **Invalid enum value is rejected**: `sre-engineer.md` uses `model_tier: observability`, which is not a valid enum value — the auditor reports this as invalid and names the expected valid values.
- **Profile mismatch is reported as a WARN, not FAIL**: `chaos-engineer.md` has `model_tier: standard` but an operational profile suggesting a higher tier — the auditor flags this as a mismatch/recommendation without hard failing it (since it's a judgment call, not a schema violation).
- **Valid entries are not false-flagged**: `analyst` (reasoning), `developer` (standard), and `agent-evaluator` (budget) all have valid model_tier values and are reported as PASS.
- **Audit findings state the expected vs actual**: for each FAIL finding, the output includes what value was found and what was expected (e.g., "found: observability, expected one of: budget | standard | reasoning").

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
