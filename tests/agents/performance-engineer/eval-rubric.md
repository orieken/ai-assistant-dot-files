# Eval Rubric: performance-engineer / input-architecture-notes.md

- **Missing index is flagged with specific column(s)**: the agent identifies that `notifications(user_id, created_at)` (or a partial index on `read_at IS NULL`) is missing — not a generic "add an index," but the specific composite key that serves the query pattern.
- **Connection pool exhaustion is quantified**: the agent computes (or estimates) the math — 10,000 concurrent polls vs. a pool of 10 = severe contention — and gives a concrete recommendation (raise pool size to match concurrency budget, or switch to a push model).
- **Long-polling architecture is challenged**: the agent explicitly recommends WebSockets or Server-Sent Events as a more connection-efficient alternative to long-polling, with a brief rationale tied to the 10,000-user estimate.
- **Unbounded result set flagged**: the `SELECT *` query with no `LIMIT` is flagged — even for unread notifications, a user could accumulate thousands. Recommends pagination or a capped response.
- **Timeout is end-to-end**: the 30-second poll timeout must be shorter than any upstream proxy/load-balancer timeout or the connection will be torn down by the infra before the app can respond cleanly. Agent either notes this or recommends an end-to-end timeout audit.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
