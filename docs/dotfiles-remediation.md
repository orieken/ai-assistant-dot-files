# Dotfiles Repo — Remediation Tasks

Hand this file to Claude Code from the root of your `ai-assistant-dot-files` repo.
Work through each task in order. Each task is self-contained and verifiable.

---

## Task 1 — Add `code-reviewer` agent to new agent set

**Problem:** The `code-reviewer` agent exists in `.claude/agents/code-reviewer.md` and
`templates/claude-feature-team/.claude/agents/code-reviewer.md` but was not included in the
recent additions. The `CLAUDE.md` orchestrator references it explicitly in the pipeline between
`developer` and `security-reviewer` — without it, unrefactored code reaches the security reviewer.

**Instructions:**
1. Read the existing `.claude/agents/code-reviewer.md` to understand the current version.
2. Verify it includes all of the following evaluation criteria:
   - Architecture & layer boundary violations (Clean Architecture)
   - Destructive database migration detection (Expand/Contract pattern)
   - Sandi Metz rules: classes ≤ 100 lines, methods ≤ 5 lines, ≤ 4 params
   - Frontend craftsmanship: semantic HTML, accessibility, no `onClick` on `<div>`
   - Shift-left reliability: explicit timeouts, idempotency on mutations, no N+1 queries
   - Fowler smells: Feature Envy, Magic Numbers, DRY violations, missing refactor pass
   - Ubiquitous Language violations against `DOMAIN_DICTIONARY.md`
3. If any of the above are missing, add them.
4. Confirm the same file exists at `templates/claude-feature-team/.claude/agents/code-reviewer.md`.
   If it doesn't, copy it there.

**Verification:** Both paths exist and contain all evaluation criteria listed above.

---

## Task 2 — Restore `deliver-feature` skill to the template

**Problem:** The scaffold script (`scaffold-team.sh`) copies `.claude/skills/` into target projects,
but the `deliver-feature` skill is missing from `templates/claude-feature-team/.claude/skills/`.

**Instructions:**
1. Read the existing `.claude/skills/deliver-feature/SKILL.md`.
2. Verify it orchestrates the full 8-agent pipeline in this order:
   ```
   Analyst → Architect* → Developer → Code-Reviewer → Security-Reviewer* → QA → Tech Writer → DevOps
   ```
   Where `*` = conditional agents (invoke-when rules apply).
3. If the pipeline order is wrong or agents are missing, update it.
4. Confirm the file exists at `templates/claude-feature-team/.claude/skills/deliver-feature/SKILL.md`.
   If it doesn't, copy it there.

**Verification:** `templates/claude-feature-team/.claude/skills/deliver-feature/SKILL.md` exists and
reflects the correct 8-step pipeline.

---

## Task 3 — Fix pipeline order in `templates/claude-feature-team/CLAUDE.md`

**Problem:** The orchestrator `CLAUDE.md` in the feature team template may not include
`code-reviewer` in the correct position (between `developer` and `security-reviewer`).

**Instructions:**
1. Read `templates/claude-feature-team/CLAUDE.md`.
2. Find the pipeline definition (the table and the numbered sequence).
3. Ensure `code-reviewer` appears as step 4, between `developer` (step 3) and
   `security-reviewer` (step 5). The full correct sequence is:
   ```
   1. analyst           → analysis.md
   2. architect*        → architecture-notes.md
   3. developer         → implementation-notes.md
   4. code-reviewer     → code-review-report.md
   5. security-reviewer*→ security-report.md
   6. qa-engineer       → qa-report.md
   7. tech-writer       → docs-report.md
   8. devops-engineer   → devops-report.md
   ```
4. Update the agent table, pipeline sequence, workspace convention section, and delivery
   summary pipeline checklist to all reflect this 8-step order.
5. Add a "When to invoke Code Reviewer" rule: always — no conditions, it runs after every
   developer pass and loops back if changes are requested.

**Verification:** Pipeline table, numbered sequence, workspace file list, and delivery summary
all show `code-reviewer` at step 4.

---

## Task 4 — Add `DOMAIN_DICTIONARY.md` enforcement to `analyst` agent

**Problem:** The analyst agent in `.claude/agents/analyst.md` does not include the
DDD Ubiquitous Language enforcement steps that the template version has. Specifically:
- Reading `DOMAIN_DICTIONARY.md` at the start of analysis
- Updating `DOMAIN_DICTIONARY.md` when new business concepts are introduced
- Mapping any synonym from the feature spec to the correct canonical term

**Instructions:**
1. Read `.claude/agents/analyst.md`.
2. Read `templates/claude-feature-team/.claude/agents/analyst.md` — it has the stronger version.
3. Add the following to the analyst's process (after reading the feature file, before exploring
   the codebase):
   - Read `DOMAIN_DICTIONARY.md` (or create it from `DOMAIN_DICTIONARY.template.md` if absent)
   - If the feature introduces new entities, value objects, or domain events: add them to
     `DOMAIN_DICTIONARY.md` before producing `analysis.md`
   - If the feature spec uses a synonym for an existing term: map it to the canonical term
     in the analysis and flag it
4. Add a `## Domain Language` section to the analyst's output format (in `analysis.md`) that
   lists all domain terms used in the feature and their canonical definitions.
5. Apply the same update to `templates/claude-feature-team/.claude/agents/analyst.md`.

**Verification:** Both analyst agent files include `DOMAIN_DICTIONARY.md` read + update steps
and the output format includes a `## Domain Language` section.

