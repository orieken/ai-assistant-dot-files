# Eval Rubric: analyst / input-feature-spec.md

- **Three distinct edge cases, not one lumped together**: acceptance criteria or edge cases separately cover
  (1) non-enumeration of registered vs. unregistered email, (2) token expiry at 30 minutes, and (3) single-use
  enforcement — as three distinguishable items, not one vague "handle security concerns" bullet.
- **Bounded Context correctly identified**: names the owning context (something like "Identity & Access")
  and explicitly flags the Notifications capability as a Context Crossing, since the spec says it "does not
  exist yet."
- **New dependency flagged with a reliability angle**: the notification service is called out under "New
  Dependencies" and/or "Non-Functional Requirements" with at least a mention that it's a new external
  integration needing a resilience decision — not just listed as a bullet with no follow-through.
- **Feature Flag / Toggle strategy explicitly named**: per the agent's own mandatory rule, `analysis.md`
  states a concrete flag strategy for trunk-based delivery, not a generic "we'll use feature flags" without
  specifics.
- **Acceptance criteria are behavior-focused, not implementation-focused**: written as observable
  given/when/then outcomes (e.g., "given an unregistered email... the response is identical") rather than
  describing internal mechanics like specific function names or SQL.

## How to Grade
For each bullet above, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no
supporting quote, mark it FAIL and say what's missing.
