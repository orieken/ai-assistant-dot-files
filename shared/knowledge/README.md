# Knowledge Items (KIs)

Knowledge Items are searchable, reusable records of patterns, decisions, and fixes that `context-engineer`
surfaces during Proactive RAG (see `docs/runbooks/context-engineering.md`, Principle #2) — before the
analyst or developer re-derive something that's already been solved.

## Location convention
- `shared/knowledge/` — portable KIs that apply across any project using this framework (this directory).
- `.claude/knowledge/` — project-specific KIs, local to a single codebase, not portable.

## Format
Each KI is a single markdown file with frontmatter:

```markdown
---
name: kebab-case-slug
tags: [tag-one, tag-two]
domain: bounded-context-or-area
created: YYYY-MM-DD
---

Body: the pattern, decision, or fix — what it is, why it exists, when it applies.
```

## How KIs get used
1. `context-engineer` invokes the `search-ki` skill (tag/domain match against the active feature or task)
   during manifest creation rather than grepping this directory ad hoc.
2. Matches are listed in `context-manifest.md` under "4. Knowledge Items & ADRs (To Load)" with a reason
   they're relevant.
3. Downstream agents (analyst, developer) read the manifest and treat the referenced KI as authoritative —
   they should not re-solve a problem a KI already documents.

## Adding a new KI
Use the `create-ki` skill (`shared/skills/create-ki/SKILL.md`) — it searches for an existing KI on the same
topic first (via `search-ki`) so you update rather than duplicate, then writes the file in the format above.

