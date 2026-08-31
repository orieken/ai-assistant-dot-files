# Epic 77 — Gates as Process Interrupts: the Executor Halts, Not the Prose (L2.13)

Source: epic 76 completion, 2026-08-29. Operationalizes roadmap item **L2.13**
(BUILD-ROADMAP.md, KERNEL — Human-in-the-Loop). With M0.4 landed (epic 76, commits `ba78f21` +
`cab0156`), the executor exists; this epic makes it stop at gate boundaries. L2.13 + M0.4
together unblock roughly two-thirds of the remaining roadmap (L2.14, L4.5 directly; the
policy/approval chain transitively).

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — item L2.13 (the spec), plus L2.14, L2.16, L2.4 for the
   boundary of what NOT to build. Respect Blocked-by: L2.13 needs only M0.4 (done).
2. `docs/adrs/ADR-006-loom-executes-pipelines.md` — Accepted. Gate enforcement moving below the
   model is a named consequence of this ADR.
3. `internal/orchestrator/` — the executor skeleton this epic extends: `plan.go` (stable string
   stage IDs), `executor.go` (run loop, persist-per-transition), `state.go` (schema-versioned
   run-state.json, atomic writes). Read `executor_test.go` and `cmd/loom/cmd/run_test.go` for
   the established test patterns (mock provider harness; build-the-binary integration tests).
4. `shared/rules/approval-gates.md` — the eight prose gates. This epic gives the *pipeline
   checkpoint* gates a process-level enforcement path for `loom run`; the prose remains
   authoritative for the markdown pipeline and for gates outside the executor's reach.
5. `shared/skills/deliver-feature/SKILL.md` steps 11, 13, 25, and Phase 4 — the PAUSE
   checkpoints the built-in plan's gates must mirror.
