---
name: performance-profile
description: Diagnose why something is slow, not just verify that it is.
triggers:
  keywords: ["slow", "profile", "SLA", "latency", "degrading"]
  intentPatterns: ["Why is * slow?", "Profile *", "The SLA is failing for *", "Performance is degrading on *"]
standalone: true
---

## When To Use
Triggered to diagnose performance issues or latency spikes. Do NOT use for premature optimization before a test fails its SLA.

## Context To Load First
Relevant failing SLAs, logs, endpoints, or target code.

## Process
1. Establish the baseline: "What is the current measured time and what is the SLA?"
2. Identify the layer: browser, server, database, or network? Ask the user: "Where does the time appear to be spent?" If unknown, start at browser.
3. Generate the profiling command for the identified layer.
4. Interpret the output: what is the hottest path?
5. Apply the diagnostic checklist for that layer.
6. Produce Performance Profile Report with ranked hypotheses and the one experiment to run first.

## Output Format
`.claude/feature-workspace/performance-profile-[description].md`

```markdown
# Performance Profile: [Feature/Endpoint/Page]

## Baseline
- Current measured time: [Xms]
- SLA target: [Yms]
- Gap: [Xms over target]

## Layer Identified
Browser | Server | Database | Network

## Profiling Commands Run
```bash
[exact commands used]
```

## Hottest Path
[What the profiler identified as the primary bottleneck]

## Diagnostic Checklist Results
- [ ] [Check]: [result]

## Ranked Hypotheses
1. [Most likely cause]: [evidence] — [experiment to confirm]
2. [Second hypothesis]: [evidence] — [experiment to confirm]

## Recommended First Experiment
[One specific change to make, with expected impact and how to measure it]

## Expected Outcome
If hypothesis 1 is correct, making [change] should reduce time from [X]ms to [Y]ms.
Verify with: [exact command or test to rerun]
```

## Guardrails
- Never recommend an optimization before running a profiler.
- Only recommend one experiment at a time — concurrent changes make causation impossible to establish.
- Never recommend premature abstraction before ruling out algorithmic fixes.

## Standalone Mode
Provides profiling command generation and diagnostic checklists even when profiling tools aren't installed. Always produces a structured hypothesis list.
