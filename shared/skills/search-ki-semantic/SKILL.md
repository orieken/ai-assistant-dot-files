---
name: search-ki-semantic
description: Semantic search over the framework KI corpus using LLM-as-retriever (AOS Phase 3). Loads the full KI + ADR corpus into context and judges relevance holistically — catches paraphrases and conceptual matches that lexical search-ki misses. Use when search-ki returns nothing for a concept you believe should be documented.
triggers:
  keywords: ["search-ki-semantic", "semantic ki search", "find knowledge semantically"]
  intentPatterns: ["Semantic search KIs for *", "/search-ki-semantic *"]
standalone: true
---

## When To Use

Use when:
- `search-ki` returned empty results but the concept feels like it should be documented.
- The query is conceptual or paraphrase-heavy ("how do we handle cascading failures?" vs "circuit-breaker KI").
- You want a holistic relevance judgment across the full corpus, not a tag/keyword pre-filter.

Do NOT use when:
- You know the exact KI name or tag — use `search-ki` directly (faster).
- The corpus exceeds 200 KIs — check `shared/memory-registry.json` item count first. Above 200, context budget may be exhausted; consider `search-ki` with explicit tags instead.
- Searching application source code — use `Grep`/`Glob`.

## Adapter

Implements the `llm-as-retriever` adapter shape documented in
[`shared/rag/adapters/llm-as-retriever.md`](../../rag/adapters/llm-as-retriever.md) for the
`framework-ki` corpus.

## Context To Load First

1. `shared/memory-registry.json` — verify `framework-ki` source has `retrievalBackend: "llm-as-retriever"` and item count ≤ 200
2. All KI files: `shared/knowledge/*.md`, `.claude/knowledge/*.md` (if directory exists)
3. All ADR files: `docs/adrs/*.md`

## Process

1. **Load registry**: read `shared/memory-registry.json`. If `framework-ki` source is missing or `retrievalBackend` ≠ `"llm-as-retriever"`, halt and tell the user the registry is misconfigured.

2. **Check corpus size**: count KI + ADR files. If > 200, warn the user that context budget may be tight and suggest falling back to `/search-ki` with explicit tag filters. Proceed regardless.

3. **Load full corpus**: read every KI and ADR file in parallel using Glob + Read. Do not pre-filter.

4. **Judge relevance**: for each item, determine:
   - `score`: 1.0 = directly answers the query | 0.7 = partial / one dimension of query | 0.4 = tangentially related | 0.0 = not relevant
   - `excerpt`: 1-2 sentences explaining *why* this item matches (or note `not relevant`)
   - Discard items with score < 0.4

5. **Sort and cap**: return top 10 items sorted by descending score.

6. **Format output**:

```
## Semantic KI Search: "<query>"

### Results (N found, top 10 shown)

1. [<title>](<path>) — score: 0.9
   <excerpt>

2. [<title>](<path>) — score: 0.7
   <excerpt>

...

### No match found (if 0 results)
No KIs or ADRs matched "<query>" at the 0.4 relevance threshold.
Consider: (a) the concept may not be documented yet — run /create-ki to document it; (b) try /search-ki with a different keyword; (c) check DOMAIN_DICTIONARY.md for canonical terminology.
```

## Guardrails

- Never emit content copies — only `path`, `title`, `excerpt` (≤ 2 sentences), and `score`.
- If the corpus load hits a context-budget error, truncate to the first 150 items alphabetically and note the truncation in the output.
- `search-ki` (lexical) remains the default invocation path for `context-engineer`'s Proactive RAG step and `query-memory`'s KI delegation. This skill is the opt-in semantic alternative — it does not replace lexical search.

## Standalone Mode

Operates purely using local file reading tools (Read, Glob). No network calls, no external index.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 3 Runtime layer. CC BY 4.0.*
