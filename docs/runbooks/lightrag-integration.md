# Runbook: LightRAG Integration (Deferred — Not Built Yet)

`shared/memory-registry.json`'s `lightrag` backend entry is `"status": "disabled"`. This runbook exists so
that decision is a documented, reversible choice — not a gap nobody noticed — the same pattern as
`docs/runbooks/scaling-cross-feature-learning.md`'s "don't build the index below until you've actually felt
the grep approach strain."

**Nothing under this runbook is implemented.** No `scripts/memory/build-memory-corpus.sh` exists, no LightRAG
dependency is installed, no code calls it. Writing that now, with no consumer to test it against, is the
exact kind of premature abstraction `shared/rules/design-principles.md`'s Anti-Pattern Radar already warns
about.

## Why lexical search is still the right choice today

`search-ki` and `query-memory` both do tag/domain pre-filter + full-body LLM judgment reads — no embeddings,
no vector index. That's a deliberate choice (see `search-ki`'s own Guardrails section), not an oversight:
for a KI corpus in the tens of entries, reading the whole corpus is cheap and the judgment-based ranking is
*more* precise than nearest-neighbor search, which can miss a KI that's the right answer despite sharing no
vocabulary with the query (exactly the `subagent-isolation-is-a-hard-boundary` example `search-ki`'s own docs
use).

## When to actually build this
Don't build it preemptively. Build it when one of these is true, not before:
- The combined KI + ADR + feature-archive corpus grows into the hundreds of entries, and pre-filter-then-read
  visibly stops being affordable (context-budget warnings from `context-engineer`, or `query-memory` runs
  taking noticeably long).
- You have a concrete case where lexical search demonstrably misses something a semantic/graph-based
  retrieval would have caught — not a hypothetical, an actual missed match you can point to.
- `docs/runbooks/scaling-cross-feature-learning.md`'s own index (if built) also starts straining — the two
  problems are related (both are "the corpus got big"), so it's worth revisiting both at once.

## What to build, when that day comes

1. **`scripts/memory/build-memory-corpus.sh`** — generates a curated LightRAG input corpus from the sources
   `shared/memory-registry.json` marks `indexable: true`: `shared/knowledge/`, `docs/adrs/`,
   `docs/features/*/retrospective.md`, `docs/lessons-learned/`, `shared/DOMAIN_DICTIONARY.md`,
   `shared/TEAM_TOPOLOGY.md`. Exclude caches, generated platform configs, secrets, and any transient
   `.claude/feature-workspace/` state — none of that is durable memory.
2. **Indexing** — run the corpus through LightRAG's own ingestion. Configuration should be local-first (no
   hosted service dependency) and live in a project-local config file, never committed with credentials.
3. **Update `shared/memory-registry.json`** — flip `lightrag.status` to `"active"`, list `query-memory` (and
   optionally `search-ki`) under `lightrag.usedBy`.
4. **Update `query-memory`** (and optionally `search-ki`) to actually call the LightRAG backend when the
   registry says it's active — but **every result retrieved through it must be verified against its source
   markdown file before being presented as fact**. This is non-negotiable: LightRAG becomes a faster way to
   *find* the right markdown, never a replacement for reading it. A hallucinated or stale embedding match
   presented as ground truth would be worse than the lexical search it replaced.
5. **Add a verification step** to `scripts/health-check.sh` or `scripts/memory/validate-memory.sh` (build the
   latter at this point too — no reason to have it before there's a corpus/index to validate) confirming the
   generated corpus doesn't contain anything the exclusion list in step 1 was supposed to catch.

## What NOT to do
- Don't half-integrate this — a `lightrag.status: "active"` with no actual corpus/index behind it is worse
  than leaving it disabled; the registry would be lying about what's actually available.
- Don't skip the markdown-verification step under time pressure — that's the one guardrail that keeps this
  addition from undermining the "markdown is canonical" principle the whole memory system is built on.
- Don't build this because it sounds more sophisticated than lexical search — build it because a concrete,
  observed pain point from the "When to actually build this" section above is real.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
