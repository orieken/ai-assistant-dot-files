---
name: team-topology-check
description: Cross-references a feature's Bounded Context and Context Crossings (from analysis.md or architecture-notes.md) against TEAM_TOPOLOGY.md to flag Conway's-Law-shaped smells — a stale Collaboration that should have become an X-as-a-Service contract, or a stream-aligned team bypassing a platform team's service. Also runs standalone against git history to catch ownership drift.
triggers:
  keywords: ["team-topology-check", "team topology", "conway's law", "check ownership"]
  intentPatterns: ["/team-topology-check *", "Check team ownership for *", "Does this crossing make sense", "Are we bypassing a platform team"]
standalone: true
---

## When To Use
- Whenever `analyst` or `architect` flags a Context Crossing in `analysis.md` / `architecture-notes.md` —
  `architect` invokes this automatically as part of its "Bounded Context" step (see `shared/agents/architect.md`).
- Standalone, on demand: "is this crossing actually fine, or are we quietly rebuilding a distributed monolith
  at the team level" for any pair of bounded contexts, not just an in-flight feature.
- Periodic audit: run without a specific feature to check `TEAM_TOPOLOGY.md` against actual git history for
  ownership drift.

Do NOT use for single-team, single-bounded-context features with no crossings — there's nothing to check.

## Process

### 1. Load the Registry
Read `TEAM_TOPOLOGY.md`. If it doesn't exist, or every row still has an unfilled `_(fill in)_` team name,
report that plainly and stop — there's no registry to check against yet, and fabricating team ownership
would be worse than saying so.

### 2. Identify the Crossing
From the feature's `analysis.md` ("Bounded Context" → "Context Crossings") or `architecture-notes.md`
("Bounded Context" → "Crossings" / "Integration Pattern"), or from direct user input, identify the two (or
more) bounded contexts involved.

### 3. Look Up Each Side
For each bounded context in the crossing, find its row in `TEAM_TOPOLOGY.md`: Team, Team Type, Primary
Interaction Mode.

### 4. Check for the Two Named Mismatches
- **Stale Collaboration**: both sides are Stream-aligned teams, the declared Interaction Mode is
  Collaboration, and the crossing is not a brand-new, still-being-discovered boundary (check
  `docs/features/*/analysis.md` via the same bounded-context grep `context-engineer` already uses — if this
  pair has crossed in 2+ prior deliveries, the boundary is no longer new). Flag: recommend defining a
  Consumer-Driven Contract and evolving the declared mode to X-as-a-Service.
- **Bypassed platform team**: one side is Stream-aligned and reaches directly into another Stream-aligned
  context's internals, when a Platform team's row exists for a service that should sit between them (e.g.,
  two feature teams both reimplementing something `TEAM_TOPOLOGY.md` says a platform team already owns as
  X-as-a-Service). Flag: recommend routing through the platform team's service instead.
- If neither applies, say so explicitly — most crossings are fine, don't manufacture a finding.

### 5. Optional: Ownership Drift Check (Standalone Mode Only)
When run without a specific feature, for each row in `TEAM_TOPOLOGY.md`: `git log --format='%an' -- <paths owned by this bounded context>`
over the last N commits (best-effort — this repo doesn't track "teams" as a git concept, so treat author
names as a proxy, not ground truth). Flag if a large share of recent changes to a context's files came from
outside what the registry implies, as a prompt to update `TEAM_TOPOLOGY.md` rather than as a hard failure.

## Output Format

```markdown
# Team Topology Check: [Feature Name or "Standalone Audit"]

## Crossing(s) Checked
- [Context A] ([Team], [Type]) <-> [Context B] ([Team], [Type]) — declared mode: [Collaboration/X-as-a-Service/Facilitating]

## Findings
- **[Stale Collaboration / Bypassed Platform Team / None]**: [explanation, or "This crossing looks fine — no action needed"]

## Recommendation
[Concrete next step — e.g., "Define a Pact contract for X <-> Y and update TEAM_TOPOLOGY.md's Interaction Mode to X-as-a-Service", or "No change needed"]
```

## Guardrails
- **Never** invent a team name or type not actually present in `TEAM_TOPOLOGY.md` — if a context has no row,
  say so and stop for that context.
- **Never** report a mismatch just because two teams talk to each other — Collaboration is the *correct*
  mode for a genuinely new or evolving boundary. Only flag it as stale using the 2+ prior deliveries signal
  in step 4.
- This is a judgment tool, not a lint rule — it produces a recommendation, not an automatic block. Nothing
  here should ever be wired up as a hard CI gate the way `verify-dependencies` is.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
