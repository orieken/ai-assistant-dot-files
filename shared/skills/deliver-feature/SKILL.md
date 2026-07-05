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
3. **Check for an existing `.claude/feature-workspace/pipeline-state.json`.** If one exists for this feature: stop and invoke `resume-pipeline` instead of continuing here — do not blindly clean the workspace out from under an in-progress or crashed run. If the user explicitly asked to start over ("start fresh", "restart delivery"), archive the old state file to `.claude/feature-workspace/.history/pipeline-state.json.<timestamp>` and proceed. If no state file exists, create the feature workspace: `.claude/feature-workspace/` — clean any prior artifacts.
4. **Create the feature archive directory**: `docs/features/<feature-name>/` — this is where all final artifacts are persisted.
5. **Initialize `.claude/feature-workspace/pipeline-state.json` and `.claude/feature-workspace/pipeline-trace.json`** — see "Checkpointing & Pipeline State" and "Pipeline Tracing" below.

### Phase 1: Discovery and Design
6. **Invoke context-engineer** -> produces `context-manifest.md` in `.claude/feature-workspace/`. This scopes the bounded context, pins the specific files analyst/developer must read, lists relevant KIs/ADRs, and estimates the token budget. If it flags a budget WARNING, tell the user which files it recommends cutting before continuing. **Checkpoint**: record in `pipeline-state.json`.
7. **Invoke validate-artifact** against `shared/contracts/context-manifest-contract.md`. If FAIL: send back to context-engineer with the specific violations listed; re-validate. Repeat until PASS. **Checkpoint** on PASS.
8. **Invoke analyst** -> reads `context-manifest.md` first, then produces `analysis.md`.
9. **Invoke validate-artifact** against `shared/contracts/analysis-contract.md`. If FAIL: send back to analyst with the specific violations listed; re-validate. Repeat until PASS. **Checkpoint** on PASS.
10. **PAUSE**: show summary to user. Wait for confirmation before continuing.
11. **Invoke architect** (if analysis.md has Architectural Flags != "None") -> produces `architecture-notes.md`.
12. **Invoke validate-artifact** against `shared/contracts/architecture-contract.md` (only if architect was invoked). If FAIL: send back to architect; re-validate. **PAUSE** if an RFC was written — human must acknowledge before developer starts. **Checkpoint** on PASS or SKIP.
13. **Invoke performance-engineer** (if analysis.md has Performance SLAs or Non-Functional Requirements with latency/throughput targets) -> produces `performance-report.md`. **Checkpoint**.
14. **Invoke validate-artifact** against `shared/contracts/performance-contract.md` (only if performance-engineer was invoked). If FAIL: send back to performance-engineer; re-validate. **Checkpoint** on PASS or SKIP.
15. **Invoke data-engineer** (if analysis.md has Data Model Changes != "None") -> produces `data-engineering-notes.md`. **Checkpoint**.
16. **Invoke validate-artifact** against `shared/contracts/data-engineering-contract.md` (only if data-engineer was invoked). If FAIL: send back to data-engineer; re-validate. **Checkpoint** on PASS or SKIP.

### Phase 2: Implementation and Review
17. **Invoke developer** -> reads `context-manifest.md` first, then produces `implementation-notes.md`.
18. **Invoke validate-artifact** against `shared/contracts/implementation-contract.md`. If FAIL: send back to developer; re-validate. **Checkpoint** on PASS.
19. **Invoke code-reviewer** -> produces `code-review-report.md`.
20. **Invoke validate-artifact** against `shared/contracts/review-contract.md`. If FAIL (structural): send back to code-reviewer; re-validate. If verdict is CHANGES REQUESTED (qualitative, independent of the structural check): before sending back to developer, back up the current `implementation-notes.md` and `code-review-report.md` to `.claude/feature-workspace/.history/` (see Rollback), then repeat from step 17 until APPROVED and structurally valid. **Checkpoint** on final PASS+APPROVED, including the retry count.
21. **Invoke accessibility-engineer** (if the feature involves UI components, templates, or user-facing HTML) -> produces `accessibility-report.md`. **Checkpoint**.
22. **Invoke validate-artifact** against `shared/contracts/accessibility-contract.md` (only if accessibility-engineer was invoked). If FAIL: send back to accessibility-engineer; re-validate. **Checkpoint** on PASS or SKIP.
23. **Invoke security-reviewer** (if security surface exists — auth, user input, API endpoints, tokens, trust boundaries) -> produces `security-report.md`.
24. **Invoke validate-artifact** against `shared/contracts/security-contract.md` (only if security-reviewer was invoked). If FAIL: send back to security-reviewer; re-validate. If Critical findings exist: block pipeline, alert user. **Checkpoint** on PASS or SKIP.

