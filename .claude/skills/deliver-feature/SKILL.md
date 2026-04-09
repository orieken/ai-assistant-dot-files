---
name: deliver-feature
description: The main pipeline orchestrator — kicks off the full agent sequence for a feature and persists all artifacts to docs/features/<feature-name>/.
triggers:
  keywords: ["deliver", "build", "pipeline", "feature"]
  intentPatterns: ["Deliver *", "Implement *", "Build *", "Start delivery on *", "/deliver-feature *"]
standalone: true
---

## When To Use
When the user asks to implement a feature, build a specific feature markdown file, or explicitly runs the `/deliver-feature` command. This delegates work to the full agent sequence.

Do NOT use when the user only wants a single agent's output (e.g., just an analysis or just a code review). Use the specific agent or skill instead.

## Context To Load First
1. The feature file (passed as argument or in `features/`)
2. `ARCHITECTURE_RULES.md`
3. `DOMAIN_DICTIONARY.md`
4. `CLAUDE.md`
5. `docs/features/README.md` (for artifact persistence conventions)

## Process

### Phase 0: Setup
1. **Read the feature file** — confirm it follows `features/TEMPLATE.md` structure. If not, stop and ask the user to run `/new-feature` first.
2. **Derive the feature name** — kebab-case from the feature file name (e.g., `features/user-auth.md` becomes `user-auth`).
3. **Create the feature workspace**: `.claude/feature-workspace/` — clean any prior artifacts.
4. **Create the feature archive directory**: `docs/features/<feature-name>/` — this is where all final artifacts are persisted.

### Phase 1: Discovery and Design
5. **Invoke analyst** -> produces `analysis.md`. **PAUSE**: show summary to user. Wait for confirmation before continuing.
6. **Invoke architect** (if analysis.md has Architectural Flags != "None") -> produces `architecture-notes.md`. **PAUSE** if RFC was written — human must acknowledge before developer starts.
7. **Invoke performance-engineer** (if analysis.md has Performance SLAs or Non-Functional Requirements with latency/throughput targets) -> produces `performance-report.md`.
8. **Invoke data-engineer** (if analysis.md has Data Model Changes != "None") -> produces `data-engineering-report.md`.

### Phase 2: Implementation and Review
9. **Invoke developer** -> produces `implementation-notes.md`.
10. **Invoke code-reviewer** -> produces `code-review-report.md`. If verdict is CHANGES REQUESTED: send back to developer. Repeat until APPROVED.
11. **Invoke accessibility-engineer** (if the feature involves UI components, templates, or user-facing HTML) -> produces `a11y-report.md`.
12. **Invoke security-reviewer** (if security surface exists — auth, user input, API endpoints, tokens, trust boundaries) -> produces `security-report.md`. If Critical findings exist: block pipeline, alert user.

### Phase 3: Verification and Shipping
13. **Invoke qa-engineer** -> produces `qa-report.md`. Tests must be green.
14. **Invoke sre-engineer** -> produces `observability-report.md`.
15. **Invoke tech-writer** -> produces `docs-report.md`.
16. **Invoke devops-engineer** -> produces `devops-report.md`.

### Phase 4: Persistence and Delivery
17. **Write delivery summary** -> produces `delivery-summary.md` in `.claude/feature-workspace/`.
18. **Persist all artifacts** — copy every produced artifact from `.claude/feature-workspace/` to `docs/features/<feature-name>/`.
19. **Create feature archive index** — write `docs/features/<feature-name>/README.md` listing all artifacts with descriptions and links.
20. **Update feature index** — add the new feature entry to `docs/features/README.md`.
21. **PAUSE**: show the user the full `docs/features/<feature-name>/` listing. Confirm the documentation is complete.
22. **Ship to Friday** — ask: "Ship to Friday?" On confirmation ("ship" or "yes"): POST Cucumber JSON to Friday.

