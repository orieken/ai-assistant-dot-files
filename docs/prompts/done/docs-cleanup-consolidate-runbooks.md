# Docs Cleanup — Consolidate docs/RUNBOOKS.md into docs/runbooks/README.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/RUNBOOKS.md` is a thin summary document that only covers `debug-environment` and
`debug-tests`. The real operational index has grown in `docs/runbooks/` which is the authoritative
directory. `RUNBOOKS.md` is redundant and should be consolidated into `docs/runbooks/README.md`
or replaced with a short redirect.

Known consumers to update:
- `docs/README.md` — may reference `RUNBOOKS.md`
- `shared/agents/documentation-manager.md` — may reference `RUNBOOKS.md`

## Scope

**Op 1 — Diff the two files:**

Read `docs/RUNBOOKS.md` and `docs/runbooks/README.md` in full. Identify any content in
`docs/RUNBOOKS.md` not already covered in `docs/runbooks/README.md`.

**Op 2 — Merge any missing content:**

If `docs/RUNBOOKS.md` contains content absent from `docs/runbooks/README.md`, add it to the
README under an appropriate heading.

**Op 3 — Replace or delete docs/RUNBOOKS.md:**

Option A (preferred if any external tooling might still reference the path):
Replace `docs/RUNBOOKS.md` with a one-paragraph redirect:
```markdown
# Runbooks

Operational runbooks have moved to [`docs/runbooks/`](runbooks/README.md).
```

Option B (if no external references are found):
Delete it: `git rm docs/RUNBOOKS.md`

Before choosing, grep for references:
```bash
grep -rn "RUNBOOKS.md\|docs/RUNBOOKS" . --include="*.md" --include="*.sh" --include="*.yaml" --include="*.json"
```

Use Option B only if grep returns no references outside `docs/TODO.md` itself.

**Op 4 — Update known consumers:**

- Read `docs/README.md` — if it references `RUNBOOKS.md`, update the link to `runbooks/README.md`.
- Read `shared/agents/documentation-manager.md` — if it references `RUNBOOKS.md`, update accordingly.

**Op 5 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Consolidate `docs/RUNBOOKS.md`" item as `[x]`.

## Guardrails

- Conventional commit: `docs(runbooks): consolidate RUNBOOKS.md into docs/runbooks/README.md`
- Stage only the files you modified — do not `git add -A`.
- Do not rewrite the runbooks content itself; this is structure and reference cleanup only.

## Escalation

Stop and report if `docs/RUNBOOKS.md` contains significant operational content absent from
`docs/runbooks/README.md` that you are unsure how to categorize — flag for human review.

## Report

On completion, confirm:
- Whether any content was merged from `RUNBOOKS.md` into `docs/runbooks/README.md`
- Which option was chosen (redirect vs. delete) and why
- Which references were updated
- Commit hash
