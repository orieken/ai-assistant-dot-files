# Java Porting Guide (`shared/mcp-patterns/porting-guides/java.md`)

This guide explains how to port the Go reference MCP patterns in `shared/mcp-patterns/go/` into a Java MCP server (e.g. using Spring AI MCP SDK or Java MCP SDK).

## 1. Interface & Record Mapping

### `Retriever` Interface
```java
public record Reference(
    String title,
    String path,
    String summary,
    List<String> tags,
    String domain,
    double relevance
) {}

public interface Retriever {
    List<Reference> retrieve(String query, List<String> tags, String domain);
}
```

### `KICorpusRetriever`
- Walk markdown files using `java.nio.file.Files.walk(Path start)`.
- Frontmatter parsing: use Jackson YAML or line scanner for frontmatter metadata.
- Lexical scoring matches Go's `scoreRelevance`:
  - Tag match: +2.0
  - Domain match: +1.5
  - Title token match: +1.0
  - Summary/body token match: +0.3

## 2. File System Walk Utilities

Replace Go's `filepath.Walk` + `walkutil.SkipUninterestingDir` with `java.nio.file.Files.walkFileTree`:
- Skip directories: `node_modules`, `.git`, `dist`, `build`, `vendor`, `.venv`, `.target`.

## 3. SQLite FTS5 BM25 Engine (`search_docs`)

- Package: `org.xerial:sqlite-jdbc`.
- FTS5 Table creation & Query matching:
  ```sql
  CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(
    path UNINDEXED, title, body, tokenize = 'porter unicode61'
  );
  ```

## 4. Testing Recommendations

- Testing framework: **JUnit 5** + **Mockito** (per `java-conventions.md`).
- File System testing: use `@TempDir Path tempDir`.
