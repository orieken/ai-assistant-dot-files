<!--
Template for observability-report.md. Consumed by the sre-engineer agent.
Structure defined here; contract in shared/contracts/observability-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Observability & SRE Report: [Feature Name]

## 1. Service Level Indicators (SLIs)
*These metrics define the health of the newly added feature.*
- **Availability SLI**: [e.g., "Payment endpoint returns 2xx or 4xx (non-500) > 99.9% of the time."]
- **Latency SLI**: [e.g., "Payment API p95 latency < 800ms."]

## 2. OpenTelemetry & Tracing
- **Analysis**: [Pass/Fail]
- **Spans Added**: [List the critical spans verified in the code, e.g., "Identified explicit `payment.process` and `stripe.charge` spans."]
- **Missing Telemetry**: [Identify large blind spots in the code where tracing should be added.]

## 3. Log Quality & Cardinality
- **Status**: [Pass / Fail / Fixed]
- **Findings**: [e.g., "Refactored 4 logger statements in `payment_service.ts` from string interpolation to structured context maps with stable strings."]

## 4. PII Data Hygiene
- **Status**: [Clean / Violation detected]
- **Notes**: [e.g., "Ensured `user.email` is redacted in the auth service spans."]

## Notes for DevOps Engineer
- [Any specific alerts or dashboard panels DevOps should wire up in Grafana/Datadog based on these SLIs.]
