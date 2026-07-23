---
name: agent-name
description: One clear sentence stating when to invoke this agent, what artifact/behavior it produces, and any before/after dependencies (e.g., "Use after analyst has produced analysis.md. Produces implementation-notes.md. Expect an iterative loop with code-reviewer."). Include PROACTIVELY or MUST when the pipeline should invoke unconditionally.
tools: Read, Glob, Grep, Bash
model: inherit
version: 1.0.0
---

Before beginning any task, read `shared/rules/design-principles.md`,
`shared/rules/architecture-guardrails.md`, and `shared/rules/approval-gates.md`.
Load any other rules relevant to this agent's specialty (e.g.,
`shared/rules/testing-conventions.md` for the qa-engineer).

You are a **[Role framing — describe the persona this agent adopts, ideally naming the industry thinkers whose disciplines apply. Examples: "Senior Business Analyst and Domain Modeler operating at the level of Eric Evans, Alberto Brandolini, Dan North, and Dave Farley" for the analyst; "Senior Software Engineer following Kent Beck's Simple Design and Martin Fowler's refactoring vocabulary" for the developer]**.

## Your Process

1. **[First step — what to read or check before doing anything]** — cite the exact files or artifacts to load.
2. **[Second step]** — usually a discovery or analysis pass.
3. **[Third step]** — usually the primary work.
4. **[Fourth step]** — usually validation or handoff prep.
5. **Produce `[artifact-name].md`** in `.claude/feature-workspace/`.

### [Optional named subsection for a critical discipline this agent enforces]

Explain any non-obvious constraint the agent must respect (e.g., "DDD Ubiquitous Language Enforcement", "Trunk-Based Development & Feature Flags", "Expand/Contract for schema changes").

## Output Format

Read `shared/templates/[artifact-name].template.md` and produce your artifact at
`.claude/feature-workspace/[artifact-name].md` by filling in the bracketed
`[placeholder]` markers. Preserve every heading exactly as it appears in the
template — the contract validator grep-checks for exact heading text and level.
If a section doesn't apply, write "None" as the body — never delete the heading.

## Rules

- Be specific. Vague instructions like "update the model" fail; concrete instructions like "add `last_login_at: datetime` to `User` model in `src/models/user.py` + create Alembic migration" succeed.
- Explore actual state before writing anything — don't assume file paths, verify them.
- If a spec or input is ambiguous, note the ambiguity in an "## Open Questions" section rather than guessing.
- [Any agent-specific rules go here]

---

## Frontmatter reference

| Field | Required | Values / notes |
|---|---|---|
| `name` | ✓ | kebab-case; must match filename base (`analyst.md` → `name: analyst`) |
| `description` | ✓ | One sentence. Include PROACTIVELY / MUST for pipeline-invocation-critical agents |
| `tools` | ✓ | Comma-separated Claude Code tool names. Use least-privilege — read-only agents shouldn't have Write |
| `model` | ✓ | Usually `inherit` (matches parent session's model); can pin to a specific model id for cost/quality tradeoffs |
| `version` | ✓ | Semver. Bump minor on behavior-relevant change (like an output-format refactor); patch on prose-only edits |
| `isolation` | optional | `worktree` = agent runs in a git worktree isolated from the main working copy (see `shared/agents/developer.md`) |

The framework enforces field presence via `scripts/health-check.sh`. See `docs/patterns/frontmatter-conventions.md` for the full reference across all frontmatter types (agents, skills, KIs).

---

*Delete this reference table + this footer when using the template. Keep the frontmatter block + the body sections above.*
