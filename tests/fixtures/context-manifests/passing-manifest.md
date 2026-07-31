# Context Manifest: Rate Limit Middleware (fixture — passing)

## 1. Scope and Boundaries
- **Target Component**: rate-limiting
- **Relevant Layers**: Application Use Cases, Interface Adapters
- **Bounded Context**: API Gateway

## 2. Pinpoint Files (To Keep Open)
- [auth.middleware.ts](file:///src/middleware/auth.middleware.ts#L1-L80) -- Interface -- JWT validation hook; rate-limit checks attach here
- [redis-client.ts](file:///src/cache/redis-client.ts) -- Infrastructure -- Sliding window counter backing store
- [api.router.ts](file:///src/routes/api.router.ts#L1-L45) -- Interface -- Affected route registrations

## 3. Global Rules and Constraints
- [CLAUDE.md](file:///CLAUDE.md)
- [architecture-guardrails.md](file:///shared/rules/architecture-guardrails.md)

## 4. Knowledge Items & ADRs (To Load)
- [ki-redis-sliding-window.md](file:///shared/knowledge/ki-redis-sliding-window.md) -- Documents the ZADD+ZREMRANGEBYSCORE pattern used here
- [ADR-007-rate-limiting-strategy.md](file:///docs/adrs/ADR-007-rate-limiting-strategy.md) -- Defines the chosen sliding-window approach and Redis key schema

## 5. Prior Deliveries in This Bounded Context
No prior deliveries found in this bounded context.

## 6. Prune Recommendations (To Close)
- [ ] [user.model.ts](file:///src/models/user.model.ts) -- No data model changes; different architecture layer

## 7. Token Budget
- **Estimated total tokens for pinned files**: ~4,200
- **Target agent tier**: Developer: ≤80% of a 200k-token context window (160,000 tokens)
- **Status**: OK
- **Cut recommendations (if WARNING)**: N/A
