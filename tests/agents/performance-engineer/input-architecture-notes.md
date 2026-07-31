# Architecture Notes: Real-Time Notification Feed

## Decision
Implement a per-user notification feed using long-polling (`GET /notifications/poll`).
Each poll waits up to 30 seconds for a new notification, then returns 200 (with payload) or 204 (timeout).

## Data Model
Table `notifications`:
- `id UUID PRIMARY KEY`
- `user_id UUID REFERENCES users(id)`
- `payload JSONB`
- `created_at TIMESTAMPTZ DEFAULT now()`
- `read_at TIMESTAMPTZ`

## Query Pattern
`SELECT * FROM notifications WHERE user_id = $1 AND read_at IS NULL ORDER BY created_at DESC`

No index defined yet beyond the primary key.

## Concurrency Model
Each open poll holds a database connection for up to 30 seconds. No connection pool sizing adjustment planned yet.

## Scale Estimate
10,000 active users → 10,000 concurrent open connections. Connection pool default: 10.
