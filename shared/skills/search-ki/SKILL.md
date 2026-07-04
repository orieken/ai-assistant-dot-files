---
name: search-ki
description: Searches Knowledge Items by tag, domain, or keyword across shared/knowledge/ (portable), .claude/knowledge/ (project-specific), and docs/adrs/ (why a past decision was made) before an agent does independent analysis. The mechanism context-engineer's Proactive RAG step and analyst's feedback-loop step both call into.
triggers:
  keywords: ["search-ki", "find knowledge item", "has this been solved before", "search knowledge"]
  intentPatterns: ["Search KIs for *", "Has * been solved before", "/search-ki *"]
standalone: true
---

## When To Use
- `context-engineer` invokes this during manifest creation (Proactive RAG, see `shared/agents/context-engineer.md` step 5) — before letting the analyst reason independently, check whether the pattern is already documented.
- Any agent or human can invoke it directly: "has anyone solved rate-limiting on the login endpoint before?"

Do NOT use to search application source code — use `Grep`/`Glob` directly for that. This only searches the
KI corpus (`shared/knowledge/`, `.claude/knowledge/`) and `docs/adrs/`. Do NOT use when the question could be
answered by the feature archive or `DOMAIN_DICTIONARY.md` too, not just KIs/ADRs — use `query-memory` for a
search across every registered memory source; it delegates the KI/ADR portion to this skill, so nothing here
needs to change to support that.

## Context To Load First
1. `shared/memory-registry.json` — confirms the KI/ADR retrieval backend is `lexical` (no embeddings) before
   proceeding; this skill is that backend's implementation
2. `shared/knowledge/*.md` and `.claude/knowledge/*.md` (if the latter directory exists)
3. `docs/adrs/*.md`

## Process
This is judgment-based semantic search, not a grep — you're an LLM reading the corpus, not a regex engine.
Tag/domain matching below is a *pre-filter* to keep the read-through cheap as the corpus grows, never the
final relevance decision. A KI can be the right match even when it uses none of the query's words.

1. **Parse the query** into candidate tags/domain/keywords. If the caller gave explicit tags (e.g. from
   context-engineer's bounded-context mapping), use those directly; otherwise extract likely tags from the
   free-text question.
2. **Pre-filter cheaply**: scan KI frontmatter in `shared/knowledge/` and `.claude/knowledge/` for `tags:`/
   `domain:` overlap with the candidates from step 1. If the corpus is small (roughly under ~30 KIs total),
   skip this step and just read everything — the pre-filter exists to avoid reading a large corpus in full,
   not to replace judgment on a small one.
3. **Read the candidate KIs' full bodies** (not just frontmatter) and judge relevance by actual meaning —
   a KI titled `subagent-isolation-is-a-hard-boundary` is the right answer to "why can't my agent see what
   the last one did" even though neither phrase appears in the query. This is the step a plain grep can't
   do; don't skip straight from tag-matching to "no match."
4. **Read `docs/adrs/` the same way** — ADRs aren't KIs but often answer "why did we choose X" questions a
   KI search is really asking.
5. **Rank by actual relevance to the question**, not by tag-match count — a KI with one matching tag but
   directly on-point beats one with three matching tags that's tangential. Cap at the 5 most relevant.

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
  stretching a tenuous connection just because a tag happened to overlap.
- **Read-only**: this skill never writes to `shared/knowledge/` or `.claude/knowledge/` — use `create-ki`
  for that.
- Cap results at 5 — a long list defeats the point of a high-signal search.
- **No embeddings/vector search, by design**: this framework's skills keep the intelligence in the markdown
  and treat API calls as plumbing (see `SKILL_TEMPLATE.md`'s Standalone Mode guidance) — adding an
  embeddings dependency for a corpus this size would trade a real architectural change for a problem lexical
  search doesn't actually have yet. If the KI corpus grows into the hundreds and pre-filter-then-read stops
  being affordable, that tradeoff is worth revisiting then, not now. `shared/memory-registry.json` tracks
  this decision explicitly — its `lightrag` backend entry stays `"status": "disabled"` until that day comes
  (see `docs/runbooks/lightrag-integration.md`).
- **Duplicate/stale KIs are `memory-engineer`'s job, not this skill's** — if search results feel noisy or
  redundant, that's a signal to run a memory sweep, not a reason to change how this skill ranks results.

## Standalone Mode
Pure local file reads (frontmatter parsing for the pre-filter, full reads for judgment). No external calls,
no embeddings API, no vector index to keep in sync.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
