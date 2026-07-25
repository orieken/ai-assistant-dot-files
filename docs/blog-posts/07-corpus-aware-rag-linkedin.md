Adding vector RAG to a codebase is the fashionable default. But chunking everything into a vector database fails for software developer tools.

Why? Because code repositories contain four fundamentally different artifact types:

1. **Knowledge Items & ADRs** → LLM-as-Retriever (summaries fit in context window)
2. **Project Documentation** → Lexical BM25 Search (exact technical terms matter)
3. **Feature Retrospectives** → Vector Semantic Search (great for finding past bug patterns)
4. **Source Code** → Deterministic Grep / AST tools (exact symbol definitions)

In ADR-002, we established a Corpus-Aware Retrieval Strategy. Matching the retrieval mechanism to the artifact category delivers dramatically higher precision than any one-size-fits-all vector database.

Full technical breakdown: TODO_DEVTO_URL
