# Task: Add rate limiting to the public API

## What needs to happen
Add per-user rate limiting (100 req/min) to the public API endpoints under `/api/v1/`.
Use Redis for the sliding-window counter. Return 429 with `Retry-After` header when exceeded.

## Known context
- `src/api/middleware/auth.middleware.ts` — existing auth middleware (likely the right insertion point)
- `src/cache/redis-client.ts` — existing Redis abstraction
- Previous delivery: `docs/features/auth-overhaul/` (may have relevant rate-limiting ADR notes)
- ADR-007 covers the Redis usage decision
