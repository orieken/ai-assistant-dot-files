# Domain-Driven Design

Not a speculative addition — this is already how the framework operates. `analyst.md` explicitly channels
Eric Evans (domain modeling, ubiquitous language, bounded contexts) and runs Event Storming Lite;
`architect.md` does Strategic Domain Design with bounded-context mapping; `shared/DOMAIN_DICTIONARY.md`
*is* a ubiquitous-language artifact, enforced ("MUST update `DOMAIN_DICTIONARY.md`" is a literal rule in
`analyst.md`). This file collects the building blocks that practice already leans on, so they're
documented once instead of re-derived per agent.

## Entity

**Context**: A domain object with identity that persists over time, independent of its attribute
values — two entities with identical attributes but different identity are still different entities.

**Structure**: Identity is typically an ID field; equality is based on that ID, not on structural
comparison. This is the DDD framing behind Clean Architecture's Domain Layer (see
`clean-architecture-layers.md`) — an Entity is what actually lives there.

**Example**: A `User` entity with the same name and email as another `User` is still a different `User`
if their IDs differ — unlike a Value Object, where identical attributes mean identical objects.

## Value Object

**Context**: An immutable object defined entirely by its attributes — no identity, no lifecycle.

**Structure**: Two Value Objects with the same attributes are interchangeable. This is the DDD framing
behind this framework's own immutability convention (`CLAUDE.md`'s Models section: "Prefer immutability:
`readonly` (TS) · value receiver (Go) · `frozen=True` (Python) · `record` (Java)") — that convention
exists because most well-designed Models are Value Objects, not Entities.

**Example**: A `Money` value object (`amount`, `currency`) — two `Money(10, "USD")` instances are
equivalent and interchangeable; neither has an identity independent of its value.

**Trade-offs**: Requires creating a new instance for every change instead of mutating in place — more
object churn, in exchange for eliminating an entire class of "who else is holding a reference to this and
will be surprised when it changes" bugs.

## Aggregate / Aggregate Root

**Context**: A cluster of Entities and Value Objects treated as a single unit for the purposes of data
changes, with one Entity designated as the Aggregate Root — the only member external code is allowed to
reference directly.

**Structure**: All invariants that must hold across the cluster are enforced by the root. External code
never reaches into an aggregate to modify an internal member directly; it goes through the root.

**Example**: An `Order` aggregate root owns `LineItem` entities. Code that wants to add a line item calls
`order.addLineItem(...)`, which enforces invariants (e.g. "can't add items to a shipped order") — it
never constructs a `LineItem` and pushes it into the order's collection directly.

**Trade-offs**: Drawing aggregate boundaries too large creates contention (everyone's fighting over one
lock/version); too small and you lose the ability to enforce a real invariant atomically. Getting this
boundary right is one of the harder judgment calls in DDD — get the boundary from the actual
transactional consistency requirement, not from what looks tidy in a diagram.

## Repository

**Context**: An abstraction for retrieving and persisting Aggregates, hiding storage details from the
domain and use-case layers entirely.

**Structure**: Defined as an interface in the use-case layer (per this repo's own interface-placement
convention), implemented in the adapter layer. This is the DDD-flavored name for the same abstraction
`clean-architecture-layers.md`'s Use Case Layer section already describes — `OrderRepository` there is a
Repository in the DDD sense.

## Domain Service

**Context**: Business logic that doesn't naturally belong to any single Entity or Value Object — usually
because it coordinates across more than one Aggregate.

**Structure**: Stateless, operates on Entities/Value Objects passed to it rather than owning its own
persistent state.

**Example**: A `TransferMoneyService` that debits one `Account` aggregate and credits another doesn't
belong on either `Account` individually — the operation spans both.

**Trade-offs**: Overusing Domain Services is how domain logic quietly drains out of entities into a pile
of stateless procedures — the classic Anemic Domain Model anti-pattern
(`shared/rules/design-principles.md`'s Anti-Pattern Radar names this explicitly: "Domain entities have
only getters/setters; all logic is in 'Service' classes"). Reach for a Domain Service only when logic
genuinely can't belong to one Entity, not as a default dumping ground.

## Domain Event

**Context**: Something that happened in the domain that other parts of the system may care about — this
framework already has a literal Domain Events table (`shared/DOMAIN_DICTIONARY.md` section 4):
`SpecReady`, `AnalysisComplete`, `ArchitectureApproved`, `CodeReviewApproved`, `SecurityCleared`,
`PipelineComplete`, `ArtifactsPersisted`, `ShippedToFriday`.

**Structure**: Named in past tense (something that *already happened*, not a command to do something).
Published by the aggregate/process that caused it; consumed by whatever needs to react — the same
producer/consumer decoupling the Observer pattern describes (see `gang-of-four-patterns.md`).

**Trade-offs**: See Observer's trade-offs — the same debuggability cost (control flow isn't visible at
the call site) applies to domain events generally.

## Bounded Context

**Context**: An explicit boundary within which a particular domain model is defined and internally
consistent — this repo has five, listed in `shared/DOMAIN_DICTIONARY.md` section 1: Agent Orchestration,
Craftsmanship Governance, Test Automation, Feature Delivery, Documentation Knowledge Base.

**Structure**: The same term can mean different things in different bounded contexts, and that's fine —
DDD's whole point is that a model only needs to be consistent *within* its context, not universally.

**Trade-offs**: Drawing bounded contexts too finely creates translation overhead at every boundary;
too coarsely and you're back to one tangled model trying to serve incompatible needs. `TEAM_TOPOLOGY.md`
is where this framework tracks context ownership and crossings in practice.

## Context Map

**Context**: Documents how bounded contexts relate to and integrate with each other.

**Structure**: In this repo, `TEAM_TOPOLOGY.md`'s Context Crossings and `architect.md`'s "Team Topology
Fit" check (which flags a stale Collaboration mode or a bypassed Platform team) are the operational form
of a context map — checked automatically rather than only living in a diagram nobody updates.

## Anti-Corruption Layer

**Context**: A translation layer that prevents an external or legacy model from leaking into and
corrupting the internal domain model.

**Structure**: Structurally the same as the Adapter pattern (see `gang-of-four-patterns.md`) and the
Adapter Layer in Clean Architecture — the DDD name emphasizes *why* the translation exists (protecting
model integrity across a context boundary), not just that it exists.

**Example**: Integrating a legacy invoicing system that models "customer" completely differently from
this domain's `Customer` aggregate — the anti-corruption layer translates at the boundary so the legacy
system's shape never touches the domain layer directly.

## Ubiquitous Language

**Context**: The shared vocabulary between code, documentation, and domain experts — enforced in this
repo, not just aspirational: `shared/DOMAIN_DICTIONARY.md`'s own rule is "When a new concept is
introduced (via a new agent, skill, or rule), it MUST be added here before merging," and `analyst.md`
is required to update it whenever a feature introduces a new term or uses a synonym for an existing one.

**Trade-offs**: Maintaining a living dictionary is ongoing discipline, not a one-time document. The
alternative — letting terminology drift per-file — is exactly what `shared/DOMAIN_DICTIONARY.md`'s
"Synonyms to AVOID" column exists to prevent (e.g. never `PageObject` for `BasePage`, never `IHttpAdapter`
called `HttpService`).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
