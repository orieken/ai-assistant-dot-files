# Policy Evaluator

Part of the `shared/orchestration/` runtime. The evaluator is the bridge between the declarative
policy files in `.claude/policies/` and the `FeatureDeliveryWorkflow`'s stage boundaries.

---

## Role

At every checkpoint stage in `FeatureDeliveryWorkflow`, the runtime calls the evaluator before
proceeding past the gate. The evaluator:

1. Loads all `.policy.yaml` files from `.claude/policies/` (project-local) or the path configured
   in `.claude/delivery-policy.yaml` under `policyDir`.
2. Filters to policies whose `matcher.gate` matches the current gate ID.
3. Evaluates each matching policy's `condition` against the current pipeline context.
4. Determines the final action (see Conflict Resolution below).
5. Emits a `policy.evaluated` telemetry event for every policy — whether it matched or not.
6. Returns one of three signals to the workflow runtime:
   - `proceed` — the pipeline advances without prompting the human operator
   - `halt` — the pipeline stops and reports the reason; human must intervene
   - `require-human` — the pipeline pauses and asks the human operator for confirmation

---

## Inputs

```yaml
evaluatorInputs:
  gate: string                       # gate ID (e.g. "git-commit", "fitness-function-wiring")
  pipelineContext:                   # structured view of current pipeline state
    diffLines: integer               # total lines changed in current diff
    diffType: string                 # "docs-only" | "test-additions" | "source" | "mixed"
    filePaths: [string]              # list of changed file paths
    testsPass: boolean               # did all configured tests pass?
    fitnessFunction.allPass: boolean # did all fitness functions pass?
    codeReviewer.verdict: string     # latest code-reviewer verdict if present
    codeReviewer.behaviorChange: boolean
    codeReviewer.criticals: integer
    securityReviewer.criticals: integer
    dryRunPass: boolean              # for fitness-function-wiring gate
  policyDir: path                    # resolved policy directory (default: .claude/policies/)
  telemetryStream: path              # where to append events (default: .claude/telemetry/events.jsonl)
```

---

## Evaluation algorithm

```
1. If policiesEnabled == false in .claude/delivery-policy.yaml → return require-human (kill-switch)
2. Load policies from policyDir (ignore files that fail schema parse; log warning)
3. Filter: keep policies where matcher.gate == current gate AND enabled == true
4. Further filter: discard policies targeting non-eligible gates (always-human list):
   ship-to-friday, db-migration, db-contract-phase, external-api, deploy
5. For each candidate policy:
   a. Evaluate condition against pipelineContext
   b. Emit policy.evaluated event (conditionMet: true|false)
6. Collect actions from policies where conditionMet == true
7. Apply conflict resolution (see below)
8. If no policies matched → return require-human (default)
9. Return resolved signal
```

---

## Conflict resolution

When multiple policies match the same gate with conflicting actions:

| Conflicting actions | Resolution |
|---|---|
| `auto-approve` + `require-human` | `require-human` wins unconditionally |
| `auto-approve` + `auto-reject` | `auto-reject` wins; log `policy.conflict` event |
| `auto-approve` + `escalate` | `escalate` wins |
| Multiple `auto-approve` | `auto-approve` (unanimous) |
| Multiple `require-human` | `require-human` |

A `policy.conflict` telemetry event is emitted whenever two policies disagree on a gate. The
conflict event lists both policy names and their actions so operators can diagnose and deduplicate.

---

## Telemetry events emitted (non-negotiable)

Every call to the evaluator emits at least one event. There are no silent auto-approvals.

### `policy.evaluated`

```json
{
  "event": "policy.evaluated",
  "policyName": "<name from policy file>",
  "gate": "<gate-id>",
  "conditionMet": true | false,
  "action": "auto-approve | auto-reject | require-human | escalate | no-op",
  "reason": "<policy action.reason or 'no-match'>",
  "timestamp": "<ISO-8601>",
  "pipelineRunId": "<run-id>",
  "pipelineStage": "<stage-id>"
}
```

### `policy.conflict`

```json
{
  "event": "policy.conflict",
  "gate": "<gate-id>",
  "conflictingPolicies": ["<name-a>", "<name-b>"],
  "actions": ["auto-approve", "require-human"],
  "resolution": "require-human",
  "timestamp": "<ISO-8601>",
  "pipelineRunId": "<run-id>"
}
```

### `policy.skipped`

```json
{
  "event": "policy.skipped",
  "policyName": "<name>",
  "gate": "<gate-id>",
  "reason": "gate-not-eligible | disabled | parse-error",
  "timestamp": "<ISO-8601>",
  "pipelineRunId": "<run-id>"
}
```

---

## Integration points in FeatureDeliveryWorkflow

The `FeatureDeliveryWorkflow` (`shared/workflows/feature-delivery-workflow.md`) has explicit stage
boundary comments noting where Phase 4 policy hooks evaluate. The evaluator is called at:

- `code-review` stage boundary — `gate: git-commit`
- `fitness-function-wiring` gate (inline in devops stage) — `gate: fitness-function-wiring`
- `out-of-boundary-write` check (inline in development stage) — `gate: out-of-boundary-write`

The workflow does NOT call the evaluator for non-eligible gates (1, 3, 4, 5, 8). Those gate
decisions always return `require-human` unconditionally, bypassing the evaluator entirely.

---

## Backward-compatibility guarantee

If `.claude/policies/` does not exist or is empty, the evaluator returns `require-human` for every
gate — identical to v3.2 behavior. Teams that upgrade to v3.3 without creating any policy files
see **zero behavior change**.

---

## Dry-run mode

```
/orchestrate --workflow feature-delivery --spec <file> --dry-run-policies
```

In dry-run mode, the evaluator evaluates all policies against historical telemetry from
`.claude/telemetry/events.jsonl` instead of live pipeline context. The `policy.evaluated` events
are written to `.claude/telemetry/events-dryrun.jsonl` and never mutate the pipeline state. Use
this to validate new policies before enabling them.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) AOS Phase 4
Policy Layer. Licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
