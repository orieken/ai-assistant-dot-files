# Stability Patterns (Michael Nygard, *Release It!*)

`architect.md` explicitly names Nygard as an influence ("the stability patterns of Michael Nygard —
circuit breakers, bulkheads, timeouts"), but only Circuit Breaker ever got written down, and only in its
API-testing-specific form (`sunday-framework-patterns.md`'s Resilience Primitives). This file covers the
general, not-API-specific versions, plus the other named patterns that never got documented at all.

## Circuit Breaker

**Context**: Stops calling a dependency once it's crossed a failure threshold, so a genuinely-down
dependency fails fast instead of getting hammered by every caller still retrying it. The general form of
what `sunday-framework-patterns.md` documents for the Sunday framework specifically.

**Structure**: Three states — closed (calls go through normally), open (calls fail immediately without
attempting the dependency), half-open (a trial call checks whether the dependency has recovered). This
is the same open/closed/half-open state machine `gang-of-four-patterns.md`'s State pattern entry names
this as a concrete instance of.

**Trade-offs**: A tripped circuit breaker means legitimate requests fail fast even after the dependency
has actually recovered, until the half-open trial succeeds. That's the intended trade — a short window
of unnecessary failures is far cheaper than an unbounded pile-up against a struggling dependency.

## Bulkhead

**Context**: Partitions resources (connection pools, thread pools, capacity) per dependency, so one
failing or slow dependency can't exhaust resources that other, healthy dependencies need.

**Structure**: Named for ship compartmentalization — a hull breach floods one compartment, not the whole
ship. In software: a separate connection pool per downstream service means a slow payment provider can't
starve the connection pool a healthy inventory service also depends on.

**Trade-offs**: Partitioning resources means each partition is individually smaller than one shared pool
would be, so a partition can be under-provisioned for a legitimate traffic spike even while the whole
system has spare capacity elsewhere. The alternative — one shared pool — is what turns "one dependency is
slow" into "the whole system is down."

## Timeout

**Context**: Every network call must have an explicit timeout — this is already a hard guardrail in this
repo, not just a Nygard-ism: `shared/rules/architecture-guardrails.md` #5 states it directly ("Every
network call MUST have an explicit timeout defined"), and `shared/rules/go-conventions.md` repeats it as
a Go-specific ALWAYS.

**Structure**: A call without a timeout doesn't fail — it hangs, holding whatever resource (thread,
connection, request-handling capacity) it occupies indefinitely. A slow dependency without a timeout
doesn't degrade the caller; it can eventually stop it entirely.

**Trade-offs**: Picking too aggressive a timeout produces false-positive failures on a dependency that's
just briefly slow, not actually down. Picking too generous a timeout defeats the purpose. There's no
universal right number — it has to come from the actual latency budget of the calling code, which is
exactly why `analyst.md`'s own contract requires explicit SLA/timeout thresholds to be defined per
external or long-running call, not left as an implementation detail decided later.

## Fail Fast

**Context**: Validate what you can up front and reject immediately, rather than doing partial work and
failing deep inside a call chain where the failure is expensive to unwind and hard to diagnose.

**Structure**: Input validation, schema checks, and precondition checks belong at the boundary, not
scattered through business logic — the same principle behind Sunday's mandatory Schema Validation
(`sunday-framework-patterns.md`): a malformed response fails immediately at the validation boundary
instead of producing a confusing downstream error three calls later.

## Steady State

**Context**: A system should be able to run indefinitely without manual intervention to clean up its own
accumulated garbage — logs that rotate themselves, caches that evict, connections that get reclaimed.

**Structure**: Anything that grows unbounded over the life of a running process (an in-memory cache with
no eviction policy, a log file with no rotation) is a Steady State violation waiting to become an
incident, even if it looks fine in every short-lived test run.

**Trade-offs**: Building in cleanup (rotation, eviction, expiration) is easy to skip when a feature is
first built and everything looks fine in dev. It's specifically the multi-week-uptime failure mode that
short-lived local testing can't surface — which is exactly why it's easy to miss until it's a production
incident.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
