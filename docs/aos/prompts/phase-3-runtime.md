# AOS Migration — Phase 3: Runtime (v3.2)

You are executing Phase 3 of the AOS migration. Scope: 13 operations spanning orchestration, RAG, Learning/Forgetting engines, and the trinity-native workflow refactor. This is the heaviest phase — pace accordingly.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prerequisite

**Phase 2 (v3.1) must be complete.** Verify:
- 10 audit-relationship counter agents exist under `shared/agents/*-auditor.md` + evaluators/reviewers/validators
- 7-8 opposing-force skills exist under `shared/skills/`
- `shared/hooks/` directory with README + schema + examples
- `validate-artifact` supports opt-in auditor invocation
- `docs/aos/governance-pairs.md` exists
- `CHANGELOG.md` has v3.1.0 entry

If not complete, halt and execute `docs/aos/prompts/phase-2-governance.md` first.

## Source of truth

`docs/aos/migration-plan.md` — Phase 3 section (Ops 3.1–3.13). Also read:
- `docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md` — governs the RAG layer's design (LLM-as-retriever for framework corpus, BM25 for installed docs, vector for feature archive, defer for source)
- `docs/aos/AOS_Governance_Design_Pack/06-LightRAG-Strategy.md` — the vector-store design guidance
- `docs/aos/AOS_Governance_Design_Pack/07-Memory-Engineering-Roadmap.md` — Learning/Forgetting engine intent
- `saturday-mcp/mcp-add-plan.md` "Retrofit Complete" — the trinity-native workflow refactor pattern applied to a real MCP server; Ops 3.11–3.13 apply the same pattern to the framework itself
- `saturday-mcp/internal/tools/retriever.go` + `bm25_retriever.go` — the retrieval adapter shape that the framework's `shared/rag/` layer should mirror

## Backward-compatibility guarantee (non-negotiable)

A team on v3.1 that upgrades to v3.2 and does not invoke `/orchestrate`, does not use `--semantic` flags, does not enable Learning/Forgetting engines, and does not opt into the trinity-refactored `FeatureDeliveryWorkflow` MUST see zero behavior change. Every runtime layer is opt-in. `deliver-feature` skill continues to work identically.

## Scope: 13 ops

### Ops 3.1-3.4 — `shared/rag/` layer (per ADR-002)

**Op 3.1**: Design `shared/rag/`:
- `README.md` — three-corpus model
- `retriever.interface.md` — pluggable adapter contract `Retrieve(query, corpus) → []Reference`
- Three adapter shapes documented per ADR-002: `llm-as-retriever`, `bm25`, `vector`
- Commit: `feat(rag): scaffold shared/rag/ interface (AOS Phase 3 Op 3.1)`

**Op 3.2**: Implement framework-corpus retrieval:
- Keep `search-ki` skill lexical (unchanged)
- Add `search-ki-semantic` skill using LLM-as-retriever
- Add `--semantic` optional flag to `query-memory`
- Commit: `feat(rag): LLM-as-retriever for framework corpus (AOS Phase 3 Op 3.2)`

