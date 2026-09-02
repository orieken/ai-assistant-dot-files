# Epic 87 — A Policy Evaluator That Runs (L2.16)

Source: session discussion, 2026-09-02, following epic 86. Implements roadmap item **L2.16**, whose
blocker **L2.13** shipped in epic 77. It **blocks L4.3** (budget governor), whose other blocker
L3.8 shipped in epic 84 — so L2.16 is the last thing standing in front of it.

Epic 86 made three gaps visible and named this item behind all of them:

- `shared/rules/approval-gates.md` — a hard-constraint file — now states that "no silent
  auto-approvals" has **no audit trail** and that nothing enforces it.
- `shared/orchestration/policy-evaluator.md` records that its `policy.evaluated`,
  `policy.conflict` and `policy.skipped` events were specified and never emitted.
- `--dry-run-policies` has nothing to read, because no history was ever written.

Making a gap visible and leaving it is worse than not having looked.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context (read before any phase)

1. `docs/roadmaps/BUILD-ROADMAP.md` — **L2.16**, **L4.3** (what it unblocks), **L2.13** (the gate
   enforcement this must not weaken).
2. `shared/policies/policy-schema.md` — the condition vocabulary, the gate-ID table, and the
   always-human column. Note that the "Scalar checks" block is **invalid YAML**: `diffType` and
   `filePaths` each appear twice in one map. No parser has ever read this file.
3. `shared/orchestration/policy-evaluator.md` — the evaluation order, conflict resolution, and the
   sections epic 86 rewrote to say nothing is emitted.
4. `internal/orchestrator/gate.go`, `executor.go` `checkGate` — where a gate halt happens today, and
   the property that must survive: nothing an agent returns can approve a gate.
5. `internal/orchestrator/vocabulary.go` — the event enum. A policy decision that is actually
   emitted becomes a real member here, and the fitness function will require it documented.
6. `internal/orchestrator/loop.go` `conditions()` — epic 82's precedent: named conditions resolved
   in Go rather than an expression language. This epic extends that reasoning rather than reversing
   it.

## Design decisions (fixed — do not relitigate in-phase)

