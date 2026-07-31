# Context Manifest: Notification Feed Redesign (fixture — WARNING, no cuts → FAIL)

## 1. Scope and Boundaries
- **Target Component**: notifications
- **Relevant Layers**: Domain Entities, Application Use Cases, Interface Adapters, Infrastructure
- **Bounded Context**: Notifications

## 2. Pinpoint Files (To Keep Open)
- [notification.entity.ts](file:///src/domain/notification.entity.ts) -- Domain -- Core entity
- [notification.repository.ts](file:///src/repositories/notification.repository.ts) -- Interface -- Repository interface
- [notification.repository.impl.ts](file:///src/infrastructure/notification.repository.impl.ts) -- Infrastructure -- Postgres implementation
- [ws-handler.ts](file:///src/notifications/ws-handler.ts) -- Interface -- WebSocket push handler
- [notification.service.ts](file:///src/application/notification.service.ts) -- Application -- Orchestrates fanout
- [user-preferences.service.ts](file:///src/application/user-preferences.service.ts) -- Application -- Reads notification settings
- [event-bus.ts](file:///src/infrastructure/event-bus.ts) -- Infrastructure -- EventEmitter backing fanout
- [notification.controller.ts](file:///src/api/notification.controller.ts) -- Interface Adapter -- REST endpoints
- [schema.sql](file:///migrations/20260715_notification_feed.sql) -- Migration -- Adds the notifications table
- [notification.factory.ts](file:///src/factories/notification.factory.ts) -- Factory -- Test fixture builder

## 3. Global Rules and Constraints
- [CLAUDE.md](file:///CLAUDE.md)
- [architecture-guardrails.md](file:///shared/rules/architecture-guardrails.md)
- [design-principles.md](file:///shared/rules/design-principles.md)
- [testing-conventions.md](file:///shared/rules/testing-conventions.md)

## 4. Knowledge Items & ADRs (To Load)
- [ki-websocket-fanout.md](file:///shared/knowledge/ki-websocket-fanout.md) -- Single-instance EventEmitter risk
- [ADR-012-notification-architecture.md](file:///docs/adrs/ADR-012-notification-architecture.md) -- Defines chosen fanout approach

## 5. Prior Deliveries in This Bounded Context
- [email-notifications](docs/features/email-notifications/) -- 2026-03-14 -- no retrospective.md exists for this one

## 6. Prune Recommendations (To Close)
None

## 7. Token Budget
- **Estimated total tokens for pinned files**: ~138,000
- **Target agent tier**: Developer: ≤80% of a 200k-token context window (160,000 tokens)
- **Status**: WARNING (exceeds 80% developer tier budget)
- **Cut recommendations (if WARNING)**: None
