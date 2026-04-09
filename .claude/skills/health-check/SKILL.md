---
name: health-check
description: Validates the ai-assistant-dot-files installation — checks symlinks, file references, cross-platform configs, and reports any broken or missing components.
triggers:
  keywords: ["health-check", "health", "check installation", "verify setup", "check setup"]
  intentPatterns: ["Check my setup", "Is everything installed?", "Run health check", "Verify my installation", "/health-check"]
standalone: true
---

## When To Use
When the user wants to verify their ai-assistant-dot-files installation is complete and functional. Useful after running `./install`, after a `git pull`, or when agents or skills are not loading as expected.

Do NOT use for debugging application code — use `/debug-environment` instead.
Do NOT use for checking code quality — use `/complexity-check` or `/design-review` instead.

## Context To Load First
1. The `install` script — to understand expected symlink targets
2. `.claude/agents/` — to enumerate all agents
3. `.claude/skills/` — to enumerate all skills

## Process

1. **Check symlink integrity** — verify that `~/.claude/agents/` and `~/.claude/skills/` exist and contain valid symlinks pointing back to this repository:
   - For each file in `.claude/agents/`, confirm `~/.claude/agents/<name>` exists and is a valid symlink
   - For each directory in `.claude/skills/`, confirm `~/.claude/skills/<name>` exists and is a valid symlink
   - Report broken symlinks (target missing or link does not exist)

2. **Check agent file references** — for each agent in `.claude/agents/`:
   - Verify the agent references files that exist (e.g., `DOMAIN_DICTIONARY.md`, `CLAUDE.md`, `ARCHITECTURE_RULES.md`)
   - Verify the frontmatter is well-formed (has `name`, `description`, `tools`, `model` fields)
   - Report any agents referencing non-existent files

3. **Check skill file references** — for each skill in `.claude/skills/`:
   - Verify `SKILL.md` exists in the skill directory
   - Verify the frontmatter is well-formed (has `name`, `description`, `triggers`, `standalone` fields)
   - If a `check.sh` or `run.sh` exists, verify it is executable

4. **Check cross-platform config files**:
   - `CLAUDE.md` (root) — exists and is non-empty
   - `ARCHITECTURE_RULES.md` — exists and is non-empty
   - `DOMAIN_DICTIONARY.md` — exists (warn if only template exists)
   - `.cursor/rules/` — directory exists with at least one rule file
   - `.windsurfrules` — exists and is non-empty
   - `.openai.md` — exists and is non-empty
   - `.github/copilot-instructions.md` — exists and is non-empty
   - `.gemini/antigravity/instructions.md` — exists and is non-empty

5. **Check documentation structure**:
   - `docs/README.md` — exists
   - `docs/features/README.md` — exists
   - `docs/adrs/README.md` — exists
   - `docs/patterns/README.md` — exists
   - `docs/runbooks/README.md` — exists

6. **Check settings.json** — verify `.claude/settings.json` is valid JSON and contains expected hook definitions.

7. **Produce the health report** — display inline to the user.

## Output Format

```markdown
# Installation Health Check

Date: [YYYY-MM-DD]
Repository: [path to this repo]

## Overall Status
HEALTHY | DEGRADED | BROKEN

## Symlinks
| Component | Expected | Status |
|---|---|---|
| ~/.claude/agents/ | [N] agents | PASS / [N] missing |
| ~/.claude/skills/ | [N] skills | PASS / [N] missing |

### Missing Symlinks (if any)
- ~/.claude/agents/[name] — run `./install` to fix
- ~/.claude/skills/[name] — run `./install` to fix

## Agents ([N] total)
| Agent | Frontmatter | File References | Status |
|---|---|---|---|
| analyst | PASS | PASS | HEALTHY |
| [name] | FAIL — missing model | PASS | DEGRADED |

## Skills ([N] total)
| Skill | SKILL.md | Frontmatter | Scripts | Status |
|---|---|---|---|---|
| deliver-feature | PASS | PASS | N/A | HEALTHY |
| [name] | PASS | FAIL | not executable | DEGRADED |

## Cross-Platform Configs
| Platform | File | Status |
|---|---|---|
| Claude | CLAUDE.md | PASS |
| Cursor | .cursor/rules/ | PASS |
| Windsurf | .windsurfrules | PASS |
| OpenAI | .openai.md | PASS |
| Copilot | .github/copilot-instructions.md | PASS |
| Gemini | .gemini/antigravity/instructions.md | PASS |

## Documentation Structure
| Directory | Status |
|---|---|
| docs/features/ | PASS |
| docs/adrs/ | PASS |
| docs/patterns/ | PASS |
| docs/runbooks/ | PASS |

## Recommended Fixes
1. [Specific fix with command to run]
2. [Another fix]

— or "No issues found."
```

## Guardrails
- Never modify any files during the health check — this is a read-only diagnostic.
- Never suppress warnings. If something is missing, report it even if the system still functions.
- Always recommend running `./install` as the first fix for symlink issues.
- If the repository path cannot be determined, ask the user to confirm it.

## Standalone Mode
Works entirely offline. All checks are local filesystem operations. No external services required.
