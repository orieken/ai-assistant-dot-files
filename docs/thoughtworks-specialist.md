# Engineering Guidelines (Thoughtworks-Aligned)

## 1. Optimize for Fast Feedback
Deliver the smallest working increment as early as possible.
- Prefer incremental, testable outputs over large implementations
- Ensure every step produces runnable code

## 2. Always Produce Production-Ready Code
Code is not done until it is deployable.
- Include tests, error handling, and configuration
- Avoid prototypes or throwaway implementations

## 3. Evolve the Design Continuously
Design for change, not completeness.
- Avoid premature abstraction
- Refactor incrementally as understanding improves

## 4. Default to Simplicity
Choose the simplest solution that works.
- Minimize dependencies and complexity
- Avoid over-engineering and unnecessary generalization

## 5. Make Intent Explicit
Code should clearly communicate purpose.
- Use meaningful names
- Document assumptions and non-obvious decisions

## 6. Build Observability by Default
Systems must be operable.
- Include logging, metrics, and traceability where relevant
- Expose meaningful signals for debugging and monitoring

## 7. Treat Data as a First-Class Concern
Data design is critical.
- Define clear schemas and contracts
- Validate inputs and outputs explicitly

## 8. Design for Decoupling
Components should be easy to change or replace.
- Use clear interfaces and boundaries
- Avoid tight coupling between modules

## 9. Test Behavior, Not Implementation
Tests should validate outcomes.
- Focus on business logic and edge cases
- Avoid brittle, implementation-specific tests

## 10. Make Trade-offs Explicit
Every decision has a cost.
- Document key trade-offs and assumptions
- Prefer reversible decisions when possible



### Stories

- Thin vertical slices of user-facing functionality
- Comparable in scope size to maintain flow stability
- Independently deployable
- Deliver observable user value
- Organized around rollout events (feature flags, migrations, releases)

## Testing
Maintain a balanced test pyramid:
- Mostly unit tests
- Some integration tests
- Few end-to-end tests

Rules:
- Every story updates or adds tests.
- Prefer fast, deterministic tests.
- Avoid E2E-heavy designs.

### Domain-Driven Design

- Model around clear bounded contexts.
- Align code structure with domain language.
- Make aggregates and invariants explicit.
- Keep domain logic isolated from infrastructure concerns.
- Use Aggregates, Entities, Value Objects, Repositories, Services, Modules, Events

### Planning File Hygiene

- **Completed stories/epics**: keep only the description and task checklist (marked `[x]`). Strip all code snippets, fixture examples, and key-assertion blocks — the source of truth is the codebase.
- **Pending stories/epics**: code snippets and design specs are allowed; they serve as implementation guidance until the story is implemented.
- The planning files are navigation and status aids, not documentation mirrors of the code.

Strip all code snippets, fixture examples, and key-assertion blocks — the source of truth is the codebase.