# Adapter: LLM-as-Retriever

**Corpus**: `framework-ki` (Knowledge Items + ADRs)
**Files**: `shared/knowledge/*.md`, `.claude/knowledge/*.md`, `docs/adrs/*.md`

## Rationale

The framework's KI/ADR corpus is small (~30-200 items) and fits in a context window. Reading it in full
and letting the LLM reason about relevance is simpler, more accurate, and cheaper to maintain than
maintaining an embedding index. No index drift, no embedding cost, no vector store dependency.

This is the "simplest thing that works" per Kent Beck's Simple Design rule: passes the tests (retrieves
relevant KIs), reveals intention (algorithm is obvious from reading this file), no duplication (reuses
the LLM's existing reading capacity), fewest elements (zero infrastructure dependencies).

## Implementation

The `search-ki` skill is this adapter's reference implementation (lexical mode). The `search-ki-semantic`
skill (added in Op 3.2) is the same adapter operating in semantic mode — it loads the full corpus into
context and lets the LLM judge relevance holistically, rather than pre-filtering by tag/keyword.

### Algorithm (semantic mode)

1. Load `shared/memory-registry.json` — confirm the `framework-ki` source entry has `retrievalBackend: "llm-as-retriever"`.
2. Load all KI files from `shared/knowledge/` and `.claude/knowledge/` (if present) and `docs/adrs/`.
3. For each item, output a relevance judgment (relevant | not relevant) and a 1-2 sentence excerpt explaining why it matches (or doesn't). Scoring: 1.0 = directly answers the query, 0.7 = partial relevance, 0.4 = tangentially related, 0.0 = not relevant.
4. Return the top 10 items with score ≥ 0.4, sorted descending by score.

### Algorithm (lexical mode — `search-ki` default)

1. Tag/domain pre-filter: compare query terms against KI `tags:` frontmatter to reduce the read set.
2. Full-content read of the pre-filtered set.
3. LLM relevance judgment within that subset.

The semantic mode skips step 1 — it reads the full corpus. This is safe as long as the corpus stays under
~200 items (design limit enforced by `health-check` + `memory-registry.json`'s KI count check).

## Cost Profile

- **No index to maintain** — zero drift risk.
- **Scales with corpus size** — at ~200 KIs each ~800 tokens = ~160K tokens context budget. Acceptable for
  the framework corpus. At 500+ KIs this approach would need revisiting.
- **Latency** — one LLM call with a large context. Acceptable for interactive use.

## References

- Interface: [`../retriever.interface.md`](../retriever.interface.md)
- ADR: [`../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md`](../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md)
- Skills: `shared/skills/search-ki/SKILL.md` (lexical), `shared/skills/search-ki-semantic/SKILL.md` (semantic)
