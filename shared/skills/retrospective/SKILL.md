---
name: retrospective
description: Runs a structured retrospective on a completed feature delivery by analyzing all pipeline artifacts in docs/features/<feature-name>/.
triggers:
  keywords: ["retro", "retrospective", "review delivery", "postmortem"]
  intentPatterns: ["Run a retrospective on *", "How did * go?", "Review the delivery of *", "/retrospective *"]
standalone: true
---

## When To Use
When a feature delivery is complete and the user wants to analyze what went well, what went poorly, and what to improve. Reads all persisted artifacts from `docs/features/<feature-name>/`.

Do NOT use for active incident response — use `/on-call` instead.
Do NOT use for root cause analysis of a specific bug — use `/five-whys` instead.

## Context To Load First
1. `docs/features/<feature-name>/` — all pipeline artifacts
2. `docs/features/<feature-name>/delivery-summary.md` — pipeline run status
3. `CLAUDE.md` — project constraints for benchmarking

## Process

1. **Identify the feature** — accept a feature name (kebab-case) or path to the feature docs directory. Verify `docs/features/<feature-name>/` exists and contains a `delivery-summary.md`.

2. **Read all pipeline artifacts** in the feature directory:
   - `analysis.md` — scope and acceptance criteria
   - `architecture-notes.md` — structural decisions (if present)
   - `implementation-notes.md` — what was built
   - `code-review-report.md` — review findings and cycles
   - `security-report.md` — security findings (if present)
   - `a11y-report.md` — accessibility findings (if present)
   - `qa-report.md` — test results
   - `observability-report.md` — OTel coverage (if present)
   - `docs-report.md` — documentation updates
   - `devops-report.md` — deployment changes
   - `delivery-summary.md` — final synthesis

3. **Analyze delivery metrics**:
   - **Scope accuracy**: Did the acceptance criteria in analysis.md match what was actually delivered?
   - **Review cycles**: How many code-reviewer CHANGES REQUESTED loops occurred?
   - **Security findings**: Were there Critical findings? How many review cycles to resolve?
   - **Complexity scores**: Did any implemented code exceed complexity thresholds?
   - **Test coverage**: Did QA achieve the 85% threshold?
   - **Architectural decisions**: Were fitness functions defined and wired?
   - **Documentation completeness**: Were all relevant docs updated?

4. **Categorize findings** into:
   - **What went well** — things the pipeline caught and handled correctly
   - **What went poorly** — things that required rework, were missed, or caused delays
   - **What to improve** — actionable process changes for next time

5. **Produce the retrospective** at `docs/features/<feature-name>/retrospective.md`.

## Output Format

Write `docs/features/<feature-name>/retrospective.md`:

```markdown
# Retrospective: [Feature Name]

Date: [YYYY-MM-DD]
Delivery status: Complete | Complete with notes | Blocked

## Delivery Metrics

| Metric | Value | Threshold | Status |
|---|---|---|---|
| Review cycles | [N] | <= 2 | PASS / WARN / FAIL |
| Security findings (Critical) | [N] | 0 | PASS / FAIL |
| Security findings (total) | [N] | — | INFO |
| Cyclomatic complexity (max) | [N] | < 7 | PASS / FAIL |
| Test coverage | [N%] | >= 85% | PASS / FAIL |
| Fitness functions added | [N] | >= 1 per arch decision | PASS / WARN |
| Acceptance criteria met | [N/M] | M/M | PASS / FAIL |
| Accessibility violations | [N] | 0 | PASS / FAIL |

## What Went Well
- [Specific positive observation with evidence from the artifacts]
- [Another positive observation]

## What Went Poorly
- [Specific problem]: [evidence from artifacts] -> [impact]
- [Another problem]: [evidence] -> [impact]

## What To Improve
- [Actionable recommendation]: [which pipeline step it affects] -> [expected benefit]
- [Another recommendation]: [step] -> [benefit]

## Process Recommendations
- [ ] [Specific change to the pipeline or agent configuration]
- [ ] [Specific change to the spec-writing process]
- [ ] [Specific change to the review criteria]

## Patterns Identified
[Recurring themes across this and prior deliveries, if other retrospectives exist in docs/features/]
```

## Guardrails
- Never fabricate metrics — only report what is evidenced in the artifacts.
- If an artifact is missing (e.g., no security-report.md), note it as "not assessed" rather than assuming it passed.
- The retrospective is for the process, not for individuals — focus on pipeline steps and handoffs.
- If prior retrospectives exist in other feature directories, reference recurring patterns.
- Never modify any existing pipeline artifacts — the retrospective is an additive document only.

## Standalone Mode
Works entirely offline. Reads only local markdown files. No external services required.
