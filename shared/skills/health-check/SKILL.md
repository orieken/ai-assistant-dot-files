---
name: health-check
description: Validates the ai-assistant-dot-files installation — symlinks, agent/skill frontmatter, platform config drift, domain dictionary orphans, inter-agent contracts, changelog/version consistency, Knowledge Item frontmatter, and memory registry integrity. Wraps scripts/health-check.sh for everything scriptable and adds AI judgment on top for anything that requires reading prose. Also detects and reports the presence of opt-in AOS layers (telemetry, evaluation, hooks, orchestration, rag) and counter agents when they exist — never fails on their absence.
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

4. **Detect optional AOS layers** (v3.0+): check for the presence of each opt-in AOS layer and any
   counter agents. AOS layers are optional per the migration plan (`docs/aos/migration-plan.md`) —
   a missing layer is never a failure, only an "absent" observation:
   - `shared/telemetry/` — layer present if the directory exists with a `README.md`
   - `shared/evaluation/` — same shape
   - `shared/hooks/` — same shape (Phase 2)
   - `shared/orchestration/` — same shape (Phase 3)
   - `shared/rag/` — same shape (Phase 3)
   - Counter agents under `shared/agents/` matching `*-auditor.md`, `*-evaluator.md`, `*-reviewer.md`, or `*-validator.md`.
     DO NOT include the review-shaped producers whose names happen to end in `-reviewer` — those
     (`code-reviewer`, `security-reviewer`, `accessibility-engineer`-style producers) are producers-in-
     role-name, not AOS counter agents. Count counter agents mapped to the 15 pairs in `docs/aos/governance-pairs.md`:
     `memory-auditor`, `context-auditor`, `knowledge-auditor`, `prompt-evaluator`, `agent-evaluator`,
     `rule-auditor`, `pattern-reviewer`, `tool-validator`, `documentation-auditor`, `retrieval-evaluator`, `privacy-auditor`.
   - Validate hook configs under `shared/hooks/examples/` and `.claude/hooks/` against `shared/hooks/hooks-schema.md`.

5. **Produce the health report** — synthesize the script's output plus your judgment calls into the format
   below; don't just paste the raw script output. The AOS Layers section always appears at the bottom of
   the report, even when everything is absent — it's an inventory, not a pass/fail check.

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

## AOS Layers
Opt-in per the AOS migration plan — absence is never a failure, only an inventory observation.

| Layer | Present? | Notes |
|---|---|---|
| `shared/telemetry/` | Yes / No | [If present: brief list of files.] |
| `shared/evaluation/` | Yes / No | [Same shape] |
| `shared/hooks/` | Yes / No | [Present — Phase 2: hooks-schema + examples] |
| `shared/orchestration/` | Yes / No | [Same shape — Phase 3] |
| `shared/rag/` | Yes / No | [Same shape — Phase 3] |

Counter agents present under `shared/agents/`: [comma-separated list, or "None"]
Landed in v3.1: 11 counter agents (memory-auditor + 10 Phase 2 counter-auditors).
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
