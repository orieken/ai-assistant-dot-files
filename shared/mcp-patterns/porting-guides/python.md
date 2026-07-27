# Python Porting Guide (`shared/mcp-patterns/porting-guides/python.md`)

This guide explains how to port the Go reference MCP patterns in `shared/mcp-patterns/go/` into a Python MCP server (e.g. using FastMCP or official `mcp` library).

## 1. Interface & Data Classes

### `Retriever` Interface
```python
from dataclasses import dataclass
from typing import Protocol, Sequence, Optional

@dataclass
class Reference:
    title: str
    path: str
    summary: str
    tags: Optional[list[str]] = None
    domain: Optional[str] = None
    relevance: float = 0.0

class Retriever(Protocol):
    async def retrieve(self, query: str, tags: Optional[list[str]] = None, domain: Optional[str] = None) -> Sequence[Reference]:
        ...
```

### `KICorpusRetriever`
- Walk markdown files using `pathlib.Path.rglob("*.md")`.
- Frontmatter parsing: use `python-frontmatter` or standard line parsing.
- Lexical scoring matches Go's `scoreRelevance`:
  - Tag match: +2.0
  - Domain match: +1.5
  - Title token match: +1.0
  - Summary/body token match: +0.3

## 2. File System Walk Utilities

Replace Go's `filepath.Walk` + `walkutil.SkipUninterestingDir` with `pathlib` filtering:
```python
SKIPPED_DIRS = {"node_modules", ".git", "dist", "build", "vendor", ".venv"}

def is_skipped_dir(path: Path, root: Path) -> bool:
    if any(part in SKIPPED_DIRS for part in path.parts):
        return True
    if any(part.startswith(".") for part in path.relative_to(root).parts[:-1]):
        return True
    return False
```

## 3. SQLite FTS5 BM25 Engine (`search_docs`)

- Package: stdlib `sqlite3` (built-in, includes FTS5 support on Python 3.9+).
- Database handle: managed via `sqlite3.connect(db_path)`.
- FTS5 Virtual Table & Query:
  ```python
  cursor.execute("""
      CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(
          path UNINDEXED, title, body, tokenize = 'porter unicode61'
      )
  """)
  ```

## 4. MCP Server Registration Pattern

Using `mcp.server.fastmcp`:
```python
from mcp.server.fastmcp import FastMCP

mcp = FastMCP("project-mcp")

@mcp.tool()
async def analyze_complexity(project_path: str, max_complexity: int = 7, max_lines: int = 30) -> str:
    analyzer = ComplexityAnalyzer()
    result = await analyzer.analyze(project_path, max_complexity, max_lines)
    return result.to_json()
```

## 5. Testing Recommendations

- Testing framework: **pytest** with `pytest-asyncio` (per `python-conventions.md`).
- Fixtures: use `tmp_path` fixture for standard file system setup.
