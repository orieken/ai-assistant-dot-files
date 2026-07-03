# GitHub Copilot Verification Protocol

**Update 2026-07-02**: Tests 1 and 2 are confirmed — see
[results/copilot-2026-07-02.md](results/copilot-2026-07-02.md). Test 3 (`.test.ts`) came back inconclusive:
the original expected guidance ("behavior over implementation, descriptive names, mocking") was written
generically enough that a capable reviewer would say the same thing whether or not
`testing.instructions.md` actually loaded — not a distinctive enough signal. Replaced with a sharper
question below.

Confirms `.github/copilot-instructions.md` (repo-wide) and the `.github/instructions/*.instructions.md`
(path-scoped) files both load and combine as GitHub's docs claim. Path-scoped instructions are documented as
supported in **Copilot Chat in VS Code, Visual Studio, and Copilot cloud agent** — not necessarily every
Copilot surface, so run this in VS Code's Copilot Chat if you have a choice.

## Setup
Nothing to install — both files are already generated and committed. Open this repo in VS Code (or your
Copilot-enabled editor) with Copilot Chat available.

## Test 1 — Repo-wide instructions (no file needed)
**CONFIRMED 2026-07-02.** Open Copilot Chat with no file open and ask:

> What are the approval gates in this project?

**Expected and observed**: Copilot references the approval gates from `.github/copilot-instructions.md` —
this file always applies regardless of what's open.

## Test 2 — Path-scoped instructions combine with repo-wide (`.go` / `.vue`)
**CONFIRMED 2026-07-02.** Open `tests/platform-verification/fixtures/sample.go` and ask:

> Review this file for any issues.

**Expected and observed**: Go-specific guidance (typed interfaces instead of `interface{}`, parameterized
queries, explicit timeouts, structured logging, Clean Architecture dependency direction) appeared, matching
`.github/instructions/go-backend.instructions.md`'s actual content closely enough to confirm it was read,
not just generic advice.

Repeat with `tests/platform-verification/fixtures/sample.vue` — expect and confirm Vue-specific guidance
(Composition API with `<script setup>`, Tailwind instead of custom CSS, extracting business logic to a
composable/service) matching `vue-frontend.instructions.md`.

## Test 3 — Does `testing.instructions.md` actually attach to `.test.*` files? (needs a sharper prompt)
Don't just ask "review this file" against `sample.test.ts` — its planted violations (vague name, testing
private-looking state, no mock) are generic enough that any capable reviewer flags them with or without
`testing.instructions.md` loaded, which is exactly what made the first run of this test inconclusive. Ask
instead, with the file open:

> What testing framework or pattern does this project's rules say I should use for this file? Name it specifically.

**Expected if `testing.instructions.md` is genuinely attaching**: it should name the **Saturday Framework**,
**Site-Centric pattern** (`BaseSite`/`BasePage`/`BaseElement`/`BaseFlow`), **Vitest**, the custom `api`
fixture, `BaseApiClient`, `CircuitBreaker`, or the **85% coverage threshold** — something specific from the
actual file content. If it can't name any of these, the scoped file isn't attaching for `.test.*` files in
whatever surface you're using, even though `.go`/`.vue` clearly did.

## Test 4 — Confirm which surface actually reads path-scoped files
If Test 2 or Test 3 don't show scoped guidance, try the same fixture + prompt in a different Copilot surface
(VS Code Copilot Chat vs. Copilot code review on a PR vs. GitHub.com chat) — GitHub's docs say support
varies by surface. Note which surface you tested in either case.

## Report back
```
- [ ] Test 3 (.test.ts, sharper prompt): did it name a specific Saturday/Sunday Framework term? Y/N — which one?
- Copilot surface(s) tested (VS Code Chat / code review / GitHub.com / other): ___
- Anything unexpected: ___
```
