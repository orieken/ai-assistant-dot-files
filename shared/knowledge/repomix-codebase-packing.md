---
name: repomix-codebase-packing
tags: [context-engineering, codebase-packing, cli-tools, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

repomix (`github.com/yamadashy/repomix`) concatenates a directory or file set into a single,
structured file optimised for LLM consumption. The output includes a file tree header, each
file's content with clear delimiters, a token count, and optional XML/Markdown/plain-text
formatting. Useful when the bounded context is small enough that sending the whole thing
verbatim is cheaper than selective file pinning.

## When to reach for it

- Bounded context is a single self-contained package or module (≤ ~5,000 code lines)
- The task requires holistic understanding — a large refactoring, a design review — where
  selective pinning would miss implicit relationships
- Handing context to a subagent that needs a complete picture with no navigation overhead

## When NOT to use it

- The bounded context is large — repomix packs everything, which can blow the token budget
  faster than selective pinning
- You only need a few specific files — use selective Read + pin instead
- The repo has generated files, lock files, or binaries that inflate the pack

## Usage

```bash
# Pack a specific directory
repomix path/to/subsystem

# Pack with ignore patterns (always add these for framework projects)
repomix path/to/subsystem \
  --ignore "**/*.lock,**/node_modules,**/.claude/feature-workspace"

# Output token count only (useful for budget check before sending)
repomix path/to/subsystem --output-tokens
```

## Installation

```bash
npm install -g repomix
# or
npx repomix <path>  # no install needed
```

## Guardrails

- Always run `--output-tokens` first to check the pack size before using it as context.
  A pack over 50,000 tokens is likely too large for a productive context window.
- Use `--ignore` to exclude framework workspace directories (`.claude/feature-workspace/`),
  test fixtures, and generated files — they add noise without signal.
- repomix output is a snapshot. If files change mid-session, the pack is stale.

## See also

- `shared/skills/context-engineer/SKILL.md` — alternative to selective pinning for small,
  self-contained bounded contexts
- `shared/knowledge/tokei-token-budget.md` — use tokei on the directory first to estimate
  whether selective pinning or repomix-packing is the more token-efficient approach
