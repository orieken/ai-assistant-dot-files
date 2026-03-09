---
name: chaos-experiment
description: Designs and runs Game Days fault injection tests for resilience verification.
triggers:
  keywords: ["chaos", "game day", "fault injection", "resilience verify"]
  intentPatterns: ["Design a chaos experiment for *", "Game day for *", "Test the circuit breaker on *"]
standalone: true
---

## When To Use
Triggered specifically by the `chaos-engineer` or a developer wanting to test resilience patterns (Timeouts, Circuit Breakers, Bulkheads).

## Context To Load First
The target service architecture or the code implementing the resilience pattern.

## Process
1. Define the blast radius securely.
2. Select the fault injection method (e.g., `toxiproxy`, network dropping via Chaos Mesh, or application-level mock faults).
3. Provide the exact bash commands or test harness code to inject the fault.
4. Tell the user how to observe the system's reaction via their APM/logs.
5. Provide the recovery/rollback commands to revert the fault.

## Output Format
`.claude/feature-workspace/chaos-experiment-design.md`

```markdown
# Chaos Experiment Design: [Component]

## Hypothesis
Under [Condition], [Component] will respond by [Graceful Degradation].

## Injection Setup
`bash script or code mock to disrupt the network/service`

## Observation Checklist
- [ ] Metric A spikes to X
- [ ] Log pattern "Circuit Open" appears
- [ ] Error rate holds steady at Y% without cascading

## Rollback Command
`bash script to restore normal operations`
```

## Guardrails
- **Always** provide the rollback command FIRST in case of an emergency.
- Chaos experiments should target observability as much as functionality (can we actually *see* the failure occurring?).
