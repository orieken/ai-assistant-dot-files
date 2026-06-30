---
name: api-contract-verify
description: Generates and verifies Pact-style consumer-driven contracts. Ensures APIs never break consumers.
triggers:
  keywords: ["contract", "consumer-driven", "pact", "API verification"]
  intentPatterns: ["Verify the API contract for *", "Generate consumer contract for *", "Does this API change break consumers?"]
standalone: true
---

## When To Use
Triggered when an API change is requested, or when writing integration tests between microservices or frontend/backend boundaries.

## Context To Load First
1. `.claude/feature-workspace/openapi-[resource].yaml` (if available)
2. Consumer code (frontend fetch clients or downstream service calls)
3. Producer code (the API endpoints)

## Process
1. Identify the Consumer's strict requirements: What exact JSON paths and types does the consumer *actually* read?
2. Formulate the Contract Expectation: Write the minimal JSON shape required for the consumer to not crash.
3. Compare the expectation against the Producer's OpenAPI spec or implementation.
4. If there is a mismatch (e.g., Producer removed a field the Consumer uses), flag it as a Breaking Contract.
5. Output the Verification Report.

## Output Format
`.claude/feature-workspace/api-contract-verification.md`

```markdown
# Consumer-Driven Contract Verification: [API Endpoint]

## Consumer
- [Name of upstream service or frontend app]

## Strict Requirements
- `response.data.id` (string)
- `response.data.status` ("ACTIVE" | "INACTIVE")
- [List minimal fields strictly required to avoid crash/error]

## Producer Verification
- [X] Field `id` exists and is string.
- [ ] Field `status` changed to boolean in producer API. ❌ CONTRACT BROKEN.

## Verdict
✅ SAFE TO DEPLOY / ❌ BREAKING CHANGE
```

## Guardrails
- Contracts are defined by what the consumer *actually uses*, not what the producer *can return*.
- If a producer removes a field the consumer never reads, it is NOT a breaking change.
