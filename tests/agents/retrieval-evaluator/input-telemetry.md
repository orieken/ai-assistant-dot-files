# Retrieval Telemetry Report (under audit)

## Query Log (last 30 days)

| Query | Matches | KIs returned | Result |
|---|---|---|---|
| "redis rate limiting" | 2 | ki-redis-sliding-window.md, ki-redis-patterns.md | OK |
| "expand contract migration" | 0 | none | ZERO MATCH |
| "JWT validation middleware" | 1 | ki-auth-jwt.md | OK |
| "circuit breaker retry" | 0 | none | ZERO MATCH |
| "WebSocket fanout multi-instance" | 0 | none | ZERO MATCH |
| "OpenTelemetry span adapter" | 1 | ki-otel-adapter-pattern.md | OK |

## memory-registry.json metadata
- `ki-redis-sliding-window.md`: tags = [redis, rate-limiting] ✓
- `ki-auth-jwt.md`: tags = [auth, jwt] ✓
- `ki-otel-adapter-pattern.md`: tags = [otel, observability] ✓
- No KI exists with tags: [expand-contract, migration, database]
- No KI exists with tags: [circuit-breaker, resilience, retry]
- No KI exists with tags: [websocket, fanout, pub-sub]
