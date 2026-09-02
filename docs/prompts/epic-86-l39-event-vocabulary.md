# Epic 86 — One Event Vocabulary (L3.9)

Source: session discussion, 2026-09-01, following epic 85. Implements roadmap item **L3.9**, whose
blocker **L3.8** shipped in epic 84. Epics 84 and 85 both had to write "that is L3.9's problem" into
their scope boundaries; this closes it.

The framework documents **six** event types, specifies **nine more** across its spec files, and
emits **none of the fifteen** to the file the schema names. `event-recorder.md` instructs producers
to refuse any type not documented — while 60% of the specified surface is undocumented. Meanwhile
the executor emits a real, Go-owned event log that none of this prose mentions.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L3.9**, and **L3.5** (which will build on this vocabulary).
2. `internal/orchestrator/timeline.go` — the `EventKind` enum and `Event` struct that already exist
   and are already the only things emitting events. This is the source of truth being formalised,
   not a new one being invented.
3. `shared/telemetry/README.md`, `event-schema.md`, `event-recorder.md` — the layer being retired,
   and the three-way distinction epic 84 wrote into the README.
4. `cmd/gen-schemas/` and `internal/state/schema.go` — the generation-plus-drift-test pattern
   established for typed state. Follow it rather than inventing a second one.
5. `shared/orchestration/policy-evaluator.md`, `shared/policies/policy-schema.md`,
   `shared/orchestration/audit-composition-pattern.md` — where the nine unemitted types are
   specified, and where the `event` / `event_type` key mismatch lives.

## Verified before drafting

- All nine extra types (`policy.evaluated`, `policy.conflict`, `policy.skipped`, `audit.fail`,
  `audit.retry`, `audit.halt`, `contract.retry`, `workflow.completed`, `workspace.migrated`) appear
  **only in prose**. Nothing emits them.
- `policy-schema.md` really does use the key `"event"` where the schema requires `"event_type"`.
- `event-recorder.md` is a **spec document in `shared/telemetry/`, not an installed skill**, so
  retiring it changes no installation and breaks no user's setup.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| The enum stays in `internal/orchestrator` | `EventKind` is **not** moved to `internal/telemetry`, despite L3.9's target-files line naming a new `internal/telemetry/events.go` | The executor emits these events, and `internal/telemetry` is the OTel adapter. Moving the enum there would make `internal/orchestrator` import it, which the L3.8 guardrail fitness function forbids and should keep forbidding. The roadmap line was written before that boundary existed |
| Source of truth | The Go enum. The JSON Schema and the documentation table are **generated** from it | Exactly the pattern typed state already uses, including the drift test. Two hand-maintained lists is how the contradiction happened |
| `events.jsonl` and `event-recorder.md` are retired | Both go. `run-events.jsonl` is the event log | The file never had an emitter and the recorder never had a caller. Keeping a prose layer that describes a file nobody writes is what made "60% of the surface undocumented" possible — there was nothing to check the prose against |
| The nine unemitted types | **Not in the enum.** They become a documented "specified, not emitted" list, each naming the roadmap item that would build its emitter | An enum where two-thirds of entries fire from nowhere is the same trap in a new location. The design intent survives as a roadmap pointer, which is honest about its status in a way a type name is not |
| The `event` / `event_type` key mismatch | Not fixed — **removed**. Those blocks stop claiming to be emitted and point at the item that would emit them | Fixing a key in a payload nothing writes is busywork that makes the prose look more real than it is |
| Scope of the retirement | Prose only. No agent, skill, or executor behaviour changes | Nothing calls the recorder, so there is nothing to migrate. If a phase finds a real caller, STOP and escalate — that would mean the premise is wrong |
| What replaces the AOS opt-in guarantee | Nothing needs to. `run-events.jsonl` is local, project-scoped, and written only by `loom run` | The opt-in guarantee protected users from telemetry they did not ask for. A local audit log beside a run's own state, written only when they run the executor, does not engage that concern — and the network exporter it *would* engage is already opt-in (L3.8) |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, max 4 parameters.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**64.2%** as of epic 85). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- Deleting prose is the point of this epic, but deleting a *reference* without redirecting it is
  not. Every file that points at `event-recorder.md` or `events.jsonl` must end this epic pointing
  somewhere true.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — The enum is the source of truth — UNBLOCKED

1. Give every `EventKind` a documented meaning in Go — a description and which fields of `Event` it
   populates — beside the constant, so the documentation cannot drift from the constant it
   describes.
