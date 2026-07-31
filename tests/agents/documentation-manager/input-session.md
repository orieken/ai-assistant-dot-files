# Session Summary: WebSocket notification feed (ad-hoc, never went through deliver-feature)

## What happened
Engineer added a WebSocket notification feed directly to `src/notifications/ws-handler.ts` without an analysis or architecture phase. Changes were made over two hours in a single unreviewed session.

## Key decisions made
1. Used raw `ws` package instead of Socket.IO — "Socket.IO overhead not justified for this use case"
2. Chose in-process event emitter (`EventEmitter`) over Redis pub/sub for fanout — "only one server instance in staging, will revisit for multi-instance"
3. Did not add OpenTelemetry spans to the WebSocket handler — "will add later"
4. Added a hardcoded 30-second heartbeat interval in the handler directly

## Files changed
- `src/notifications/ws-handler.ts` (new file, 240 lines)
- `src/server.ts` (attached WebSocket upgrade handler)
- `package.json` (added `ws` dependency)

## What was NOT done
- No ADR written for the ws vs Socket.IO decision
- No KI written for the EventEmitter fanout approach
- OTel deferred (violates architecture-guardrails.md §8)
- Heartbeat interval hardcoded (violates architecture-guardrails.md §3 spirit — not a secret but a magic number)

## Notes
The EventEmitter approach will silently break when a second server instance is added. This is a known future risk that was not documented.
