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

## V. Testing Principles (Kent Beck's TDD, Dan North's BDD, Dave Farley)

- **Red-Green-Refactor**: Passing tests are insufficient without the subsequent "Refactor" step to ensure clean code. The refactor step happens *before* moving to the next feature.
- **Test Behavior, Not Implementation**: Tests should describe *WHAT* the system does, not *HOW* it does it. Refactoring internal logic should never break a test unless the behavior changed.
- **BDD (Dan North)**: Scenarios must describe behavior observable from *outside* the system. "Given the CartService has been initialized" is wrong. "Given the cart has 3 items" is right. Given/When/Then is a communication tool first, a test structure second.
- **Acceptance Tests vs Unit Tests (Dave Farley)**: Acceptance tests verify that the system does what the business needs (written in terms of *what*, never *how*). Unit tests provide design feedback and verify internal contracts. Test isolation is non-negotiable.
- **The "Boy Scout" Rule**: Always leave the code cleaner than you found it. Proactively clean up minor technical debt or formatting issues in the files you touch.

## VI. Simple Design Rules (Kent Beck)

Apply these four rules, strictly in this priority order, whenever you write code:

1. **Passes the Tests**: The system works. If it doesn't work, nothing else matters.
   - *Example (TypeScript)*: Write the failing `expect(calculator.add(2, 2)).toBe(4)` BEFORE you write the implementation.
2. **Reveals Intention**: The code communicates its purpose clearly to the next reader.
   - *Example (Go)*: Rename `x := Calculate(u, items)` to `invoiceTotal := CalculateUserInvoice(user, cartItems)`.
3. **No Duplication**: The code handles a specific concept in exactly one place. Do not Repeat Yourself (DRY).
   - *Example (Python)*: Extract identical data validation logic from `create_user` and `update_user` into a shared `_validate_user_data()` function.
4. **Fewest Elements**: Minimize the number of classes, methods, and variables. Remove abstractions that aren't earning their keep.
   - *Example (TypeScript)*: Don't create an `IUserProcessorStrategyFactory` if a single `processUser(user)` function is sufficient for the current requirements.

## VII. Refactoring Catalog (Martin Fowler)

Refactoring is a named, structured operation, not a vague "cleanup." Use these 10 most applicable operations:

1. **Extract Function**
   - *Smell*: A function is too long (> 30 lines) or requires comments to explain its blocks.
   - *Before/After (TypeScript)*:
     ```typescript
     // Before
     function printOwing(invoice: Invoice) {
       let outstanding = 0;
       console.log("***********************");
       console.log("**** Customer Owes ****");
       console.log("***********************");
       // ... calculate outstanding ...
       console.log(`name: ${invoice.customer}`);
       console.log(`amount: ${outstanding}`);
     }
     
     // After
     function printOwing(invoice: Invoice) {
       const outstanding = calculateOutstanding(invoice);
       printBanner();
       printDetails(invoice.customer, outstanding);
     }
     ```
2. **Inline Function**
   - *Smell*: A function's body is as clear as its name, making the abstraction an unnecessary indirection.
3. **Extract Variable**
   - *Smell*: An expression is hard to read or understand.
   - *After*: Extract `const isMacOs = platform.toUpperCase().indexOf("MAC") > -1;`
4. **Rename Variable / Extract Variable (Intention-Revealing Name)**
   - *Smell*: A variable name doesn't explain its purpose. `const a = 5;` -> `const retryCount = 5;`
5. **Replace Conditionals with Polymorphism**
   - *Smell*: A switch statement or long if/else chain checks the type or properties of an object. (See Strategy pattern).
6. **Introduce Parameter Object**
   - *Smell*: A group of parameters naturally go together and appear in multiple function signatures.
   - *After*: Replace `(startDate, endDate)` with a `DateRange` object.
7. **Decompose Conditional**
   - *Smell*: A complex conditional statement hides the business rule.
   - *After*: Replace `if (date.isBefore(SUMMER_START) || date.isAfter(SUMMER_END))` with `if (isNotSummer(date))`
8. **Replace Magic Number/String with Symbolic Constant**
   - *Smell*: A literal value is used with unexplained meaning.
9. **Separate Query from Modifier**
   - *Smell*: A function returns a value BUT also has side effects (mutates state).
   - *After*: Split into two functions: one that mutates state (returns void) and one that queries state.
10. **Move Function / Field**
    - *Smell*: A function or field is used more often by another class/module than the one it currently resides in (Feature Envy).

## VIII. Fitness Functions (Neal Ford)

Evolutionary architecture requires continuous validation. **Fitness Functions** are objective functions that guide development by enforcing measurable constraints on the architecture.

Every significant architectural decision should produce a fitness function to ensure the constraint isn't violated over time.

Examples relevant to the Saturday/Sunday ecosystem:
1. **Cyclomatic Complexity Limits**: `eslint complexity rule (max 6)` fails the build if code is too complex to maintain.
2. **Dependency Direction**: `eslint-plugin-boundaries` ensures that the Domain layer (Entities) never imports from the Infrastructure layer.
3. **Test Coverage Thresholds**: `nyc` or `jest --coverage` fails the build if domain logic drops below 85% test coverage.
4. **Performance Thresholds**: Playwright tests that fail if critical UI workflows (like login or checkout) take longer than 2000ms.
5. **Security Baselines**: `npm audit` or `tfsec` fails the pipeline if high-severity vulnerabilities or misconfigurations are introduced.

## IX. Legacy Code Protocol (Michael Feathers)

When modifying untested or "legacy" code, you must follow this protocol:

1. **Identify the Seam**: Find a place where you can alter behavior without editing in that place (e.g., extracting a method, injecting a dependency).
2. **Write Characterization Tests**: BEFORE changing any code, write tests that assert the *current actual behavior* of the system (even if it's flawed). This sets the baseline.
3. **Introduce Testability**: Use the seam to get the code under test safely.
4. **Refactor and Add Behavior SEPARATELY**: **NEVER** refactor code and add new behavior in the same commit. Refactor first to make the change easy, commit, then add the new feature.
