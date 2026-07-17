# Sunday Framework Patterns (API Testing)

Declarative, resilient API test automation. Terms match `shared/DOMAIN_DICTIONARY.md` exactly.

## Mental Model

Sunday organizes API testing into three primary concepts and three supporting ones. Read top-down:

- **API Clients** — the top-level abstraction. One `BaseApiClient` subclass per API domain
  (`UserApiClient`, `OrderApiClient`); each exposes named business operations, never generic HTTP
  verbs. Tests never see URL strings.
- **Assertions** — how tests express what they expect:
  - **Fluent Matchers** (`toHaveStatus`, `toBeSuccessful`, `toRespondWithin`, `toMatchSchema`) —
    declarative status/timing/shape checks with purpose-built failure messages.
  - **Schema Validation** — Zod schemas validate response bodies; a right-status/wrong-shape response
    still fails the test.
- **Resilience** — `CircuitBreaker` and `ExponentialBackoffStrategy` for failure behaviors, replacing
  ad hoc retry loops. `architecture-guardrails.md` #5 forbids `for`/`while`+`sleep` retry code
  outright — this layer is the only sanctioned alternative.

Supporting:

- **HTTP Adapter** — `IHttpAdapter` sits behind every `BaseApiClient`, hiding the concrete HTTP
  library so tests never depend on it directly.
- **Test Data** — `Model`s (Zod-typed request/response shapes) are the domain objects tests reason
  about; `Factory`s produce valid instances of them without tests hardcoding UUIDs and emails inline.
- **Test Doubles** — `Mock Server` (Playwright `page.route()`, MSW, WireMock, or an equivalent
  provider stub) replaces real upstream services when the test needs to exercise a failure path or
  isolate from real infrastructure.

The rest of this file is one entry per concept, in the order above.

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

## Model

**Context**: A typed representation of what crosses the API boundary — an `OrderCreateRequest`, an
`OrderResponse`. Sunday-specific: **Zod schemas serve as both the Model and the validator**. One
`OrderSchema` produces the TypeScript type (`z.infer<typeof OrderSchema>`) used to construct request
payloads AND the runtime validator that Schema Validation calls on responses. That double duty is
Sunday's central bet on Zod — one source of truth for what an Order looks like across every test.

**Structure**: One Zod schema per API domain type. Schemas live alongside the `BaseApiClient` they
serve, or in a shared `schemas/` directory if multiple clients reference the same type. Never inlined
per test — a schema inlined in a `.test.ts` file will silently drift from the schema used elsewhere.

**Example**: `OrderSchema = z.object({ id: z.string().uuid(), totalCents: z.number().int().positive(),
lineItems: z.array(LineItemSchema) })`. A test creating an order derives its payload type from
`z.infer<typeof OrderCreateRequestSchema>`. A test asserting on a response calls
`expect(response.body).toMatchSchema(OrderSchema)`. Both use the same schema file.

**Trade-offs**: A schema per domain type is real ongoing cost — schemas need to change when the API
changes. In exchange, one file becomes the definition of "what an Order looks like" across every test,
every client, and every mock (see `Mock Server` below — the same Zod schema can validate the mock's
outbound shape too, catching mocks that drift from the real contract).

**Related**: `Schema Validation` (the assertion mechanism that uses the schema), `Factory` (produces
valid instances of a Model), the `api-contract-verify` skill (Pact-style consumer-driven contracts,
the counterbalance for when Zod schemas locally pass but drift from what a real consumer expects).

## Factory

**Context**: Produces valid `Model` instances with realistic defaults, using this framework's existing
per-language factory-library convention — `fishery` + `@faker-js/faker` for TypeScript (see
`shared/rules/typescript-conventions.md`). Sunday-specific: often produces request payloads for POST
and PUT tests, guaranteeing the payload conforms to its Zod schema without every test hardcoding UUIDs
and email strings inline.

**Structure**: One Factory per Model, exposing named creation scenarios per this framework's
`CLAUDE.md` Factories convention (`createOrder`, `createOrderWithLineItems`, `createExpiredOrder`).
Static methods, not instantiation.

**Example**: `OrderFactory.create({ status: 'shipped' })` returns a valid `Order` with a fresh UUID,
a realistic `totalCents`, one line item by default, and the specified status. A test asserting on
shipped-order behavior doesn't hand-write any of that.

**Trade-offs**: One more class per Model. Pays off exactly when Models do — one file changes when the
schema changes; every test using the Factory picks up the new default automatically.

**Related**: `Model` (what a Factory produces), `BaseApiClient` (typically consumes Factory-produced
Models as request payloads).

## API Test Coverage Matrix

**Context**: A shared baseline of scenarios every API endpoint under test should have coverage for,
regardless of language or specific test framework. This is Sunday's answer to "how do you know when
you've written enough API tests" — not a mandate, a baseline. Missing one shouldn't fail a build
automatically, but the absence should be a conscious decision, not an oversight.

**Structure**: Per HTTP method, the recommended scenario set:

