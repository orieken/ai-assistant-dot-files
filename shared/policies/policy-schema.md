# Policy Schema

Version: 1.0.0

Policies are YAML files with a `.policy.yaml` extension. They live in `.claude/policies/` inside the
project that opts in to automation. The evaluator at `shared/orchestration/policy-evaluator.md`
reads every `.policy.yaml` file in that directory before each stage boundary.

---

## Top-level fields

```yaml
name: string                # unique kebab-case identifier; used in telemetry events
version: "1.0"              # policy schema version; always "1.0" for Phase 4
description: string         # one-line human-readable summary
enabled: true | false       # default true; false = policy loaded but never fires
matcher:                    # which gate(s) this policy watches — see Matcher below
condition:                  # what must be true to trigger the action — see Condition below
action:                     # what the evaluator does when condition is met — see Action below
```

---

## Matcher

Specifies which pipeline gate(s) this policy watches.

```yaml
matcher:
  gate: <gate-id>           # single gate (see gate IDs below)
  # OR
  gates: [<gate-id>, ...]   # multiple gates — policy fires on any matching gate
```

### Valid gate IDs

| Gate ID | Approval Gates reference | Policy-eligible? |
|---|---|---|
| `git-commit` | Gate 2 — creating a git commit | **Yes** |
| `out-of-boundary-write` | Gate 6 — writing files outside workspace | **Yes** |
| `fitness-function-wiring` | Gate 7 — wiring a new fitness function | **Yes** |
| `ship-to-friday` | Gate 1 — POST to Friday dashboard | No — always human |
| `db-migration` | Gate 3 — SQL migration against remote DB | No — always human |
| `db-contract-phase` | Gate 4 — DROP/RENAME in contracting phase | No — always human |
| `external-api` | Gate 5 — third-party API mutation | No — always human |
| `deploy` | Gate 8 — triggering a deployment | No — always human |

Policies targeting non-eligible gates are **silently ignored** by the evaluator and logged as
`policy.skipped` in telemetry (reason: `gate-not-eligible`). The evaluator never auto-approves a
gate marked "No — always human."

---

## Condition

A condition is a YAML map of key-value checks. All checks must pass for the condition to be true
(implicit AND). Use a list under `any:` for OR logic.

### Scalar checks

```yaml
condition:
  diffLines:
    lessThan: 200           # diff size in lines
  testsPass: true           # all configured tests pass (boolean)
  codeReviewer.verdict:
    equals: "APPROVED"      # code-reviewer produced APPROVED verdict
  securityReviewer.criticals:
    equals: 0               # zero critical security findings
  fitnessFunction.allPass: true   # all fitness functions pass
  diffType:
    equals: "docs-only"     # diff contains only documentation files
  diffType:
    equals: "test-additions" # diff adds test files without touching source
  filePaths:
    allMatch: "docs/**"     # all changed files match a glob
  filePaths:
    noneMatch: "**/security/**"   # no changed files match a glob
  dryRunPass: true          # CI dry-run validation passed
```

### OR logic

```yaml
condition:
  any:
    - codeReviewer.verdict:
        equals: "APPROVED"
    - diffLines:
        lessThan: 50
```

### NOT / inversion

```yaml
condition:
  not:
    filePaths:
      anyMatch: "**/auth/**"
```

---

## Action

```yaml
action:
  type: auto-approve        # one of: auto-approve | auto-reject | require-human | escalate
  reason: "string"          # human-readable reason written to the telemetry event
  escalateTo: "string"      # optional; used with type: escalate — names who or what to notify
```

| Action | Meaning |
|---|---|
| `auto-approve` | Evaluator approves the gate; pipeline proceeds without human prompt |
| `auto-reject` | Evaluator rejects; pipeline halts and reports reason to operator |
| `require-human` | Evaluator explicitly requires human confirmation (overrides auto-approve policies) |
| `escalate` | Evaluator pauses pipeline and surfaces findings to `escalateTo` target |

**Conflict resolution**: if multiple policies match the same gate and their actions conflict
(`auto-approve` vs `require-human`), `require-human` always wins. The evaluator logs a
`policy.conflict` telemetry event listing the conflicting policy names.

---

## Full example

```yaml
name: auto-approve-doc-only-commits
version: "1.0"
description: "Auto-approve git commits when the diff is documentation-only and tests pass"
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
  reason: "Docs-only diff, tests green — no review required"
```

---

## Telemetry events emitted

Every evaluation (match or no-match) emits a `policy.evaluated` event:

```json
{
  "event": "policy.evaluated",
  "policyName": "auto-approve-doc-only-commits",
  "gate": "git-commit",
  "conditionMet": true,
  "action": "auto-approve",
  "reason": "Docs-only diff, tests green — no review required",
  "timestamp": "2026-08-01T12:00:00Z",
  "pipelineRunId": "<run-id>"
}
```

Non-matching policies emit `conditionMet: false` and `action: "no-op"`. Skipped non-eligible gate
policies emit event type `policy.skipped`.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 4
Policy Layer. Licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md).*