### Phase 3: Verification and Shipping
25. **Invoke qa-engineer** -> produces `qa-report.md`. Tests must be green.
26. **Invoke validate-artifact** against `shared/contracts/qa-contract.md`. If FAIL: send back to qa-engineer; re-validate. **Checkpoint** on PASS.
27. **Invoke sre-engineer** -> produces `observability-report.md`.
28. **Invoke validate-artifact** against `shared/contracts/observability-contract.md`. If FAIL: send back to sre-engineer; re-validate. **Checkpoint** on PASS.
29. **Invoke tech-writer** -> produces `docs-report.md`. **Checkpoint**.
30. **Invoke validate-artifact** against `shared/contracts/docs-contract.md`. If FAIL: send back to tech-writer; re-validate. **Checkpoint** on PASS.
31. **Invoke devops-engineer** -> produces `devops-report.md`. **Checkpoint**.
32. **Invoke validate-artifact** against `shared/contracts/devops-contract.md`. If FAIL: send back to devops-engineer; re-validate. **Checkpoint** on PASS.

### Phase 4: Persistence and Delivery
33. **Write delivery summary** -> produces `delivery-summary.md` in `.claude/feature-workspace/`.
34. **Persist all artifacts** — copy every produced artifact from `.claude/feature-workspace/` to `docs/features/<feature-name>/`.
35. **Create feature archive index** — write `docs/features/<feature-name>/README.md` listing all artifacts with descriptions and links.
36. **Update feature index** — add the new feature entry to `docs/features/README.md`.
37. **Count total deliveries** — count `docs/features/*/delivery-summary.md` (including the one just written). If the count is evenly divisible by 5, auto-invoke `/retrospective` for the feature just delivered, producing `docs/features/<feature-name>/retrospective.md` before the next step. This is a single-delivery retrospective, distinct from `pipeline-retrospective` (which analyzes trends across many `pipeline-trace.json` files and is invoked manually, not on this cadence).
38. **PAUSE**: show the user the full `docs/features/<feature-name>/` listing. Confirm the documentation is complete.
39. **Ship to Friday** — ask: "Ship to Friday?" On confirmation ("ship" or "yes"): POST Cucumber JSON to Friday. Set `pipeline-state.json` phase to `complete`.

## Human Checkpoints
- After context-engineer passes contract validation (step 7): if token budget is WARNING, confirm pruning before analyst starts
- After analyst passes contract validation (step 10): confirm scope before any code is written
- After architect RFC (step 12): confirm architectural direction before developer starts
- After code-review CHANGES REQUESTED loop (step 20): confirm all findings resolved
- After security Critical finding (step 24): explicit "fix confirmed" before QA starts
- After artifact persistence (step 38): confirm documentation is complete
- Before shipping to Friday (step 39): explicit "ship" confirmation

## Checkpointing & Pipeline State

After every step marked **Checkpoint** above, write/update both `.claude/feature-workspace/pipeline-state.json`
(resumability — see below) and `.claude/feature-workspace/pipeline-trace.json` (timing/performance history —
see "Pipeline Tracing" below). They're updated together but serve different consumers: `pipeline-state.json`
is read by `resume-pipeline` to continue an interrupted run; `pipeline-trace.json` is read by
`pipeline-retrospective` and `agent-scorecard` to analyze trends across many runs.

