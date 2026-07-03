# Cursor Verification Protocol

Confirms that `.cursor/rules/*.mdc` actually loads as expected — the two always-apply rules, the
glob-triggered rules, and the 24 per-agent persona files. Run this in Cursor with this repo open as the
workspace.

Report back using the checklist at the bottom — you don't need to run every test, but the more you run, the
more confidence we have.

## Setup
Nothing to install — `.cursor/rules/*.mdc` is already generated and committed. Just open this repo in Cursor.

## Test 1 — Always-apply rules (no file needed)
Open a new Cursor Chat (no file open, or any file) and ask:

> What are the approval gates in this project? List them.

**Expected**: Cursor should list the gates from `approval-gates.mdc` (creating a git commit, shipping to
Friday, database migrations, external API calls, deploying, etc.) even with no relevant file open — this
file is `alwaysApply: true`.

Then ask:

> What specialized personas/agents are available in this project?

**Expected**: Cursor should list agent names from `agent-roster.mdc` (analyst, architect, developer,
code-reviewer, security-reviewer, etc. — 24 total) — also `alwaysApply: true`.

## Test 2 — Auto-Attach rules (glob-triggered)
Open `tests/platform-verification/fixtures/sample.go` in the editor, then ask Cursor Chat:

> Review this file for any issues.

**Expected**: Since `.go` matches `go-backend.mdc`'s and `architecture.mdc`'s globs, Cursor should flag
several of these (it doesn't need to catch all of them, but should catch at least 2-3):
- `interface{}` instead of a typed interface
- Raw SQL string concatenation (SQL injection risk) instead of a parameterized query
- No explicit timeout on the `http.Get` call
- Ignored error return values (`_`)
- No repository interface — the function talks to `database/sql` directly

If Cursor's review doesn't mention *any* of these, the rule likely isn't Auto-Attaching correctly for `.go`
files — worth reporting.

Repeat with `tests/platform-verification/fixtures/sample.vue` (should trigger `vue-frontend.mdc` +
`architecture.mdc`/`design-principles.mdc` — expect flags on Options API instead of Composition API, custom
CSS instead of Tailwind, business logic in the component) and
`tests/platform-verification/fixtures/sample.test.ts` (should trigger `testing.mdc` — expect flags on the
vague test name, testing internal state instead of behavior, no mocking).

## Test 3 — Manual persona invocation
With any file open, ask Cursor Chat:

> @developer.mdc — review the sample.go fixture and tell me what you'd change.

**Expected**: Cursor should respond in the voice/process of the `developer` persona from
`shared/agents/developer.md` (TDD-first, named refactoring log, self-review checklist) — not a generic
code review. This confirms the per-agent persona files are both present and manually invocable.

## Report back
```
- [ ] Test 1 (always-apply): approval gates listed correctly? Y/N
- [ ] Test 1 (always-apply): persona roster listed correctly? Y/N
- [ ] Test 2 (.go Auto-Attach): flagged go-backend/architecture issues? Y/N — which ones?
- [ ] Test 2 (.vue Auto-Attach): flagged vue-frontend issues? Y/N — which ones?
- [ ] Test 2 (.test.ts Auto-Attach): flagged testing issues? Y/N — which ones?
- [ ] Test 3 (manual persona): did @developer.mdc actually change Cursor's response style/process? Y/N
- Cursor version used: ___
- Anything unexpected: ___
```
