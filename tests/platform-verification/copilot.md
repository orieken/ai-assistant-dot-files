# GitHub Copilot Verification Protocol

Confirms `.github/copilot-instructions.md` (repo-wide) and the new `.github/instructions/*.instructions.md`
(path-scoped) files both load and combine as GitHub's docs claim. Path-scoped instructions are documented as
supported in **Copilot Chat in VS Code, Visual Studio, and Copilot cloud agent** — not necessarily every
Copilot surface, so run this in VS Code's Copilot Chat if you have a choice.

## Setup
Nothing to install — both files are already generated and committed. Open this repo in VS Code (or your
Copilot-enabled editor) with Copilot Chat available.

## Test 1 — Repo-wide instructions (no file needed)
Open Copilot Chat with no file open and ask:

> What are the approval gates in this project?

**Expected**: Copilot should reference the approval gates from `.github/copilot-instructions.md` (commits,
deploys, migrations, external API calls, etc.) — this file always applies regardless of what's open.

## Test 2 — Path-scoped instructions combine with repo-wide
Open `tests/platform-verification/fixtures/sample.go` and ask:

> Review this file for any issues.

**Expected**: Because `.go` matches `go-backend.instructions.md`'s `applyTo` pattern, Copilot's answer should
include Go-specific guidance (typed interfaces instead of `interface{}`, parameterized queries, explicit
timeouts, structured logging) **in addition to** general repo-wide guidance — GitHub's docs say both
combine. If you only see generic feedback with no Go-specific rules mentioned, the path-scoped file likely
isn't being picked up in whatever Copilot surface you're using.

Repeat with `tests/platform-verification/fixtures/sample.vue` (expect Vue-specific guidance — Composition
API, Tailwind, no business logic in components) and `tests/platform-verification/fixtures/sample.test.ts`
(expect testing-specific guidance — behavior over implementation, descriptive names, mocking).

## Test 3 — Confirm which surface actually reads path-scoped files
If Test 2 doesn't show Go/Vue/testing-specific guidance, try the same fixture + prompt in a different
Copilot surface (VS Code Copilot Chat vs. Copilot code review on a PR vs. GitHub.com chat) — GitHub's docs
say support varies by surface. Note which surface you tested in either case; this tells us where the
scoped-instructions investment actually pays off.

## Report back
```
- [ ] Test 1 (repo-wide): approval gates mentioned without any file open? Y/N
- [ ] Test 2 (.go scoped): Go-specific guidance appeared? Y/N — which points?
- [ ] Test 2 (.vue scoped): Vue-specific guidance appeared? Y/N — which points?
- [ ] Test 2 (.test.ts scoped): testing-specific guidance appeared? Y/N — which points?
- Copilot surface(s) tested (VS Code Chat / code review / GitHub.com / other): ___
- Anything unexpected: ___
```
