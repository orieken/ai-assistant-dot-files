# API Design Patterns

`sunday-framework-patterns.md` covers how to *test* an API well. The `openapi` skill ("Deliberate API
contract design before a line of implementation code is written") already implies real design
principles, worked out through a structured interview, but they've never been written down as a
standalone reference — every design decision the skill produces is currently rediscovered per-endpoint
rather than checked against a documented pattern.

## Resource-Oriented Design

**Context**: Endpoints model resources (nouns), not actions (verbs) — the HTTP verb carries the action,
the path carries the resource.

**Example** (from `openapi/SKILL.md`'s own worked design-decision table): `POST /sessions`, not
`POST /login` — "session is the resource" is the stated rationale. The verb `login` describes an action;
`/sessions` describes what's actually being created.

**Trade-offs**: Resource-oriented design occasionally produces a path that feels less immediately
readable than an action-named one (`/sessions` vs. `/login`) in exchange for consistency: every resource
in the API follows the same create/read/update/delete shape, instead of every endpoint inventing its own
verb-based naming convention.

## Idempotency Keys

**Context**: Every mutation must specify an idempotency strategy — a guardrail in `openapi/SKILL.md`
itself, not optional per-endpoint judgment.

**Structure**: A client-supplied key that lets a retried request (network timeout, client retry logic)
be recognized as "the same request again," so the server can return the original result instead of
performing the mutation a second time.

**Related**: This is what makes `stability-patterns.md`'s Timeout and Circuit Breaker safe to actually
use on mutating calls — a retry after a timeout is only safe to attempt automatically if the operation
being retried is idempotent. Without an idempotency key, a `CircuitBreaker`-wrapped retry on
`POST /orders` risks creating the order twice.

## Status Code Discipline

**Context**: Failure cases each need a distinct status code — from `openapi/SKILL.md`'s process step 4
("What can go wrong? List the failure cases — each needs a distinct status code") and its own explicit
guardrail: never use 200 for error responses.

**Example**: A 429 with a `Retry-After` header (from `openapi/SKILL.md`'s own design-decision table)
communicates lockout duration without revealing account state — the status code and headers carry
information the response body deliberately doesn't, avoiding the account-enumeration risk STRIDE's
Information Disclosure category flags (see `security-patterns.md`).

## Pagination by Default

**Context**: Every collection endpoint must have pagination — another explicit `openapi/SKILL.md`
guardrail, not a per-endpoint judgment call.

**Trade-offs**: Adding pagination to a collection endpoint that "will never have many rows" feels like
unneeded ceremony right up until it does have many rows, and the endpoint is already live with clients
depending on an unpaginated shape — retrofitting pagination onto a live API is a breaking change; building
it in from the start isn't.

## User Enumeration Prevention

**Context**: Never return error messages (or status codes, or timing) that reveal whether an account
exists — an explicit `openapi/SKILL.md` guardrail, and a direct instance of STRIDE's Information
Disclosure category (see `security-patterns.md`).

**Example**: A login endpoint that returns "invalid password" for a wrong password but "no account
found" for a nonexistent email leaks which emails are registered. The fix is a single generic message
("invalid credentials") for both cases — the API design decision *is* the security control here, not a
separate concern layered on afterward.

## Schema-First Contract

**Context**: The OpenAPI spec is designed and approved *before* implementation starts, not written after
the fact to describe whatever the implementation happened to produce.

**Structure**: `openapi/SKILL.md`'s process ends with "Draft the OpenAPI spec section and show it for
approval" before any code is written, and its output includes an explicit "Sunday Framework Mapping"
section — the client class extending `BaseApiClient`, the Zod schema that will validate the response
shape (see `sunday-framework-patterns.md`) — so the design phase and the test-tooling phase are
connected by construction, not reconciled after the fact.

**Trade-offs**: Designing the contract up front means occasionally discovering during implementation that
the designed shape doesn't quite fit — that's a real cost. It's cheaper than the alternative: an API
whose actual contract only becomes visible by reading the implementation, with every consumer (including
the test suite) reverse-engineering it independently and potentially disagreeing about what it actually
is.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
