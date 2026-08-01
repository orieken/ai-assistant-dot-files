# Policy Layer

The policy layer enables **graduated automation** for the `FeatureDeliveryWorkflow` and other AOS
pipelines. Policies are **strictly opt-in per-project** — the absence of any policy file guarantees
identical behavior to v3.2.

---

## What policies do

Each policy declares a *matcher* (which pipeline gate it watches), a *condition* (what must be true to
trigger), and an *action* (`auto-approve`, `auto-reject`, `require-human`, `escalate`). The
`FeatureDeliveryWorkflow` calls the policy evaluator at every stage boundary; if no matching policy
exists, the stage boundary defaults to `require-human` — the same behavior as all prior versions.

A team that upgrades to v3.3 but places no `.claude/policies/` files in their project sees **zero
behavior change**.

---

## Where to put project policies

Policies live in the *project* (not in this framework repo) at:

```
.claude/policies/<your-policy-name>.policy.yaml
```

The evaluator loads every `.policy.yaml` file found in that directory at pipeline startup. File order
is undefined — write policies that do not depend on evaluation order.

---

## Audit trail (non-negotiable)

Every policy decision — whether it auto-approves, rejects, escalates, or falls through to a human —
MUST emit a telemetry event to `.claude/telemetry/events.jsonl`. There are no silent auto-approvals
in this framework. Teams that disable telemetry lose their audit trail and should not enable
auto-approve policies.

Event shape: see `shared/telemetry/event-schema.md`, event type `policy.evaluated`.

---

## Schema

The declarative policy format is documented in `shared/policies/policy-schema.md`.

## Examples

`shared/policies/examples/` contains three reference policies demonstrating common patterns:
- `auto-approve-refactor.policy.yaml` — Tier A auto-proceed on pure-refactor commits
- `auto-approve-doc-changes.policy.yaml` — auto-proceed when diff is docs-only
- `auto-approve-test-additions.policy.yaml` — auto-proceed when adding new tests
- `require-human-review-security.policy.yaml` — inversion: force human regardless of other policies

---

## Emergency override

To disable all policies for a project without deleting them:

```yaml
# .claude/delivery-policy.yaml
policiesEnabled: false
```

This is a global kill-switch. All gate decisions revert to `require-human`. The policy files are
preserved — re-enable by removing or flipping the flag.

---

## Gate classification

Not every approval gate is policy-eligible. Gates 1, 3, 4, 5, and 8 from
`shared/rules/approval-gates.md` are permanently human-only regardless of any policy you write —
the evaluator ignores policies targeting those gates and always returns `require-human`. See
`docs/aos/policy-authoring-guide.md` for the full classification and rationale.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 4
Policy Layer. Licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md).*
