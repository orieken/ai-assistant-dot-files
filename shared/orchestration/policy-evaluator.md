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
5. Was specified to emit a `policy.evaluated` event for every policy, matched or not. Nothing emits
   it — see "Telemetry events: specified, not emitted" below.
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
   b. (No event is emitted — see below.)
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

## Telemetry events: specified, not emitted

This section used to define three event payloads — `policy.evaluated`, `policy.conflict`,
`policy.skipped` — and to state that every call emits at least one, so there are no silent
auto-approvals.

**Nothing emits them.** The `.claude/telemetry/events.jsonl` file they targeted had no verified
writer and was retired in roadmap **L3.9**; the payloads used the key `event` where that schema
required `event_type`, which is the kind of mismatch a format nobody writes never has to resolve.
This evaluator is itself prose a model follows, not code that runs.

The payload definitions are removed rather than corrected. Fixing a key in a message nothing sends
makes the specification look more real than it is, which is how a documented guarantee comes to be
believed.

**What this means for the guarantee.** The rule stands: a policy-based gate is a human gate whose
prompt is replaced by a policy decision, and gates marked Always Human are never delegated. What
does not exist is any record of those decisions. Building an evaluator that runs as code and emits
a real audit trail is roadmap **L2.16**; until it lands, "no silent auto-approvals" is a property of
whoever follows these instructions, not one anything can demonstrate after the fact.

See `shared/telemetry/README.md` for the full list of specified-but-unemitted types and the roadmap
item behind each.

## Dry-run mode

```
/orchestrate --workflow feature-delivery --spec <file> --dry-run-policies
```

Dry-run mode was specified to evaluate policies against historical telemetry from
`.claude/telemetry/events.jsonl`. That file was retired in roadmap L3.9 and never had a writer, so
there is no history to replay and dry-run has nothing to read. It lands with **L2.16**, alongside
the evaluator that would produce the history in the first place.

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) AOS Phase 4
Policy Layer. Licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md).*
