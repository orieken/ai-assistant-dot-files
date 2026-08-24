# Agent-Skill Name Pair Convention

When the same `name` exists in both `shared/agents/*.md` and `shared/skills/*/SKILL.md`, the pair
is **intentional** if one of the two sub-patterns below applies. A pair that fits neither is an
**accidental duplicate** and should be resolved via the `forgetting-engine` capability inventory flow.

Evaluated: 2026-08-02. Current confirmed pairs: `spec-writer`, `context-engineer`.

---

## Sub-pattern A: Delegation Wrapper

**Description**: The skill acts as a trigger-routing entry point. It delegates all substantive
work to the identically-named agent. The skill body is thin — it detects mode, says "invoke the
`<name>` agent," and documents output format by reference.

**Why this exists**: Skills have first-class keyword/intent-pattern routing and slash-command
registration. Agents have rich, stateful prompt bodies. Combining them gives the feature both
slash-command discoverability (skill) and full agent-quality reasoning (agent).

**Current example**: `spec-writer`
- Skill: detects Write Mode vs. Review Mode, then says "Invoke the `spec-writer` agent."
  Output format docs say "delegates to `shared/agents/spec-writer.md`."
- Agent: the actual interview logic, 300+ lines of Q&A scaffolding, readiness critique rubric,
  and spec output template.

**Identifying a Delegation Wrapper**: The skill's Process section contains an explicit
"Invoke the `<name>` agent" step. The skill's Output Format section defers to the agent file
rather than defining its own template.

**When to add a new one**: When an agent already implements a complete workflow and you need it
to be slash-command-accessible and auto-triggerable, but the agent's prompt body is too complex
to duplicate in a skill. Write a thin skill that routes to the agent; do not copy-paste the
agent's logic.

---

## Sub-pattern B: Parallel Implementation (Pipeline vs. Ad-hoc)

**Description**: The skill and the agent implement the same conceptual job independently.
The agent runs as an orchestrated subagent inside `deliver-feature` or a similar pipeline skill.
The skill auto-triggers in ad-hoc sessions where no pipeline is running. Both share the same name
because they are the same capability for two different invocation contexts.

**Why this exists**: Claude Code has two invocation paths: the main context (where skills run)
and subagent contexts (where agents run via the `Task` tool). An auto-trigger keyword that fires
in the main context cannot seamlessly hand off to a subagent — the skill must be capable of
running the workflow directly. The same job needs both a pipeline-safe subagent form and an
ad-hoc, main-context skill form.

**Current example**: `context-engineer`
- Agent: invoked as a subagent by `deliver-feature` before analyst/developer steps. Produces
  `context-manifest.md` with full prune logic, KI retrieval, and token budget estimation.
- Skill: auto-triggers when the user says "optimize context", "build context manifest", etc.
  in any ad-hoc session. Runs the same 6-step manifest process directly, calling `search-ki`
  and `query-memory` as sub-skills. Does NOT explicitly invoke the agent.

**Identifying a Parallel Implementation**: The skill has its own Process section with steps
that mirror the agent's workflow but run in the main context. The skill does not say "invoke the
`<name>` agent." The agent is referenced in `deliver-feature` or similar pipeline skills.

**When to add a new one**: When an existing agent implements a workflow that is also valuable
outside the pipeline (in ad-hoc sessions, triggered by keywords). Write a skill with the same
name that implements the same steps as main-context skill instructions. Keep both in sync when
the workflow changes — they are two representations of the same capability, not a canonical plus
a copy.

---

## Detecting an Accidental Pair

A pair is **accidental** (and should be resolved) when:

1. The skill's Process section neither delegates to the agent nor mirrors the agent's workflow
   with independent main-context steps.
2. The descriptions diverge significantly — one is a superset/refactoring of the other.
3. The skill and agent have overlapping keyword triggers but different output formats and no
   clear division of invocation context.

Run `/forgetting-engine capability-inventory` to surface all pairs with an ASSESSMENT
(WRAPPER, PARALLEL, or COLLISION) and get a human-reviewed resolution proposal.

---

## Rules

- **Never accidentally create a pair.** Before naming a new skill or agent, check for an existing
  capability with the same name in the other catalog.
- **If an intentional pair diverges in workflow**, update both files in the same commit.
- **Delegating skills must NOT duplicate the agent's logic.** The entire point of sub-pattern A
  is one canonical prompt body, not two.
- **Parallel implementation skills must NOT invoke the agent** (they would create a nested
  subagent from the main context, which is expensive and breaks the skill's standalone guarantee).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
