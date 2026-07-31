# Context Manifest: rate-limit-feature (under audit)

## Scoped Files
- `src/middleware/auth.middleware.ts` — JWT validation; rate-limit hooks go here
- `src/cache/redis-client.ts` — used for sliding window counter
- `src/routes/api.router.ts` — endpoint definitions; read to identify affected routes
- `src/models/user.model.ts` — **never opened** during implementation; pinned speculatively
- `docs/adrs/ADR-007-rate-limiting-strategy.md` — confirmed; describes Redis sliding window
- `docs/adrs/ADR-003-enterprise-memory-sync.md` — **unrelated to rate limiting**; appears to be a copy-paste error
- `shared/knowledge/ki-redis-patterns.md` — KI link; file exists

## Token Budget
Claimed: 18,400 tokens
Actual (sum of file sizes): 11,200 tokens

## Prior Deliveries Referenced
- `deliver-feature/rate-limit-v1/implementation-notes.md` — prior delivery correctly referenced

## KI Links
- `ki-redis-patterns.md` — valid
- `ki-jwt-validation.md` — **file does not exist under shared/knowledge/ or .claude/knowledge/**
