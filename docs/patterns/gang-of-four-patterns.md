# Gang of Four Patterns

`CLAUDE.md`'s decision table already covers *when* to reach for each pattern in one line. This file
doesn't repeat that table — it adds the structure, concrete example, and trade-off each one-liner
doesn't have room for. Applied where they solve a real problem in this codebase; never applied to look
clever (`CLAUDE.md`: "Apply Gang of Four patterns when they solve a real problem. Never apply them to
look clever").

## Factory

**When** (per `CLAUDE.md`): object creation is complex or varies by type.

**Structure**: A dedicated class or function that centralizes `new` for a domain object, so construction
logic doesn't leak across every call site. Per this repo's own convention: "Factories are the only place
`new` is called on complex domain objects outside of tests" — a `new` on a complex domain object
anywhere else is itself a signal a factory is missing.

**Example**: `UserFactory.createAdmin()`, `UserFactory.createGuest()` — named creation scenarios, not one
generic constructor with a pile of optional flags.

**Trade-offs**: Adds a class for something a constructor could technically do. Worth it the moment
construction has more than one meaningfully different scenario — the named methods document intent a
constructor's positional arguments can't.

## Builder

**When**: object has many optional parameters.

**Structure**: A fluent chain (`.withX().withY().build()`) that replaces a constructor with a long
optional-parameter list, or worse, a boolean-parameter pileup (see `CLAUDE.md`'s own flag: "Boolean
parameter on a public method" is something to flag and split into two functions — Builder is often the
alternative to reaching for one).

**Trade-offs**: More code than a constructor for objects with only 1-2 optional fields — don't reach for
this until the parameter list is actually unwieldy.

## Strategy

**When**: behavior varies by context; eliminates conditionals.

**Structure**: Extract each branch of an `if`/`switch` into its own class implementing a shared
interface, then select the implementation at runtime instead of branching on a type code every time the
behavior is needed. Directly related to this repo's own complexity-reduction ladder
(`shared/rules/design-principles.md`'s complexity section step 3-4: replace `if/else` chains with lookup
tables, replace `switch` on type with polymorphism).

**Example**: A discount calculation that branches on `customerType` in five different places is five
bugs waiting to happen when a sixth customer type is added; a `DiscountStrategy` interface with one
implementation per type means adding a type is one new class, not five edited branches.

**Trade-offs**: More files than one function with a `switch`. Pays for itself the moment the same
type-code branch would otherwise need to be duplicated in more than one place.

## Observer

**When**: event-driven decoupling between producers and consumers.

**Structure**: A producer emits an event without knowing who's listening; consumers subscribe
independently. Keeps the producer from needing a direct reference to every consumer it might ever need
to notify.

**Trade-offs**: Debugging an event chain is harder than reading a direct function call — the control flow
isn't visible at the call site. Reach for this when producers and consumers genuinely need to evolve
independently, not as a default way to connect two pieces of code that could just call each other
directly.

## Adapter

**When**: wrapping third-party or legacy dependencies.

**Structure**: A thin translation layer between an external API's shape and the interface this codebase's
inner layers actually depend on — the same role `IHttpAdapter` plays for the Sunday framework (see
`sunday-framework-patterns.md`) and the Adapter layer plays for Clean Architecture generally (see
`clean-architecture-layers.md`).

**Trade-offs**: An extra indirection for a dependency you're confident will never change. The payoff
shows up the one time it does change — the adapter is the only thing that needs to.

## Decorator

**When**: adding behavior (logging, caching, instrumentation) without touching the class.

**Structure**: Wraps an existing object with the same interface, adding behavior before/after delegating
to the wrapped instance. This is the mechanism that keeps observability code out of domain
logic — an OTel-instrumented adapter decorates the plain adapter rather than the domain entity importing
OTel directly (`shared/rules/architecture-guardrails.md` #8).

**Trade-offs**: Stacking many decorators makes the actual call chain harder to trace at a glance. Keep
decorator chains shallow; if you're nesting more than two or three, something's probably trying to be a
use case instead.

## Command

**When**: encapsulating operations; supporting undo/redo or queuing.

**Structure**: Wraps a single operation (and its parameters) as an object with a uniform `execute()` (and
optionally `undo()`), so operations can be queued, logged, retried, or reversed without the caller
needing to know what the operation actually does.

**Trade-offs**: Overkill for a one-off function call. Reach for it specifically when you need one of the
things a plain function call can't give you: queuing, undo, or a uniform audit log of "what operations
happened," not just as a wrapper for its own sake.

## Template Method

**Not in `CLAUDE.md`'s table** — added here because this repo's own Saturday framework already uses it,
unnamed. Included for recognition, not as a new recommendation to reach for.

**When**: an algorithm has a fixed overall shape, but specific steps need to vary by subclass.

**Structure**: A base class defines the skeleton of an operation and calls out to methods subclasses
override for the parts that vary; the sequence itself stays fixed in the base class.

**Example**: `BasePage`/`BaseFlow` (see `saturday-framework-patterns.md`) are Template Method in
practice — the base classes define how a page action or a multi-step flow is structured and sequenced;
each concrete `LoginPage` or `CompleteCheckoutFlow` overrides the specific steps, not the shape.

**Trade-offs**: Subclasses are constrained to the base class's fixed sequence — if a subclass genuinely
needs a different order of operations, Template Method is the wrong pattern; that's usually a sign
Strategy (a swappable whole algorithm, not just its steps) fits better.

## Composite

**Not in `CLAUDE.md`'s table** — same as Template Method, named here for recognition.

**When**: part-whole hierarchies need to be treated uniformly, whether a caller is holding one leaf
object or a whole tree of them.

**Structure**: A common interface shared by both individual objects and groups of them, so client code
doesn't need to know which one it's holding.

**Example**: `BaseElement` (see `saturday-framework-patterns.md`) is Composite in practice — a complex
`BaseElement` like a data table is built from smaller `BaseElement`s, and both the simple and composite
cases expose the same interaction interface to the `BasePage` that uses them.

## State

**Not in `CLAUDE.md`'s table** — named here for recognition; two existing mechanisms are already
concrete instances of it.

**When**: an object's behavior needs to change based on its internal state, without a pile of
conditionals checking that state everywhere the behavior is needed. `CLAUDE.md`'s own complexity ladder
names this directly (step 4: "Replace `switch` on type with polymorphism (Strategy, State)").

**Structure**: Each state is its own class implementing a shared interface; transitioning state means
swapping which implementation is active, rather than mutating a flag that every method has to branch on.

**Example**: `CircuitBreaker`'s closed/open/half-open states (see `stability-patterns.md`) are a State
machine — each state has different behavior for "should this call go through." `pipeline-state.json`'s
phase tracking (`shared/skills/deliver-feature/SKILL.md`) is the same shape applied to an entire feature
delivery run: the current phase determines what `resume-pipeline` does next.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
