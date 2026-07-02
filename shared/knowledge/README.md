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
1. `context-engineer` searches this directory (and `.claude/knowledge/`) by tag/domain match against the
   active feature or task.
2. Matches are listed in `context-manifest.md` under "Relevant Knowledge Items (KIs) & ADRs" with a reason
   they're relevant.
3. Downstream agents (analyst, developer) read the manifest and treat the referenced KI as authoritative —
   they should not re-solve a problem a KI already documents.

## Adding a new KI
There's no dedicated authoring skill yet (tracked as an open item — see
`docs/features/context-engineering-framework/TODO.md`, Epic 14). For now, add a markdown file directly to
this directory following the format above.
