# Eval Rubric: pattern-reviewer / input-pattern-doc.md

- **Stale class name is flagged with version**: `BasePage` is flagged as stale because it was renamed to `AbstractPage` in v2.1 — the auditor names the old class, the new class, and the version where the rename occurred.
- **Broken file path is flagged with the correct path**: `SiteManager.ts` is referenced at the old path (`src/saturday-core/src/managers/`) but was moved to `src/saturday-core/src/orchestration/` in v2.0 — the auditor names both paths.
- **Valid content is not false-flagged**: the `BaseSite` snippet and the `SiteManager` narrative description are not flagged — only the stale class name and broken path are findings.
- **Findings cite the exact location in the doc**: each finding references the specific line or section of the pattern doc where the stale content appears (not just "there's a stale snippet somewhere").
- **Audit suggests a remediation**: for each STALE finding, the output proposes a concrete fix (update the class name, update the file path) — not just "this is wrong."

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
