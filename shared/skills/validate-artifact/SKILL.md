---
name: validate-artifact
description: Validates a pipeline artifact against its inter-agent contract in shared/contracts/ — checks required sections are present (and a few contract-specific content rules) before the artifact is allowed to move to the next agent. Invoked automatically between every contract-bound handoff in deliver-feature.
triggers:
  keywords: ["validate-artifact", "validate contract", "check contract", "contract check"]
  intentPatterns: ["/validate-artifact *", "Validate this artifact against its contract", "Check *.md against its contract"]
standalone: true
---

## When To Use
Invoked automatically by `deliver-feature` immediately after context-engineer, analyst, architect,
developer, code-reviewer, security-reviewer, qa-engineer, or sre-engineer produce their artifact — before
the pipeline proceeds to the next step. Can also be run standalone: "validate analysis.md against its
contract".

Do NOT use for agents without a contract yet: performance-engineer, data-engineer, accessibility-engineer,
tech-writer, devops-engineer. Their outputs are consumed directly without a structural gate (tracked as an
open item in `docs/features/context-engineering-framework/TODO.md`, Epic 5).

## Context To Load First
1. The contract file for the artifact being validated (see mapping below)
2. The artifact file itself

## Contract Mapping
| Agent | Artifact | Contract |
|---|---|---|
| context-engineer | `.claude/feature-workspace/context-manifest.md` | `shared/contracts/context-manifest-contract.md` |
| analyst | `.claude/feature-workspace/analysis.md` | `shared/contracts/analysis-contract.md` |
| architect | `.claude/feature-workspace/architecture-notes.md` | `shared/contracts/architecture-contract.md` |
| developer | `.claude/feature-workspace/implementation-notes.md` | `shared/contracts/implementation-contract.md` |
| code-reviewer | `.claude/feature-workspace/code-review-report.md` | `shared/contracts/review-contract.md` |
| security-reviewer | `.claude/feature-workspace/security-report.md` | `shared/contracts/security-contract.md` |
| qa-engineer | `.claude/feature-workspace/qa-report.md` | `shared/contracts/qa-contract.md` |
| sre-engineer | `.claude/feature-workspace/observability-report.md` | `shared/contracts/observability-contract.md` |

## Process
1. Read the contract's "Required Sections" list.
2. Grep the artifact for each required heading — exact string and heading level (`##` vs `###`) must match.
3. Apply the contract's "Validation Rule" content checks (e.g., "Overall Status must contain APPROVED or CHANGES REQUESTED", "Test Results must show Failed: 0").
4. Record every miss — missing heading or failed content rule — as a violation.
5. Report `PASS` (zero violations) or `FAIL` (one or more violations).

## Output Format
```markdown
# Artifact Validation: [artifact filename]
Contract: [contract filename]
Status: PASS | FAIL

## Missing Sections (if any)
- [heading]

## Content Rule Violations (if any)
- [rule] — [what was found instead]

## Present Sections
- [heading] ✓
```

## Guardrails
- **Structural only**: this checks presence of headings and the small set of content rules each contract
  spells out explicitly. It never judges whether the content is *correct* — that stays with the human
  PAUSE checkpoints and the downstream agents (code-reviewer, security-reviewer, qa-engineer).
- **Exact match**: a heading at the wrong level (e.g., `### Summary` when the contract requires `## Summary`)
  is a FAIL, not a warning — a level mismatch breaks downstream `grep`-based parsing just as much as a
  missing heading would.
- **Never silently pass**: every miss must be listed explicitly in the output, even if the orchestrator
  chooses to proceed anyway (it shouldn't, except by explicit human override).
- **On FAIL**: the orchestrator sends the artifact back to the producing agent with the specific violations
  listed, and re-validates after the fix. This is a structural loop, distinct from code-reviewer's
  qualitative CHANGES REQUESTED loop — both can fire independently.

## Standalone Mode
No external tools required — this is a Read + Grep pass over two local markdown files. Works identically
inside or outside the `deliver-feature` pipeline.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
