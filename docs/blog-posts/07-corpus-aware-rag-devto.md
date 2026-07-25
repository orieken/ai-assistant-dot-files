---
title: "Corpus-Aware RAG: Why One-Size-Fits-All Retrieval Is Wrong for AI Coding"
published: false
description: "Why chunking everything into a vector database fails for developer tools—and how a tiered retrieval strategy fixed it."
tags: ai, architecture, rag, devtools, software-engineering
canonical_url:
cover_image:
---

Adding Retrieval-Augmentation Generation (RAG) has become the default knee-jerk answer for AI developer tools.

Got a codebase? Embed it into a vector database. Got project documentation? Chunk it into 512-token passages, compute cosine similarity, and dump the top 5 matches into the context window.

Then you ask the AI to find the definition of a specific configuration interface, and the vector search returns 5 semi-relevant README paragraphs because semantic distance in high-dimensional vector space is terrible at exact symbol lookups.

In `ai-assistant-dot-files`, we authored **ADR-002: Corpus-Aware Retrieval Strategy** to challenge the assumption that one retrieval mechanism fits all code artifact types.

Here is why one-size-fits-all RAG fails in software engineering—and how a tiered, corpus-aware approach delivers vastly higher precision.

---

## The Four Artifact Categories in Software Repositories

A software project does not contain one homogeneous text corpus. It contains four distinct types of artifacts, each with completely different scale, structure, and query dynamics:

| Artifact Type | Scale | Best Retrieval Mechanism | Why Embeddings Fail Here |
|---|---|---|---|
| **Knowledge Items (KIs) & ADRs** | 30–200 files | **LLM-as-Retriever / Summary Scoring** | Small enough to fit summaries into context; semantic vector search misses explicit tag matching. |
| **Project Documentation & Guides** | 50–500 files | **BM25 / Lexical Search** | Technical terms, class names, and CLI flags require exact keyword precision. |
| **Feature Archive & Retrospectives** | 100–1000 files | **Vector Embedding Search** | Natural language queries like "have we solved race conditions in payment webhooks before?" benefit from semantic similarity. |
| **Source Code & AST Schemas** | 100–10,000 files | **Grep / Glob / AST Tools** | Exact symbol definitions (`func ProcessPayment`) cannot tolerate fuzzy vector approximations. |

---

## 1. Small Knowledge Corpora: LLM-as-Retriever

When your project maintains 50 Knowledge Items (KIs) or Architecture Decision Records (ADRs), creating an embedding pipeline adds massive infrastructure complexity for zero gain.

Instead, we feed a 1-line title and summary of every KI directly to our `context-engineer` agent at startup. The LLM acts as its own retriever in-context, selecting exact relevant KI paths (`shared/knowledge/ki-012.md`) before implementation begins.

## 2. Technical Documentation: Lexical BM25 Precision

If a developer queries `validate-artifact --schema agent-frontmatter`, vector search often returns generic documentation about frontmatter rules because "agent" and "frontmatter" appear everywhere.

Lexical search (BM25 or Ripgrep) prioritizes exact token matches. For technical documentation containing CLI parameters, environment variables, and interface names, lexical precision outperforms dense vector embeddings almost every time.

## 3. Feature History: Vector Semantic Search

Vector databases excel when searching unstructured past human experiences.

When asking, *"What went wrong the last time we refactored external API retry loops?"*, vector search across past `retrospective.md` and `delivery-summary.md` files surfaces relevant historical lessons across unrelated modules.

## 4. Source Code: Exact Tooling (Grep / Glob)

Source code retrieval should never be delegated to a vector database. Tools built into modern AI environments (such as Claude Code's native Grep and Glob, or LSP definitions) provide 100% deterministic symbol resolution.

---

## Generative Engine Optimization (GEO) for Code & Docs

ADR-001 (*Adopt RAG-Friendly Docs Structure*) established that retrieval quality is a function of document structure. We mandate:

- **H2/H3 Headings as Self-Contained Units**: Every section must make sense if extracted in isolation.
- **Explicit Metadata Frontmatter**: Every document carries `domain`, `tags`, and `updated` timestamps.
- **No Pronoun Ambiguity**: Replacing "it handles this" with "the `SecurityAuditor` handles token verification."

---

## Takeaways for AI Engineers

1. **Stop chunking your entire codebase into a single vector DB.**
2. **Match the retrieval mechanism to the artifact category.**
3. **Optimize your documentation structure (GEO) before tuning retrieval algorithms.**

---

## Image Prompt

> Hero image prompt: A high-tech infographic diagram visualizing four distinct data streams (Knowledge Items, Docs, Retrospectives, Code) routing into four specialized retrieval filters (Context Summary, Lexical BM25, Vector Search, AST Grep). Dark futuristic dark mode UI with glowing cyan, magenta, and emerald accents.
