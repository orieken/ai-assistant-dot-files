# Clean Architecture Layer Patterns

The dependency-direction rule this whole framework enforces (`shared/rules/architecture-guardrails.md`
#1): inner layers never import from outer layers. This file expands each layer with structure and
concrete examples; the guardrail itself stays the single source of truth for the rule.

## Domain Layer (Entities)

**Context**: Pure business entities and value objects. The innermost layer — depends on nothing else in
the system.

**Structure**: No database, HTTP, or framework imports of any kind. A domain model exposes computed
properties and business logic as methods, but is never responsible for its own persistence (see
`CLAUDE.md`'s Models section: "Models are never responsible for their own persistence").

**Example**: A TypeScript `Invoice` entity can compute `isOverdue()` from its own fields, but cannot
import a TypeORM decorator, a Prisma client, or an Express request type — the concrete violation
`architecture-guardrails.md` names explicitly.

**Trade-offs**: Keeping the domain layer pure means some logic that "would be easier" with direct
database access has to go through a repository interface instead. That friction is the point — it's
what makes the domain layer testable without a database and portable if the persistence layer changes.

## Use Case Layer (Application Business Rules)

**Context**: Orchestrates domain objects to fulfill one specific application behavior — "place an
order," not "the `Order` entity."

**Structure**: Defines the interfaces it needs from the outside world (a repository, a payment gateway)
without importing any concrete implementation or framework/library
(`shared/rules/architecture-guardrails.md` #1: "UseCases cannot import Adapters or
Frameworks/Libraries"). Per this repo's own interface-placement convention, those interfaces are
defined here, in the consumer package — not in the adapter that will eventually implement them.

**Example**: `PlaceOrderUseCase` depends on an `OrderRepository` interface it defines itself; the
concrete Postgres-backed implementation lives in the adapter layer and is injected in, never imported
directly.

**Trade-offs**: Requires defining an interface for every external dependency the use case touches, even
when there's only one real implementation today. That's deliberate — see Architecture Guardrail #1's
broader point: an interface at this boundary is what makes swapping the implementation (or mocking it in
a test) a non-event.

## Adapter Layer (Interface Adapters)

**Context**: Concrete implementations of the interfaces the inner layers defined — repositories, HTTP
clients, OpenTelemetry instrumentation.

**Structure**: This is the only layer allowed to import both a framework/library and an inner-layer
interface simultaneously. It translates between the outside world's data shapes and the domain's own
model.

**Example**: A `PostgresOrderRepository implements OrderRepository` lives here, importing whatever SQL
client the project uses — something the use case layer that defined `OrderRepository` is never allowed
to do itself.

**Trade-offs**: Adapters are where the most "boilerplate-feeling" translation code lives (mapping a
database row to a domain entity, a domain entity to an HTTP response body). Resist the temptation to
skip the adapter and let a use case talk to the ORM directly — that's exactly the coupling this whole
pattern exists to prevent, and it's usually irreversible once a few use cases have taken the shortcut.

## Framework Layer (Frameworks & Drivers)

**Context**: The outermost layer — the actual framework and driver code (Express, Vue, Playwright,
whatever's wired through dependency injection at the composition root).

**Structure**: Wires concrete adapters into use cases at application startup. This is the only place
that should know about the full object graph; everything inside it interacts through interfaces.

**Observability boundary** (`shared/rules/architecture-guardrails.md` #8): OpenTelemetry instrumentation
is only allowed in the adapter layer or an interceptor — never inside domain entities or use cases. A
domain entity that emits a trace span has taken on an outer-layer concern it shouldn't know about.

**Trade-offs**: None, really — this layer is supposed to be the "glue," and gluing things together is
its entire job. The failure mode to watch for isn't over-engineering here, it's under-engineering
elsewhere: business logic that should live in a use case creeping into a route handler because it was
faster to write it there once.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
