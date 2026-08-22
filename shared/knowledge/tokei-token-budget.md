---
name: tokei-token-budget
tags: [context-engineering, token-budget, cli-tools, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

tokei (`github.com/XAMPPRocky/tokei`) is a fast, accurate line-count CLI. It reports code
lines, comment lines, and blank lines per file and per language — and outputs JSON. This
replaces context-engineer's rough heuristic (`total lines × 8 ÷ 4`) with a measurement that
matches how LLMs actually consume tokens: comment-heavy files consume significantly more tokens
per line than terse code.

## Integration point — context-engineer Step 5 (Estimate Token Budget)

Instead of the heuristic, run:
```bash
tokei --output json <file1> <file2> ...
```

From the JSON output, use the `code` field (not `lines` or `total`) as the line count for
token estimation — comment lines average ~2× the tokens of code lines, so using code lines
with the standard multiplier (`lines × 4.5 tokens`) gives a tighter estimate for
mixed-comment files.

If tokei is not installed, fall back to the existing heuristic and note:
```
context-debt: tokei not installed — token estimate uses total-line heuristic
```

## Installation

```bash
brew install tokei
# or
cargo install tokei
```

## Why not `wc -l`

`wc -l` counts total lines including blanks and comments. tokei separates them. For a file
that is 40% comments (common in well-documented code), `wc -l`-based estimates overshoot by
~30%. The error compounds when many files are pinned.

## Guardrails

- tokei measures the file as it exists on disk — not what will be sent after any preprocessing
  (stripping comments, truncating). The estimate is an upper bound.
- JSON output key is `code`, not `lines` — don't confuse them.

## See also

- `shared/skills/context-engineer/SKILL.md` — Step 5 (Estimate Token Budget) is the primary
  integration point
- `shared/knowledge/git-churn-risk-signal.md` — pairs with tokei for Step 4b prioritization
