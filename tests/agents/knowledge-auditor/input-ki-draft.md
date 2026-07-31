---
title: Redis Sliding Window Rate Limiting
slug: redis-sliding-window-rate-limit
domain: infrastructure
tags: [redis, rate-limiting, middleware]
created: 2026-07-15
last-referenced: 2026-07-15
---

# Redis Sliding Window Rate Limiting

Use a sorted set with `ZADD` + `ZREMRANGEBYSCORE` + `ZCARD` to implement a sliding window
counter. The key pattern is `rate:{userId}:{windowSizeMs}`.

```typescript
async function isRateLimited(userId: string): Promise<boolean> {
  const now = Date.now();
  const window = 60_000;
  const key = `rate:${userId}:${window}`;
  await redis.zadd(key, now, `${now}`);
  await redis.zremrangebyscore(key, 0, now - window);
  const count = await redis.zcard(key);
  return count > RATE_LIMIT_MAX;
}
```

See also: ADR-007-rate-limiting-strategy.