| Decision | Choice | Rationale |
|---|---|---|
| No expression language | A typed Go evaluator over the declared schema. **Not** CEL, not Rego, despite the roadmap naming both | The condition vocabulary is closed and small — about ten checks over facts the executor already holds. CEL buys arbitrary boolean logic at the cost of a dependency, a cost-limiting surface, and an authorization decision written in a language most readers of a `.policy.yaml` will not know. Epic 82 reached the same conclusion for loop conditions; the reasoning did not change because the subject did |
| Extending the vocabulary is a code change | Adding a condition check means adding a typed struct and a test, not editing a config | That is the honest trade for the decision above, and it is the right one for authorization: a new way to skip a human review should not be addable by editing YAML |
| **This epic does not auto-approve anything** | The executor evaluates matching policies at each gate, records the decision, reports what *would* have happened — and still halts for a human | `loom run` has never auto-approved, and the first run that skips a barrier should not also be the first evidence the evaluator is correct. This closes the audit gap, makes `--dry-run-policies` real, and earns the right to honor a decision later. Honoring it is a small follow-up item this epic raises |
| The always-human list is compiled in | The five non-eligible gates are a Go constant, not a YAML field | A kill-switch that lives in the same prose the model may misread is the defect being removed. A policy targeting one of them is rejected **at load time**, not skipped at evaluation time |
| Load-time rejection over silent skipping | `policy-schema.md` says policies targeting non-eligible gates are "silently ignored". They now fail to load, loudly | Silently ignoring a policy someone wrote to auto-approve a deployment is how a person concludes it worked. This is a behaviour change from the spec, and the spec is what changes |
| The decision becomes a real event | `policy.evaluated` joins the event vocabulary and is emitted on `run-events.jsonl` | This is what makes "no silent auto-approvals" demonstrable rather than asserted, and it is only legitimate to add to the enum now because something finally emits it (L3.9's rule) |
| Invalid YAML fails | The schema doc's own duplicate-key block is fixed, and duplicate keys are a load error | L2.16's done-when names this. A parser that tolerates the file that has never been parsed would prove nothing |
| Conflict resolution stays as specified | `require-human` beats `auto-approve`; conflicts are recorded | The rule is right. What was missing was anything that applied it |

## Shared guardrails (all phases)

- Conventional Commits; commit at the end of each phase, then PAUSE for human review before the
  next phase. Never `git push`.
- Go work follows `shared/rules/go-conventions.md` + `shared/rules/design-principles.md`:
  complexity < 7 (`golangci-lint run ./...` must pass — that IS the build gate), table-driven
  tests, no `interface{}`/`any` in exported signatures, interfaces at the consumer, max 4
  parameters (introduce a parameter object rather than a fifth).
- Every phase ends by running `go build ./...`, `go test ./...`, `golangci-lint run ./...`, and
  `scripts/health-check.sh` + `scripts/check-parity.sh` when `shared/` content changed. Coverage
  must stay ≥ the CI ratchet floor (**64.6%** as of epic 86). Integration tests driving the built
  binary are invisible to the coverage tool — if the floor is threatened, add in-process tests,
  never lower the floor.
- **Nothing in this epic may make a gate easier to pass.** If a phase finds itself weakening
  `checkGate`, STOP and escalate — evaluating is not approving, and the two must stay separable in
  the code as well as in the prose.
- If a decision point arises that this file doesn't settle, STOP and escalate rather than
  guessing; record the open question in the phase report.

---

## Phase A — Load and validate policies — UNBLOCKED

1. New `internal/policy/`: parse `.claude/policies/*.policy.yaml` into typed structs — matcher,
   condition, action — rejecting unknown fields and duplicate keys.
2. The always-human gate list as a compiled constant. A policy targeting one fails to load, naming
   the gate and the rule.
3. Load-time validation: unknown gate IDs, unknown condition checks, unknown action types, and a
   missing `name` are all errors with messages that say what to fix.
4. Fix `shared/policies/policy-schema.md`'s invalid YAML block, and add a fixture built from it.
5. Tests: the previously-invalid example now parses; a duplicate key fails; a policy targeting
   `deploy` fails at load with a message naming the always-human rule; every example under
   `shared/policies/examples/` loads.

**Done when**: a policy targeting an always-human gate is rejected at load time, and the invalid
YAML example fails to parse — L2.16's stated done-when, in full. **Commit** (`feat(policy): load and
validate policies against a compiled gate list`), report, PAUSE.

## Phase B — Evaluate, and record what was decided — BLOCKED BY Phase A

1. Typed evaluation of every declared condition check against a `GateContext` the executor builds
   from run state — diff size, test results, review verdict, security criticals, changed paths.
   Where a fact is not available, the check is **unknown**, and an unknown check never satisfies a
   condition.
2. Conflict resolution as specified: `require-human` wins; the conflict is recorded.
3. `policy.evaluated` joins the event vocabulary — it is emitted now, which is what earns it a place
   — and is written for every policy considered, matched or not.
4. The executor calls the evaluator at each gate and **still halts**. The CLI reports what a policy
   would have done, so the decision is visible before anyone trusts it.
5. Make `--dry-run-policies` real: evaluate against a completed run's recorded state and print the
   decisions, with no mutation.
6. Tests: each condition check, including its unknown case; conflicting policies resolve to
   require-human; a matched auto-approve policy does **not** approve anything; the decision reaches
   the timeline; dry-run mutates nothing.

**Done when**: every policy decision a run makes is recorded and reviewable, and no gate is skipped
because of one. **Commit** (`feat(policy): evaluate policies at gates and record every decision`),
report, PAUSE.

## Phase C — Docs, and the honesty pass — BLOCKED BY Phase B

1. `shared/rules/approval-gates.md`: replace epic 86's "no decision is recorded, and none ever has
   been" with what is now true — decisions are recorded and auditable, and **auto-approval still
   does not happen**. Do not overstate: the gate still stops for a human.
2. `shared/orchestration/policy-evaluator.md` and `shared/policies/`: describe the evaluator that
   exists, the load-time rejection, and the compiled always-human list. Remove the "silently
   ignored" behaviour.
3. Roadmap: L2.16 **SHIPPED**, stating plainly that honoring an auto-approve decision is
   deliberately not included and raising the follow-up item for it. Update **L4.3** with what it
   inherits.
4. `cmd/loom/README.md`, `README.md`, `docs/ARCHITECTURE.md`, `shared/DOMAIN_DICTIONARY.md` if a
   term appears. Bump `shared/VERSION`.
5. Run `scripts/health-check.sh` and `scripts/check-parity.sh` plus the full Go gate.

**Done when**: the docs describe a policy evaluator that runs, and are explicit that it advises
rather than approves. **Commit** (`docs(policy): an evaluator that runs, and does not yet approve`),
report, PAUSE — epic complete.

---

## Explicitly out of scope (do not build, even if it feels adjacent)

- **Honoring an auto-approve decision** — the follow-up item Phase C raises
- CEL, Rego, or any expression language
- The budget governor (**L4.3**), which consumes this
- New condition checks beyond those `policy-schema.md` already declares
- Policy evaluation in the markdown pipeline

## Report format (end of every phase)

```
## Epic 87 Phase <X> Report
- Roadmap item: L2.16 — a policy evaluator that runs
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
