# Eval Rubric: privacy-auditor / input-pipeline-artifact.md

- **Hardcoded API key is a BLOCKER**: `sk_test_4eC39HqLyjWDarjtT1zdp7dc` appearing in the implementation notes is flagged as a BLOCKER (not just a warning) — even test keys must not appear in committed artifacts because they train engineers to treat secrets as documentation.
- **PII logging is flagged**: `customer_email` and `customer_name` logged at INFO level is identified as a PII data boundary leak — the auditor names the specific fields and the logger call.
- **Remediation requires both removal and rotation**: the output does not stop at "remove it" — it also states the key should be rotated (since it appeared in a committed artifact).
- **Correct env var usage is acknowledged**: `STRIPE_WEBHOOK_SECRET` is correctly read from `process.env` — the auditor notes this as correct practice, not a finding.
- **Hardcoded key in config file is also caught**: the auditor flags `src/config/stripe.config.ts` in addition to the implementation notes — both locations where the key appears are called out.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
