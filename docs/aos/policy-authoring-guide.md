# Policy Authoring Guide

v3.3 — AOS Phase 4

This guide explains how to write, test, and maintain project-level policies that enable
**graduated automation** in the `FeatureDeliveryWorkflow`. Read `shared/policies/README.md` first
for the architectural overview; this document is the hands-on reference.

---

## 1. The 8 Approval Gates — Policy-Eligibility Classification

`shared/rules/approval-gates.md` defines 8 non-negotiable gates. Phase 4 classifies each gate
against the `automated-delivery-design.md` Tier A/B/C model:

| # | Gate | Risk class | Policy-eligible? | Automation tier | Rationale |
|---|---|---|---|---|---|
| 1 | Shipping to Friday (POST to dashboard) | External side-effect | **No — always human** | None | Modifies shared external metrics; cannot be undone via git revert |
| 2 | Creating a Git Commit | Repository history | **Yes** | Tier A | Revertable; risk is proportional to diff size and reviewer verdict |
| 3 | DB Migrations (Expand/Migrate) | Stateful infrastructure | **No — always human** | None | Modifies live database state; infrastructure blast radius too high for automation |
| 4 | DB Contracting Phase (DROP/RENAME) | Data loss | **No — always human** | None | Irreversible data destruction; no safe automation path |
| 5 | Posting to External APIs | External side-effect | **No — always human** | None | Third-party mutations; no rollback guarantee |
| 6 | Writing Files out of Boundary | System structure | **Yes** | Tier A | Safe when target paths are within a known allowlist |
| 7 | Wiring a New Fitness Function | CI/CD pipeline | **Yes** | Tier A | Safe when dry-run passes and the wired function is a new test (not modifying existing checks) |
| 8 | Deploying to Environment | Production stability | **No — always human** | None | Deployment failures can cause downtime; human judgment essential |

**Summary**: Gates 2, 6, and 7 are policy-eligible (Tier A). Gates 1, 3, 4, 5, and 8 are
permanently human-only — the evaluator ignores policies targeting them and always returns
`require-human`.

---

## 2. Schema Quick Reference

Full schema: `shared/policies/policy-schema.md`. The essentials:

```yaml
name: <kebab-case unique identifier>
version: "1.0"
description: <one sentence>
enabled: true

matcher:
  gate: <gate-id>               # single gate
  # OR
  gates: [<gate-id>, ...]       # multiple gates

condition:
  <key>:
    <operator>: <value>         # lessThan, equals, greaterThan, allMatch, anyMatch, noneMatch

action:
  type: auto-approve | auto-reject | require-human | escalate
  reason: <string written to telemetry event>
```

---

## 3. Writing Your First Policy

**Scenario**: your team writes documentation daily and wants to skip the manual "approve commit"
prompt for doc-only changes.

Step 1 — Create the policy file:

```
.claude/policies/auto-approve-docs.policy.yaml
```

Step 2 — Fill in the YAML:

```yaml
name: auto-approve-docs
version: "1.0"
description: "Auto-approve commits when the entire diff is docs-only and tests pass"
enabled: true

matcher:
  gate: git-commit

condition:
  diffType:
    equals: "docs-only"
  testsPass: true
  diffLines:
    lessThan: 500

action:
  type: auto-approve
  reason: "Docs-only diff, tests green — no behavioral risk"
```

Step 3 — Dry-run validate (see Section 5 below).

Step 4 — Commit the policy file and it becomes active for the next pipeline run.

---

## 4. Inversion Policies (Requiring Humans)

Use `require-human` to *force* human review regardless of what other policies say.
`require-human` always beats `auto-approve` in conflict resolution — this is the safety default.

Example: always require human review for security/auth file changes:

```yaml
name: require-human-on-auth
version: "1.0"
description: "Any change touching auth/ always requires human confirmation"
enabled: true

matcher:
  gate: git-commit

condition:
  filePaths:
    anyMatch: "**/auth/**"

action:
  type: require-human
  reason: "Auth paths require human review per security policy"
```

