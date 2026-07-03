# Eval Rubric: security-reviewer / input-vulnerable-code.ts

Qualitative criteria beyond `expected-patterns.txt`'s keyword check — judges whether the *reasoning* is
correct, not just whether the right words appear.

- **Correct severity on SQL injection**: classified as CRITICAL or HIGH, not Medium/Low — this is a
  pre-auth, unauthenticated injection point.
- **Specific SQL injection fix**: recommends parameterized queries / prepared statements by name, not a
  generic "sanitize the input" or "validate the input" instruction with no concrete mechanism.
- **Correct severity on the hardcoded Stripe key**: classified as CRITICAL or HIGH (a live secret key
  committed to source), and the fix names moving it to an environment variable or secrets vault — not just
  "don't hardcode secrets" without a concrete remediation.
- **User enumeration reasoning is correct**: identifies *why* the differing 404 ("Username not found") vs.
  401 ("Incorrect password") responses leak account existence, and recommends unifying both into a single
  generic response (e.g., "Invalid username or password") — not just flagging that two error codes exist.
- **Plaintext password logging is flagged as its own distinct finding**, not folded into the enumeration or
  injection findings — the `console.log` line leaks the raw password itself, a separate exposure.
- **No false negative on the response body**: notes that `stripeKey` is returned directly in the JSON
  response to the client, not just that it's hardcoded in source — this is a second exposure of the same
  secret (returned to any caller, not just visible to whoever reads the repo).

## How to Grade
For each bullet above, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no
supporting quote, mark it FAIL and say what's missing.
