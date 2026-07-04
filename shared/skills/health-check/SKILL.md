---
name: health-check
description: Validates the ai-assistant-dot-files installation — symlinks, agent/skill frontmatter, platform config drift, domain dictionary orphans, inter-agent contracts, changelog/version consistency, Knowledge Item frontmatter, and memory registry integrity. Wraps scripts/health-check.sh for everything scriptable and adds AI judgment on top for anything that requires reading prose.
triggers:
  keywords: ["health-check", "health", "check installation", "verify setup", "check setup"]
  intentPatterns: ["Check my setup", "Is everything installed?", "Run health check", "Verify my installation", "/health-check"]
standalone: true
---

## When To Use
After running `install.sh` (it runs this automatically at the end unless `--dry-run`), after a `git pull`,
or whenever agents/skills aren't loading as expected, or you want an overall health snapshot of the
framework itself.

Do NOT use for debugging application code in a project you're building features for — use
`/debug-environment` instead. Do NOT use for checking code quality — use `/complexity-check` or
`/design-review` instead.

## Context To Load First
1. `scripts/health-check.sh` — the deterministic backbone this skill wraps
2. `shared/agents/CHANGELOG.md`, `shared/contracts/`, `DOMAIN_DICTIONARY.md` — for the checks the script
   performs, in case a finding needs more context than the script's one-line output gives

## Process

1. **Run the script**: `bash scripts/health-check.sh --verbose`. This performs 9 checks, all scriptable:
   - Symlinks (`.claude/{agents,skills,rules}` -> `shared/` equivalents) resolve correctly
   - Every agent's frontmatter has `name`, `description`, `tools`, `model`, `version`
   - Every skill's `SKILL.md` has `name`, `description`, `triggers`, `standalone`
   - Platform configs match `shared/` (delegates to `scripts/check-parity.sh` rather than duplicating it)
   - Domain dictionary terms that appear nowhere else in `shared/`/`docs/` (best-effort — see Guardrails)
   - Every contract-bound agent (`shared/contracts/`) has its contract file present
   - Every agent's current version appears in `shared/agents/CHANGELOG.md`
   - Every Knowledge Item (`shared/knowledge/`, `.claude/knowledge/`) has valid frontmatter
   - `shared/memory-registry.json` is valid, every non-optional path it declares exists, and no two KIs
     share an exact frontmatter `name:` (a deterministic subset of what `memory-engineer` judges more fully)

2. **Add judgment the script can't**: for each `WARN` the script reports (it never hard-fails on these,
   since they need a human/AI read), decide whether it's real:
   - A domain term with zero references *inside this repo* may still be intentional — some
     `DOMAIN_DICTIONARY.md` terms (e.g. Saturday/Sunday framework classes like `BaseSite`, `BaseApiClient`)
     describe patterns that show up in *generated project code*, not in this repo's own `shared/`/`docs/`.
     Read the term's description before recommending removal.
   - A version/changelog mismatch warning might mean the changelog entry uses different wording than an
     exact string match caught — read the actual changelog section for that agent before concluding it's
     truly undocumented.

3. **If the user asked for a fix**: re-run with `--fix` — it regenerates configs on drift and recreates
   broken/missing symlinks. It does not touch anything else (contracts, changelog, KI frontmatter, domain
   dictionary) since those need human judgment about the *right* fix, not just a mechanical one.

4. **Produce the health report** — synthesize the script's output plus your judgment calls into the format
   below; don't just paste the raw script output.

## Output Format

```markdown
# Installation Health Check

Date: [YYYY-MM-DD]
Repository: [path to this repo]

## Overall Status
HEALTHY (0 fails) | DEGRADED (warns only) | BROKEN (1+ fails)

## Results
| Check | Pass | Warn | Fail |
|---|---|---|---|
| Symlinks | [N] | — | [N] |
| Agent frontmatter | [N] | — | [N] |
| Skill frontmatter | [N] | — | [N] |
| Platform config drift | [N] | — | [N] |
| Domain dictionary terms | [N] | [N] | — |
| Inter-agent contracts | [N] | — | [N] |
| Changelog/version consistency | [N] | [N] | — |
| Knowledge Item frontmatter | [N] | — | [N] |

## Failures (if any)
- [component] — [what's wrong] — [exact fix: run `scripts/health-check.sh --fix`, or manual steps if not auto-fixable]

## Warnings Worth a Human Look
- [term/agent] — [the script's warning] — [your judgment: real issue or expected/benign, and why]

## Recommended Next Steps
1. [Specific command or edit]
— or "No issues found."
```

## Guardrails
- **Never** modify any files yourself beyond what `scripts/health-check.sh --fix` already does — this skill
  reports and (optionally) triggers the script's narrow auto-repair, it doesn't freelance additional fixes.
- **Never** suppress a warning without stating why you believe it's benign — "probably fine" isn't a
  judgment, it's a guess. Read the term/entry in question first.
- **Always** run the underlying script rather than re-deriving its checks by hand — it's the single source
  of truth for what "healthy" means here, and re-deriving invites drift between the skill and the script.

## Standalone Mode
`scripts/health-check.sh` is pure local filesystem operations — no external services required. The skill's
judgment layer is local reasoning over the script's output.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
