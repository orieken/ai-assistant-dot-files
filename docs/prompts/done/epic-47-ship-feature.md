# Epic 47 — Automated Release & PR Skill (`ship-feature`)

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 2 (ranked priority #3 there);
re-confirmed open and promoted to top remaining standalone epic by
`docs/audits/framework-gap-audit-2026-07-31.md` § F5.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

The delivery pipeline automates everything from spec to docs, then stops: no skill automates
branch creation, commit formatting, PR assembly, or release tagging. What exists today:

- `release-manager` **agent** — analyzes git history, determines semver bumps, drafts release
  notes. It *plans*; nothing *executes*.
- `shared/rules/approval-gates.md` — gate #2 (git commit requires "commit"/"approve commit"),
  gate #5 (external API mutations require "send"/"approve request"). PR creation via `gh` is an
  external mutation.
- `deliver-feature` persists artifacts to `docs/features/<feature-name>/` (spec, analysis,
  retrospective, qa-report) — the raw material for a PR description.
- Conventional Commits format is defined in root `CLAUDE.md` § Git Commits.

## Scope

Build `shared/skills/ship-feature/SKILL.md`. The skill:

1. **Branch**: creates `feat/<feature-name>` (or `fix/`) from main if not already on a feature
   branch; never commits to main directly.
2. **Commit**: stages explicit paths (never `git add -A`), formats a Conventional Commit, and
   **halts at gate #2** — presents the staged diff summary and waits for "commit".
3. **PR**: compiles the PR body from `docs/features/<feature-name>/` artifacts (spec summary,
   acceptance criteria coverage from qa-report.md, retrospective link), then **halts at gate #5**
   — presents the full PR title + body and waits for "send" before running `gh pr create`.
4. **Release** (optional `--release` mode): delegates version-bump analysis to the
   `release-manager` agent, then halts for approval before `git tag` + `gh release create`.

Constraints:

- Frontmatter must satisfy `shared/contracts/skill-frontmatter-contract.md` (name, description,
  triggers, standalone). `standalone: true` — must work outside a `deliver-feature` run (degrade
  gracefully when `docs/features/<name>/` artifacts don't exist: fall back to git-log-derived PR
  body).
- Every gate halt must state which gate it is and quote the exact approval word it's waiting for.
- No hardcoded remote/repo names — derive from `git remote` and `gh repo view`.

Commit sequence (one commit per op):

1. `feat(skills): add ship-feature skill (Epic 47)` — the SKILL.md
2. `docs(skills): wire ship-feature into deliver-feature as optional final phase (Epic 47)` —
   mention in `shared/skills/deliver-feature/SKILL.md` as an optional post-devops step, opt-in
3. `docs: update inventory counts for ship-feature skill (Epic 47)` — README / AGENT_REFERENCE
   skill counts (run `scripts/check-inventory-drift.sh` to find every stated count; it FAILs the
   health check if you miss one)

After each commit: `bash scripts/health-check.sh` must stay green (watch the skill-frontmatter
section) and `bash scripts/check-inventory-drift.sh` must pass.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- If gate semantics are ambiguous for any step (e.g., is `git tag` gate #2 or gate #5 territory?),
  halt and propose a gate mapping rather than inventing one — gates CANNOT be loosened by this
  skill.
- If `deliver-feature` integration would make ship-feature non-optional for existing users, halt —
  the backward-compat guarantee (new capability = opt-in) applies.
- If `gh` absence handling gets complicated, degrade to printing the `gh` command for the human to
  run, and say so in the report.

## Report (under 150 words)

```
Commits:
  <sha> <message>
  ...
Gate mapping implemented: branch=<none|gate>, commit=#2, pr=#5, tag/release=<mapping>
Standalone fallback (no feature artifacts): <how the PR body degrades>
health-check: <pass/fail>  check-inventory-drift: <pass/fail>
Open questions: <anything punted>
```

Go.
