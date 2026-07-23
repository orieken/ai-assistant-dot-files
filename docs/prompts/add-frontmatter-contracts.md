# Add frontmatter contracts for agents / skills / KIs

Create formal contract files so `validate-artifact` can grep-check agent, skill, and Knowledge Item frontmatter the same way it checks pipeline artifacts (`analysis.md`, `architecture-notes.md`, etc.). Currently `scripts/health-check.sh` field-presence-checks frontmatter but there's no contract file the validator can point to.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo. Commits go here directly.

## Prior context

The framework has `shared/contracts/` with 13 contract files for pipeline artifacts (analysis-contract.md, architecture-contract.md, etc.). Each declares required section headings, and `validate-artifact` grep-checks them. Zero contract exists for frontmatter today.

Field-presence enforcement lives in `scripts/health-check.sh` (steps 2, 3, 8 for agents/skills/KIs respectively). That check is functional but the shape is implicit rather than referenceable — contributors have to read the script to know what fields are required.

The consolidated reference is `docs/patterns/frontmatter-conventions.md` — read this first for the exact required-fields tables.

## Scope

### Op 1: create `shared/contracts/agent-frontmatter-contract.md`

Match the shape of existing contract files (`shared/contracts/analysis-contract.md` is the reference). Sections:
- `**Produced by**: humans authoring agent files under shared/agents/`
- `**Consumed by**: install.sh (symlink), Claude Code / Cursor loaders, orchestrator agents citing by role`
- `## Required Fields` — table with field name, type, notes
- `## Validation Rule` — describe what validate-artifact should check (field presence, and for `version`, valid semver)

Content sourced verbatim from `docs/patterns/frontmatter-conventions.md` section 1.

### Op 2: create `shared/contracts/skill-frontmatter-contract.md`

Same shape, content from `docs/patterns/frontmatter-conventions.md` section 2. Includes the `triggers` sub-structure requirement.

### Op 3: create `shared/contracts/ki-frontmatter-contract.md`

Same shape, content from `docs/patterns/frontmatter-conventions.md` section 3. Includes the `created` ISO-date format rule.

### Op 4: update `shared/skills/validate-artifact/SKILL.md`

Extend the Contract Mapping table with three new rows:

```
| humans (agent author) | shared/agents/*.md | shared/contracts/agent-frontmatter-contract.md |
| humans (skill author) | shared/skills/*/SKILL.md | shared/contracts/skill-frontmatter-contract.md |
| humans (KI author) | shared/knowledge/*.md, .claude/knowledge/*.md | shared/contracts/ki-frontmatter-contract.md |
```

Update the skill's Process section to note that validate-artifact can be invoked on frontmatter files (not just pipeline artifacts).

Bump the skill's version frontmatter (if it has one).

### Op 5: update `docs/patterns/frontmatter-conventions.md`

In the "Gaps and follow-ups" section, mark the "Formal contracts" bullet as done and cross-reference the three new contract files.

## Discipline

- **One commit per op preferred** (5 commits total) — OR one commit for Ops 1-3 (all three contracts are one epic) + one for Op 4 + one for Op 5. Either shape is fine.
- Conventional Commits: `feat(contracts): add agent-frontmatter-contract` for the contracts, `feat(validate-artifact): support frontmatter contracts` for Op 4, `docs(patterns): mark frontmatter contracts as done` for Op 5.
- **NEVER `git add -A`** — pre-existing untracked directories exist (`docs/audits/`, `docs/blog-posts/`, `docs/aos/AOS_Governance_Design_Pack.zip`). Stage explicit paths only.
- Run `scripts/health-check.sh` after Op 4 to verify nothing regressed.

## Escalation criteria

Stop and report if:
- An existing agent, skill, or KI would fail validate-artifact under the new contract (means the contract is stricter than reality — either loosen the contract or fix the offender before shipping)
- The validate-artifact contract mapping structure doesn't cleanly accept "frontmatter" as an artifact type (may need a small refactor first — flag it)

## Report format (under 200 words)

```
STATUS: complete | stopped-at-op-N
Commits: <sha> <message>, ...

Files created:
- shared/contracts/agent-frontmatter-contract.md
- shared/contracts/skill-frontmatter-contract.md
- shared/contracts/ki-frontmatter-contract.md

Files updated:
- shared/skills/validate-artifact/SKILL.md
- docs/patterns/frontmatter-conventions.md

Health-check status: pass / <details>

Any existing files that would newly fail under the contracts (should be zero):
- <list, or "None">
```

Go.
