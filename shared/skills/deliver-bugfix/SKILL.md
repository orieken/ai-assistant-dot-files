---
name: deliver-bugfix
description: Lightweight pipeline for fixing a known, reproducible bug — reproduce-first discipline, developer fix, code-reviewer, QA verify. Reuses existing agents and approval gates. Escalates to /deliver-feature when the fix scope exceeds a single bounded context.
triggers:
  keywords: ["deliver-bugfix", "fix bug", "bug fix", "reproduce-first", "bugfix"]
  intentPatterns: ["Fix bug *", "Deliver bugfix for *", "Reproduce and fix *", "/deliver-bugfix *"]
standalone: true
---

## When To Use

When you have a **known, reproducible bug** — a symptom is confirmed, the location is
approximately known, and the fix should not require architectural changes.

**Use `/deliver-feature` instead when any of these apply:**

- The root cause is unknown and needs `/analyst` exploration
- The fix touches more than one bounded context (≥3 files in separate domains)
- A contract or schema change is required
- The scope feels like a small feature that was missing, not a regression

**Escalation trigger**: at any phase, if the developer discovers the fix requires a contract
change, a database migration, or changes in ≥3 files across distinct bounded contexts — stop,
write a summary, and route to `/deliver-feature` instead. Never scope-creep a bugfix into a
feature silently.

## Context To Load First

1. Error logs, stack trace, or reproduction steps provided by the reporter
2. Relevant source files (grep for the affected function or module)
3. `docs/features/<related-feature>/` — the original delivery that introduced the bug, if known
4. `docs/incidents/<date>-<slug>.md` — the incident record that triggered this bugfix, if one exists
5. `CLAUDE.md` — complexity limits and naming rules

## Phases

### Phase 1 — Reproduce

**Characterization test (Michael Feathers rule, already in `CLAUDE.md`)**: Write a test that
**fails** against the current code, capturing the *actual buggy behavior*. Commit nothing yet —
the red test is evidence the bug is real and reproduced.

```
Bug slug: <kebab-case-description>
Workspace: docs/features/<bug-slug>/
```

1. Write the failing characterization test.
2. Run it — confirm it fails with the expected error/output (not some unrelated failure).
3. Note: do NOT fix the bug here. If the test passes without changes, the bug is not reproduced —
   investigate further before proceeding.

**Commit** (after user confirms the red test is correct):
`test(<scope>): characterization test for <bug-slug> (deliver-bugfix phase 1)`

### Phase 2 — Fix

Invoke the `developer` agent: implement the minimal fix that makes the failing test pass. The fix
must:
- NOT add new behavior beyond making the failing test pass
- NOT introduce new dependencies or change API signatures
- Preserve all passing tests (the full suite must remain green after the fix)

```
Developer produces: docs/features/<bug-slug>/implementation-notes.md
```

**Approval gate**: commit requires `shared/rules/approval-gates.md` Gate #2 confirmation.

### Phase 3 — Review

Invoke the `code-reviewer` agent against the fix:
- APPROVED → proceed to Phase 4
- CHANGES REQUESTED → loop back to Phase 2; do not proceed with unresolved review comments

The code-reviewer also verifies the fix did not silently broaden scope.

```
Code-reviewer produces: docs/features/<bug-slug>/code-review-report.md
```

### Phase 4 — Verify

Run the full test suite. Both of these must pass before closing:
- The characterization test from Phase 1 is now **green**
- The full suite has **no new failures** (no regressions introduced by the fix)

If `security-reviewer` is warranted (fix touches auth, sessions, inputs, or secrets), invoke it
now. Otherwise skip — the full-pipeline threat model is not required for a targeted bugfix.

```
Verification produces: docs/features/<bug-slug>/qa-report.md (minimal — test results only)
```

### Phase 5 — Record

Write the minimal artifacts to `docs/features/<bug-slug>/`:
- `implementation-notes.md` (from Phase 2)
- `code-review-report.md` (from Phase 3)
- `qa-report.md` (from Phase 4)
- A CHANGELOG entry in the project's CHANGELOG.md: `fix(<scope>): <description>` line under
  the Unreleased section

**Link to incident** (if triggered from an incident): add a "Fixed by" line to
`docs/incidents/<date>-<slug>.md` pointing to `docs/features/<bug-slug>/`.

The `docs/features/<bug-slug>/` archive is first-class input for `/retrospective`, `/extract-lessons`,
and the feature-archive retrieval source — bugfixes and features share the same archive,
so production feedback loops see both.

## What This Pipeline Skips vs `/deliver-feature`

| Stage | deliver-feature | deliver-bugfix |
|---|---|---|
| spec-writer | Required | **Skipped** — the bug report IS the spec |
| product-owner | Required | **Skipped** — bugs don't need ROI review |
| context-engineer | Optional | **Skipped** — fix is targeted; context is loaded manually |
| analyst | Required | **Skipped** — replaced by characterization test |
| architect | Optional | **Skipped** — escalate to deliver-feature if needed |
| performance-engineer | Optional | **Skipped** |
| data-engineer | Optional | **Skipped** — escalate if migration required |
| developer | Required | Required |
| code-reviewer | Required | Required |
| accessibility-engineer | Optional | **Skipped** — unless fix touches UI |
| security-reviewer | Optional | **Only if fix touches auth/sessions/inputs** |
| qa-engineer | Required | Minimal — test run only (no new test-writing phase) |
| sre-engineer | Optional | **Skipped** |
| tech-writer | Required | **Minimal** — CHANGELOG entry only |
| devops-engineer | Required | **Skipped** — bugfix deployment is same as normal deploy |

## Output

```
docs/features/<bug-slug>/
  implementation-notes.md
  code-review-report.md
  qa-report.md
CHANGELOG.md  (new entry under Unreleased)
docs/incidents/<date>-<slug>.md  (updated "Fixed by" link, if applicable)
```

## Guardrails
- NEVER add new behavior in the same run. If the fix accidentally adds a feature, split the commit.
- NEVER skip the characterization test. A fix without a red test first is not reproduce-first.
- NEVER proceed to Phase 2 if the Phase 1 test passes on the unfixed code — the bug is not
  reproduced; investigate further.
- If the fix scope expands (new files, contract changes, migration needed), halt and escalate to
  `/deliver-feature` — do not silently absorb the scope increase.
- Approval gates from `shared/rules/approval-gates.md` apply unchanged. Incident urgency does not
  waive Gate #2 (commit) or Gate #8 (deploy).

## Standalone Mode
Orchestrates existing agents (`developer`, `code-reviewer`) and standard test tooling.
No external services required beyond the project's own test runner.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
