# Team Topology

Maps each Bounded Context (see `DOMAIN_DICTIONARY.md`, "Core Domains") to the team that owns it, following
Matthew Skelton and Manuel Pais's *Team Topologies*. `architect` reads this file when a feature crosses a
bounded context boundary, and the `team-topology-check` skill (`shared/skills/team-topology-check/SKILL.md`)
audits it standalone. Update this file whenever a bounded context changes owners or a new one is added to
`DOMAIN_DICTIONARY.md`.

## Team Types (Skelton & Pais's four fundamental types)
- **Stream-aligned**: owns a single, valuable stream of work end to end — usually one bounded context.
- **Platform**: provides internal services (APIs, tooling, infrastructure) that reduce the cognitive load of
  stream-aligned teams. Interacted with as X-as-a-Service, not Collaboration, once the service stabilizes.
- **Enabling**: temporarily helps a stream-aligned team close a capability gap (e.g., a security or
  performance specialist embedding for a few weeks). Should not become a permanent dependency.
- **Complicated-subsystem**: owns a part of the system that needs deep specialist knowledge (e.g., a codec,
  a pricing engine) that would overload a stream-aligned team's cognitive load if owned directly.

## Interaction Modes
- **Collaboration**: two teams work closely and communicate frequently — appropriate while a boundary is
  still being discovered, not as a permanent state. If two teams are still in Collaboration mode a year
  after a boundary was drawn, that's a Shotgun Surgery / Distributed Monolith smell (see
  `shared/rules/design-principles.md`, Anti-Pattern Radar).
- **X-as-a-Service**: one team consumes another's well-defined API/contract with minimal ongoing
  communication. This is the target end-state for most stable bounded-context crossings — see `architect`'s
  Consumer-Driven Contracts (CDC) guidance.
- **Facilitating**: one team helps another learn or adopt something, expected to end once the gap closes.

## Registry

| Bounded Context | Team | Team Type | Primary Interaction Mode | Notes |
|---|---|---|---|---|
| Agent Orchestration | _(fill in)_ | Stream-aligned | — | Owns `shared/agents/`, `shared/skills/`, pipeline orchestration |
| Craftsmanship Governance | _(fill in)_ | Platform | X-as-a-Service | Owns `shared/rules/`, fitness functions, CI gates consumed by every other context |
| Test Automation | _(fill in)_ | Complicated-subsystem | X-as-a-Service | Saturday/Sunday Framework internals — deep specialist knowledge, consumed via `BaseSite`/`BaseApiClient` as a stable API |
| Feature Delivery | _(fill in)_ | Stream-aligned | — | Owns the `deliver-feature` pipeline and `docs/features/` archive |
| Documentation Knowledge Base | _(fill in)_ | Enabling | Facilitating | Helps other contexts write KIs/ADRs/runbooks; should stay lightweight, not become a bottleneck |

Replace the `_(fill in)_` placeholders with your actual team names once you adopt this framework in a
multi-team org. The row shapes above (type + interaction mode) are a reasonable starting default for this
framework's own bounded contexts — adjust if your org's structure differs.

## What "mismatch" means in practice
`architect` and `team-topology-check` flag two situations, both under the existing Anti-Pattern Radar
umbrella (`shared/rules/design-principles.md`):
1. **Stale Collaboration**: a Context Crossing between two stream-aligned teams that has existed for a while
   but is still ad hoc (no defined contract/API) — recommend Consumer-Driven Contracts and evolving to
   X-as-a-Service.
2. **Bypassed platform team**: a stream-aligned team reaching directly into another stream-aligned team's
   context instead of going through the declared platform team's service — this is exactly the Distributed
   Monolith shape, just named at the team level instead of the code level.

This file does not replace `DOMAIN_DICTIONARY.md`'s bounded-context definitions — it adds an ownership/team
layer on top. Bounded context names here MUST match `DOMAIN_DICTIONARY.md` exactly (same drift risk as
`docs/features/_index-by-context.md` in `docs/runbooks/scaling-cross-feature-learning.md`).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
