---
name: ollama-local-embeddings
tags: [context-engineering, embeddings, semantic-search, local, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

Ollama (`ollama.com`) runs open-weight models locally — no API key, no network call, no cost
per token. When paired with `nomic-embed-text`, it provides a local embedding model suitable
for semantic search over the KI corpus, session history, or any small document collection.

The primary use case in this framework is as a backend for `search-ki-semantic` in environments
where calling an external embedding API is not desirable — offline teams, air-gapped CI, or
cost-sensitive setups where the KI corpus is large enough to make per-query API calls add up.

## When to reach for it

- The team runs `search-ki-semantic` frequently and wants to avoid external API latency/cost
- The environment is air-gapped or internet-restricted
- ctx is configured for semantic search (`ctx setup --semantic`) and needs a local embedding
  model rather than the ctx-managed cloud embedding

## Setup

```bash
# Install ollama
brew install ollama
# or: curl -fsSL https://ollama.com/install.sh | sh  (Linux)

# Start the ollama server
ollama serve

# Pull the embedding model (runs once, ~274MB)
ollama pull nomic-embed-text
```

## Generating embeddings

```bash
# Via the REST API (ollama serve must be running)
curl http://localhost:11434/api/embeddings \
  -d '{"model": "nomic-embed-text", "prompt": "context engineering framework"}'
```

The response contains a 768-dimensional vector. Use a lightweight vector store (e.g.
`sqlite-vss`, `hnswlib`, or a flat cosine-similarity scan for small corpora) to index and
retrieve.

## Integration with ctx semantic search

If ctx is installed and configured:
```bash
# Tell ctx to use the local ollama embedding endpoint
ctx setup --semantic --embedding-url http://localhost:11434/api/embeddings
ctx index
```

Verify with `ctx search "some domain term" --semantic` — results should appear without
any network calls to an external embedding service.

## Alternative embedding models via ollama

| Model | Size | Notes |
|---|---|---|
| `nomic-embed-text` | ~274MB | Default recommendation — good quality, small |
| `mxbai-embed-large` | ~670MB | Higher quality, larger |
| `all-minilm` | ~46MB | Smallest option; lower quality but very fast |

## Guardrails

- `ollama serve` must be running before any embedding call. In a pipeline context, start it
  as a background process and add a readiness check before the first embedding call.
- Embedding dimensions differ across models — `nomic-embed-text` produces 768-dim vectors;
  switching models requires re-indexing all documents.
- Ollama is not suitable for high-concurrency production use. For a shared team setup, a
  dedicated embedding service is more appropriate.
- This is an offline fallback, not an upgrade. If an external embedding API is already in use
  and performing well, there is no benefit to switching.

## See also

- `shared/skills/search-ki-semantic/` — the skill this most directly enables in offline mode
- `shared/knowledge/ctx-session-history-search.md` — ctx's `--semantic` flag can be pointed
  at the local ollama endpoint
