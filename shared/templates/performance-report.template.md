---
feature: "<feature-name>"
bounded_context: "<owning-bounded-context>"
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

<!--
Template for performance-report.md. Consumed by the performance-engineer agent.
Structure defined here; contract in shared/contracts/performance-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Performance & Reliability Report: [Feature Name]

## 1. Idempotency Guarantees
- **Status**: [Pass / Fail / N/A]
- **Notes**: [e.g., "The POST /checkout endpoint must require an Idempotency-Key header and check a distributed cache before processing."]

## 2. Timeout & Circuit Breaker Mandates
- **Status**: [Pass / Fail]
- **Mandates**:
  - [e.g., "The call to the Stripe API MUST be wrapped in a CircuitBreaker with a hard 3000ms timeout."]
  - [e.g., "Database queries inside `UserRepository` MUST explicitly pass a 1000ms Context timeout."]

## 3. N+1 Query Prevention
- **Status**: [Pass / Fail]
- **Findings**: [Identify any loops or iterative accesses in the proposed architecture that will cause N+1 database queries. Demand the use of a DataLoader or explicit SQL JOINs.]

## 4. Hot Path Caching
- **Analysis**: [Identify slow aspects of this feature's structure]
- **Strategy**: [e.g., "The user settings payload should be cached in Redis with a 5-minute TTL since it is read on every page load but rarely updated."]

## Notes for Developer
- [Actionable, specific instructions the Developer must follow to fulfill these constraints while writing the code.]