2. Extend `cmd/gen-schemas` to generate `shared/schemas/telemetry/run-event.schema.json` and a
   markdown table of the event types from the enum. Never hand-edit either.
3. The fitness function L3.9 asks for: a test that fails when a type is added without a
   description, and the existing drift test pattern for the generated artifacts.
4. Tests: every emitted `EventKind` is described; generated output matches the committed files; an
   event with an unknown kind fails to validate against the generated schema.

**Done when**: adding a sixteenth event type without documenting it fails the build — L3.9's stated
done-when, adjusted to the fifteen types that turned out to be real. **Commit** (`feat(telemetry):
generate the event vocabulary from one Go enum`), report, PAUSE.

### Phase B decisions (settled 2026-09-01, before the phase)

| Decision | Choice | Rationale |
|---|---|---|
| Historical records are left alone | `shared/agents/CHANGELOG.md`, `docs/human-tasks.md`'s v3.0.0 release verification, `docs/audits/`, and the `docs/aos/` design pack are **not edited** | A changelog entry and a dated verification log ("events.jsonl NOT created — PASS") are correct as records of what was true then. Editing them to match today makes the history lie. Only *live instructions* get redirected |
| `event-schema.md` is deleted, not rewritten | Both it and `event-recorder.md` go. `shared/telemetry/README.md` is the one surviving file | Keeping a schema document for a file nothing will ever write is close to the thing being retired. The README carries the honest content: the run event timeline, the generated vocabulary, the L3.8 traces, and the specified-but-unemitted list |
| `shared/telemetry/` survives as a directory | Not folded into `shared/schemas/telemetry/` | The AOS design pack made telemetry a first-class top-level concern deliberately; quietly reversing that as a side effect of a cleanup is the wrong way to revisit it |

## Phase B — Retire the layer that never emitted — BLOCKED BY Phase A

1. Delete `shared/telemetry/event-recorder.md` and `shared/telemetry/event-schema.md`. Rewrite
   `shared/telemetry/README.md` around what exists: the generated vocabulary, the run event
   timeline, and the OTel traces from L3.8 — plus the specified-but-unemitted list.
2. Record the nine specified-but-unemitted types in one place, each pointing at the roadmap item
   that would emit it. Remove their payload examples from `policy-schema.md`,
   `policy-evaluator.md`, and `audit-composition-pattern.md`, replacing each with a pointer — the
   `event` / `event_type` mismatch disappears with them.
3. Redirect every remaining **live instruction** referencing `event-recorder` or
   `.claude/telemetry/events.jsonl` — `deliver-feature`, `extract-lessons`, `scheduler`,
   `retrospective`, `policies/README.md`, `evaluation/README.md`, `DOMAIN_DICTIONARY.md`,
   `docs/ARCHITECTURE.md`, `docs/runbooks/parallel-delivery.md`. `retrospective` in particular
   still mines `gate_decision` per feature and should read `artifact.corrected` instead, the way
   `extract-lessons` was taught to in epic 85. A dangling reference is a worse outcome than the
   contradiction being fixed.
4. If any of them turns out to have a real caller or a real emitter, STOP and escalate.

**Done when**: no file describes a telemetry file nothing writes, and every reference resolves.
**Commit** (`refactor(telemetry): retire the event layer that never had an emitter`), report, PAUSE.

## Phase C — Roadmap and docs — BLOCKED BY Phase B

1. Roadmap: L3.9 **SHIPPED**, with what it did and did not do. Update **L3.5** to say it inherits a
   generated vocabulary rather than needing to invent one. Check L3.9's blocks/blocked-by against
   reality, the way epics 84 and 85 both had to.
2. `README.md`, `cmd/loom/README.md`, `docs/ARCHITECTURE.md`, `shared/DOMAIN_DICTIONARY.md` if a
   term changes meaning — "Run Event Timeline" in particular now has a generated schema.
3. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the telemetry story is one vocabulary, one log, and one generated schema, described
in the present tense only where it is true. **Commit** (`docs(telemetry): one event vocabulary`),
report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Emitters for any of the nine unemitted types — each belongs to its own item
- The episodic store (**L3.5**), which consumes this vocabulary
- Changing what the executor emits, or when
- Metrics or logs signals; L3.8 shipped traces only and that stands

## Report format (end of every phase)

```
## Epic 86 Phase <X> Report
- Roadmap item: L3.9 — one event vocabulary
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
