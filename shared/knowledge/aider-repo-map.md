---
name: aider-repo-map
tags: [context-engineering, codebase-overview, symbol-graph, cli-tools, optional-dependency]
domain: context-engineering
created: 2026-08-21
---

Aider's repo-map (`github.com/paul-gauthier/aider`) uses tree-sitter to build a compact
symbol graph of a codebase — every class, function, and their call relationships — in a few
hundred tokens. It is significantly more signal-dense than a file tree and far less verbose
than full file contents. The algorithm is open-source and can be extracted as a standalone
script without using aider itself as a coding assistant.

## When to reach for it

- The bounded context is not yet known at the start of a context-engineer run — the task is
  described functionally ("fix the payment flow") and the relevant files are not obvious
- Large monorepos where reading the file tree gives no structural intuition
- Pre-flight for refactor-engineer campaigns where understanding call relationships matters

## Usage

```bash
# Generate a repo-map for the current directory
aider --map-tokens 1024 --no-stream --message "" --yes 2>/dev/null

# Or use the standalone script (no aider session, just the map)
python -c "
from aider.repomap import RepoMap
rm = RepoMap(root='.')
files = rm.get_repo_map([], list_all_files=True)
print(files)
"
```

The standalone script approach avoids launching an interactive aider session and is suitable
for pipeline use in context-engineer.

## Installation

```bash
pip install aider-chat
# tree-sitter language grammars are downloaded on first use
```

Or use `uv` (if the project uses uv):
```bash
uv tool install aider-chat
```

## Guardrails

- The repo-map is a heuristic — it surfaces symbols it considers "important" based on usage
  frequency, not correctness. Validate its output against actual file reads before relying on
  it for scoping decisions.
- aider-chat is a heavy install (~200MB with tree-sitter grammars). If tokei + ast-grep already
  give sufficient scope signal, skip this tool.
- The repo-map algorithm changes between aider versions — don't treat its output as stable
  across upgrades.
- Extracting the standalone script requires importing from `aider.repomap` — a private internal
  module with no stability guarantee. Wrap the import in a try/except and fall back gracefully.

## See also

- `shared/skills/context-engineer/SKILL.md` — Step 1 (Identify Target Component Scope) when
  the bounded context is initially unknown
- `shared/knowledge/ast-grep-structural-search.md` — more targeted alternative once you know
  the relevant interface or method name
