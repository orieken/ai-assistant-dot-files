# Retriever Adapter Interface

## Contract

```
Retrieve(query: string, corpus: CorpusID) → []Reference
```

### Types

```typescript
type CorpusID =
  | "framework-ki"       // shared/knowledge/ + docs/adrs/
  | "project-docs"       // <project>/docs/ (BM25)
  | "project-features"   // <project>/docs/features/ (vector)
  | "project-source"     // <project>/src/ (DEFERRED)

type Reference = {
  path: string        // repo-relative path to the canonical markdown or source file
  title: string       // human-readable label (KI name, ADR title, feature name)
  excerpt: string     // 1-2 sentence context excerpt — never a content copy
  score: number       // 0.0-1.0 relevance; interpretation is backend-specific
  corpus: CorpusID    // which corpus this result came from
}
```

### Behavior Contract

1. **Returns references, not content.** Callers must read the referenced file to get content. This prevents stale embedded copies and respects the "always verify against canonical markdown" principle from `docs/aos/AOS_Governance_Design_Pack/06-LightRAG-Strategy.md`.

2. **Graceful empty result.** If no relevant results exist, returns `[]` — never throws. Callers must handle the empty case.

3. **Top-K bounded.** Returns at most 10 results per call. Callers that need broader coverage call `Retrieve` multiple times with refined queries.

4. **Corpus isolation.** A `Retrieve` call targets exactly one corpus. Multi-corpus search is implemented by the caller making sequential calls and merging by score.

5. **No side effects.** `Retrieve` is a pure read operation. It never writes to the index, never updates `memory-registry.json`, never emits telemetry. Callers may emit telemetry around `Retrieve` calls if desired.

6. **Score is comparable within a backend, not across backends.** BM25 scores and cosine similarities live on different scales. The merger pattern is: sort each corpus's results by descending score, interleave round-robin, deduplicate by path.

### Implementations

| Backend | File | Corpus |
|---|---|---|
| LLM-as-retriever | `adapters/llm-as-retriever.md` | `framework-ki` |
| BM25 | `adapters/bm25.md` | `project-docs` |
| Vector (sqlite-vec) | `adapters/vector.md` | `project-features` |
| Deferred | `adapters/source-retrieval.deferred.md` | `project-source` |

### Extension

To add a new corpus:
1. Add a `CorpusID` variant above.
2. Create an `adapters/<name>.md` file documenting the implementation.
3. Register the backend in `shared/memory-registry.json` under a new source entry with `retrievalBackend: "<value>"`.
4. No existing adapter or caller needs to change.
