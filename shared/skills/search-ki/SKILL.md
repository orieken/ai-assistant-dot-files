---
name: search-ki
description: Searches Knowledge Items by tag, domain, or keyword across shared/knowledge/ (portable) and .claude/knowledge/ (project-specific) before an agent does independent analysis. The mechanism context-engineer's Proactive RAG step and analyst's feedback-loop step both call into.
triggers:
  keywords: ["search-ki", "find knowledge item", "has this been solved before", "search knowledge"]
  intentPatterns: ["Search KIs for *", "Has * been solved before", "/search-ki *"]
standalone: true
---

## When To Use
- `context-engineer` invokes this during manifest creation (Proactive RAG, see `shared/agents/context-engineer.md` step 5) — before letting the analyst reason independently, check whether the pattern is already documented.
- Any agent or human can invoke it directly: "has anyone solved rate-limiting on the login endpoint before?"

Do NOT use to search application source code — use `Grep`/`Glob` directly for that. This only searches the
KI corpus (`shared/knowledge/`, `.claude/knowledge/`) and `docs/adrs/`.

## Context To Load First
1. `shared/knowledge/*.md` and `.claude/knowledge/*.md` (if the latter directory exists)
2. `docs/adrs/*.md`

## Process
1. **Parse the query** into candidate tags/domain/keywords. If the caller gave explicit tags (e.g. from
   context-engineer's bounded-context mapping), use those directly; otherwise extract likely tags from the
   free-text question.
2. **Scan KI frontmatter** in both `shared/knowledge/` and `.claude/knowledge/` for `tags:` or `domain:`
   matches. A KI matches if any candidate tag/domain overlaps.
3. **Fall back to full-text grep** across KI bodies (not just frontmatter) if the tag/domain scan finds
   nothing — a KI's body may use different wording than its own tags.
4. **Scan `docs/adrs/`** the same way — ADRs aren't KIs but often answer "why did we choose X" questions a
   KI search is really asking.
5. **Rank results**: exact tag match > domain match > full-text body match. Cap at the 5 most relevant.

## Output Format
```markdown
## KI Search: "[query]"

### Matches
- [KI Name](shared/knowledge/<file>.md) -- [tags] -- [one-line why it matches]
- [ADR Name](docs/adrs/<file>.md) -- [one-line why it matches]

### No Match
[If nothing matches: "No existing KI or ADR covers this — consider running create-ki after this task if the
solution turns out to be reusable."]
```

## Guardrails
- **Never** fabricate a match — if nothing in the corpus is actually relevant, say so plainly rather than
  stretching a tenuous connection.
- **Read-only**: this skill never writes to `shared/knowledge/` or `.claude/knowledge/` — use `create-ki`
  for that.
- Cap results at 5 — a long list defeats the point of a high-signal search.

## Standalone Mode
Pure local file reads (frontmatter parsing + grep). No external calls.
