# GitHub Copilot Verification Protocol

**Update 2026-07-02: all tests below are confirmed** — see
[results/copilot-2026-07-02.md](results/copilot-2026-07-02.md) for the full report, including how Test 3's
first attempt came back inconclusive (the original expected guidance was written generically enough that a
capable reviewer would say the same thing whether or not `testing.instructions.md` actually loaded) before
the sharper re-run below definitively confirmed it. Kept here for regression testing after future changes.

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

## Test 3 — Does `testing.instructions.md` actually attach to `.test.*` files?
**CONFIRMED 2026-07-02**, on the second attempt — the first "review this file" prompt was inconclusive
because `sample.test.ts`'s planted violations are generic enough that any capable reviewer flags them
without the scoped file loaded at all. The disambiguating question, asked with the file open:

> What testing framework or pattern does this project's rules say I should use for this file? Name it specifically.

**Observed**: correctly named the **Sunday Framework** (Vitest for unit tests, Playwright for
integration/E2E), citing `testing.instructions.md`, and — unprompted — correctly distinguished that a UI/E2E
file would instead need the **Saturday Framework**'s Site-Centric pattern, naming all four classes
(`BaseSite`/`BasePage`/`BaseElement`/`BaseFlow`). No generic AI response invents these specific
day-of-week-themed framework names or class names — definitive confirmation the scoped file is genuinely
read for `.test.*` files, not just `.go`/`.vue`.

## Test 4 — Confirm which surface actually reads path-scoped files
Not needed this run — VS Code Copilot Chat confirmed all three scoped files directly. Keep this test in
reserve if a future check in a different surface (code review, GitHub.com chat) comes back negative.

## Report back
All tests confirmed as of 2026-07-02 — nothing outstanding for this protocol. Re-run after future changes to
`.github/instructions/*.instructions.md` generation to catch regressions.