---

## Task 5 — Add `DOMAIN_DICTIONARY.md` enforcement to `developer` agent

**Problem:** The developer agent in `.claude/agents/developer.md` does not reference
`DOMAIN_DICTIONARY.md`. The template version does — it instructs the developer to use
exactly the terms defined in the Ubiquitous Language and never invent technical synonyms.

**Instructions:**
1. Read `.claude/agents/developer.md`.
2. Add the following under "Before Writing Code":
   - Read `DOMAIN_DICTIONARY.md`
   - Use exactly the terms defined — class names, method names, variable names must match
     the Ubiquitous Language
   - Do not invent technical synonyms for business concepts (e.g. if the dictionary says
     `User`, do not use `Client` or `AccountHolder`)
3. Apply the same update to `templates/claude-feature-team/.claude/agents/developer.md`.

**Verification:** Both developer agent files include `DOMAIN_DICTIONARY.md` read step and the
Ubiquitous Language enforcement rule under "Before Writing Code".

---

## Task 6 — Add `DOMAIN_DICTIONARY.template.md` to the feature team template

**Problem:** `templates/claude-feature-team/DOMAIN_DICTIONARY.template.md` exists but the
analyst references creating `DOMAIN_DICTIONARY.md` from it if the file doesn't exist. If
a team scaffolds a new project with `scaffold-team.sh`, they won't get the template and
the analyst will have nothing to bootstrap from.

**Instructions:**
1. Read `templates/claude-feature-team/DOMAIN_DICTIONARY.template.md`.
2. Verify it includes sections for: Core Domains, Entities, Value Objects, Domain Events,
   and Operations/Actions — each with clear instructions and examples.
3. If any sections are missing or examples are unclear, improve them.
4. Verify that `scaffold-team.sh` copies this file into the target project. Read the script
   and check if it's included. If not, add a copy step for it.

**Verification:** `DOMAIN_DICTIONARY.template.md` is complete and `scaffold-team.sh` deploys
it to target projects.

---

## Task 7 — Add `features/TEMPLATE.md` to the feature team template

**Problem:** `templates/claude-feature-team/features/TEMPLATE.md` may not exist or may be
missing the `## Test Approach` field that the analyst uses to detect the correct test framework
(Saturday/Cucumber vs Saturday/Playwright vs Sunday/API).

**Instructions:**
1. Check if `templates/claude-feature-team/features/TEMPLATE.md` exists.
2. If it doesn't exist, create it. It should include:
   - `## Summary` — what the feature does and why
   - `## Acceptance Criteria` — Given/When/Then format, BDD-style
   - `## Out of Scope` — explicit exclusions
   - `## Non-Functional Requirements` — performance SLAs, security notes, scaling concerns
   - `## Test Approach` (optional override) — `Saturday/Cucumber` | `Saturday/Playwright` |
     `Sunday/API` — leave blank for auto-detection
   - `## Open Questions` — ambiguities to resolve before analysis
3. If it exists, verify `## Test Approach` is present.

**Verification:** `templates/claude-feature-team/features/TEMPLATE.md` exists and contains
all sections listed above including `## Test Approach`.

---

## Task 8 — Verify `scaffold-team.sh` deploys all required files

**Problem:** The scaffold script may not copy all files that agents depend on, causing
silent failures when agents look for files that don't exist.

**Instructions:**
1. Read `scaffold-team.sh`.
2. Verify it copies ALL of the following into the target project:
   - `CLAUDE.md` (orchestrator)
   - `.claude/agents/` (all agent files including `code-reviewer`)
   - `.claude/skills/deliver-feature/SKILL.md`
   - `features/TEMPLATE.md`
   - `DOMAIN_DICTIONARY.template.md`
3. If any are missing from the copy steps, add them.
4. Add a verification step at the end that lists what was deployed:
   ```bash
   echo "Deployed agents: $(ls $TARGET_DIR/.claude/agents/ | wc -l) agents"
   echo "Deployed skills: $(ls $TARGET_DIR/.claude/skills/ | wc -l) skills"
   ```

**Verification:** Running `scaffold-team.sh` on a fresh empty directory produces a project
with all 8 agents, the deliver-feature skill, the feature template, and the domain dictionary
template.

---

## Final Check

After completing all tasks, run this verification checklist:

```bash
# From the repo root — all of these should exist
ls .claude/agents/analyst.md
ls .claude/agents/architect.md
ls .claude/agents/developer.md
ls .claude/agents/code-reviewer.md        # Task 1
ls .claude/agents/security-reviewer.md
ls .claude/agents/qa-engineer.md
ls .claude/agents/tech-writer.md
ls .claude/agents/devops-engineer.md

ls .claude/skills/deliver-feature/SKILL.md  # Task 2

ls templates/claude-feature-team/.claude/agents/code-reviewer.md  # Task 1
ls templates/claude-feature-team/.claude/skills/deliver-feature/SKILL.md  # Task 2
ls templates/claude-feature-team/features/TEMPLATE.md  # Task 7
ls templates/claude-feature-team/DOMAIN_DICTIONARY.template.md  # Task 6

# Grep checks
grep -l "DOMAIN_DICTIONARY" .claude/agents/analyst.md   # Task 4
grep -l "DOMAIN_DICTIONARY" .claude/agents/developer.md # Task 5
grep -l "code-reviewer" templates/claude-feature-team/CLAUDE.md  # Task 3
```

All commands should return file paths with no errors.
