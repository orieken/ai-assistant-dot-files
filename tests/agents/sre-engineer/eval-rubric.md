# Eval Rubric: sre-engineer / input-implementation-notes.md

- **Span wraps the critical path**: a trace span is recommended (or required) around the signature verification + `markAsPaid` call, not just at the HTTP handler level — the span should capture the meaningful unit of work.
- **Low-cardinality log fields called out**: the agent flags that logging `orderId` and `eventType` as structured fields is correct, but warns against logging full Stripe event payloads (high cardinality, potential PII/token leakage).
- **Missing retry observability identified**: because there's no retry logic yet, the agent calls out the gap — Stripe will retry on 500s, but there is no metric or log to distinguish a first attempt from a Stripe retry. Recommends a counter or a flag.
- **SLI proposed for the webhook path**: suggests a concrete SLI (e.g., "webhook processing latency p99 < 2 s" or "5xx rate < 0.1%") not just a generic "add metrics."
- **No OTel instrumentation inside the domain layer**: the agent confirms (or flags as a violation if present) that span creation is in the adapter/interceptor layer, not inside `order-service.ts` itself, per architecture-guardrails.md §8.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
