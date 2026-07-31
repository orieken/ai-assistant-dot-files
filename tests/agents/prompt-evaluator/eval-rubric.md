# Eval Rubric: prompt-evaluator / input-prompt-draft.md

- **Hardcoded secret is a BLOCKER**: the connection string `postgres://admin:hunter2@localhost:5432/myapp_dev` contains a plaintext password — the auditor flags this as a BLOCKER finding and requires replacement with an environment variable placeholder.
- **Fabricated/internal URL is flagged**: `https://company-internal.example.com/db-standards` is an internal URL embedded in a prompt — the auditor flags this as undiscoverable by the agent at runtime and recommends removing or replacing with a relative path.
- **Tightly-coupled example is flagged**: the SQL example references `ADR-011` by number — this is a specific internal reference that will break when moved to other projects. The auditor recommends making examples generic.
- **`standalone` field is not relevant here** (this is an agent prompt, not a SKILL.md) — the auditor does not flag phantom issues beyond what's actually in the input.
- **Consistent voice is evaluated**: the prompt mixes imperative ("ALWAYS use", "NEVER use") with instructional prose — the auditor notes whether the voice is consistent throughout.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
