# shared/rag/ — Corpus-Aware Retrieval

**Status**: AOS Phase 3 (v3.2) — opt-in. Teams that don't invoke semantic retrieval skills see no change from v3.1.

## Three-Corpus Model

The framework retrieves across three structurally different corpora, each matched to the right algorithm:

| Corpus | Default Backend | Location | Rationale |
|---|---|---|---|
| Framework KIs + ADRs | LLM-as-retriever | `shared/knowledge/`, `docs/adrs/` | Small (~30-200 items); fits in context; no index maintenance |
| Installed project `docs/` | BM25 via sqlite-fts5 | `<project>/docs/` | Prose-heavy structured docs; lexical precision; no ML cost |
| Installed project feature-archive | Vector via sqlite-vec | `<project>/docs/features/` | Semantic "have we built this before?" query; paraphrases matter |
| Installed project source | Lean on client (Claude Code Grep/Glob) | `<project>/src/` | Modern LLM clients already do this well; **deferred** |

This graduated strategy is defined in [ADR-002](../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md).

## Adapter Interface

All retrieval implementations satisfy the contract in [`retriever.interface.md`](retriever.interface.md):

```
Retrieve(query: string, corpus: CorpusID) → []Reference
```

Implementations swap without touching callers. The three concrete adapters are documented in `adapters/`.

## Opt-In Path

| Capability | How to enable |
|---|---|
| Semantic search over framework KIs | Invoke `/search-ki-semantic` or `/query-memory --semantic` |
| BM25 over installed project docs | Install `saturday-mcp` M1 (ships `search_docs`, `search_adrs` MCP tools) |
| Vector over feature archive | Install `saturday-mcp` M2 (ships `search_features` MCP tool) |
| Source retrieval | **Not yet implemented** — see `adapters/source-retrieval.deferred.md` |

Teams that use only `/search-ki` and `/query-memory` (without `--semantic`) continue with lexical retrieval unchanged.

## Index Management

Installed-project indexes live in `.claude/rag/` (per-install, gitignored). Rebuilt via:
- **Install-time**: initial index pass during `install.sh --full`
- **On-demand**: `/reindex` skill (to be added in a future op)

The framework KI corpus requires no index — it fits in context and is read directly by the LLM-as-retriever adapter.

## File Map

```
shared/rag/
├── README.md                          ← this file
├── retriever.interface.md             ← adapter contract
└── adapters/
    ├── llm-as-retriever.md            ← framework corpus (default)
    ├── bm25.md                        ← installed-project docs
    ├── vector.md                      ← installed-project feature-archive
    └── source-retrieval.deferred.md   ← rationale for deferral
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 3 Runtime layer — Oscar Rieken, CC BY 4.0.*
