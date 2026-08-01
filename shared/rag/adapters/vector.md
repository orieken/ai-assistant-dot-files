# Adapter: Vector (Feature Archive Semantic Search)

**Corpus**: `project-features`
**Files**: `<project>/docs/features/*/` (feature deliveries and retrospectives)
**Implementation home**: `saturday-mcp` (M2 milestone)

## Rationale

The "have we built something like this before?" query is inherently semantic. A feature spec may say
"user preference sync" while a prior delivery called it "settings persistence across devices". Keyword
matching misses this. Cosine similarity on embeddings captures it.

The feature archive is the right scope for vector search because:
- Bounded size — one subdirectory per delivered feature, grows linearly with deliveries.
- High-value semantic queries — analyst asking about prior work is where vector search pays off most.
- Manageable churn — feature deliveries are complete documents, not actively-edited source files.

BM25 handles `docs/adrs/` and `docs/patterns/` (keyword retrieval is sufficient for those structured
bodies); vector handles `docs/features/` (semantic similarity is the meaningful query there).

## Implementation

This adapter is implemented as an **MCP tool in `saturday-mcp`** (Milestone 2).

### MCP Tool

| Tool | Corpus | Query type |
|---|---|---|
| `search_features` | `docs/features/*/` | Cosine similarity (sqlite-vec) |

### Vector Store

- **Backend**: `sqlite-vec` — a SQLite extension for vector similarity search. Embedded, no server.
- **Storage**: `.claude/rag/features.db` (gitignored)
- **Embedding model**: configurable via `.claude/rag/config.yaml` (default: `claude-haiku-4-5` for
  embedding generation at index time; or a local embedding model if the user has one configured)
- **Dimensions**: 1024 (or match the configured model's output dimension)

### Algorithm

1. Open `.claude/rag/features.db`. If missing, return empty + surface "feature index not built" warning.
2. Generate embedding for `query` using the configured embedding model.
3. Execute sqlite-vec similarity query: `SELECT path, title, distance FROM features_vec ORDER BY distance LIMIT 10`.
4. Map rows to `Reference` objects; compute excerpt by reading the feature's `README.md` first sentence.
5. Score: `1.0 - (distance / max_distance_in_result_set)` to normalize to 0.0-1.0.

### Index Management

- **Initial build**: `install.sh --full` triggers embedding pass over `docs/features/`
- **Incremental update**: embed new feature on `docs/features/<name>/` write (hook candidate)
- **Rebuild**: `/reindex` skill

## LightRAG as v2

The design pack (`docs/aos/AOS_Governance_Design_Pack/06-LightRAG-Strategy.md`) names LightRAG as the
preferred graph-aware vector backend. sqlite-vec is the simpler v1 baseline — no Python runtime
dependency, no graph traversal overhead, straightforward upgrade path. If telemetry reveals that
sqlite-vec's flat cosine similarity misses relationship-aware queries (e.g., "features that depend on
the auth system"), upgrade to LightRAG via a new ADR. The adapter interface stays unchanged.

## References

- Interface: [`../retriever.interface.md`](../retriever.interface.md)
- ADR: [`../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md`](../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md)
- Design pack: `docs/aos/AOS_Governance_Design_Pack/06-LightRAG-Strategy.md`
- Implementation: `saturday-mcp/internal/tools/retriever.go` (Milestone 2)