## Human Checkpoints
- After analyst (step 5): confirm scope before any code is written
- After architect RFC (step 6): confirm architectural direction before developer starts
- After code-review CHANGES REQUESTED loop (step 10): confirm all findings resolved
- After security Critical finding (step 12): explicit "fix confirmed" before QA starts
- After artifact persistence (step 21): confirm documentation is complete
- Before shipping to Friday (step 22): explicit "ship" confirmation

## Output Format

### Working Artifacts (temporary)
All agents write to `.claude/feature-workspace/` during execution.

### Persisted Artifacts (permanent)
After pipeline completion, all artifacts are copied to `docs/features/<feature-name>/`:

```
docs/features/<feature-name>/
  README.md                  <- index of all artifacts with links
  analysis.md                <- analyst output
  architecture-notes.md      <- architect output (if invoked)
  performance-report.md      <- performance-engineer output (if invoked)
  data-engineering-report.md <- data-engineer output (if invoked)
  implementation-notes.md    <- developer output
  code-review-report.md      <- code-reviewer output
  a11y-report.md             <- accessibility-engineer output (if invoked)
  security-report.md         <- security-reviewer output (if invoked)
  qa-report.md               <- qa-engineer output
  observability-report.md    <- sre-engineer output
  docs-report.md             <- tech-writer output
  devops-report.md           <- devops-engineer output
  delivery-summary.md        <- final synthesis
```

### Delivery Summary Format

```markdown
# Delivery Summary: [Feature Name]

## Pipeline Run
| Agent | Status | Key Output |
|---|---|---|
| analyst | PASS | [N acceptance criteria, N architectural flags] |
| architect | PASS / SKIPPED | [N structural decisions, RFC: yes/no] |
| performance-engineer | PASS / SKIPPED | [N SLAs verified, N recommendations] |
| data-engineer | PASS / SKIPPED | [N migrations, expand/contract phase] |
| developer | PASS | [N files created, N modified, N refactoring ops] |
| code-reviewer | PASS | [Design score: C/Co/Cu/Cr — APPROVED] |
| accessibility-engineer | PASS / SKIPPED | [N violations found, N fixed] |
| security-reviewer | PASS / SKIPPED | [N findings, N critical fixed] |
| qa-engineer | PASS | [N tests, N passed, SLAs verified: yes/no] |
| sre-engineer | PASS | [N spans added, N alerts configured] |
| tech-writer | PASS | [N docs updated] |
| devops-engineer | PASS | [N CI changes, N env vars] |

## Artifacts Persisted
Location: docs/features/<feature-name>/
Files: [count] artifacts written

## Friday
Status: Shipped | Pending | Skipped

## Artifacts
- docs/features/<feature-name>/analysis.md
- docs/features/<feature-name>/architecture-notes.md
- [all produced artifacts listed with full paths]
```

### Feature Archive Index Format

```markdown
# Feature: [Feature Name]

Delivered: [YYYY-MM-DD]
Status: Complete | Complete with notes | Blocked

## Artifacts

| Document | Agent | Description |
|---|---|---|
| [analysis.md](./analysis.md) | analyst | Technical analysis and task breakdown |
| [architecture-notes.md](./architecture-notes.md) | architect | Structural decisions and fitness functions |
| ... | ... | ... |
| [delivery-summary.md](./delivery-summary.md) | orchestrator | Final pipeline synthesis |

## Summary
[2-3 sentence plain English summary from delivery-summary.md]
```

## Guardrails
- Never skip the analyst — it is always first
- Never let the developer start without analysis.md
- Never send CHANGES REQUESTED code to the security reviewer or QA
- Never ship to Friday without explicit "ship" or "yes" from the user
- Never persist artifacts to docs/features/ until the delivery summary is written
- Pipeline can be resumed from any checkpoint — check which workspace artifacts already exist
- The feature archive in docs/features/<feature-name>/ is append-only — never delete prior delivery artifacts

## Standalone Mode
All agents run locally. Friday POST is the only external call — non-blocking if Friday is not running. Artifact persistence to docs/features/ works entirely offline.
