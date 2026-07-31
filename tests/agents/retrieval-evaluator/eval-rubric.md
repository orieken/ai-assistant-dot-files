# Eval Rubric: retrieval-evaluator / input-telemetry.md

- **All three zero-match queries are flagged**: "expand contract migration", "circuit breaker retry", and "WebSocket fanout multi-instance" each returned zero KI matches — the auditor identifies all three, not just one.
- **Missing-KI vs bad-metadata distinction is made**: for each zero-match query, the auditor proposes whether the gap is a missing KI (no KI exists on the topic) or bad metadata (a KI exists but tags don't surface it) — and states which it is based on the corpus snapshot.
- **`create-ki` recommendation is concrete**: for at least one zero-match query, the auditor proposes a specific KI title and suggested tags — not just "someone should create a KI."
- **Successful queries are not flagged**: the three queries that returned results (redis rate limiting, JWT validation, OTel span) are acknowledged as working — no false positives.
- **Priority is assigned**: findings are ordered by impact (e.g., expand-contract is a critical architectural pattern referenced in guardrails — its absence is higher priority than a WebSocket KI).

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
