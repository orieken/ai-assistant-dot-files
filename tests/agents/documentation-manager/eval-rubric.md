# Eval Rubric: documentation-manager / input-session.md

- **ADR candidate is proposed for the ws vs Socket.IO decision**: the agent recommends creating an ADR (not just a comment) for the raw `ws` package choice — this is a structural decision affecting future maintainers.
- **KI candidate is proposed for the EventEmitter fanout approach**: the single-instance limitation and the known multi-instance breakage risk is proposed as a KI with the specific risk documented, not buried in a comment.
- **OTel gap is flagged as a guardrail violation**: the missing OpenTelemetry spans are identified as violating `architecture-guardrails.md §8` (OTel not in domain layer, but it must be in the adapter layer) — not just noted as "will add later."
- **Hardcoded heartbeat is flagged**: the 30-second constant is called out as a magic number violation requiring extraction to a named constant or config.
- **Candidate records are framed for human review — not auto-applied**: the agent produces proposals (ADR Candidate, KI Candidate) and does NOT claim to have written the ADR or KI. The output makes clear these require human approval.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
