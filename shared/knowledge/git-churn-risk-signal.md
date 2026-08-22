---
name: git-churn-risk-signal
tags: [context-engineering, git, risk-signal, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

Git churn — the frequency with which a file has been modified — is a free, always-available
risk signal requiring no extra tooling. High-churn files are statistically more likely to
contain the current bug or conflict with the current change; stable files that rarely change
can be summarised rather than pinned in full, freeing token budget for higher-signal content.

## Integration point — context-engineer Step 4b (Prefer High-Fidelity References)

Run against the set of candidate files before deciding pin depth:

```bash
# Files modified more than N times in the last 90 days
git log --since="90 days ago" --diff-filter=M --name-only --pretty="" \
  -- <file1> <file2> ... | sort | uniq -c | sort -rn
```

**Decision rule:**
- Count ≥ 10 in 90 days → **pin in full** — high-churn; active development area, high risk
- Count 3–9 → **pin in full** — actively touched, worth the token cost
- Count 0–2 → **pin summary** — stable; extract the interface/signature only unless it is a
  direct dependency of the target component

Add a note in `context-manifest.md`:
```
## Churn Signal
high-churn (≥10 changes / 90 days): [file list]
stable (0-2 changes / 90 days): [file list] — summary-pinned
```

## Integration point — deliver-bugfix Phase 0

After running context-engineer, check churn on the buggy file and its immediate callers.
A file that has changed 15 times in 30 days is a different debugging context than one that
hasn't changed in a year.

```bash
git log --since="90 days ago" --diff-filter=M --name-only --pretty="" \
  -- <buggy-file> | wc -l
```

## No install required

This uses only `git log` — no extra dependency. Works in any git repository. The
`--diff-filter=M` flag restricts to modifications (not renames, deletions, or additions),
giving a cleaner signal for "how often is this file actively edited."

## Guardrails

- Churn reflects edit frequency, not correctness. A stable file can still be wrong;
  a high-churn file can still be the right abstraction.
- Initial commits (adding a file) count as additions (`A`), not modifications (`M`). A brand-
  new file with zero modifications is not the same as a proven stable file — check its age
  before treating it as low-risk.
- The 90-day window is a convention, not a hard rule. For longer-lived codebases with slow
  release cycles, extend to 180 days.

## See also

- `shared/skills/context-engineer/SKILL.md` — Step 4b (Prefer High-Fidelity References)
- `shared/knowledge/tokei-token-budget.md` — pairs with churn signal for Step 5 budget sizing
