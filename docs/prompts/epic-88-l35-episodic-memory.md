# Epic 88 — Episodic Memory (L3.5)

Source: session discussion, 2026-09-02, following epic 87. Implements roadmap item **L3.5**, whose
blocker **L3.8** shipped in epic 84. It **blocks L4.4** (prompt registry), which blocks **L4.6**
(lessons reach prompts) — the largest remaining dependency chain on the roadmap.

The framework's "organizational memory" is semantic only: markdown KIs and ADRs. There is no record
of what was *attempted* — what failed, what a retry changed, what a human corrected. Four epics have
now produced exactly that data per run, and it dies with the run's workspace.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L3.5**, **L4.4** (what consumes this), **L3.13** (agent
   quality metrics, also blocked on it), **L3.4**/**L3.1** (why the corpus is out of scope).
2. `internal/orchestrator/timeline.go` and `vocabulary.go` — the event record and its generated
   vocabulary. This epic reads that vocabulary; it does not invent one (L3.9 shipped it).
3. `internal/orchestrator/state.go` — `RunState`, `StageRecord`, and the three record types four
   epics added: `Usage` (L3.8), `Correction` (L4.5), `PolicyRecord` (L2.16).
4. `shared/mcp/internal/tools/bm25_retriever.go` — the existing sqlite usage and driver import.
   `modernc.org/sqlite` is already a direct dependency; this epic adds no new one.
5. `shared/rag/retriever.interface.md` — read it to understand what this epic is **not** doing.

## What already exists, and must be adopted rather than rebuilt

L3.5's Problem paragraph predates four epics that quietly built most of its raw material:

| Data | Where it already is | From |
|---|---|---|
| Stage transitions, timings | `run-events.jsonl`, timestamped by the process doing the work | M0.4, L2.12 |
| Retries and loop rounds | `loop.iterated` / `loop.exhausted`, plus `Iteration` on each record | L2.17 |
| Gate decisions and approvals | `gate.waiting` / `gate.approved` / `gate.invalidated` | L2.13, L2.14 |
| Human corrections | `artifact.corrected` plus retained diffs | L4.5 |
| Token counts and cost | `Usage` on every `StageRecord` | L3.8 |
| Policy decisions | `PolicyRecord` in run state | L2.16 |
| Routing decisions | `SkipReason` per stage | L3.0 |

The item says `pipeline-trace.json` "does not exist anywhere in the repo" and that nothing can learn
from execution. The first half is still true; the second is no longer. **This epic is a store and a
query surface over data that exists**, not a new collection mechanism. Any phase that finds itself
adding a new emitter should stop and ask why.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Ingest, don't co-write | A separate step reads a run's `run-events.jsonl` + `run-state.json` and populates sqlite. The executor is **not** modified | The JSONL is already the durable append-only source of truth. Co-writing puts a database inside the run loop, adds a locking failure mode that could disturb a delivery, and duplicates a record the executor already writes reliably. A learning store that can break a run is one people switch off. Ingest is also idempotent and retroactive: every run that already happened can be ingested |
| Project-local | `.claude/memory/episodes.db` | Every other framework record is project-local, never uploaded, never aggregated — feature workspaces, knowledge, run state. A global store would silently pool data from every repo loom is run in, including client code, into one file outside those repos. That is a privacy posture change and it is not a side effect this epic gets to make |
| Query surface only | A `loom memory` command. **No** `episodic` CorpusID in `retriever.interface.md` | Its adapters are markdown specs with no running backends, and the planner that would consume the corpus is L3.1. Adding an interface entry nothing implements is precisely the defect epics 84–87 kept finding and cleaning up. L3.4 adds the corpus when there is a retriever to add it to |
| The schema mirrors the vocabulary | Event kinds are stored as the generated vocabulary's strings, not re-enumerated | L3.9 made one enum the source of truth. A second list in a database schema would drift from it, which is the exact thing that item removed |
| Re-ingest is safe | Ingesting a run twice replaces its rows rather than duplicating them, keyed by run identity | An append-only *file* becomes a queryable *store*; those have different idempotency needs. A retrospective that double-counts a retry is worse than one that has no data |
| Cost and corrections are first-class columns | Not JSON blobs | The done-when is a query, and the two things anyone will actually ask about are "what did this cost" and "what did a human have to fix" |
| No new event types | If a query needs a fact nothing records, that is a finding to report, not a licence to add an emitter | The vocabulary's rule (L3.9): a type joins by acquiring an emitter, never by being wanted |

### Additional decisions (settled 2026-09-02, before Phase A)

**The executor's records are not archived.** `run-events.jsonl` and `run-state.json` live only in
`.claude/feature-workspace/<feature>/`, the temporary directory. `docs/features/<name>/` receives
the artifacts and `pipeline-trace.json` but not these — so the episodic source dies when a
workspace is cleaned.

| Decision | Choice | Rationale |
|---|---|---|
| Ingest runs automatically | `loom run` ingests when it completes or halts at a gate, plus a `loom memory ingest` command for anything else | A store nobody populates is a feature that exists only in its tests, and the payoff here is months away when L4.4 wants history that was never captured. Ingest failure must **never** fail a run — same posture as L4.5's baseline capture. This touches the CLI, not `internal/orchestrator`, so the "do not modify the executor" guardrail holds |
| The records get archived | `run-events.jsonl` and `run-state.json` join what is persisted into `docs/features/<name>/` | They then live in git beside the artifacts they describe: version-controlled, reviewable in a PR, and enough to rebuild the store if `episodes.db` is lost or gitignored. Making an untracked binary the sole home of a run's history is how history quietly disappears |
| The store is a cache, not the record | The archived JSONL is the durable record; sqlite is a queryable projection of it | This is why re-ingest must be idempotent, and it is what makes deleting `episodes.db` a recoverable event rather than a loss |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer, max 4
  parameters (introduce a parameter object rather than a fifth).
- **Parameterized queries only**, per `shared/rules/go-conventions.md`. No string-built SQL.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**66.0%** as of epic 87).
- **The executor is not modified.** If a phase needs to change `internal/orchestrator` to make
  ingest work, STOP and escalate — that means the ingest premise is wrong.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — The store and the ingest — UNBLOCKED

1. New `internal/memory/`: a sqlite schema keyed by run and stage — runs, stages, events,
   corrections, policy decisions — with cost and correction counts as columns rather than blobs.
2. Ingest one run from its `run-events.jsonl` and `run-state.json`. Idempotent: re-ingesting
   replaces that run's rows.
3. A run identity that survives re-ingest and does not collide across features.
4. Tests: a fixture run ingests; re-ingesting changes no counts; a partial run (halted at a gate,
   no `run.completed`) ingests cleanly, because most runs a human looks at are halted ones; a
   corrupt trailing line is skipped rather than fatal, matching how the timeline is already read.

5. Archive `run-events.jsonl` and `run-state.json` into `docs/features/<name>/`, and ingest
   automatically at the end of a run — reporting a failure without failing the run.

**Done when**: a completed or halted run's history is in sqlite, re-ingesting it is a no-op, and the
records it was built from are archived beside the artifacts they describe.
**Commit** (`feat(memory): an episodic store ingested from the run timeline`), report, PAUSE.

### Phase B decisions (settled 2026-09-02, before the phase)

| Decision | Choice | Rationale |
|---|---|---|
| Named queries only | Fixed subcommands, each parameterized by construction. **No** `--sql` escape hatch | Injection-proof by construction, discoverable through `--help`, and it keeps the schema private so later phases can change it. The cost is real and accepted: a question nobody anticipated needs a code change, so the store is only as useful as the questions already imagined — which is an argument for choosing the first few carefully, not for opening the database |
| Backfill from the archive | `loom memory ingest` with no argument walks `docs/features/*/` and ingests every run whose records are archived | This is what makes the store genuinely rebuildable after deletion, and it imports deliveries that happened before this epic. Idempotency already makes repeated runs safe. Without it, archiving the records buys much less than it should |
| "Retried more than twice" means iteration count > 2 | At least three attempts. The output header states the reading | The done-when is ambiguous — a retry could be the stage record's `Iteration` or a count of `loop.iterated` events, and "more than twice" could mean three attempts or three retries after the first. Pinning it in the output means nobody has to guess which reading they are looking at |

## Phase B — The query surface — BLOCKED BY Phase A

1. `loom memory` subcommands: ingest a run, and query. The queries must include L3.5's own
   done-when — "every run where code-reviewer retried more than twice" — plus cost per run, and
   which agents humans corrected most.
2. Parameterized throughout. Output readable in a terminal, with `--json` for anything else.
3. Report honestly when a question cannot be answered from what is recorded, naming the missing
   fact rather than returning empty — the lesson from epic 87's UNKNOWN.
4. Tests: the done-when query returns the right runs against a fixture with known retry counts;
   queries over an empty store return nothing rather than erroring; `--json` parses.

**Done when**: "show me every run where code-reviewer retried more than twice" is answerable by
query — L3.5's stated done-when, in full. **Commit** (`feat(memory): query the episodic store`),
report, PAUSE.

## Phase C — Docs and boundaries — BLOCKED BY Phase B

1. Roadmap: L3.5 **SHIPPED**, stating what it stores, what it deliberately does not (the corpus),
   and that the executor was not touched. Update **L4.4** and **L3.13** with what they inherit.
2. `shared/skills/pipeline-trace/SKILL.md`: it owns a `pipeline-trace.json` schema whose timings are
   model-written estimates. Say plainly which questions the episodic store now answers with measured
   data, and what remains its own.
3. `shared/rag/retriever.interface.md`: record that an `episodic` corpus is intended and belongs to
   **L3.4**, without adding it.
4. `cmd/loom/README.md`, `README.md`, `docs/ARCHITECTURE.md`, `shared/DOMAIN_DICTIONARY.md`, and
   `shared/VERSION`.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the docs say what is stored, what is queryable, and what is still an estimate.
**Commit** (`docs(memory): what the episodic store knows`), report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- The `episodic` CorpusID or any retriever adapter (**L3.4**)
- The planner that would consume it (**L3.1**)
- Prompt-variant generation or anything that acts on the data (**L4.4**)
- Any new event type or emitter
- Cross-project aggregation, and any upload of anything
- Modifying the executor

## Report format (end of every phase)

```
## Epic 88 Phase <X> Report
- Roadmap item: L3.5 — episodic memory
- Blockers verified: <list, with evidence (commit SHAs / files)>
- Commits: <sha> <subject>
- Build/lint/test: go build PASS|FAIL · golangci-lint PASS|FAIL · go test PASS|FAIL (counts) · health-check PASS|FAIL|n/a
- Done-when criterion: <restate it> — MET | NOT MET (why)
- Escalations / open questions: <list or "none">
- Next phase blocked by: <what must land first>
```

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering
Framework by Oscar Rieken — licensed under
[CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
