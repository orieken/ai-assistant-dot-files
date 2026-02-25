# Saturday Framework: Comprehensive Architecture Rules

This document codifies the core software craftsmanship principles used within the Saturday Framework ecosystem. 

All AI subagents, orchestrators, and coding assistants **MUST** read and adhere to these technical laws when interacting with or generating code for the ecosystem.

## I. Clean Architecture (The 4 Layers)

We adhere strictly to Uncle Bob's Clean Architecture. Dependencies must ALWAYS point inwards toward the Domain Layer. The inner layers must NEVER know about the outer layers.

1. **Entities (Domain Layer)**
   - Contains enterprise-wide business rules and domain models.
   - **Constraint**: Zero external dependencies. Strictly typed. No ORM decorators or framework-specific logic here.

2. **Use Cases (Application Layer)**
   - Contains application-specific business rules. Orchestrates the flow of data to and from the domain entities.
   - **Constraint**: Must define interfaces for external dependencies (e.g., `IUserRepository`) that are fulfilled by the outer layers.

3. **Interface Adapters (Controllers/Presenters/Gateways)**
   - Converts data from the format most convenient for the Use Cases (DTOs) to the format most convenient for the external agency (e.g., Web, Database).
   - **Constraint**: This is where HTTP routing handlers and DB implementations live.

4. **Frameworks & Drivers (Infrastructure/UI)**
   - External agencies like databases, UI frameworks (Vue), third-party libraries (Playwright), and web servers.
   - **Constraint**: Code here is kept to an absolute minimum. It acts solely as the glue to the adapter layer.

## II. Domain-Driven Design (DDD "Lite")

We embrace the strategic and tactical patterns of DDD.

- **Bounded Contexts**: Structure your code by domain (e.g., `src/billing`, `src/users/auth`) rather than by technical concern (`src/controllers`, `src/models`, `src/services`).
- **Ubiquitous Language**: Variables, classes, and methods MUST perfectly match the business terminology. Do not invent technical synonyms for business concepts.
- **Repositories**: Abstract all data persistence behind Repository interfaces.
- **Value Objects**: Use immutable Value Objects for concepts like Currency, Email Addresses, or Coordinates rather than primitive types.

## III. Gang of Four (GoF) Design Patterns

Leverage these proven object-oriented design patterns to solve recurring problems:

### Creational Patterns
- **Factory Method**: Use when delegating object creation to subclasses or when object creation logic becomes complex.
- **Builder**: Use to construct complex objects step-by-step (e.g., building HTTP request payloads or complex Test Fixtures).
- **Singleton**: **Use sparingly.** Only acceptable for stateless managers or logging/observability configurations.

### Structural Patterns
- **Adapter**: Explicitly required when connecting the Core Logic to external execution clients (e.g., creating a `PlaywrightHttpAdapter` that implements `IHttpAdapter`).
- **Facade**: Use to provide a simplified interface to a complex body of code (e.g., `BaseSite` acts as a facade over various `BasePage` objects).
- **Decorator**: Use to dynamically add responsibilities to objects. Heavily used for `Filters` (state-based guards in E2E automation).

### Behavioral Patterns
- **Strategy**: Use to define a family of algorithms, encapsulate them, and make them interchangeable (e.g., `ExponentialBackoffStrategy` vs. `LinearBackoffStrategy` for API resilience).
- **Observer / PubSub**: Use for decoupled event handling (e.g., emitting OpenTelemetry traces when specific domain events occur).
- **State**: Use when an object's behavior must change dynamically based on its internal state.

## IV. Micro-Rules & Code Smells

- **Return Early (Guard Clauses)**: Eliminate nested `if` statements. Handle errors and invalid edge cases at the very top of a function.
- **No Magic Numbers or Strings**: Extract all hardcoded values into named, capitalized constants at the top of the file or in a constants module.
- **Immutability First**: Prefer pure functions and functional transformations (e.g., `map`, `filter`, `reduce`) over mutating state.
- **Cyclomatic Complexity**: Strictly enforce a McCabe complexity ceiling of `< 7` per method.
- **Function Length**: No function may exceed 30 Lines of Code (LOC). If it does, refactor it using extract-method.

## V. Testing Principles (Kent Beck's TDD)

- **Red-Green-Refactor**: Passing tests are insufficient without the subsequent "Refactor" step to ensure clean code.
- **Test Behavior, Not Implementation**: Tests should describe *WHAT* the system does, not *HOW* it does it. Refactoring internal logic should never break a test unless the behavior changed.
- **The "Boy Scout" Rule**: Always leave the code cleaner than you found it. Proactively clean up minor technical debt or formatting issues in the files you touch.