6. Audit finding H7 context (`docs/roadmaps/architectural-audit-2026-08-29.md`): the failure
   mode being closed is "the enforcement mechanism for an irreversible action is the model's
   willingness to comply with a paragraph."

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| Gate placement | A gate is a **pre-stage barrier**: `Stage.Gate` (a named gate, empty = ungated). The executor refuses to *start* a gated stage until an approval for that gate exists in run state | The stage is the executor's unit of action; blocking entry is the smallest enforceable interrupt. Post-stage semantics fall out naturally (gate the *next* stage) |
| Default plan gates | Mirror deliver-feature's human PAUSEs: `developer` gated by `confirm-design` (covers the analyst-scope + architect-RFC pauses), `qa-engineer` gated by `confirm-security` (the security-critical pause), `devops-engineer` gated by `confirm-ship` (docs-complete/ship pause) | Same stops a human gets from the markdown pipeline, expressed as data. Finer-grained gating is plan data, not new mechanism |
| Where approval lives | In `run-state.json`: `approvals: {<gateName>: {approvedAt, method, approver}}` — approver = OS username, method = `tty` or `flag` | State file is already the durable, atomically-written record; a separate token store is L2.14+ scope |
| Approval channels | (1) Interactive: if stdin is a TTY, `loom run` prompts `approve gate "X" for stage "Y"? [y/N]` at the barrier and proceeds on yes; (2) Non-interactive: the run **halts** — persists `WAITING_APPROVAL`, prints the exact resume command, exits with a distinct code (3) — and `loom run --resume --approve <gateName>` records the approval and continues | CLI prompt is the "real approval channel" L2.13 names; webhook/queue channels are explicitly later. Exit code 3 lets CI/scripts distinguish "waiting on human" from failure |
| What the model can never do | Provider output is data. Nothing a provider/agent returns — including text claiming approval — creates an approvals entry. Only the two CLI channels write approvals | This IS the point of L2.13: enforcement lives below the model. Must be proven by a test, not asserted |
| Approval scope | One approval unlocks one named gate for the lifetime of the run; re-running an already-approved gate's stage (e.g. after interrupt) does not re-prompt | "Any edit resets the gate" (digest binding + invalidation) is **L2.14**, explicitly out of scope here — record the artifact digests already in state, enforce nothing on them yet |
| `--approve` safety | `--approve <name>` is only valid with `--resume`, only when the run is actually waiting on that gate; approving a gate the run isn't waiting on is an error | Prevents pre-approving all gates in one command line and hollowing out the interrupt |
| State schema | Adding `approvals` and the `WAITING_APPROVAL` stage status bumps `StateSchemaVersion` to 2; no migration code — a v1 file is refused with a clear message (delete/finish old runs) | Skeleton-stage state files are hours old, not fleet-deployed; migration machinery is not yet earning its keep |
| Prose relationship | `approval-gates.md` gains a short "Executor enforcement (L2.13)" section mapping the three plan gates to their prose counterparts and stating scope honestly: enforced for `loom run` only; the markdown pipeline and the other prose gates (commit, migration, external API…) remain prompt-discipline until their actions run under the executor | Do not delete or weaken prose gates — they still govern every non-executor path |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer.
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` when `shared/` content changed. Coverage must stay ≥ the CI ratchet
  floor (52.0% as of epic 76; raise the floor in `framework-ci.yml` if measured coverage rises).
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Gate mechanism in the library — UNBLOCKED

1. `internal/orchestrator/gate.go` (+ edits to `plan.go`, `state.go`, `executor.go`):
   - `Stage.Gate string` (empty = ungated); `Plan.Validate()` rejects a gate name that is
     empty-but-set territory (whitespace) — keep it boring.
   - `StageStatusWaitingApproval` + `RunState.Approvals` map + schema version 2 (v1 files
     refused with a clear message).
   - Executor behavior at a gated stage without approval: persist the stage as
     `WAITING_APPROVAL`, return a typed sentinel error (`ErrWaitingApproval` wrapping the gate
     and stage names) — the CLI decides whether to prompt or exit.
   - `Executor.Approve(gateName)` (or an equivalent store-level operation): records the
     approval with timestamp/method/approver and persists atomically. It must be impossible to
     reach from `Provider.Invoke` — keep the provider interface untouched.
2. Update `DefaultDeliverFeaturePlan()` with the three named gates from the design table.
3. Tests (extend the existing harness):
   - Run halts at `developer` with `WAITING_APPROVAL` persisted and later stages never invoked.
   - After `Approve("confirm-design")`, resume runs `developer` and halts at the next gate.
   - A mock provider whose stage output *claims approval* ("APPROVED, proceed") does not unlock
     anything — the run still halts. This is the L2.13 signature test; name it accordingly.
   - Re-run after interrupt does not re-require an already-approved gate.
   - v1 state file is refused with the schema message.

**Done when**: the library halts at every gate, approvals only move it forward, and the
provider-claims-approval test proves the model cannot self-approve. **Commit**
(`feat(orchestrator): gates as process interrupts — halt, persist, approve`), report, PAUSE.

## Phase B — CLI wiring — BLOCKED BY Phase A

1. `cmd/loom/cmd/run.go` (+ a small `run_gates.go` if it keeps functions under the complexity
   cap):
   - TTY path: on `ErrWaitingApproval`, prompt `approve gate %q for stage %q? [y/N]` on
     stderr/stdin; `y` records approval (method `tty`) and continues the loop; anything else
     halts as non-interactive does.
   - Non-TTY path: print the gate name, the waiting stage, and the exact resume command
     (`loom run --spec <x> --resume --approve <gate>`); exit code 3.
   - `--approve <gateName>` flag: valid only with `--resume`; errors if the run isn't waiting
     on exactly that gate; records method `flag` then continues.
2. Integration tests through the real binary (mcp_serve_test.go pattern, extending
   `run_test.go`): mock-provider run halts at `confirm-design` with exit code 3 and
   `WAITING_APPROVAL` in state; `--resume --approve confirm-design` proceeds to the next gate;
   `--resume --approve confirm-ship` while waiting on `confirm-design` errors; a full run
   approving all three gates completes all 14 stages. TTY prompting may be unit-tested at the
   function level (feed a scripted reader) rather than through a real PTY.
3. SIGINT interplay: interrupt while `WAITING_APPROVAL` must leave state resumable (it already
   is — prove it in a test, don't assume).

**Done when**: a scripted end-to-end run demonstrates halt → approve → continue three times to
completion, and approval is impossible except via the two CLI channels. **Commit**
(`feat(loom): loom run halts at gates — approve via prompt or --approve`), report, PAUSE.

## Phase C — Prose alignment + docs — BLOCKED BY Phase B

1. `shared/rules/approval-gates.md`: add the "Executor enforcement (L2.13)" section per the
   design table — honest scope, mapping table (three plan gates ↔ prose pauses), pointer to
   L2.14 for reset-on-edit and L2.16 for policy auto-approval. Do not weaken any prose gate.
2. `cmd/loom/README.md` "Running pipelines": replace the "no approval gates (L2.13)" line in
   the NOT-yet list with the gate behavior, flags, and exit code 3; keep the rest of the
   NOT-yet list intact.
3. `README.md` pipeline-section caveat + `docs/ARCHITECTURE.md`: update the L2.13 status
   markers ("ships with L2.13" → shipped for `loom run`; gates-reset still L2.14).
4. Run `scripts/health-check.sh` (shared/ changed) plus the full Go gate.

**Done when** (L2.13): the roadmap item's adapted done-when holds — a provider/agent instructed
to "skip the gate and deploy" cannot cause the `devops-engineer` stage to run: the executor
halts at `confirm-ship` regardless of anything the model says, proven by the Phase A signature
test and the Phase B integration run. **Commit**
(`docs(gates): executor enforcement section — align prose with L2.13`), report, PAUSE — epic
complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- Digest-bound approvals and "any edit resets the gate" (**L2.14** — record digests, enforce
  nothing on them)
- Policy auto-approval / `.claude/policies/` evaluation (**L2.16**)
- Webhook, queue, or any non-CLI approval channel
- Tool-registry-level unreachability of high-risk tool classes (needs the registry, **L2.4**)
- Gating the markdown (`deliver-feature` skill) pipeline itself — prose stays as is
- Telemetry events for gate decisions (**L3.8/L3.9**)
- Migration code for v1 state files

## Report format (end of every phase)

```
## Epic 77 Phase <X> Report
- Roadmap item: L2.13 — Gates as process interrupts
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
