# Token Footprint Audit

Date: 2026-08-07
Repository: `/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`

## Overview

The dominant token-usage issue in this repository is not source-code verbosity. It is the large, near-identical instruction payload duplicated across multiple AI client entry points, plus repeated rule trees that increase both repository size and model context load.[cite:2][cite:3]

## Highest-impact reductions

| Priority | Current issue | Concrete refactor | Likely token reduction |
|---|---|---|---:|
| 1 | Six root instruction files are each about 76 KB and roughly 1,050 lines | Replace them with one canonical `AGENTS.md` and tiny platform adapters that point to it | 75–85% of always-loaded instructions |
| 2 | Root prompts include 95–96 headings and 27 tables each | Convert them into a concise runtime policy; move reference material to on-demand skills/docs | 60–75% per active session |
| 3 | The same rules exist in `shared/`, `.roo/rules/`, `.clinerules/`, and templates | Treat `shared/rules/` as the single source of truth and generate target formats at install time | Major repo reduction and less drift |
| 4 | Large agent roster is loaded with general instructions | Keep a short routing table in default context and store full role contracts in agent definitions | 50–70% of roster tokens |
| 5 | Repeated imperative phrasing such as always, never, and must | Express invariants once and use compact checklists | Lower repetition and fewer contradictions |

The six primary entry-point files total roughly 456 KB before any project-specific context: `AGENTS.md`, `.openai.md`, `.cursorrules`, `.windsurfrules`, `.github/copilot-instructions.md`, and `.junie/guidelines.md` all contain approximately the same 1,047–1,054 lines of content.[cite:3]

## Canonical prompt architecture

Make `AGENTS.md` the single canonical instruction source, then generate intentionally minimal client-specific files.

```text
shared/
  policy/
    core.md
    routing.md
    quality-gates.md
  rules/
    typescript-conventions.md
    python-conventions.md
    architecture-guardrails.md
    ...

AGENTS.md
.cursor/rules/*.mdc
.openai.md
.cursorrules
.windsurfrules
.github/copilot-instructions.md
