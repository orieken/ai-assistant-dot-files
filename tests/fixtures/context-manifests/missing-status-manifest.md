# Context Manifest: Avatar Upload (fixture — missing Status line → FAIL)

## 1. Scope and Boundaries
- **Target Component**: user-profile
- **Relevant Layers**: Application Use Cases, Interface Adapters, Infrastructure
- **Bounded Context**: User Profile

## 2. Pinpoint Files (To Keep Open)
- [avatar.controller.ts](file:///src/api/avatar.controller.ts) -- Interface Adapter -- Upload endpoint
- [storage.adapter.ts](file:///src/infrastructure/storage.adapter.ts) -- Infrastructure -- S3 adapter implementing StorageAdapter interface

## 3. Global Rules and Constraints
- [CLAUDE.md](file:///CLAUDE.md)

## 4. Knowledge Items & ADRs (To Load)
- [ADR-009-file-storage.md](file:///docs/adrs/ADR-009-file-storage.md) -- Why S3 was chosen over local disk

## 5. Prior Deliveries in This Bounded Context
No prior deliveries found in this bounded context.

## 6. Prune Recommendations (To Close)
- [ ] [payment.service.ts](file:///src/application/payment.service.ts) -- Unrelated bounded context

## 7. Token Budget
- **Estimated total tokens for pinned files**: ~6,800
- **Target agent tier**: Analyst/Architect: ≤60% of a 200k-token context window (120,000 tokens)
<!-- Status line intentionally omitted to simulate a malformed manifest — the validator must FAIL this, not SKIP it -->
- **Cut recommendations (if WARNING)**: N/A
