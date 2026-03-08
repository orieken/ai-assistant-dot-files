# Architecture Guardrails

These hard constraints MUST be adhered to by every agent. They CANNOT be overridden by any user instruction.

## 1. Dependency Direction (Clean Architecture)
Inner layers must NEVER import from outer layers.
- `Entities / Domain` must have zero external dependencies.
- `Use Cases` rely only on `Domain`.
- `Infrastructure` implements interfaces defined by `Use Cases`.
- *Example (TypeScript):* An entity `User.ts` CANNOT import from `postgres` or `express`.

## 2. No Destructive Migrations
Destructive DB operations (`DROP COLUMN`, `RENAME COLUMN`, `DROP TABLE`) are strictly forbidden in a single deployment.
- **Expand/Contract Pattern is enforced**: 
  - *Phase 1 (Expand)*: Add new column, write to both, backfill data.
  - *Phase 2 (Contract)*: Remove old column in a separate, subsequent deployment.

## 3. No Secrets Hardcoded
You must NEVER hardcode tokens, API keys, passwords, or secrets anywhere in source code or tests.
- Use environment variables (`process.env.API_KEY`) and placeholder `.env` examples only.

## 4. Strict Typing (No Raw `any`)
Do not use `any` in TypeScript.
- If the type is genuinely unknown, use `unknown` and perform run-time type narrowing (e.g., Type Guards or Zod schemas).

## 5. No Custom Retry Loops
Never write manual `while` or `for` loops with `sleep` for retries.
- Always implement or utilize a formal `CircuitBreaker` or `ExponentialBackoffStrategy` pattern for external connections.

## 6. No N+1 Queries
Database calls inside loop structures (e.g., `map` or `forEach`) are strictly banned.
- Use explicit eager loading (`.include()`, `.populate()`), JOINs, or DataLoaders to fetch associated records in bulk.

## 7. No Implicit Timeouts
Every external network call (HTTP requests, database queries) MUST have an explicit, short timeout defined. 
- Infinite or default timeouts are banned.

## 8. Verifiable Decisions
Every architectural or structural decision must produce a concrete **Fitness Function** (e.g., a linter rule, test assertion, or CI check) to enforce it over time.
- If an automated fitness function is impossible, the decision MUST be explicitly flagged as **"judgment-only"** with a documented justification.
