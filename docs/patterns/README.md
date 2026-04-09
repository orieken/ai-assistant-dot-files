# Reusable Patterns

This directory holds documentation for design patterns and architectural patterns used across the ecosystem. Each pattern gets its own markdown file with context, rationale, and usage guidance.

---

## Categories

### Design Patterns (Gang of Four)

Standard object-oriented patterns applied where they solve a real problem. The repository enforces these through code review and complexity checks rather than applying them speculatively.

Patterns in active use include Factory, Builder, Strategy, Observer, Adapter, Decorator, and Command. See the root `CLAUDE.md` for the decision table on when each pattern applies.

### Saturday Framework Patterns (E2E Testing)

The Site-Centric architecture replaces the traditional Page Object Model. These patterns define how UI automation is structured:

- **BaseSite** -- Top-level site abstraction managing cross-application journeys.
- **BasePage** -- Page-level abstraction encapsulating page state and actions.
- **BaseElement** -- Element-level abstraction for reusable UI component interactions.
- **BaseFlow** -- Multi-step workflow orchestration across pages.
- **Filters** -- Declarative element selection and filtering.
- **SiteManager** -- Cross-application journey orchestration.
- **TabManager** -- Multi-tab browser context management.

### Sunday Framework Patterns (API Testing)

Patterns for declarative, resilient API test automation:

- **BaseApiClient** -- Abstract base for all domain-specific API clients.
- **IHttpAdapter** -- Interface hiding HTTP execution details from test logic.
- **Fluent Matchers** -- Custom assertion matchers (`toHaveStatus`, `toBeSuccessful`, `toRespondWithin`).
- **Schema Validation** -- Zod-based response validation via `validateSchema()`.
- **Resilience Primitives** -- `CircuitBreaker` and `ExponentialBackoffStrategy` for retry logic.

### Clean Architecture Layer Patterns

Patterns enforcing dependency direction and layer separation:

- **Domain Layer** -- Pure entities and value objects with no external dependencies.
- **Use Case Layer** -- Application-specific business rules orchestrating domain objects.
- **Adapter Layer** -- Implementations of interfaces defined by inner layers (repositories, HTTP clients, OTel instrumentation).
- **Framework Layer** -- Framework and driver code (Express, Vue, Playwright) wired through dependency injection.

---

## Contributing a Pattern

Each pattern document follows this structure:

1. **Name** -- Pattern name and category.
2. **Context** -- When and why this pattern is used.
3. **Structure** -- Key abstractions, interfaces, and relationships.
4. **Example** -- A concrete usage example from the codebase.
5. **Trade-offs** -- What this pattern costs and what it buys.
6. **Related Patterns** -- Links to complementary or alternative patterns.