```json
{
  "featureName": "user-auth",
  "featureFile": "features/user-auth.md",
  "startedAt": "2026-07-02T10:00:00Z",
  "updatedAt": "2026-07-02T10:45:00Z",
  "currentPhase": 2,
  "lastCompletedStep": 15,
  "completedAgents": [
    {
      "agent": "context-engineer",
      "step": 6,
      "artifact": "context-manifest.md",
      "checksum": "sha256:<hash>",
      "status": "PASS",
      "completedAt": "2026-07-02T10:05:00Z"
    },
    {
      "agent": "analyst",
      "step": 8,
      "artifact": "analysis.md",
      "checksum": "sha256:<hash>",
      "contractStatus": "PASS",
      "contractRetries": 0,
      "completedAt": "2026-07-02T10:15:00Z"
    }
  ]
}
```

- Compute the checksum as `sha256` of the artifact's current file content.
- Before overwriting an artifact that already exists in `.claude/feature-workspace/` (a re-run of the same agent, e.g. after a validate-artifact FAIL or a CHANGES REQUESTED loop), copy the existing version to `.claude/feature-workspace/.history/<artifact-name>.<unix-timestamp>.md` first, so it can be restored by a rollback.
- Skipped agents (conditional agents whose trigger condition was false) get a `completedAgents` entry with `"status": "SKIPPED"` and no artifact/checksum, so a later resume doesn't try to re-evaluate the skip condition against a possibly-changed `analysis.md`.
- If an artifact on disk doesn't match the checksum recorded for its step (someone hand-edited a workspace file outside the pipeline), treat that step as **not** completed — re-run the agent rather than trusting stale state.

## Rollback

If an agent's artifact turns out to be wrong (not just a validate-artifact FAIL, which self-heals via the retry loop, but a case where a human or a later agent determines an *earlier* artifact was flawed):

1. Identify the artifact to roll back to its previous version, and find its latest entry in `.claude/feature-workspace/.history/`.
2. Restore that history file over the current artifact.
3. In `pipeline-state.json`, remove (or mark `"stale": true` on) every `completedAgents` entry for that agent and every agent after it in the pipeline — they consumed content that no longer exists.
4. Re-run the pipeline starting at the rolled-back agent's step.
5. This is a structural/human-triggered rollback, distinct from the automatic validate-artifact and CHANGES REQUESTED retry loops, which don't require rollback because they re-run in place before anything downstream has consumed the bad artifact.

For resuming an interrupted run or replaying from a specific phase, use the `resume-pipeline` skill rather than restarting `deliver-feature` from Phase 0 — it reads `pipeline-state.json` and continues from `lastCompletedStep + 1`, or from the start of an explicitly requested phase ("resume delivery on user-auth from phase 2" / `--from-phase 2`).

## Pipeline Tracing

Alongside `pipeline-state.json`, maintain `.claude/feature-workspace/pipeline-trace.json` — a timing and
iteration-count record consumed by `pipeline-retrospective` and `agent-scorecard` (see
`shared/skills/pipeline-trace/SKILL.md` for the full schema and query usage). Minimal shape:

```json
{
  "featureName": "user-auth",
  "startedAt": "2026-07-02T10:00:00Z",
  "completedAt": "2026-07-02T11:30:00Z",
  "totalDurationSeconds": 5400,
  "agents": [
    {
      "agent": "analyst",
      "agentVersion": "1.0.0",
      "step": 8,
      "startedAt": "2026-07-02T10:05:00Z",
      "completedAt": "2026-07-02T10:15:00Z",
      "durationSeconds": 600,
      "status": "PASS",
      "iterations": 1,
      "contractRetries": 0,
      "budgetUtilization": 0.42
    },
    {
      "agent": "code-reviewer",
      "agentVersion": "1.0.0",
      "step": 19,
      "startedAt": "2026-07-02T10:40:00Z",
      "completedAt": "2026-07-02T11:10:00Z",
      "durationSeconds": 1800,
      "status": "APPROVED",
      "iterations": 3,
      "changesRequestedCount": 2,
      "budgetUtilization": null
    }
  ]
}
```

