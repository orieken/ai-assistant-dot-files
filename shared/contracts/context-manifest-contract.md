# Contract: context-manifest.md

**Produced by**: context-engineer
**Consumed by**: analyst, architect, developer, and every other pipeline agent (all read it first, per their
own "Context To Load First" sections)

## Required Sections (exact heading text and level)
- `## 1. Scope and Boundaries`
- `## 2. Pinpoint Files (To Keep Open)`
- `## 3. Global Rules and Constraints`
- `## 4. Knowledge Items & ADRs (To Load)`
- `## 5. Prior Deliveries in This Bounded Context`
- `## 6. Prune Recommendations (To Close)`
- `## 7. Token Budget`

## Validation Rule
- `## 7. Token Budget` must contain a `**Status**` line with value `OK` or `WARNING` — never blank and
  never any other value. If `WARNING`, a `**Cut recommendations**` line must also be present (per
  `context-engineer`'s own guardrail: never report a budget without cut recommendations to back it up).
- `## 5. Prior Deliveries in This Bounded Context` must not be empty — it must contain either at least one
  feature reference or the literal phrase "No prior deliveries found in this bounded context" (per
  `context-engineer`'s own guardrail: never skip this section silently).

This is a structural check only. It does not verify the Pinpoint Files are actually the *right* files, or
that the token estimate is accurate — that's judgment, caught downstream if analyst or developer finds the
manifest missing something relevant (see `context-audit` for after-the-fact waste analysis).
