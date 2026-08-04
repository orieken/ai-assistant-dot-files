# Epic 59 — Retrieval Write-Time Enrichment (project-as-RAG, part 1 of 2)

Source: retrieval-optimization discussion 2026-07-31 (follow-up to
`docs/audits/framework-gap-audit-2026-07-31.md`), building on ADR-001 (RAG-friendly docs
structure) and ADR-002 (graduated corpus-aware retrieval). Sibling: Epic 60
(`epic-60-retrieval-index-freshness-eval.md`) — index freshness + retrieval eval. This epic has
no dependency on Epic 60, AOS Phase 3, or any vector/BM25 machinery; everything here improves
today's lexical tier AND every future tier, because later backends inherit the enriched corpus.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

The core insight: an installed project isn't a corpus indexed after the fact — the delivery
pipeline *generates* the corpus. Every `deliver-feature` run writes spec, analysis,
architecture-notes, qa-report, and retrospective through `validate-artifact` contract checks.
That write moment is when full knowledge is available; capturing retrieval metadata there is
deterministic and free, versus reconstructing it later with embeddings. Relevant existing pieces:

- `shared/contracts/*.md` (16 contracts) — enforce required *sections* per artifact;
  `shared/skills/validate-artifact/SKILL.md` grep-checks them at every pipeline handoff.
- KIs already carry retrieval frontmatter (name, tags, domain — see
  `shared/knowledge/README.md`); pipeline artifacts carry none. That asymmetry is the gap.
- `shared/DOMAIN_DICTIONARY.md` — canonical vocabulary; `check-ubiquitous-language` skill
  enforces it in code, so code + docs + queries can share one controlled vocabulary.
- `shared/skills/search-ki/SKILL.md` and `shared/skills/query-memory/SKILL.md` — the lexical
  retrieval layer that benefits immediately.
- `shared/skills/summarize-artifact/SKILL.md` — produces ~200-word artifact summaries for
  context decay in `deliver-feature`, currently ephemeral (never persisted).

## Scope: 4 ops

**Op 1 — Artifact retrieval frontmatter (contracts + template + validator).**
Define a standard frontmatter block for pipeline artifacts: `feature`, `bounded_context`,
`domain_terms` (from DOMAIN_DICTIONARY), `files_touched`, `issue_refs`, `linked_adrs`,
`linked_kis`. Add it to `shared/templates/*.template.md`, document it in each contract, and teach
`validate-artifact` to check it — **WARN-level first, not FAIL** (existing installed projects
must not break; promote to FAIL in a later release once adoption is visible). Update
`shared/contracts/README.md` if one exists, else the contracts' shared preamble.
Commit: `feat(contracts): retrieval frontmatter on pipeline artifacts, WARN-level (Epic 59 Op 1)`

**Op 2 — Domain-dictionary query expansion in the lexical tier.**
Update `search-ki` and `query-memory` SKILL.md instructions: before searching, expand the query
with synonyms/aliases from `DOMAIN_DICTIONARY.md` (and the installed project's own
`DOMAIN_DICTIONARY.md` when run inside one). Document the limits honestly — this catches
vocabulary drift ("checkout" vs "purchase flow"), not deep semantics; that stays ADR-002's
vector tier.
Commit: `feat(skills): domain-dictionary query expansion in search-ki + query-memory (Epic 59 Op 2)`

**Op 3 — Persist artifact summaries as retrieval surrogates.**
Extend `summarize-artifact` + `deliver-feature` so summaries land as a persistent sidecar
(propose the exact convention — e.g. `docs/features/<name>/summary.md` — after reading how the
feature archive is laid out today). Register the surrogate layer in
`shared/memory-registry.json` (or note it under the feature-archive source). Summaries are what
future BM25/vector tiers should index first — say so in the file's own header.
Commit: `feat(skills): persist artifact summaries as retrieval surrogates (Epic 59 Op 3)`

**Op 4 — Machine-followable citation links.**
Whatever link syntax KIs/ADRs already use for cross-references (READ FIRST — adopt, don't
invent), standardize it across pipeline artifacts via the Op 1 frontmatter (`linked_adrs`,
`linked_kis`) plus in-prose convention, and document the traversal story (feature → ADR → KI) in
`docs/patterns/` (new short pattern doc or extension of an existing retrieval-related one). This
is the cheap approximation of graph-RAG that keeps LightRAG deferred.
Commit: `docs(patterns): citation-link convention for artifact traversal (Epic 59 Op 4)`

After every op: `bash scripts/health-check.sh` green; after Op 1 run a `validate-artifact` dry
check against an existing archived feature's artifacts to prove WARN (not FAIL) behavior.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If any op requires changing artifact *sections* (not just adding frontmatter), halt — that
  breaks the existing contract-validation story for archived features.
- If Op 3's persistence would change `deliver-feature` behavior for teams that never invoke
  summarize-artifact, halt — opt-in guarantee applies.
- If KIs and ADRs turn out to use *conflicting* link syntaxes, halt and propose one, with a
  migration note — don't silently pick.

## Report (under 150 words)

```
Commits:
  <sha> <message>  (x4)
Frontmatter fields shipped: <list>
validate-artifact behavior on legacy artifacts: WARN confirmed (evidence: <what you ran>)
Link syntax adopted: <syntax, and where it came from>
Surrogate convention: <path pattern>
health-check: <pass>
```

Go.