- `budgetUtilization` is copied from `context-manifest.md`'s Token Budget section for that agent's tier
  (fraction of the tier ceiling, not of the full window — see `shared/skills/pipeline-trace/SKILL.md` for
  the exact definition). Use `null` for agents that don't consume the manifest's Pinpoint Files directly
  rather than fabricating a number.
- `agentVersion` is that agent's `version:` frontmatter field (see `shared/agents/CHANGELOG.md`) at the time
  it ran — this is what lets `agent-scorecard` and `pipeline-retrospective` correlate a duration or
  iteration-count trend with a specific prompt edit, instead of just observing "it got slower" with no way
  to tie that to a cause.
- Record real wall-clock `startedAt`/`completedAt` for each agent invocation — never estimate or fabricate.
- If an agent re-runs (a validate-artifact retry or a CHANGES REQUESTED loop), don't create a second entry
  for it — update the same entry: add the additional elapsed time to `durationSeconds`, increment
  `iterations`, and (for code-reviewer specifically) increment `changesRequestedCount`.
- This file is persisted to `docs/features/<feature-name>/pipeline-trace.json` in Phase 4 alongside every
  other artifact, so `pipeline-retrospective` and `agent-scorecard` can read trace history across many
  past deliveries, not just the current run.

## Context Decay

Phases are numbered 0-4 (Setup, Discovery and Design, Implementation and Review, Verification and
Shipping, Persistence and Delivery). By the time an agent 2+ phases removed from an artifact's origin phase
needs it, read a `summarize-artifact` summary instead of the full file — the broad strokes still matter, the
exact wording usually doesn't, and the full file is still on disk for anyone who needs to dig in.

Concretely: `analysis.md` originates in Phase 1. `qa-engineer`, `sre-engineer`, `tech-writer`, and
`devops-engineer` run in Phase 3-4 — 2+ phases later — so they read a summary of `analysis.md`, not the full
file (both `qa-engineer.md` and `tech-writer.md` have been updated to do this explicitly; `sre-engineer.md`
and `devops-engineer.md` don't currently read `analysis.md` at all, so there's nothing to change there).
`implementation-notes.md` and `code-review-report.md` (Phase 2) are only 1 phase old from Phase 3 — read
those in full, not summarized.

