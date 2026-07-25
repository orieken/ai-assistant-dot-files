# Framework hygiene sweep — 8 small pending items

Clean up the small pending items that accumulated during the AOS Phase 1 + frontmatter + template extraction sessions. Each is a small commit. None are individually big enough to justify their own handoff prompt, but leaving them uncommitted is a real drift risk.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo. Do NOT push; human's step.

## Scope: 8 items, one commit per item

### Item 1 — Commit the `.gitignore` modification

`.gitignore` has been modified with three deliberate additions never committed:
- `.env` — local secrets
- `scripts/publish-blog-posts.py` — local blog publisher (contains API keys)
- `docs/blog-posts/publishing.md` — publishing config

Verify the diff matches this intent (`git diff .gitignore`). If yes, commit AS-IS.

Commit: `chore(gitignore): ignore .env, blog publish script, and publishing config`.

### Item 2 — Track the blog drafts under `docs/blog-posts/`

The ENTIRE `docs/blog-posts/` directory is untracked — 8 blog draft files (01-04, both -devto.md and -linkedin.md variants) plus a README.md. These are real content that got authored but never committed.

Stage explicit paths for the 8 numbered drafts + README:
```
git add docs/blog-posts/README.md
git add docs/blog-posts/01-*.md
git add docs/blog-posts/02-*.md
git add docs/blog-posts/03-*.md
git add docs/blog-posts/04-*.md
```

Do NOT stage `publishing.md` (now gitignored per Item 1) or the `memory_engineering_prompts/` subdirectory or `memory_engineering_prompts.zip` if present — leave those for a separate decision.

Verify no strays via `git status --short docs/blog-posts/` before commit.

Commit: `docs(blog): track existing blog drafts 01-04 (dev.to + LinkedIn variants)`.

### Item 3 — Delete `docs/aos/AOS_Governance_Design_Pack.zip`

The zip is redundant now that the extracted markdowns landed in commit `54f0abe`. Options:
- **Delete**: `git rm docs/aos/AOS_Governance_Design_Pack.zip` if git-tracked; `rm` if untracked
- **Gitignore**: add `docs/aos/*.zip` to `.gitignore` (keeps the file locally, prevents accidental re-adds)

Recommend delete — the source zip is available in Oscar's original delivery folder if ever needed again, and having it in-tree confuses future readers about whether the zip or the extracted markdowns are canonical.

Commit: `chore(aos): remove redundant AOS_Governance_Design_Pack.zip`.

### Item 4 — Add `shared/templates/ki.template.md`

`docs/patterns/frontmatter-conventions.md`'s "Gaps and follow-ups" section explicitly flagged the missing KI template. Contracts + schemas have landed; this is the last piece of the frontmatter epic.

Template shape:
- Follow `shared/templates/agent.template.md`'s pattern (frontmatter block + illustrative body + reference table at the bottom)
- Match required fields from `shared/schemas/ki-frontmatter.schema.json`: `name` (kebab-case), `tags` (list, kebab-case each), `domain` (kebab-case), `created` (ISO date)
- Body: two headers `## Context` + `## Pattern` — matches the shape of existing KIs like `shared/knowledge/subagent-isolation-is-a-hard-boundary.md` (read it first for style)
- Include a footer noting the template's own `name: template-ki-do-not-use` so if it ever accidentally gets loaded it'll be obvious

Also update `docs/patterns/frontmatter-conventions.md`:
- Section 3 "Knowledge Item frontmatter" → "Template" subsection currently says "None yet — see the gap section below." Change to point at the new template
- "Gaps and follow-ups" → mark "KI template" as DONE with the new path

Commit: `docs(templates): add ki.template.md + close the frontmatter epic`.

### Item 5 — Adopt `docs/prompts/done/` convention

`saturday-mcp/docs/prompts/` uses a `done/` subdirectory for completed handoffs; `ai-assistant-dot-files/docs/prompts/` currently marks completed items with `— DONE` suffix in the README table but leaves the files in the active dir. Two files that should move:
- `docs/prompts/add-frontmatter-contracts.md` (shipped as `ae1e440`–`38b14a9`)
- `docs/prompts/add-frontmatter-json-schemas.md` (shipped as `fdfdd9b`–`a522e14`)

