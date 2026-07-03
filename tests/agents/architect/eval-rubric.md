# Eval Rubric: architect / input-analysis.md

- **Concrete resilience decision for the new email dependency**: names a specific pattern (Circuit Breaker,
  Timeout, and/or Retry with backoff) for the new outbound email provider call, not just "should be
  reliable" without a mechanism.
- **Concrete integration pattern chosen for the Context Crossing**: explicitly picks direct call vs.
  event-driven (or another named pattern) for Identity & Access calling Notifications, with a stated reason
  — not left as an open question when the analysis already asked for this decision.
- **Fitness function is enforceable, not vague**: the fitness function tied to this crossing/dependency
  names an actual verification mechanism (a specific lint rule, a timeout config check, a circuit-breaker
  test) rather than "should be monitored" or "keep an eye on this."
- **RFC written**: since this is a first-time Context Crossing introducing a brand-new external dependency
  (touches a layer boundary, per the agent's own RFC-trigger rule), a lightweight RFC file is produced
  alongside `architecture-notes.md`, not skipped.
- **Team Topology Fit addressed**: the "Bounded Context" section's Team Topology Fit line is filled in (even
  if the answer is "no TEAM_TOPOLOGY.md row for this context yet") rather than left as a blank placeholder.

## How to Grade
For each bullet above, quote the specific line(s) of `actual-output.md` (and the RFC file, if produced) that
satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
