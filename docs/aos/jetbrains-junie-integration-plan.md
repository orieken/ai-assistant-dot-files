# JetBrains AI Assistant + Junie Integration Plan (Epic 44)

Drafted 2026-07-30 from live doc research (jetbrains.com/help/ai-assistant, junie.jetbrains.com/docs).

## Investigation Findings

### JetBrains AI Assistant

**Project Rules** (native scoping mechanism):
- **Path**: `.aiassistant/rules/*.md` at project root
- **Format**: Plain Markdown — no frontmatter supported. Scoping mode is set in IDE Settings > AI Assistant > Project Rules, NOT in the file itself
- **Application modes**: Always | By file patterns (e.g. `*.go`, `src/**`) | By model decision | Manually (`@rule:` / `#rule:`) | Off
- **Agent instructions**: Also reads `AGENTS.md` at repo root (same cross-tool convention already used by Gemini)
- **Per-IDE**: No — files travel with the project. Same files work in IntelliJ IDEA, WebStorm, Rider, etc.
- **Custom modes/agents/skills**: Not documented. No evidence of custom mode definitions or skill invocation.

### Junie (JetBrains Agentic Sub-Feature)

**Guideline search order** (checked in sequence):
1. `.junie/AGENTS.md` (project root — preferred modern format)
2. `AGENTS.md` (project root — framework already generates this for Gemini)
3. `.junie/guidelines.md` or `.junie/guidelines/` folder (legacy format, still supported)

**Global**: `~/.junie/AGENTS.md`
**Format**: Plain Markdown
**Note**: Since the framework already generates `AGENTS.md`, Junie already works via step 2. The new contribution is an explicit `.junie/guidelines.md` generated as a guaranteed Junie-first path, and a `~/.junie/AGENTS.md` global install option.

## Single vs. Two Platform Entries

JetBrains AI Assistant and Junie are the same IDE plugin — one product, two modes (chat vs. agentic). They share `AGENTS.md` and differ only in the AI Assistant adding `.aiassistant/rules/` for project-rules configuration. Treating them as one platform entry (`jetbrains`) with two config-path sections is the correct design.

## Tier Assignment

**jetbrains**: Tier 2 — Personas + Rules

Rationale: Rules with file-pattern scoping (similar to Copilot's Tier 2 scoped instructions), reads `AGENTS.md` for roster + global rules. No custom modes, no skill invocation, no hooks, no pipeline orchestration confirmed. Per-IDE differentiation needed: No.

## Config Path Collision

None. `.aiassistant/rules/` and `.junie/` are distinct from all existing platform paths.

## What's Generated

**`.aiassistant/rules/`** (AI Assistant project rules — 10 files):

| File | Recommended IDE mode | Content |
|---|---|---|
| `00-approval-gates.md` | Always | `shared/rules/approval-gates.md` |
| `01-design-principles.md` | Always | `shared/rules/design-principles.md` |
| `02-architecture-guardrails.md` | Always | `shared/rules/architecture-guardrails.md` |
| `03-agent-roster.md` | Always | Agent roster (from collect_agent_roster) |
| `04-testing-conventions.md` | By file patterns: `*.spec.*,*.test.*,*.feature` | `shared/rules/testing-conventions.md` |
| `05-go-conventions.md` | By file patterns: `*.go` | `shared/rules/go-conventions.md` |
| `06-typescript-conventions.md` | By file patterns: `*.ts` | `shared/rules/typescript-conventions.md` |
| `07-python-conventions.md` | By file patterns: `*.py` | `shared/rules/python-conventions.md` |
| `08-csharp-conventions.md` | By file patterns: `*.cs` | `shared/rules/csharp-conventions.md` |
| `09-java-conventions.md` | By file patterns: `*.java` | `shared/rules/java-conventions.md` |

Each file gets a leading HTML comment with the recommended IDE setting — JetBrains rules have no frontmatter standard, so this is the only portable hint.

**`.junie/guidelines.md`** (Junie legacy path — explicit install):
Same content as `AGENTS.md` (rules + craftsmanship + agent roster). Junie reads this before the root `AGENTS.md`, so it's the preferred Junie-specific path.

## What `AGENTS.md` Already Covers

The framework already generates `AGENTS.md` (labelled as "cross-tool convention, confirmed read by Gemini Antigravity"). Junie reads it at step 2 of its search order. This means Junie is already partially supported. The new `jetbrains` platform entry makes this explicit and adds the AI Assistant project-rules layer on top.

## Phase B Plan

- B1: Add `jetbrains` entry to `shared/platform-registry.json`
- B2: Add `generate_jetbrains()` to `scripts/generate-configs.sh`
- B3: Add JetBrains parity checks to `scripts/check-parity.sh`
- B4: Add `install_jetbrains()` + detection to `install.sh`
- B5: Update `README.md` supported-tools table
