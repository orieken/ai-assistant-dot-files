# Gemini Antigravity Verification Protocol

**Update 2026-07-02: Tests 1-4 below are confirmed** — see
[results/antigravity-2026-07-02.md](results/antigravity-2026-07-02.md) for the full report. Kept here for
regression testing after future changes, and because one real gap remains open (Test 5).

## What was generated, and confirmed status
| Artifact | Path | Status |
|---|---|---|
| ~~Legacy instructions file~~ | ~~`.gemini/antigravity/instructions.md`~~ | **Confirmed NOT read (2026-07-02). Removed.** |
| Cross-tool agents file | `AGENTS.md` (repo root) | **Confirmed read** — injected as `<RULE[AGENTS.md]>` in the system prompt |
| Skills (global) | `~/.gemini/config/skills/` -> `shared/skills/` (via `install.sh --global`) | **Confirmed** — this is where skills actually loaded from in the 2026-07-02 test |
| Skills (project) | `.agents/skills/` -> `shared/skills/` (via `install.sh --project`) | **Not yet exercised** — didn't exist at session start in the 2026-07-02 test, so it fell back to the global root. See Test 5. |
| Rules (project) | `.agents/rules/` -> `shared/rules/` | **Not yet exercised directly** — `AGENTS.md` already carries rules content, so this may be redundant; not contradicted either |

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

## Test 5 — Does project-level `.agents/skills/` work when it exists from the start? (OPEN)
The one thing 2026-07-02's test didn't confirm: whether `.agents/skills/`/`.agents/rules/` (the
**project**-scoped convention, as opposed to the global root that was actually exercised) works when present
*before* Antigravity's session starts.

1. Run `./install.sh --project /path/to/a/fresh/scratch/dir --platform gemini` **first**, confirming
   `.agents/skills/` and `.agents/rules/` exist in that directory before you open it.
2. Open that directory fresh in Antigravity (not this repo, and not a directory that already had a session
   running before the install completed).
3. Repeat Test 3's questions. If it lists skills from `.agents/skills/` specifically (rather than falling
   back to the global root, if that's even distinguishable), that confirms the project-scoped path too.

## Report back
```
- [ ] Test 5 (project-level skills): confirmed working from session start? Y/N
- Antigravity version used: ___
- Anything unexpected: ___
```
