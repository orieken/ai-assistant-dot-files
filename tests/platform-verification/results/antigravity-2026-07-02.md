# Gemini Antigravity Verification Results: 2026-07-02

Protocol: [antigravity.md](../antigravity.md) (pre-update version, before this run's findings were folded
back into it). Run by the user against Antigravity v2.1.1, installed via `./install.sh --project`.

## Test 1 — AGENTS.md
**PASS.** Approval gates listed, injected as `<RULE[AGENTS.md]>` in the system prompt. Matches
`shared/rules/approval-gates.md`.

## Test 2 — Legacy `.gemini/antigravity/instructions.md`
**CONFIRMED NOT READ.** Not present in the system prompt instructions at all — only `AGENTS.md` loaded.
Consequence: the file has been removed from `scripts/generate-configs.sh`'s output entirely, along with the
leftover `.gemini/antigravity/{agents,rules,skills}` symlinks that predated this investigation and pointed at
`.claude/{agents,rules,skills}` — none of that was ever confirmed read either.

## Test 3 — Skill recognition and invocation
**PASS, with an important correction to the original hypothesis.** All 48 skills were recognized — but not
from the project-level `.agents/skills/` this framework generates for `--project` installs. Since
`.agents/` didn't exist in the project root at session start, Antigravity fell back to its **global**
customizations root at `~/.gemini/config/skills/`. This matches a detail from the original codelab research
that hadn't been acted on yet (`install.sh --global` was symlinking to `~/.agents/skills`, not
`~/.gemini/config/skills/`). Fixed: `install_antigravity()` in `install.sh` now targets
`~/.gemini/config/skills/` for `--global` installs specifically.

Skill invocation itself was genuine, not descriptive: asking it to run `complexity-check` against
`sample.go` produced a real complexity/LOC evaluation (complexity 1, length 10 lines, both under threshold)
using the actual limits from `ARCHITECTURE_RULES.md`, not a generic response.

## Test 4 — Rule recognition on a fixture
**PASS.** All 5 planted issues in `sample.go` were flagged: Clean Architecture dependency direction
violation, missing HTTP timeout, `interface{}` instead of a typed return, swallowed errors, SQL injection
risk. Sourced from `AGENTS.md`'s inlined rules content.

## What changed as a result
- Removed `.gemini/antigravity/instructions.md` and the `.gemini/antigravity/{agents,rules,skills}` symlinks
  entirely — confirmed dead weight, not just "harmless fallback."
- `scripts/generate-configs.sh` no longer generates the legacy instructions file; `AGENTS.md` is now Gemini's
  only generated artifact.
- `install.sh`'s `install_antigravity()` now symlinks to the *confirmed* global skills root
  (`~/.gemini/config/skills/`) instead of the unconfirmed guess (`~/.agents/skills`).
- `shared/platform-registry.json`'s `gemini` entry rewritten from "medium confidence" hedging to confirmed
  facts, with the one still-open question (project-level `.agents/skills/` specifically) called out.
- `tests/platform-verification/antigravity.md` updated with a new Test 5 for that remaining gap.

## Still open
Whether `.agents/skills/`/`.agents/rules/` (the **project**-scoped convention, used by `--project` installs)
works when present from session start — this test's `.agents/` didn't exist yet, so it never got exercised.
See Test 5 in the updated protocol.
