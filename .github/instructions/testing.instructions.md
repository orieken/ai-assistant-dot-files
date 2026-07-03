---
applyTo: "**/*.spec.*,**/*.test.*,**/*.feature"
---
# Testing Rules

## Saturday Framework (E2E / UI Testing)
ALWAYS use the Site-Centric pattern: `BaseSite`, `BasePage`, `BaseElement`, `BaseFlow`.
NEVER use traditional Page Object Model (POM).
ALWAYS use Playwright driven by Cucumber.js for UI automation.
ALWAYS include OpenTelemetry instrumentation for every BDD scenario.

## Sunday Framework (API Testing)
ALWAYS use Vitest for unit tests and Playwright for integration/E2E API tests.
ALWAYS use the custom `api` fixture and fluent matchers (`toHaveStatus`, `toBeSuccessful`, `toRespondWithin`).
ALWAYS extend `BaseApiClient` for domain-specific API clients.
ALWAYS validate schemas with Zod (`validateSchema()`).
NEVER use custom retry loops — use `CircuitBreaker` or `ExponentialBackoffStrategy`.

## Test Quality
CRITICAL: Test coverage MUST be >= 85%.
CRITICAL: Cyclomatic complexity per function MUST be < 7.
ALWAYS practice TDD/BDD — Red-Green-Refactor.
NEVER write feature code without tests.
