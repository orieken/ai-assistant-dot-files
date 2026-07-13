# Sunday Framework Patterns (API Testing)

Declarative, resilient API test automation. Terms match `shared/DOMAIN_DICTIONARY.md` exactly.

## Declarative API Client Pattern

**Context**: The overarching pattern every component below implements a piece of — Sunday's answer to
"how do you write API tests that read like specifications instead of HTTP plumbing." The name comes
directly from this file's own opening line and from `API_FRAMEWORK_BLUEPRINT_PROMPT.md`'s "Declarative
Testing Style."

**Structure**: Three layers, each hiding the one below it from the test author. A domain-specific
`BaseApiClient` exposes only named business operations (`getUser(id)`, never `get(path)`) and hides an
`IHttpAdapter` behind it (Adapter, keeping the concrete HTTP library swappable). Fluent Matchers
(`toHaveStatus`, `toBeSuccessful`) and Schema Validation replace manual status-code/body assertions with
declarative ones. Resilience Primitives (`CircuitBreaker`, `ExponentialBackoffStrategy`) replace ad hoc
retry loops so failure-handling is a declared property of the client, not inline test logic.

**Why "declarative"**: A test written against this pattern reads as *what* is being verified
(`expect(response).toBeSuccessful()`, `expect(response.body).toMatchSchema(userSchema)`), never *how* —
no manual JSON parsing, no `if (status !== 200)`, no hand-rolled retry `while` loop. Every one of those
manual mechanics is exactly what a layer in this pattern exists to absorb.

**Trade-offs**: More upfront structure than calling an HTTP client directly in a test — a new API domain
needs its own `BaseApiClient` subclass before the first test can be written. Pays off the moment a
contract changes or a flaky endpoint needs a retry policy: one file changes, not every test that touches
that endpoint.

**Related**: Every entry in this file is a Declarative API Client component. Saturday's `Site-Centric
Pattern` is the equivalent umbrella pattern for E2E/UI testing — see `saturday-framework-patterns.md`.

## BaseApiClient

**Context**: The abstract base every domain-specific API client extends — the thing that keeps HTTP
mechanics out of test logic entirely.

**Structure**: Holds an `IHttpAdapter` instance and exposes domain methods (`getUser(id)`,
`createOrder(payload)`), never generic HTTP verbs (`get(path)`, `post(path, body)`) directly to test
code. A test should never see a URL string.

**Example**: `UserApiClient extends BaseApiClient` exposes `getById(id): Promise<ApiResponse<User>>` —
the path, headers, and serialization all live inside the client, not the test.

**Trade-offs**: One more class per API domain. In exchange, a URL or header contract change touches
exactly one file instead of every test that happens to call that endpoint.

**Related**: `IHttpAdapter` (the dependency every `BaseApiClient` hides behind).

## IHttpAdapter

**Context**: The interface hiding HTTP execution details from test logic — this is what makes
`BaseApiClient` swappable between a real HTTP call and a mock without touching a single test.

**Structure**: Defined in the consumer package (the Sunday framework's own core), not by whatever HTTP
library implements it — matches this repo's own interface-placement convention
(`shared/rules/design-principles.md`: "Define interfaces in the consumer package, not the provider").

**Trade-offs**: An extra layer versus calling `fetch`/`axios` directly. The payoff: switching the
underlying HTTP library is a one-file change, and unit-testing a `BaseApiClient` without a real network
call becomes trivial.

## Fluent Matchers

**Context**: Custom, chainable assertions (`toHaveStatus`, `toBeSuccessful`, `toRespondWithin`,
`toMatchSchema`, `toHaveHeader`) that read as intent instead of raw property checks.

**Structure**: Each matcher wraps one specific assertion concern and produces a failure message specific
to that concern (`expected status 200, got 404` — not a generic `assertion failed`).

**Example**: `expect(response).toHaveStatus(200)` instead of `expect(response.status).toBe(200)` — the
difference matters when it fails: the custom matcher's error message is written for this exact
assertion, not generic equality.

**Trade-offs**: More matcher code to maintain than relying on a general-purpose assertion library alone.
Worth it for the failure-message quality across a large test suite where debugging time matters more
than authoring time.

## Schema Validation

**Context**: Every response body gets validated against a Zod schema — a response that returns the
right status code but the wrong shape still fails the test.

**Structure**: `validateSchema()` wraps Zod's `.parse()`/`.safeParse()` with the Sunday framework's own
error formatting. Schemas live alongside the client or in a shared `schemas/` directory, never inlined
per-test.

**Trade-offs**: Requires writing (and maintaining) a schema per response shape, which is real ongoing
cost. The alternative — asserting on a handful of fields and hoping the rest is fine — is exactly the
kind of test that passes right up until a real consumer breaks on a field that quietly changed type.

## Resilience Primitives

**Context**: `CircuitBreaker` and `ExponentialBackoffStrategy` — the only sanctioned way to retry a
flaky or rate-limited call. `shared/rules/architecture-guardrails.md` #5 is explicit: no custom retry
loops with `for`/`while` and `sleep`.

**Structure**: `CircuitBreaker` wraps a call and stops retrying once a failure threshold trips, so a
truly-down dependency fails fast instead of hammering it. `ExponentialBackoffStrategy` spaces out
retries for calls that are flaky rather than down.

**Trade-offs**: A hand-rolled `for` loop with `sleep()` is faster to write once. It's also exactly the
pattern that turns a five-minute outage into a self-inflicted denial-of-service against your own
dependency — the framework-provided primitives exist because that failure mode has actually happened
before.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