**Op 3.3**: (NOTE: this differs from what the plan text calls Op 3.3 — the plan renumbered mid-execution when the corpus-aware reframing landed via `b55c0dd`; the current plan has ops 3.3 as "Design shared/rag/", already covered above in Op 3.1 of this handoff. The remaining Op 3.3 work is the installed-project docs BM25 story — but that's shipped in saturday-mcp M1 as commit `5a47441` and does not need re-implementing at the framework level.) Verify no re-work needed; commit a `docs(aos)` note in `migration-plan.md` recording that saturday-mcp handled installed-project retrieval.

**Op 3.4**: Defer installed-project source retrieval per ADR-002 (framework leans on client's Grep/Glob). Add a stub `shared/rag/source-retrieval.deferred.md` documenting the deferral rationale.

### Ops 3.5-3.6 — Learning + Forgetting engines

Phase 2 landed the SKILLS for these; Phase 3 wires the ENGINES that actually run them on triggers.

**Op 3.5**: Wire Learning engine as a hook:
- Add `shared/hooks/on-retrospective-written.yaml` — triggers `learning-engine` skill (from Phase 2 Op 2.2)
- Update `shared/skills/learning-engine/SKILL.md` if needed to formalize its "propose draft KIs, human approves" flow
- Commit: `feat(hooks): wire Learning engine on retrospective write (AOS Phase 3 Op 3.5)`

**Op 3.6**: Wire Forgetting engine as a scheduled skill:
- Add `shared/hooks/scheduled-monthly.yaml` — triggers `forgetting-engine` skill
- Update `shared/skills/forgetting-engine/SKILL.md` to formalize the "scan for staleness, propose expiration, human approves" flow
- Commit: `feat(hooks): schedule Forgetting engine monthly (AOS Phase 3 Op 3.6)`

### Ops 3.7-3.10 — `shared/orchestration/` runtime

**Op 3.7**: Design orchestration interface:
- `shared/orchestration/README.md` — the runtime's opt-in nature
- `shared/orchestration/interface.md` — how workflows plug in
- `shared/orchestration/pipeline-schema.md` — declarative pipeline definition format
- Commit: `feat(orchestration): scaffold shared/orchestration/ interface (AOS Phase 3 Op 3.7)`

**Op 3.8**: Implement orchestration wrapper for deliver-feature:
- `deliver-feature` skill continues to work unchanged
- New `/orchestrate` skill that invokes the runtime; runtime handles replay from checkpoints, parallel branches
- Commit: `feat(orchestration): implement /orchestrate wrapper (AOS Phase 3 Op 3.8)`

**Op 3.9**: Update install.sh with `--base` and `--full` modes (default `--base`)
- `--base` = current framework, no AOS layers wired
- `--full` = includes AOS layers (rag, orchestration, hooks activated)
- Commit: `feat(install): add --base and --full modes (AOS Phase 3 Op 3.9)`

**Op 3.10**: Write `docs/aos/migration-guide.md` proper (Phase 1 stubbed it) — document how to opt into each AOS layer
- Commit: `docs(aos): write full migration-guide (AOS Phase 3 Op 3.10)`

### Ops 3.11-3.13 — Trinity-native workflow refactor

Apply the Tool/Persona/Workflow trinity (established in saturday-mcp) to the framework itself. Skills that orchestrate agents become thin callers of first-class Workflows.

**Op 3.11**: Refactor `/deliver-feature` skill's orchestration logic into a `FeatureDeliveryWorkflow`:
- Skill becomes a thin caller
- Workflow owns state machine (named internal roles: analyst/architect/developer/reviewer/qa/tech-writer/devops), resumable checkpoints backed by `.claude/feature-workspace/`, explicit stage boundaries for Phase 4 policy hooks
- External invocation contract unchanged — teams keep typing `/deliver-feature <spec>`
- Commit: `refactor(deliver-feature): extract to FeatureDeliveryWorkflow (AOS Phase 3 Op 3.11)`

**Op 3.12**: Extract Red-Green-Refactor loop from `test-driven-developer` agent into `TDDWorkflow`:
- Internal roles: test-writer (unit-tester), implementer (developer), refactor-reviewer (developer + code-reviewer audit), coverage-auditor (unit-tester in audit mode)
- Agent becomes a thin caller; workflow owns the loop + retry policy + coverage gates
- Commit: `refactor(tdd): extract Red-Green-Refactor to TDDWorkflow (AOS Phase 3 Op 3.12)`

**Op 3.13**: Establish "workflow-invokes-audit-after-producer" as default composition pattern:
- Every workflow step that ends with a contract-bound artifact automatically invokes the corresponding counter agent from Phase 2 Op 2.1
- Auditor-on-failure sends artifact back to producer with the specific violations
- Config knob exists to disable per-project; default is "audits run"
- Commit: `feat(workflows): default composition invokes audit after producer (AOS Phase 3 Op 3.13)`

### CHANGELOG + verification

After all 13 ops: consolidate CHANGELOG v3.2.0 entry + run the Phase 2-shape identity check + mark Phase 3 checklist complete in `migration-plan.md`.

## Commit discipline (non-negotiable)

- ~15 commits total for Phase 3 (13 ops + CHANGELOG + verify).
- Conventional Commits.
- **NEVER `git add -A`.**
- Green health-check + tests per commit.
- Do NOT push.

## Escalation criteria

- Op 3.3's "verify no re-work needed" reveals actual re-work — halt, describe. The plan may need re-amending.
- RAG backend choice (LightRAG vs pgvector) hasn't been made in an ADR — halt, propose. Recommend LightRAG per the design pack + as a Python opt-in.
- Trinity refactor (3.11-3.13) breaks any existing test — halt. The mcp-add pattern was behavior-preserving; this MUST be too.
- Any op requires modifying a Phase 1 or Phase 2 artifact non-additively — halt.
- Install script mode addition breaks the current default install flow — halt.

## Report format (under 400 words)

```
PHASE 3 STATUS: <complete | stopped-at-op-N>

Commits landed:
  <sha> <message>
  ...

Op-by-op tally: (13 ops + support commits)
  Op 3.1 shared/rag/ interface: <landed>
  Op 3.2 LLM-as-retriever: <landed>
  Op 3.3 verify-no-rework note: <landed | rework required — describe>
  Op 3.4 defer source retrieval: <landed>
  Op 3.5 Learning engine hook: <landed>
  Op 3.6 Forgetting engine schedule: <landed>
  Op 3.7 orchestration/ interface: <landed>
  Op 3.8 /orchestrate wrapper: <landed>
  Op 3.9 install --base/--full: <landed>
  Op 3.10 migration-guide full: <landed>
  Op 3.11 FeatureDeliveryWorkflow refactor: <landed | behavior-preservation notes>
  Op 3.12 TDDWorkflow refactor: <landed>
  Op 3.13 audit-after-producer default: <landed>

Trinity refactor behavior-preservation check: <passed | drift observed — describe>

Recommended next step:
  <e.g., "human review + git push + tag v3.2.0, then execute docs/aos/prompts/phase-4-policy.md">
```

Go.
