# Epic 60 — Retrieval Index Freshness + Eval Loop (project-as-RAG, part 2 of 2)

Source: retrieval-optimization discussion 2026-07-31 (follow-up to
`docs/audits/framework-gap-audit-2026-07-31.md`), building on ADR-002. Sibling: Epic 59
(`epic-59-retrieval-write-time-enrichment.md`) — write-time enrichment. Epic 59 is independent;
this epic is best run AFTER it (the CODEMAP and eval set are more valuable over an enriched
corpus) but nothing here hard-depends on it.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

- ADR-002's own consequences section flags the index-staleness problem: any index over the
  markdown needs an install-time build + rebuild story, and cache invalidation on doc edits "is
  a real concern for the vector tier." The `shared/hooks/` layer (Phase 2) is the standing
  answer nobody has wired up yet.
- The BM25 docs-search tooling lives in **saturday-mcp** (M1, commit `5a47441` there — see
  `docs/aos/prompts/phase-3-runtime.md` Op 3.3 note), NOT in this repo. This epic writes hook
  examples and conventions AGAINST that tooling's interface; it does not build an index engine
  here.
- ADR-002 chose a judgment-only fitness function for retrieval quality, with a named
  "mechanical CI check that could be added later": every `shared/memory-registry.json` source
  carries a `retrievalBackends` value from the enum `{lexical, llm-as-retriever, bm25, vector}`.
  All six sources currently say `lexical`.
- Telemetry layer (Phase 1): `shared/telemetry/event-recorder.md` + `event-schema.md`, opt-in.
  `retrieval-evaluator` counter agent already audits KI/ADR retrievability from telemetry and
  flags zero-match queries as missing-KI or bad-metadata candidates.
- Source-code tier is deferred per ADR-002 (lean on client Grep/Glob) — the CODEMAP op below
  improves that tier *without* un-deferring it.

## Scope: 4 ops

**Op 1 — Index-freshness hooks.**
Add `shared/hooks/examples/on-artifact-written-reindex.<ext>` and
`on-ki-created-reindex.<ext>` (follow the existing examples' schema exactly): on writes under
`docs/features/`, `docs/adrs/`, or the KI dirs, trigger an index upsert via the saturday-mcp
reindex entry point (name the actual tool after reading `shared/mcp-patterns/` — if no reindex
entry point exists there, see Escalation). Document in `shared/hooks/README.md` that these
examples are the ADR-002 rebuild story's event-driven half; `/reindex`-style manual rebuild
remains the fallback. Opt-in like every hook.
Commit: `feat(hooks): index-freshness hook examples for docs/KI writes (Epic 60 Op 1)`

**Op 2 — Registry retrievalBackends fitness function.**
Implement ADR-002's named-but-deferred mechanical check: extend `scripts/health-check.sh`'s
Memory Registry section to FAIL if any registry source lacks `retrievalBackends` or uses a value
outside `{lexical, llm-as-retriever, bm25, vector}`. Pure bash/python3-stdlib — no new
dependencies (mind the ci-check container lesson from commit `6c422cb`: degrade to SKIP if a
needed module is missing, never false-FAIL).
Commit: `feat(health-check): enforce retrievalBackends enum on memory registry (Epic 60 Op 2)`

**Op 3 — CODEMAP generation for the source tier.**
Add `scripts/generate-codemap.sh`: emits `CODEMAP.md` at repo root — directory tree (depth ~2)
with a one-line purpose per directory, sourced from a package's README first line or a
documented per-language convention (Go package comment, `index.ts` header, etc.). Deterministic,
no LLM. Document it as the entry-point document for any LLM client's Grep/Glob over source
(ADR-002's deferred source tier gets better without being un-deferred). Generate this repo's own
CODEMAP.md as the demonstration. Wire a WARN into `check-inventory-drift.sh` or `health-check.sh`
if CODEMAP.md exists but is older than the newest directory change (mtime check only).
Commit: `feat(scripts): CODEMAP generator for source-tier retrieval entry point (Epic 60 Op 3)`

**Op 4 — Telemetry-sourced retrieval regression set.**
Define `shared/evaluation/retrieval-regression.md`: a format for capturing real telemetry
queries as regression cases (`query → must-appear-in-top-5 reference`), plus instructions for
`retrieval-evaluator` to (a) propose new cases from the telemetry log and (b) run the existing
set and report drift. This graduates ADR-002's judgment-only fitness function toward mechanical
— using queries people actually asked, not synthetic benchmarks. Cases are proposed, human
approves (same discipline as learning-engine's draft-KI flow). Update
`shared/agents/retrieval-evaluator.md` (version bump + CHANGELOG in same commit —
`check-agent-versions-ci.sh` enforces it).
Commit: `feat(evaluation): telemetry-sourced retrieval regression set (Epic 60 Op 4)`

After every op: `bash scripts/health-check.sh` green on a pristine repo (no hooks configured, no
telemetry) — every addition is opt-in.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- Op 1: if `shared/mcp-patterns/` exposes no reindex entry point, halt on that op — propose the
  interface as a note for saturday-mcp M2 instead of inventing a phantom tool name; ship the
  other ops.
- Op 2: if the registry's actual field shape differs from ADR-002's description (field name,
  nesting), follow the registry as-is and note the ADR discrepancy — don't change the registry
  schema in this epic.
- Op 3: if a one-line-purpose source doesn't exist for most directories, generate with `TODO`
  placeholders and WARN — don't fabricate descriptions.
- Op 4: if the telemetry event schema lacks the fields needed (query text, hits, chosen), halt
  on that op and propose the schema extension for human review — schema changes ripple.

## Report (under 150 words)

```
Commits:
  <sha> <message>  (x4, or fewer with halts noted)
Reindex entry point used/proposed: <name>
retrievalBackends check: <pass on current registry?>
CODEMAP.md: <generated, N directories, M TODO placeholders>
Regression set: <format shipped; N seed cases from existing telemetry, or 0 + why>
health-check pristine-repo: <pass>
Halted ops: <none | list + reason>
```

Go.