| Method   | Recommended baseline scenarios                                     |
|----------|--------------------------------------------------------------------|
| GET      | `happy_path`, `not_found`, `network_error`                         |
| GET list | `happy_path`, `network_error` (not_found doesn't apply to lists)   |
| POST     | `happy_path`, `server_error`, `validation_error`                   |
| PUT      | `happy_path`, `not_found`, `server_error`                          |
| PATCH    | `happy_path`, `not_found`, `server_error`                          |
| DELETE   | `happy_path`, `not_found`                                          |

Add `auth_error` for every endpoint that requires authentication (bearer, basic, apikey).
`timeout_error` is always worth considering when the endpoint has an SLA to enforce; not on the
baseline because not every endpoint needs it.

The default expected status codes and Sunday resilience-primitive tie-ins:

| Scenario           | status | typical trigger                         |
|--------------------|--------|-----------------------------------------|
| `happy_path`       | 200    | (POST: 201)                             |
| `not_found`        | 404    | resource missing                        |
| `server_error`     | 500    | upstream fault                          |
| `auth_error`       | 401    | missing/invalid credential              |
| `validation_error` | 422    | malformed payload                       |
| `network_error`    | —      | transport failure — exercises retries   |
| `timeout_error`    | —      | slow response — exercises `CircuitBreaker` |

**Example**: A `POST /orders` endpoint that requires auth should have at least four scenarios:
`happy_path` (201, order created), `server_error` (500 propagates), `validation_error` (422 on bad
payload), `auth_error` (401 without a valid token). If the client wraps calls in
`ExponentialBackoffStrategy`, add `network_error` to prove the strategy actually retries on transient
failures rather than silently swallowing them.

**Trade-offs**: The matrix defines a baseline, not a mandate. 100% matrix coverage isn't the goal —
conscious skipping is legitimate (a `GET /health` endpoint that can never 404 doesn't need
`not_found`). What matters is that the absence is a decision. Skipping because you forgot is exactly
what advisor skills like `sunday-test-advisor` (go-sunday YAML audit) exist to catch.

**Related**: `sunday-test-advisor` (go-sunday YAML audit against this matrix), `api-test-generator`
(auto-generates baseline tests from an OpenAPI spec — implicitly follows a similar matrix), `Model` /
`Factory` (produce the request payloads each scenario needs), `Mock Server` (route interception or
provider stub exercises the failure-side scenarios without needing a broken real upstream), `Resilience
Primitives` (`network_error` and `timeout_error` scenarios are how you verify these actually work).

## Mock Server (Test Doubles)

**Context**: When a Sunday test needs to isolate from a real upstream service — either because the
real service isn't reachable in this test environment, or because the test needs to exercise a
specific upstream response the real service doesn't reliably produce (429 rate limits, 500 errors,
malformed bodies, slow responses that trip a `CircuitBreaker`). "Mock Server" is the umbrella; the
concrete choice varies (Playwright `page.route()`, MSW, WireMock, custom local HTTP server).
Terminology follows Martin Fowler's Test Doubles taxonomy — most Sunday uses are technically *stubs*
(canned responses), not *mocks* (behavior verification), but "mock server" is the widely-used
industry term and this doc uses it that way.

**Structure**: Two shapes matter:
- **Route interception** — Playwright's `page.route()` or MSW intercepts the network call at the
  boundary. The client never sees a different URL; the interception replaces what would have gone over
  the wire. Best when you control the client and want to short-circuit specific requests inline with
  the test.
- **Provider stub** — a small local HTTP server (WireMock, `msw` in Node mode, a custom Express
  server) sitting at a real URL the client is pointed at. Best when the client is instantiated from a
  configured base URL and you need to fake the whole upstream service surface — including endpoints
  the test itself doesn't directly hit but a Resilience primitive might retry against.

**Example**: Testing that `OrderApiClient` correctly handles a 429 rate-limit response through
`ExponentialBackoffStrategy`. Route interception: `await page.route('**/orders', route =>
route.fulfill({ status: 429, headers: { 'Retry-After': '2' } }))` — colocated with the test, no
external server needed. Provider stub: WireMock configured to return 429 for the first request, then
200 for the second — better when the retry logic needs multiple round trips against a real socket.

**Trade-offs — validate the mock against the real contract, or you lock in fantasy**: Mocks decouple
tests from network flakiness and unlock failure-path testing, but they codify assumptions about the
upstream that can drift from reality. A mock returning `{ orderId: "abc" }` when the real API returns
`{ order_id: "abc" }` will pass forever until production breaks. Two mitigations, in decreasing
weight:
1. Validate the mock's outbound shape against the same Zod schema Schema Validation uses on real
   responses — the schema catches drift at the mock boundary too.
2. Run periodic real-upstream contract tests via the `api-contract-verify` skill (Pact-style
   consumer-driven contracts) — the counterbalance for when local schemas pass but the real API has
   evolved.

**Related**: `BaseApiClient` (the thing being isolated), `Resilience Primitives` (mocks are how you
exercise the failure paths Resilience protects against without waiting for real production incidents),
`Schema Validation` and `Model` (the schema validates mock outbound shape, not just real response
shape), the `api-contract-verify` skill (the counterbalance against mock/reality drift).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
