# Epic 73 — AOS Policy Engine Completion

Source: `docs/audits/framework-audit-2026-08-07.md` §3 item 5; completes AOS Phase 4.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

AOS Phase 4 shipped `shared/orchestration/policy-evaluator.md` (the evaluator design spec)
and `shared/policies/examples/` (3 sample policy YAML files). The feature is documented in:
- `shared/rules/approval-gates.md`: 3 Tier-A gates are policy-eligible: `git-commit` (gate 2),
  `out-of-boundary-write` (gate 6), `fitness-function-wiring` (gate 7).
- `docs/aos/policy-authoring-guide.md` (if present) — verify at execution time.

What is NOT yet present:
- `.claude/policies/` directory — the project-level activation path described in
  `approval-gates.md`. No policies are active.
- A runnable evaluator (`scripts/evaluate-policy.sh` or equivalent). The spec describes the
  contract; no implementation executes it.
- `health-check.sh` checks for policy schema validity or evaluator presence.
- Any wiring in `shared/skills/deliver-feature/SKILL.md` that calls the evaluator at gates.

The `shared/orchestration/policy-evaluator.md` spec is the authoritative source of truth for
what the evaluator must do. Read it in full before Phase A.

## Scope

**Phase A — Gap Audit (one commit, then PAUSE for user approval):**

Draft and commit as `docs(policies): policy engine gap audit + implementation contract (Epic 73 Phase A)`:

Produce `docs/aos/policy-engine-gap-audit.md` answering:

1. **What the spec says must exist** (enumerate from `policy-evaluator.md`):
   - Inputs the evaluator reads (policy file path, gate ID, context fields)
   - Decision outputs: `approve` / `require-human` / `deny`
   - Telemetry: `policy.evaluated` event
   - Which gates are Tier A and their condition criteria (from `approval-gates.md`)

2. **What already exists** (verify each):
   - `shared/policies/examples/`: list files and summarize what each policy contains
   - `docs/aos/policy-authoring-guide.md`: exists? If yes, summarize its coverage gaps.
   - Any call sites in `deliver-feature` or other skills that invoke an evaluator today

3. **Implementation contract** (decide before Phase B):
   - Evaluator technology: bash script (portable, no new dependencies) or Go binary (type-safe,
     requires compile step)? Rationale.
   - Invocation pattern: how does a skill call the evaluator? (`result=$(bash scripts/evaluate-policy.sh --gate git-commit --context ...)`)
   - Telemetry emission: where does `policy.evaluated` get written? (append to `pipeline-trace.json`?)
   - `.claude/policies/` vs `shared/policies/examples/`: what is the distinction? Project-level
     active policies live in `.claude/policies/`; examples are templates in `shared/policies/examples/`.

**Phase B — Implementation (after approval; one commit per op):**

Op 1 — `feat(scripts): evaluate-policy.sh policy evaluator (Epic 73 Op 1)`:
- `scripts/evaluate-policy.sh --gate <gate-id> [--context <json>]`
- Reads all `.yaml` files in `.claude/policies/` (if the dir exists; no-op if absent).
- For each policy, checks its `gate:` field matches `--gate` and evaluates its `conditions:`.
- Returns: `approve` | `require-human` | `deny` (stdout, exit 0).
- Emits a `policy.evaluated` JSON line to stderr (or a designated telemetry file).
- If `.claude/policies/` does not exist, returns `require-human` (safe default).
- Must be idempotent and side-effect-free (no file writes to policy files).

Op 2 — `feat(policies): bootstrap .claude/policies/ with refactor auto-approve (Epic 73 Op 2)`:
- Create `.claude/policies/auto-approve-refactor.policy.yaml` following the shape in
  `shared/policies/examples/auto-approve-refactor.policy.yaml` (adapt for this project).
- This policy covers gate `git-commit` when: diff is below a configured line threshold,
  no security/auth paths are touched, and all tests pass (condition fields TBD from spec).
- Do NOT create policies for the other Tier-A gates yet — one working policy first.

Op 3 — `feat(orchestration): wire evaluator into deliver-feature gates 2/6/7 (Epic 73 Op 3)`:
- Update `shared/skills/deliver-feature/SKILL.md`: at gates 2, 6, and 7, instruct the
  executing agent to call `bash scripts/evaluate-policy.sh --gate <gate-id>` before
  presenting the human approval prompt. If the result is `approve`, skip the human prompt
  and log `[POLICY] gate <id> auto-approved by <policy-file>`. If `deny`, halt the
  pipeline and report which policy denied and why.
- Same update for `shared/skills/deliver-bugfix/SKILL.md` at gate 2.

Op 4 — `feat(health-check): policy schema + evaluator checks (Epic 73 Op 4)`:
- FAIL-level: if `.claude/policies/` exists, every `.yaml` file must validate against
  `shared/schemas/policy-schema.json` (create the schema if it doesn't exist).
- WARN-level: if `.claude/policies/` does not exist, emit "no policies active — gates are
  always human-approved".
- FAIL-level: `scripts/evaluate-policy.sh` must be executable.
- `bash scripts/health-check.sh` green.

Op 5 — `docs(aos): policy-authoring-guide.md (Epic 73 Op 5)`:
- Create or update `docs/aos/policy-authoring-guide.md`: how to write a policy, what
  fields are valid, how conditions are evaluated, how to test a policy with
  `evaluate-policy.sh --dry-run`, and the Tier-A gate IDs.
- Cross-reference from `shared/rules/approval-gates.md` policy-eligible gate annotations.

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If Phase A finds that `policy-evaluator.md` describes an event-driven model (the evaluator
  listens on a message bus rather than being called synchronously) that cannot be implemented
  as a simple bash script, halt and propose a simplified synchronous contract before Phase B.
- If the condition evaluation in Op 1 requires parsing the diff or running tests (not just
  reading context fields), that Op becomes much larger — halt and scope it as a separate
  sub-epic rather than embedding complex logic in a single script.
- If `.claude/policies/` should NOT be committed to the repo (because it contains
  project-specific tuning that shouldn't be in a framework template repo), halt and clarify
  the boundary: should policies live in `.claude/policies/` (committed) or a gitignored
  user-local path?
- If `docs/aos/policy-authoring-guide.md` already exists and covers the Op 5 content,
  update it rather than replacing it, and reduce Op 5 to a gap-fill.

## Report (under 250 words)

```
Phase A commit: <sha>
Phase A findings:
  - Spec says must exist: <N items>
  - Already exists: <list>
  - Missing: <list>
  - Implementation contract:
    - Evaluator tech: <bash | Go binary>
    - Invocation pattern: <exact call signature>
    - Telemetry target: <file path>
    - .claude/policies/ boundary ruling: <commit | gitignore>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, evaluate-policy.sh --gate git-commit returns <approve>,
deliver-feature gate 2 auto-approve exercised.
```

Go.