This never applies to the artifact an agent is *immediately* reviewing (e.g. `code-reviewer` reading
`implementation-notes.md`'s Self-Review Checklist needs the literal checked items) — only to older artifacts
whose gist, not exact wording, is what still matters.

## Output Format

### Working Artifacts (temporary)
All agents write to `.claude/feature-workspace/` during execution.

### Persisted Artifacts (permanent)
After pipeline completion, all artifacts are copied to `docs/features/<feature-name>/`:

```
docs/features/<feature-name>/
  README.md                  <- index of all artifacts with links
  context-manifest.md        <- context-engineer output (scope, pinned files, KIs/ADRs, token budget)
  pipeline-trace.json        <- per-agent timing, status, and iteration counts (see Pipeline Tracing)
  analysis.md                <- analyst output
  architecture-notes.md      <- architect output (if invoked)
  performance-report.md      <- performance-engineer output (if invoked)
  data-engineering-notes.md  <- data-engineer output (if invoked)
  implementation-notes.md    <- developer output
  code-review-report.md      <- code-reviewer output
  accessibility-report.md    <- accessibility-engineer output (if invoked)
  security-report.md         <- security-reviewer output (if invoked)
  qa-report.md               <- qa-engineer output
  observability-report.md    <- sre-engineer output
  docs-report.md             <- tech-writer output
  devops-report.md           <- devops-engineer output
  delivery-summary.md        <- final synthesis
  retrospective.md           <- auto-generated every 5th delivery (see Phase 4); otherwise only present if the user ran /retrospective manually
```

### Delivery Summary Format

```markdown
# Delivery Summary: [Feature Name]

## Pipeline Run
| Agent | Version | Status | Contract | Key Output |
|---|---|---|---|---|
| context-engineer | [x.y.z] | PASS | PASS (N retries) | [N files pinned, N KIs/ADRs surfaced, token budget: OK/WARNING] |
| analyst | [x.y.z] | PASS | PASS (N retries) | [N acceptance criteria, N architectural flags] |
| architect | [x.y.z] | PASS / SKIPPED | PASS (N retries) / n/a | [N structural decisions, RFC: yes/no] |
| performance-engineer | [x.y.z] | PASS / SKIPPED | PASS (N retries) / n/a | [N SLAs verified, N recommendations] |
| data-engineer | [x.y.z] | PASS / SKIPPED | PASS (N retries) / n/a | [N migrations, expand/contract phase] |
| developer | [x.y.z] | PASS | PASS (N retries) | [N files created, N modified, N refactoring ops] |
| code-reviewer | [x.y.z] | PASS | PASS (N retries) | [Design score: C/Co/Cu/Cr — APPROVED] |
| accessibility-engineer | [x.y.z] | PASS / SKIPPED | PASS (N retries) / n/a | [N violations found, N fixed] |
| security-reviewer | [x.y.z] | PASS / SKIPPED | PASS (N retries) / n/a | [N findings, N critical fixed] |
| qa-engineer | [x.y.z] | PASS | PASS (N retries) | [N tests, N passed, SLAs verified: yes/no] |
| sre-engineer | [x.y.z] | PASS | PASS (N retries) | [N spans added, N alerts configured] |
| tech-writer | [x.y.z] | PASS | PASS (N retries) | [N docs updated] |
| devops-engineer | [x.y.z] | PASS | PASS (N retries) | [N CI changes, N env vars] |

Version is each agent's `version:` frontmatter field in `shared/agents/` at the time it ran — read it fresh
per run, don't cache it, since a mid-pipeline prompt edit (rare, but possible on a long-running delivery)
should be reflected accurately rather than assumed stale.

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
- Never skip context-engineer — it always runs before analyst. If it fails or is skipped, analyst and developer fall back to unscoped codebase exploration and MUST note this as a context-debt item in their output
- Never skip the analyst — it is always first
- Never let the developer start without analysis.md
- Never let analyst or developer ignore an existing context-manifest.md — its pinned files and pruning checklist take precedence over ad-hoc exploration
- Never let an artifact from a contract-bound agent (analyst, architect, developer, code-reviewer, security-reviewer, qa-engineer, sre-engineer) proceed to the next step while `validate-artifact` reports FAIL — send it back to the producing agent with the specific violations listed, and re-validate before continuing
- The structural contract check (validate-artifact) and the qualitative CHANGES REQUESTED loop (code-reviewer) are independent gates — passing one does not satisfy the other
- Never send CHANGES REQUESTED code to the security reviewer or QA
- Never ship to Friday without explicit "ship" or "yes" from the user
- Never persist artifacts to docs/features/ until the delivery summary is written
- Never overwrite an existing workspace artifact without first backing it up to `.claude/feature-workspace/.history/` — rollback depends on that history existing
- Never trust a `pipeline-state.json` entry whose checksum doesn't match the artifact currently on disk — re-run that step instead of resuming past it
- Never fabricate timing data in `pipeline-trace.json` — if a step's start/end time wasn't actually observed, omit the entry rather than estimating it
- The feature archive in docs/features/<feature-name>/ is append-only — never delete prior delivery artifacts

## Standalone Mode
All agents run locally. Friday POST is the only external call — non-blocking if Friday is not running. Artifact persistence to docs/features/ works entirely offline.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
