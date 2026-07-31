# Eval Rubric: tech-writer / input-qa-report.md

- **StorageAdapter is documented as a new interface**: the docs-report explicitly notes the new `StorageAdapter` interface and either updates API docs or flags it as a docs task — not silently skipped.
- **New env vars appear in the right place**: `AWS_S3_BUCKET` and `AWS_REGION` are added to the environment configuration section of the README/runbook (not just noted as "needed"), with example values.
- **S3 lifecycle rule gap is surfaced**: the agent picks up the QA note that a 30-day lifecycle rule is recommended but not yet implemented and either opens a chore or includes it in the ops runbook as a manual step.
- **CHANGELOG entry is specific**: the CHANGELOG entry names the feature ("user avatar upload"), the new interface, and the new env vars — not a generic "added new feature" line.
- **ADR evaluated (even if not required)**: the agent considers whether the decision to hide S3 behind `StorageAdapter` warrants an ADR entry, and either creates one or explicitly decides it doesn't rise to the threshold with a brief reason.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