Do:
- `mkdir docs/prompts/done`
- `git mv docs/prompts/add-frontmatter-contracts.md docs/prompts/done/`
- `git mv docs/prompts/add-frontmatter-json-schemas.md docs/prompts/done/`
- Update `docs/prompts/README.md` to match saturday-mcp's pattern: split into "Active Prompts" and "Completed Prompts (`docs/prompts/done/`)" tables

Commit: `docs(prompts): adopt done/ convention for completed handoffs`.

### Item 6 — Investigate `docs/audits/perplex-audit.md`

Contents unknown. Halt and ask the user:
1. Is this a Perplexity AI audit output from an ad-hoc query?
2. Should it be tracked (committed as-is), moved (into `docs/lessons-learned/` if it's a learning), or gitignored (if it's transient)?

Do NOT commit until the user answers.

If user says "commit as-is": `docs(audits): track perplex-audit.md`.
If user says "gitignore": `chore(gitignore): ignore docs/audits/`.
If user says "move to lessons-learned": `docs(lessons): promote perplex audit findings` after moving.

### Item 7 — Fix `context-engineer 2.2.0` CHANGELOG WARN

`scripts/health-check.sh` reports:
```
WARN  context-engineer 2.2.0 — not found together in CHANGELOG.md (current version may be undocumented)
```

`shared/agents/context-engineer.md`'s frontmatter says `version: 2.2.0` but `shared/agents/CHANGELOG.md` doesn't have a paired entry documenting what changed at 2.2.0. Investigate git history for `context-engineer.md` to reconstruct what the 2.2.0 bump introduced, then add a CHANGELOG entry retroactively.

If the git log is unclear on what 2.2.0 changed, note the version but describe as "undocumented — bumped without CHANGELOG entry" — honest is better than fabricated.

Commit: `docs(changelog): retroactive entry for context-engineer 2.2.0`.

### Item 8 — Suppress WARN for optional paths in `scripts/health-check.sh`

Second pre-existing WARN:
```
WARN  registry path missing (marked optional): .claude/knowledge/
```

The path is deliberately marked `optional` in `shared/memory-registry.json` — the WARN is redundant. Update `scripts/health-check.sh` to skip the WARN when the path is in the registry's `optionalPaths` array. Report as INFO or silently skip.

Verify: after the fix, `scripts/health-check.sh` should report 0 WARNs (down from 2 pre-existing). Update `docs/aos/migration-plan.md`'s Phase 1 Op 1.7 verification note that the 2 WARNs are no longer pre-existing.

Commit: `fix(health-check): don't WARN on optional paths marked in registry`.

## Discipline (non-negotiable)

- **One commit per item** — 7-8 commits total (Item 6 may be zero if user says "gitignore" and the action is just a .gitignore edit that could fold into Item 1's commit — subagent's judgment call).
- Conventional Commits.
- **NEVER `git add -A`.** Explicit paths only.
- `git status --short` after staging, before every commit.
- `scripts/health-check.sh` must be green (0 FAILs, 0 WARNs after Items 7+8) at the end.
- Do NOT push.

## Escalation criteria

Stop and report if:
- Item 6's `perplex-audit.md` is unreadable (binary, corrupt) — halt with the specifics
- Item 7's git log for context-engineer.md is completely opaque (no commit that touched frontmatter version) — halt with the specifics
- Item 8's health-check refactor turns out to require a larger restructure than a single-file edit — halt, describe

## Report format (under 250 words)

```
STATUS: complete | stopped-at-item-N

Commits landed (order):
  <sha> <message>
  ...

Item-by-item tally:
  Item 1 gitignore commit: <landed>
  Item 2 blog drafts tracked: <count of files added>
  Item 3 zip deleted: <landed | user chose gitignore instead>
  Item 4 KI template: <landed>
  Item 5 done/ convention: <landed>
  Item 6 perplex-audit: <user decision + landed>
  Item 7 context-engineer CHANGELOG: <landed>
  Item 8 health-check optional-path WARN: <landed>

Health-check final state: <n> passed, <n> warned, <n> failed
  If WARNs remain: <list — should be 0>

Recommended next step:
  <e.g., "human review + git push">
```

Go.
