# Eval Rubric: rule-auditor / input-rules-excerpt.md

- **`DOMAIN_DICTIONARY.md` path inconsistency is flagged**: `go-conventions.md` uses `docs/DOMAIN_DICTIONARY.md` while `design-principles.md` uses `DOMAIN_DICTIONARY.md` — the auditor identifies both references and flags the path inconsistency (they may resolve to the same file or may not).
- **`iac-conventions.md` un-indexed status is flagged**: the auditor identifies that `iac-conventions.md` exists on disk but is not referenced in any rules index or load order — this is a dead-path reference risk.
- **Redundant complexity rule is noted**: the `< 7` cyclomatic complexity rule appears in both `testing-conventions.md`, `design-principles.md`, and `CLAUDE.md` — the auditor flags this as redundant but not contradictory (same value, so consistent).
- **No false contradiction raised on `no raw any`**: the `any` prohibition appears in both `typescript-conventions.md` and `architecture-guardrails.md` — the auditor notes the duplication but does not flag a contradiction because both say the same thing.
- **Findings distinguish contradiction vs redundancy**: the auditor uses different severity labels for outright contradictions (FAIL) vs redundant-but-consistent rules (INFO or NOTE).

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
