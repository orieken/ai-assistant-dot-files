# Gemini Antigravity Verification Protocol

**This is the important one.** Everything generated for Antigravity is based on secondary sources (a Google
codelab, not primary reference docs — the official `antigravity.google/docs` page didn't return usable
content when fetched), so this protocol is designed to figure out which of the four generated/symlinked
artifacts actually gets read, not just confirm a single expected behavior. Please run all four tests if you
can — partial "it kind of worked" results are exactly the diagnostic signal needed here.

## What was generated, and why each exists
| Artifact | Path | Confidence |
|---|---|---|
| Legacy instructions file | `.gemini/antigravity/instructions.md` | This is what the framework generated *before* this investigation — kept as a fallback |
| Cross-tool agents file | `AGENTS.md` (repo root) | Medium — secondary sources say Antigravity reads this as its project rules file |
| Symlinked skills | `.agents/skills/` -> `shared/skills/` | Medium — a codelab described `.agents/skills/<name>/SKILL.md` as Antigravity's skill format, which looks structurally identical to this repo's own `SKILL.md` convention |
| Symlinked rules | `.agents/rules/` -> `shared/rules/` | Medium — same codelab, less detail given |

Note: `.agents/skills/` and `.agents/rules/` are only created by `install.sh --project`/`--global`, not
present in this repo's own root (they're an install-time artifact, not a generated-and-committed one — see
`docs/MIGRATION.md`'s convention). **Run `./install.sh --project /path/to/a/scratch/dir --platform gemini`
first** (or `--global` if you're comfortable symlinking into `~/.agents/`), then open *that* directory in
Antigravity — not this repo directly — unless you also want to test from `~/.gemini/` conventions.

## Test 1 — Does it read AGENTS.md?
With the installed directory open in Antigravity, start a new agent session and ask:

> What are the approval gates defined in this project's rules?

**Expected if AGENTS.md is read**: it should list the gates (commits, deploys, migrations, external API
calls) — `AGENTS.md` inlines all of `shared/rules/approval-gates.md`.

## Test 2 — Does it read the legacy instructions.md?
Ask the same question a different way:

> Does this project have a file at .gemini/antigravity/instructions.md, and if so, what does it say about
> approval gates?

**Expected if the legacy path is read**: same content as Test 1, sourced from the other file. If Test 1
already worked, this tells us whether *both* are being read (redundant but harmless) or just one.

## Test 3 — Does it recognize skills?
Ask:

> What skills or capabilities are available to you in this project?

**Expected if `.agents/skills/` is recognized**: Antigravity should list some subset of the 48 skills from
`shared/skills/` (e.g. `deliver-feature`, `complexity-check`, `threat-model`) using their `name`/`description`
frontmatter — the same format Claude Code's skills use. If it lists nothing or says it has no project-level
skills, the symlink likely isn't being picked up, or Antigravity's actual skill format differs from what the
codelab described.

If it does recognize skills, try invoking one directly:

> Run the complexity-check skill against tests/platform-verification/fixtures/sample.go

**Expected**: it should follow `shared/skills/analyze-complexity/SKILL.md`'s process (or `complexity-check`'s,
whichever name it resolves) rather than just doing a generic ad-hoc review.

## Test 4 — Does it recognize rules?
Open `tests/platform-verification/fixtures/sample.go` (copy it into the installed directory if needed) and
ask:

> Review this file for issues.

**Expected if `.agents/rules/` is read**: flags on the same issues as the Cursor/Copilot protocols (typed
interfaces, parameterized queries, timeouts, error handling) — sourced from `shared/rules/` content.

## Report back
This is the one where "none of it worked" or "only AGENTS.md worked, not skills" is genuinely useful
information, not a failure — it tells us whether to keep the current hybrid approach, simplify to just
`AGENTS.md`, or investigate further.

```
- [ ] Test 1 (AGENTS.md): approval gates listed? Y/N
- [ ] Test 2 (legacy instructions.md): also read, or only AGENTS.md, or neither? ___
- [ ] Test 3 (skills recognized): did it list any shared/skills/ entries? Y/N — which ones?
- [ ] Test 3 (skill invocation): did invoking one actually follow that skill's process? Y/N
- [ ] Test 4 (rules recognized): flagged sample.go issues? Y/N — which ones?
- Antigravity version used: ___
- Installed via --project or --global: ___
- Anything unexpected: ___
```
