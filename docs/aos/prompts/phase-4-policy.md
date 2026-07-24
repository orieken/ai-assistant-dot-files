# AOS Migration — Phase 4: Policy Layer (v3.3)

You are executing Phase 4 — the final AOS migration phase. Scope: 8 operations. Introduces the policy layer that makes graduated automation (per `docs/prompts/automate-deliver-feature.md`) possible.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prerequisite

**Phase 3 (v3.2) must be complete.** Verify:
- `shared/rag/` layer exists with adapter interface + LLM-as-retriever
- `shared/orchestration/` layer exists with `/orchestrate` runtime skill
- `shared/hooks/` has Learning + Forgetting wiring
- `install.sh` has `--base` and `--full` modes
- `FeatureDeliveryWorkflow` + `TDDWorkflow` exist (trinity refactor from Ops 3.11-3.13)
- `docs/aos/migration-guide.md` documents opt-in paths for all AOS layers

If not complete, halt and execute `docs/aos/prompts/phase-3-runtime.md` first.

## Source of truth

`docs/aos/migration-plan.md` — Phase 4 section (Ops 4.1–4.8). Also read:
- `shared/rules/approval-gates.md` — the 8 non-negotiable gates that policy MUST distinguish between (some can become policy-driven, some must stay human forever)
- `docs/prompts/automate-deliver-feature.md` — the DESIGN prompt that motivates Phase 4 (may or may not have been executed yet; the design doc `docs/aos/automated-delivery-design.md` may already exist and provide the Tier A/B/C framing to reference here)
- Phase 3's `FeatureDeliveryWorkflow` — this is where policy evaluation points live; Phase 4 wires policies INTO those points

## Backward-compatibility guarantee (non-negotiable)

A team on v3.2 that upgrades to v3.3 and does not create any `.claude/policies/` files MUST see zero behavior change. Default policy set is empty; all approvals still require human confirmation. Policies are strictly opt-in per-project.

## Scope: 8 ops

### Op 4.1 — Design `shared/policies/` (or `shared/orchestration/policies/`)

Files to create:
- `shared/policies/README.md` — policy layer's opt-in nature, audit-trail requirements, per-project storage
- `shared/policies/policy-schema.md` — declarative rule format (matcher + condition + action)
  - matcher: which pipeline stage does this policy apply to (e.g., `code-reviewer-gate`, `doc-changes-gate`)
  - condition: what conditions trigger auto-approval (e.g., `codeReviewer.approved AND securityReviewer.approved AND diffSize < 200 AND fitnessFunction.allPass`)
  - action: `auto-approve` | `auto-reject` | `require-human` | `escalate`
- Example policy stub: `shared/policies/examples/auto-approve-refactor.policy.yaml`

**One commit**: `feat(policies): scaffold shared/policies/ layer + schema (AOS Phase 4 Op 4.1)`.

### Op 4.2 — Implement policy evaluator in orchestration runtime

Update `shared/orchestration/` from Phase 3:
- Add `policy-evaluator.md` spec — reads policies from `.project-ai/policies/` (or configured path)
- Evaluates against telemetry events + agent outputs (from Phases 1 + 3)
- Emits telemetry event for every policy decision (audit trail requirement is non-negotiable — no silent auto-approvals)
- Integration point: `FeatureDeliveryWorkflow`'s stage boundaries (from Op 3.11) call the evaluator before proceeding past a gate; evaluator returns proceed/halt/require-human

**One commit**: `feat(orchestration): implement policy evaluator (AOS Phase 4 Op 4.2)`.

### Op 4.3 — Add 3 sample policies

Under `shared/policies/examples/`:
- `auto-approve-doc-changes.policy.yaml` — for gate #2 (git commit) when diff is docs-only
- `auto-approve-test-additions.policy.yaml` — for gate #7 (fitness function wiring) when the wired function is a new test
- `require-human-review-security.policy.yaml` — INVERSION example: explicitly requires human on any file matching `**/security/*` or `**/auth/*` regardless of other policies

**One commit**: `feat(policies): add three sample policies (AOS Phase 4 Op 4.3)`.

### Op 4.4 — Write `docs/aos/policy-authoring-guide.md`

Comprehensive guide covering:
- The 8 gates from `approval-gates.md` — classified per `automate-deliver-feature.md`'s Design section as policy-eligible vs. always-human
- Schema reference with examples
- Testing a policy (dry-run mode — policies can be evaluated against historical telemetry without actually mutating)
- Emergency override: how to disable all policies for a project

**One commit**: `docs(aos): write policy-authoring-guide (AOS Phase 4 Op 4.4)`.

### Op 4.5 — Update `shared/rules/approval-gates.md`

