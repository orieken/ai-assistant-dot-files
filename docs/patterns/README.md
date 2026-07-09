# Reusable Patterns

This directory documents the design and architectural patterns used across the ecosystem — the
structure, a concrete example, and the trade-offs each pattern brings, not just its name.

---

## Available Pattern Docs

- [gang-of-four-patterns.md](gang-of-four-patterns.md) — the 7 patterns in `CLAUDE.md`'s decision table
  (Factory, Builder, Strategy, Observer, Adapter, Decorator, Command) with the structure/example/
  trade-off detail that table doesn't have room for, plus 3 more (Template Method, Composite, State)
  named for recognition because this codebase's own mechanisms already use them, unnamed.
- [saturday-framework-patterns.md](saturday-framework-patterns.md) — the Site-Centric E2E architecture:
  `BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`, `Filters`, `SiteManager`, `TabManager`.
- [sunday-framework-patterns.md](sunday-framework-patterns.md) — declarative API testing:
  `BaseApiClient`, `IHttpAdapter`, Fluent Matchers, Schema Validation, Resilience Primitives.
- [clean-architecture-layers.md](clean-architecture-layers.md) — Domain, Use Case, Adapter, and
  Framework layers, expanding on the dependency-direction rule in
  `shared/rules/architecture-guardrails.md`.
- [domain-driven-design.md](domain-driven-design.md) — Entity, Value Object, Aggregate, Repository,
  Domain Service, Domain Event, Bounded Context, Context Map, Anti-Corruption Layer, Ubiquitous
  Language. Already actively practiced by `analyst.md`/`architect.md`/`DOMAIN_DICTIONARY.md`, just
  never consolidated as a reference until now.
- [enterprise-integration-patterns.md](enterprise-integration-patterns.md) — Message Channel,
  Content-Based Router, Message Translator, Publish-Subscribe Channel, Dead Letter Channel, Correlation
  Identifier, Saga. `architect.md` names Hohpe as an influence; several of this repo's own mechanisms
  (`pipeline-trace.json`, `deliver-feature`'s checkpointed pipeline) are concrete instances.
- [stability-patterns.md](stability-patterns.md) — Circuit Breaker, Bulkhead, Timeout, Fail Fast, Steady
  State (Michael Nygard, *Release It!*). `architect.md` names Nygard as an influence; the
  general/non-API-specific versions of these, complementing Sunday's API-specific Resilience Primitives.
- [twelve-factor-app.md](twelve-factor-app.md) — all 12 factors (Codebase through Admin Processes). The
  one axis nothing else here covers: how a service should actually run in production, not how its code
  is organized.
- [security-patterns.md](security-patterns.md) — STRIDE threat modeling, Secure by Default, Defense in
  Depth, Least Privilege, Paved Road / Golden Path. `security-reviewer.md` names Adam Shostack as an
  influence and already has a fully worked STRIDE table in its own contract.
- [observability-patterns.md](observability-patterns.md) — SLI, SLO, Error Budget, Low-Cardinality
  Logging, Structured Tracing, No PII in Telemetry. `sre-engineer.md`'s entire mandate; SLO/Error Budget
  are the one layer above the already-required SLI that wasn't named yet.
- [expand-contract-migrations.md](expand-contract-migrations.md) — the full pattern plus its own named
  approval gate, upgraded from "a rule that's referenced" to a documented pattern with the same
  Context/Structure/Example/Trade-offs treatment Circuit Breaker got in `stability-patterns.md`.
- [api-design-patterns.md](api-design-patterns.md) — Resource-Oriented Design, Idempotency Keys, Status
  Code Discipline, Pagination by Default, User Enumeration Prevention, Schema-First Contract. The
  `openapi` skill already enforces all of these per-endpoint; never written down as a standalone
  reference until now.

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
