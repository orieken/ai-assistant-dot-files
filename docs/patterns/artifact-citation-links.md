# Artifact Citation Links

Cross-reference convention for pipeline artifacts, Knowledge Items, and Architecture Decision Records.
Adopted as the standard in Epic 59 Op 4. Read before writing any link in a pipeline artifact — following
a single convention is what makes references machine-followable by agents and future retrieval tiers.

---

## Why machine-followable links matter

Pipeline artifacts, KIs, and ADRs are a knowledge graph, not a flat collection of files. When an agent
reads an `analysis.md` and follows its `linked_adrs` front-matter to `ADR-002`, then follows ADR-002's
"Related" section to `ADR-001`, it builds richer context than any single file provides alone. This is a
cheap approximation of graph-RAG — no LightRAG dependency, no embeddings overhead — that works today
with the existing lexical retrieval tier.

The convention must be consistent for this traversal to be reliable. Inconsistency (mixed path styles,
missing link sections, wiki-style `[[links]]` in some files and Markdown links in others) breaks
programmatic traversal and forces agents to guess.

---

## Link contexts and adopted syntax

### 1. Frontmatter citation fields (pipeline artifacts)

The retrieval frontmatter block defined in Epic 59 Op 1 (see each contract's "Retrieval Frontmatter
(WARN)" section) uses **repo-root-relative paths** as plain YAML list values:

```yaml
---
feature: my-feature
bounded_context: feature-delivery
domain_terms: [Pipeline, Artifact, Feature Spec]
files_touched: [shared/skills/deliver-feature/SKILL.md]
issue_refs: [PROJ-123]
linked_adrs:
  - docs/adrs/ADR-001-adopt-rag-friendly-docs-structure.md
  - docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md
linked_kis:
  - shared/knowledge/context-engineer-must-be-wired-into-pipeline.md
---
```

Path rules:
- **Always repo-root-relative** — never relative to the artifact's own directory. Every file in the
  project can resolve the same path regardless of where it lives.
- **Exact paths, not globs** — one entry per file, not wildcard patterns.
- `.claude/knowledge/*.md` is valid for project-local KIs.

### 2. In-prose cross-directory references

When a pipeline artifact or KI body references a file in a *different directory*, use a **Markdown link
with a repo-root-relative path**:

```markdown
See [ADR-002](docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md) for the graduated retrieval strategy.
See [context-engineer-must-be-wired-into-pipeline](shared/knowledge/context-engineer-must-be-wired-into-pipeline.md).
```

Agents rendering these files can follow the href directly. Future BM25/vector tiers can extract link
graphs from these hrefs without a custom parser.

**Do NOT use bare backtick paths** (`docs/adrs/ADR-001.md`) for cross-document references — they're
not clickable and aren't trivially extractable as typed links. Backtick paths are for file paths
mentioned in passing (not cited as references to load), such as "the template lives at
`shared/templates/analysis.template.md`."

### 3. Same-directory references (ADR-to-ADR)

Within `docs/adrs/`, ADRs use **relative Markdown links** when referencing sibling ADRs:

```markdown
[ADR-001](ADR-001-adopt-rag-friendly-docs-structure.md)
```

This is an established pattern already in use across all ADRs. Do not change it to repo-root-relative
for consistency with rule 2 — the relative path is already correct and shorter.

### 4. KI-to-KI back-references

KIs that conceptually extend or supersede another KI should include a back-reference section using
**repo-root-relative Markdown links** (rule 2 applies here — KIs are in `shared/knowledge/`, often
cited from other directories):

```markdown
See also: [docs-directory-follows-rag-friendly-structure](shared/knowledge/docs-directory-follows-rag-friendly-structure.md)
```

---

## Traversal story: feature → ADR → KI

An agent or retrieval tier can traverse the knowledge graph starting from any pipeline artifact:

```
docs/features/<name>/analysis.md
  └─ linked_adrs: [docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md]
       └─ Related: [ADR-001](ADR-001-adopt-rag-friendly-docs-structure.md)
            └─ docs/adrs/ADR-001-adopt-rag-friendly-docs-structure.md
  └─ linked_kis: [shared/knowledge/context-engineer-must-be-wired-into-pipeline.md]
       └─ (prose links to other KIs or ADRs)
```

Practical traversal rule for agents:
1. Read the artifact's YAML `linked_adrs` and `linked_kis` lists first — these are the high-signal
   references the producing agent explicitly chose to record.
2. For each linked ADR, read its "Related" section — it lists sibling ADRs that provide additional
   context.
3. For each linked KI, scan its body for Markdown links — they point to other KIs or ADRs relevant
   to the same topic.
4. Stop at depth 2 unless a specific link is clearly on-point for the current task.

This depth-limited traversal gives richer context than single-file reads without the cost of loading
the full corpus.

---

## What stays deferred

The convention here is a **cheap approximation of graph-RAG**. It does not provide:
- Semantic link ranking (which ADRs are *most* relevant to a query, not just which ones are cited)
- Link-weight scoring (is a link in `linked_adrs` frontmatter more authoritative than one buried in
  an artifact's prose body?)
- Automatic link extraction from unstructured prose

LightRAG and the vector tier in ADR-002 address those. This convention is the baseline that makes the
link graph explicit enough for the lexical tier to exploit it today, and for a vector tier to bootstrap
from when it arrives.

---

## Related

- [ADR-001](../adrs/ADR-001-adopt-rag-friendly-docs-structure.md) — established the docs structure
  this link convention traverses
- [ADR-002](../adrs/ADR-002-corpus-aware-retrieval-strategy.md) — graduated retrieval strategy;
  this pattern is the write-time enrichment for its lexical tier
- [frontmatter-conventions.md](frontmatter-conventions.md) — agent, skill, and KI frontmatter schemas
- Epic 59 retrieval frontmatter in `shared/contracts/*.md` — each contract's "Retrieval Frontmatter
  (WARN)" section defines the seven fields this convention populates