Add an optional "policy-based" gate type alongside existing human-only gates. For each of the 8 gates, add a "Policy-eligible: yes/no" annotation, matching the classification from `automate-deliver-feature.md`'s Design section.

If `automate-deliver-feature.md`'s DESIGN prompt has not been executed yet, halt and execute it first — this op depends on the classification that design produces.

**One commit**: `docs(rules): classify approval gates for policy eligibility (AOS Phase 4 Op 4.5)`.

### Op 4.6 — Update `health-check` to audit policy syntax + coverage

Extend `shared/skills/health-check/SKILL.md`:
- Detect any `.claude/policies/` directory in the running project
- Validate each policy file's syntax against `policy-schema.md`
- Report which gates each policy covers; flag any gate with contradictory policies (e.g., one says auto-approve, another says require-human — halt with a coverage conflict error)
- Never fails on absence of policies (opt-in guarantee)

**One commit**: `feat(health-check): validate policy syntax + coverage (AOS Phase 4 Op 4.6)`.

### Op 4.7 — CHANGELOG v3.3.0 entry

Update `shared/agents/CHANGELOG.md` with the v3.3.0 entry — policy layer + evaluator + samples + guide + gate classification + health-check bump. Restate the backward-compat promise for the fifth time.

**One commit**: `docs(changelog): v3.3.0 — AOS policy layer (AOS Phase 4 Op 4.7)`.

### Op 4.8 — Identity install verification

Same shape as Phases 1-3:
- `bash scripts/health-check.sh --verbose` — 0 FAILs
- `bash scripts/check-parity.sh` — PASS
- Verify: v3.3 install with no `.claude/policies/` behaves identically to v3.2 (no auto-approvals fire; every gate requires human confirmation)
- Run a `FeatureDeliveryWorkflow` in a scratch project WITH a single test policy (e.g., auto-approve-doc-changes) — verify the auto-approve fires + emits a telemetry event AND does not fire on a non-docs-only change

**One commit** (if all checks pass): `chore(release): tag AOS Phase 4 v3.3.0 candidate (AOS Phase 4 Op 4.8)` — also updates `migration-plan.md` marking Phase 4 checklist items complete + records Op 4.8 verification result.

## AOS Migration Complete

After Op 4.8, the AOS migration is fully shipped: v3.0 → v3.3, four phases, all with backward-compat preserved. The Success Criteria in `migration-plan.md`'s bottom section are the exit test — verify all 7 criteria pass before recommending the AOS-complete announcement.

## Commit discipline (non-negotiable)

- ~9 commits total for Phase 4.
- Conventional Commits.
- **NEVER `git add -A`.**
- Green health-check + tests per commit.
- Do NOT push.

## Escalation criteria

- Op 4.5 finds that `automate-deliver-feature.md`'s DESIGN prompt hasn't been executed — halt, execute it first, then come back.
- The policy schema needs a full expression language (e.g., CEL, Rego) — halt, propose. YAML with a fixed condition vocabulary should be sufficient for Phase 4; more expressive comes later.
- Any of the 8 gates in `approval-gates.md` looks like it needs re-classification vs. what `automate-deliver-feature.md`'s design specifies — halt, ask.
- Policy evaluator's telemetry event emission would require modifying telemetry event schema from Phase 1 — halt, describe. The schema was designed to be extensible; extensions should be additive.

## Report format (under 300 words)

```
PHASE 4 STATUS: <complete | stopped-at-op-N>

Commits landed:
  <sha> <message>
  ...

Op-by-op tally:
  Op 4.1 policies scaffold + schema: <landed>
  Op 4.2 policy evaluator: <landed>
  Op 4.3 sample policies: <landed>
  Op 4.4 policy-authoring-guide: <landed>
  Op 4.5 approval-gates classification: <landed | blocked on automate-deliver-feature design>
  Op 4.6 health-check policy validation: <landed>
  Op 4.7 CHANGELOG v3.3.0: <landed>
  Op 4.8 Identity check + policy test: <passed | failed with details>

AOS Migration Success Criteria (from migration-plan.md §Success Criteria):
  1. v3.3 is live: <yes | no>
  2. v2.x → v3.3 --base = zero behavior change: <verified | failed>
  3. v3.3 --full + policy = eliminates one class of just-say-yes gate: <verified>
  4. Telemetry-backed evaluations replace ad-hoc agent-scorecard: <yes | still ad-hoc>
  5. Learning engine proposes ≥ 1 KI candidate per month: <n/a until observed>
  6. Forgetting engine flags ≥ 1 stale KI per quarter: <n/a until observed>
  7. No skill/agent/rule from v2.x removed: <verified>

Recommended next step:
  <e.g., "AOS migration complete — announce internally, human review + git push + tag v3.3.0">
```

Go.
