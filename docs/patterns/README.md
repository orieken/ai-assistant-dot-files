# Reusable Patterns

This directory documents the design and architectural patterns used across the ecosystem — the
structure, a concrete example, and the trade-offs each pattern brings, not just its name.

---

## Available Pattern Docs

- [gang-of-four-patterns.md](gang-of-four-patterns.md) — Factory, Builder, Strategy, Observer, Adapter,
  Decorator, Command. `CLAUDE.md` has the one-line decision table for *when* to reach for each; this
  file has the structure/example/trade-off detail that table doesn't have room for.
- [saturday-framework-patterns.md](saturday-framework-patterns.md) — the Site-Centric E2E architecture:
  `BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`, `Filters`, `SiteManager`, `TabManager`.
- [sunday-framework-patterns.md](sunday-framework-patterns.md) — declarative API testing:
  `BaseApiClient`, `IHttpAdapter`, Fluent Matchers, Schema Validation, Resilience Primitives.
- [clean-architecture-layers.md](clean-architecture-layers.md) — Domain, Use Case, Adapter, and
  Framework layers, expanding on the dependency-direction rule in
  `shared/rules/architecture-guardrails.md`.

All terminology matches `shared/DOMAIN_DICTIONARY.md` exactly — that's the canonical definition source;
these docs go deeper into context, structure, and trade-offs, not redefine the terms.

---

## Contributing a Pattern

Each pattern entry follows this structure:

1. **Context** — when and why this pattern is used.
2. **Structure** — key abstractions, interfaces, and relationships.
3. **Example** — a concrete usage example, ideally from this codebase.
4. **Trade-offs** — what this pattern costs and what it buys.
5. **Related** — links to complementary or alternative patterns.

Group related patterns into one file by category (as the four docs above do) rather than one file per
pattern — most individual patterns are a few paragraphs, and splitting each into its own file adds
navigation overhead without adding content.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
