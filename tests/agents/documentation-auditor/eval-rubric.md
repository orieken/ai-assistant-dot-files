# Eval Rubric: documentation-auditor / input-agent-reference.md

- **Removed agent still in table is flagged**: `legacy-formatter` was removed in v1.6 but still appears in the agent table — the auditor flags this as a stale entry.
- **Missing agents are enumerated**: `sre-engineer` (added v1.8) and `visual-qa-engineer` (added v1.7) are absent from the table — both are named specifically, not just noted as "some agents missing."
- **Deprecated skill is flagged**: `/epic-planner` is listed without the replacement (`/spec-writer`) being cross-referenced, or the deprecation notice being visible to readers — the auditor flags this as stale skill documentation.
- **No false positives on valid entries**: agents and skills that are current (analyst, developer, qa-engineer, `/deliver-feature`, `/create-ki`) are not flagged as stale.
- **Audit verdict is actionable**: findings include what needs to be updated (table rows to add/remove, deprecation note to add) — not just "this doc is outdated."

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
