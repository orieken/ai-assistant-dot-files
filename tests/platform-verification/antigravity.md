# Gemini Antigravity Verification Protocol

**Update 2026-07-02: all 5 tests below are confirmed** — see
[results/antigravity-2026-07-02.md](results/antigravity-2026-07-02.md) and
[results/antigravity-test5-2026-07-02.md](results/antigravity-test5-2026-07-02.md) for the full reports.
Kept here for regression testing after future changes.

## What was generated, and confirmed status
| Artifact | Path | Status |
|---|---|---|
| ~~Legacy instructions file~~ | ~~`.gemini/antigravity/instructions.md`~~ | **Confirmed NOT read. Removed.** |
| Cross-tool agents file | `AGENTS.md` (repo root) | **Confirmed read** — injected as `<RULE[AGENTS.md]>` in the system prompt |
| Skills (global) | `~/.gemini/config/skills/` -> `shared/skills/` (via `install.sh --global`) | **Confirmed** — where skills loaded from before `.agents/skills/` existed |
| Skills (project) | `.agents/skills/` -> `shared/skills/` (via `install.sh --project`) | **Confirmed** — a fresh project with `.agents/skills/` present from session start correctly listed this framework's distinctive skill names (`numpath-alignment`, `sunday-test-advisor`, etc.), merged alongside an unrelated global/built-in skill set |
| Rules (project) | `.agents/rules/` -> `shared/rules/` | Not yet exercised directly — `AGENTS.md` already carries rules content, so this may be redundant; not contradicted either |

## Test 1 — Does it read AGENTS.md?
**CONFIRMED 2026-07-02.** With the installed directory open in Antigravity, start a new agent session and ask:

> What are the approval gates defined in this project's rules?

**Expected and observed**: lists all 8 gates, matching `shared/rules/approval-gates.md` exactly.

## Test 2 — Does it read the legacy instructions.md?
**CONFIRMED NOT READ, 2026-07-02** — the file has been removed from generation. Skip this test on future runs
unless verifying the removal didn't regress anything.

## Test 3 — Does it recognize and invoke skills?
**CONFIRMED 2026-07-02**, via the global root. Ask:

> What skills or capabilities are available to you in this project?

**Observed**: listed all 48 skills correctly (e.g. `deliver-feature`, `complexity-check`, `threat-model`).

> Run the complexity-check skill against tests/platform-verification/fixtures/sample.go

**Observed**: genuinely executed `analyze-complexity`'s actual process (heuristic complexity/LOC evaluation
against the real thresholds), not a generic ad-hoc review.

## Test 4 — Does it recognize rules on a fixture?
**CONFIRMED 2026-07-02.** Open `tests/platform-verification/fixtures/sample.go` and ask:

> Review this file for issues.

**Observed**: flagged the Clean Architecture dependency violation, missing HTTP timeout, `interface{}`
instead of a typed return, swallowed errors, and the SQL injection risk — all five planted issues.

## Test 5 — Does project-level `.agents/skills/` work when it exists from the start?
**CONFIRMED 2026-07-02** (`results/antigravity-test5-2026-07-02.md`), via `./install.sh --project
~/antigravity-test --platform gemini` followed by a fresh Antigravity session on that new directory. Also
confirmed exact fidelity to skill file content in the process: `complexity-check` invocation quoted that
skill's own literal example threshold (`gocyclo -over 6`, mathematically identical to this project's "< 7"
rule, just phrased the other way) and wrote its report to the exact path `complexity-check/SKILL.md`
specifies (`.claude/feature-workspace/complexity-report.md`) — not generic or approximate behavior.

## Report back
All five tests are confirmed as of 2026-07-02 — nothing outstanding for this protocol. Re-run any of them
after future changes to `shared/skills/`, `AGENTS.md` generation, or `install.sh`'s Antigravity handling, to
catch regressions.