Place inversion policies in your `.claude/policies/` directory alongside auto-approve policies.
The evaluator automatically applies conflict resolution — `require-human` wins.

---

## 5. Testing a Policy (Dry-Run Mode)

Before enabling a policy on live runs, validate it against historical telemetry:

```bash
/orchestrate --workflow feature-delivery --spec <file> --dry-run-policies
```

In dry-run mode:
- The evaluator reads past pipeline runs from `.claude/telemetry/events.jsonl`
- Evaluates your new policy against those historical contexts
- Writes results to `.claude/telemetry/events-dryrun.jsonl` — no pipeline state is mutated
- Prints a summary showing which historical runs would have been auto-approved vs. would have
  required human review

**Reading dry-run output**: look for `"event": "policy.evaluated"` entries in the dryrun log.
`"conditionMet": true` means the policy would have fired; `"action": "auto-approve"` means it
would have auto-approved that run.

---

## 6. Emergency Override

To disable ALL policies for a project immediately (without deleting policy files):

```yaml
# .claude/delivery-policy.yaml
policiesEnabled: false
```

This is a global kill-switch. Every gate reverts to `require-human` until the flag is removed or
flipped to `true`. No policy file is deleted. Use this during an incident or when debugging a
policy that is firing unexpectedly.

To disable a single policy without deleting it, set `enabled: false` in that policy's YAML.

---

## 7. Conflict Resolution Reference

When two policies match the same gate with different actions:

| Conflict | Winner | Telemetry event |
|---|---|---|
| `auto-approve` vs `require-human` | `require-human` | `policy.conflict` |
| `auto-approve` vs `auto-reject` | `auto-reject` | `policy.conflict` |
| `auto-approve` vs `escalate` | `escalate` | `policy.conflict` |
| Multiple `auto-approve` | `auto-approve` | none |
| Multiple `require-human` | `require-human` | none |

The `policy.conflict` event includes both policy names so you can trace and deduplicate.

---

## 8. Audit Trail Requirements

Every policy decision — including no-op decisions where condition is not met — is recorded in
`.claude/telemetry/events.jsonl`. Telemetry is a **non-negotiable** requirement for policies.
If your project does not have `.claude/telemetry/` configured, enabling auto-approve policies
is strongly discouraged — you will have no audit trail for automated decisions.

See `shared/telemetry/event-schema.md` for the full event format.

---

## 9. Governance Pairs and Policy Safety

Phase 2 counter agents (from `docs/aos/governance-pairs.md`) continue to run as audit-after-
producer steps in the workflow. Policies evaluated by the Phase 4 evaluator run *after* audit
steps complete — a policy can only auto-approve a gate if the counter agent's audit also passed.
The audit is not bypassed by a policy.

Example flow for a git-commit gate:
1. Developer produces `implementation-notes.md`
2. `code-reviewer` evaluates → produces `code-review-report.md`
3. **Audit**: `code-reviewer` verdict surfaces in `pipelineContext.codeReviewer.verdict`
4. **Policy evaluator** checks all matching policies against that context
5. If `auto-approve`: gate proceeds without human prompt + telemetry event emitted
6. If `require-human`: human is prompted (same as v3.2 behavior)

---

## 10. Reference: Sample Policies

See `shared/policies/examples/` for three reference policies:

| File | Gate | Pattern |
|---|---|---|
| `auto-approve-refactor.policy.yaml` | `git-commit` | Auto-approve pure refactors (reviewer APPROVED + no behavior change) |
| `auto-approve-doc-changes.policy.yaml` | `git-commit` | Auto-approve docs-only commits |
| `auto-approve-test-additions.policy.yaml` | `fitness-function-wiring` | Auto-approve new test fitness functions |
| `require-human-review-security.policy.yaml` | `git-commit`, `out-of-boundary-write` | Inversion: force human on security/auth paths |

Copy these to `.claude/policies/` in your project and adapt as needed.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 4
Policy Layer. Licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md).*
