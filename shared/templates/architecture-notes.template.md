<!--
Template for architecture-notes.md. Consumed by the architect agent.
Structure defined here; contract in shared/contracts/architecture-contract.md validates
that these headings survive intact. Preserve every heading exactly.
-->

# Architecture Notes: [Feature Name]

## Structural Decisions

### [Decision 1 Title]
**Decision**: [What was decided]
**Reversibility**: [Cheap / Moderate / Expensive / Essentially Permanent]
**Fitness Function**: [Property that must remain true. If a decision cannot produce a fitness function, you MUST justify why and flag it as "judgment-only"]
**Enforcement**: [Exact eslint rule / dependency-cruiser config / tsc flag / CI command]
**Responsibility**: devops-engineer wires it, architect defines it

### [Decision 2 Title]
...

## Component Placement

| Component | Package | Layer | Extends/Implements |
|---|---|---|---|
| `NewLoginFlow` | `apps/ye-olde-magic-shop` | Application | `BaseFlow` |
| `RateLimitFilter` | `@orieken/saturday-core` | Core | `Filter` decorator |

## Bounded Context
- Context:
- Crossings:
- Integration Pattern:
- Team Topology Fit: [Owning team + type from TEAM_TOPOLOGY.md, and whether the Integration Pattern matches
  the declared Interaction Mode — or "No TEAM_TOPOLOGY.md row for this context" / "No crossing, not applicable"]

## Stability Design
| Dependency | Timeout | Circuit Breaker | Bulkhead | Fail Fast | Idempotency |
|---|---|---|---|---|---|

## Observability Design
- Spans:
- Metrics:
- Logs:
- Alerts:

## Layer Boundary Checks
- [ ] No domain logic in adapter layer
- [ ] No framework imports in use case layer
- [ ] New components follow dependency direction (inward only)
- [ ] No direct HTTP calls outside `IHttpAdapter` implementations

## Anti-Pattern Check
- [ ] Checked for Distributed Monolith
- [ ] Checked for Anemic Domain Model
- [ ] Checked for God Object
- [ ] Checked for Shotgun Surgery
- [ ] Checked for Leaky Abstraction
- [ ] Checked for Premature Generalization
- [ ] Checked Team Topology Fit (stale Collaboration or bypassed Platform team, if this feature crosses a bounded context)

## Fitness Functions
These properties must remain true as the codebase evolves.
Add these to CI if not already enforced:

- **[Fitness function name]**: [What it measures] — [How to enforce: eslint rule / dependency-cruiser / tsc flag]

## Refactoring Opportunities
Adjacent code that should be cleaned up while we're in this area:
- [File/class]: [Smell identified] → [Named refactoring operation to apply]
— or "None identified"

## Developer Handoff Notes
- [Specific guidance for the developer that isn't obvious from the analysis]
- [Import paths to use for new components]
- [Patterns in the codebase to model after — with file paths]

## Open Architectural Questions
- [Ambiguity that needs human input before the developer proceeds]
— or "None"
