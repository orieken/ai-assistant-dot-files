# Architecture Guardrails

**HARD CONSTRAINTS. THESE CANNOT BE OVERRIDDEN BY ANY AGENT OR USER INSTRUCTION.**

## 1. Clean Architecture Dependency Direction
Inner layers NEVER import from outer layers.
- `Entities` (Domain) cannot import `UseCases` or `Adapters`.
- `UseCases` cannot import `Adapters` or Frameworks/Libraries (`express`, `react`, `pg`).
*Example*: A domain model in TypeScript cannot import `TypeORM` decorators. It must remain pure.

## 2. No Destructive Migrations
The Expand/Contract pattern is non-negotiable.
- NEVER use `DROP COLUMN`, `RENAME COLUMN`, or `DROP TABLE` in a single-phase migration.
- NEVER add a `NOT NULL` column without a `DEFAULT` value.

## 3. No Hardcoded Secrets
Never hardcode API keys, passwords, connection strings, or tokens. Use `.env` placeholders mapped to secure vaults.

## 4. Strict Typing
- No raw `any` types allowed in TypeScript. 
- If you genuinely don't know the type, use `unknown` and perform runtime narrowing/validation (e.g., Zod).

## 5. Failure & Reliability
- No custom retry loops with `for` or `while` and `sleep`.
- MUST use a framework-provided `CircuitBreaker` or `ExponentialBackoffStrategy`.
- Every network call MUST have an explicit timeout defined.

## 6. Performance Guarantees
- No N+1 Queries: Eager loading (`.populate`, `.include`, or DataLoaders) is required.
- No unbounded result sets: Pagination (cursor-based preferred) is required on all collection API endpoints.

## 7. Verifiable Architecture
- Every structural or architectural decision made must produce a fitness function (a CI check, linter rule, or automated test).
- If it cannot produce a fitness function, it MUST be explicitly flagged as "judgment-only" with a documented reason in the architecture notes.

## 8. Observability Boundaries
- No OpenTelemetry (OTel) instrumentation logic is allowed inside domain entities or page logic.
- Traces and spans must only be emitted from the adapter layer or interceptor layer.
