# Contract: observability-report.md

**Produced by**: sre-engineer
**Consumed by**: devops-engineer, orchestrator (deliver-feature)

## Required Sections (exact heading text and level)
- `## 1. Service Level Indicators (SLIs)`
- `## 2. OpenTelemetry & Tracing`
- `## 3. Log Quality & Cardinality`
- `## 4. PII Data Hygiene`
- `## Notes for DevOps Engineer`

## Validation Rule
`validate-artifact` checks presence of every heading above, plus:
- `## 4. PII Data Hygiene` status must be `Clean` or explicitly resolved — per the agent's own rule that PII violations must be fixed directly with `Edit`, not left as a recommendation. A status of `Violation detected` with no accompanying fix note is a FAIL.

This is a structural check only. It does not verify OTel spans actually exist in production — that's an operational concern outside this pipeline's scope.
