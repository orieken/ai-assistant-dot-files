---
name: deliver-feature
description: The main pipeline orchestrator — kicks off the full agent sequence for a feature.
triggers:
  keywords: ["deliver", "build", "pipeline", "feature"]
  intentPatterns: ["Deliver *", "Implement *", "Build *", "Start delivery on *", "/deliver-feature *"]
standalone: true
---

## When To Use
When the user asks to implement a feature, build a specific feature markdown file, or explicitly runs the `/deliver-feature` command. This delegates work to the full agent sequence.

## Context To Load First
1. The feature file (passed as argument or in `features/`)
2. `ARCHITECTURE_RULES.md`
3. `DOMAIN_DICTIONARY.md`
4. `CLAUDE.md`

## Process
1. **Read the feature file** — confirm it follows `features/TEMPLATE.md` structure. If not, stop and ask the user to run `/new-feature` first.
2. **Create the feature workspace**: `.claude/feature-workspace/` — clean any prior artifacts.
3. **Invoke analyst** → produces `analysis.md`. Pause: show summary to user. Continue.
4. **Invoke architect** (if analysis.md has Architectural Flags ≠ "None") → produces `architecture-notes.md`. Pause if RFC was written — human must acknowledge.
5. **Invoke developer** → produces `implementation-notes.md`.
6. **Invoke code-reviewer** → produces `code-review-report.md`. If verdict is CHANGES REQUESTED: send back to developer. Repeat until APPROVED.
7. **Invoke security-reviewer** (if security surface exists) → produces `security-report.md`. If Critical findings exist: block pipeline, alert user.
8. **Invoke qa-engineer** → produces `qa-report.md`. Tests must be green.
9. **Invoke tech-writer** → produces `docs-report.md`.
10. **Invoke devops-engineer** → produces `devops-report.md`.
11. **Pipeline complete**: show summary of all workspace artifacts. Ask: "Ship to Friday?" On confirmation ("ship" or "yes"): POST Cucumber JSON to Friday.

## Human Checkpoints
- After analyst: confirm scope before any code is written
- After architect RFC: confirm architectural direction before developer starts
- After code-review CHANGES REQUESTED loop: confirm all findings resolved
- After security Critical finding: explicit "fix confirmed" before QA starts
- Before shipping to Friday: explicit "ship" confirmation

## Output Format
`.claude/feature-workspace/delivery-summary.md`

```markdown
# Delivery Summary: [Feature Name]

## Pipeline Run
| Agent | Status | Key Output |
|---|---|---|
| analyst | ✅ | [N acceptance criteria, [N] architectural flags] |
| architect | ✅/⏭️ skipped | [N structural decisions, RFC: yes/no] |
| developer | ✅ | [N files created, N modified, N refactoring ops] |
| code-reviewer | ✅ | [Design score: C/Co/Cu/Cr — APPROVED] |
| security-reviewer | ✅/⏭️ skipped | [N findings, N critical fixed] |
| qa-engineer | ✅ | [N tests, N passed, SLAs verified: yes/no] |
| tech-writer | ✅ | [N docs updated] |
| devops-engineer | ✅ | [N CI changes, N env vars] |

## Friday
Status: Shipped | Pending | Skipped

## Artifacts
- `.claude/feature-workspace/analysis.md`
- `.claude/feature-workspace/architecture-notes.md`
- [all produced artifacts listed]
```

## Guardrails
- Never skip the analyst — it is always first
- Never let the developer start without analysis.md
- Never send CHANGES REQUESTED code to the security reviewer or QA
- Never ship to Friday without explicit "ship" or "yes" from the user
- Pipeline can be resumed from any checkpoint — check which workspace artifacts already exist

## Standalone Mode
All agents run locally. Friday POST is the only external call — non-blocking if Friday isn't running.
