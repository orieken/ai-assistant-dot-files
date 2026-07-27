# TypeScript Porting Guide (`shared/mcp-patterns/porting-guides/typescript.md`)

This guide explains how to port the Go reference MCP patterns in `shared/mcp-patterns/go/` into a TypeScript Node.js/Bun MCP server (e.g. using `@modelcontextprotocol/sdk`).

## 1. Interface & Class Mapping

### `Retriever` Interface
```typescript
export interface Reference {
  title: string;
  path: string;
  summary: string;
  tags?: string[];
  domain?: string;
  relevance: number;
}

export interface Retriever {
  retrieve(query: string, tags?: string[], domain?: string): Promise<Reference[]>;
}
```

### `KICorpusRetriever`
- Walk markdown files using `fast-glob` or `node:fs/promises` (`readdir({ recursive: true })`).
- Frontmatter parsing: use `gray-matter` or hand-scan lines starting between `---` delimiters.
- Lexical scoring matches Go's `scoreRelevance`:
  - Tag match: +2.0
  - Domain match: +1.5
  - Title token match: +1.0
  - Summary/body token match: +0.3

## 2. File System & Walk Utilities

Replace Go's `filepath.Walk` + `walkutil.SkipUninterestingDir` with Node.js fast-glob or fs/promises:
- Skip patterns: `['**/node_modules/**', '**/.git/**', '**/dist/**', '**/build/**', '**/vendor/**', '**/.venv/**']`
- Extension filtering: `.html`, `.htm`, `.vue`, `.jsx`, `.tsx`, `.svelte` for accessibility; `.go`, `.ts`, `.tsx`, `.js`, `.jsx`, `.py`, `.java`, `.cs` for ubiquitous language.

## 3. SQLite FTS5 BM25 Engine (`search_docs`)

- Package: [`better-sqlite3`](https://github.com/WiseLibs/better-sqlite3) (or `bun:sqlite` if on Bun).
- Schema:
  ```sql
  CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(
    path UNINDEXED,
    title,
    body,
    tokenize = 'porter unicode61'
  );
  ```
- Match Query:
  ```sql
  SELECT path, title,
         snippet(docs_fts, 2, '[', ']', '...', 20) AS summary,
         bm25(docs_fts, 1.0, 10.0, 1.0) AS relevance
  FROM docs_fts
  WHERE docs_fts MATCH ?
  ORDER BY relevance
  LIMIT 25;
  ```

## 4. MCP Server Registration Pattern

Using `@modelcontextprotocol/sdk`:
```typescript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { CallToolRequestSchema, ListToolsRequestSchema } from "@modelcontextprotocol/sdk/types.js";

export function registerTools(server: Server, tools: MCPTool[]) {
  server.setRequestHandler(ListToolsRequestSchema, async () => ({
    tools: tools.map((t) => ({
      name: t.name,
      description: t.description,
      inputSchema: t.inputSchema,
    })),
  }));

  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    const tool = tools.find((t) => t.name === request.params.name);
    if (!tool) throw new Error(`Unknown tool: ${request.params.name}`);
    return tool.execute(request.params.arguments);
  });
}
```

## 5. Testing Recommendations

- Testing framework: **Vitest** (per `typescript-conventions.md`).
- Mocking file systems: use `memfs` or Node.js `fs.mkdtemp` for temporary directories.
