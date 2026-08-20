# Docs Cleanup — Consolidate docs/CLAUDE.md into Root CLAUDE.md

Source: `docs/TODO.md` §"docs/ Root".

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Context

`docs/CLAUDE.md` is a divergent copy of the root `CLAUDE.md`. It has not been updated since
2026-03-01 and is indexed only as a reference document with no active consumers. Before deletion,
any unique content (notably a "Legacy Code & Refactoring" section referencing Michael Feathers
and characterization tests) must be reviewed and merged into the canonical root `CLAUDE.md` if
still authoritative.

## Scope

**Op 1 — Diff the two files:**

Read both `docs/CLAUDE.md` and `CLAUDE.md` in full. Identify all content in `docs/CLAUDE.md`
that is either:
- Absent from root `CLAUDE.md`, OR
- Present but with materially different/richer wording

Specifically check for the "Legacy Code & Refactoring" section — it is known to differ.

**Op 2 — Merge unique authoritative content:**

For each unique or richer section found:
- If it is still accurate and belongs in the canonical rules, add it to `CLAUDE.md` under the
  most appropriate existing heading (or a new heading if none fits).
- If it is outdated or superseded by content already in `shared/rules/`, note it in the report
  and skip it.

Do not blindly append — integrate thoughtfully, preserving the root file's structure.

**Op 3 — Delete docs/CLAUDE.md:**

After the merge is complete and committed:
```bash
git rm docs/CLAUDE.md
```

**Op 4 — Remove all references to docs/CLAUDE.md:**

Grep for references:
```bash
grep -rn "docs/CLAUDE.md" .
```

Update or remove any found references (README, index files, skill files, etc.).

**Op 5 — Mark TODO item resolved:**

In `docs/TODO.md`, mark the "Delete or consolidate `docs/CLAUDE.md`" item as `[x]`.

## Guardrails

- Commit the merge (Op 2) and the deletion (Op 3–4) as a single commit or two sequential commits.
- Conventional commit: `docs(cleanup): consolidate docs/CLAUDE.md into root CLAUDE.md`
- Do not introduce new rules or change the meaning of existing root `CLAUDE.md` content —
  this is a merge and delete, not a rewrite.
- Stage files explicitly — do not `git add -A`.

## Escalation

Stop and report if:
- `docs/CLAUDE.md` contains substantive content whose current accuracy you cannot determine
  (flag it for human review before merging).
- Any reference to `docs/CLAUDE.md` is found in a skill, agent, or generated config — those
  consumers may need a separate update.

## Report

On completion, confirm:
- Which sections (if any) were merged from `docs/CLAUDE.md` into root `CLAUDE.md`
- Which sections were skipped and why
- Whether any references to the old path were found and updated
- Commit hash(es)
