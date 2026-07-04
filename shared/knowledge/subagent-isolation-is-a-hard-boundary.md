---
name: subagent-isolation-is-a-hard-boundary
tags: [subagents, isolation, context-engineering, pipeline]
domain: agent-orchestration
created: 2026-07-02
---

Spawning a subagent (any of `shared/agents/*.md`) is a clean-slate context, not a shared-memory call. The
subagent sees only what its own definition and the specific artifact/task handed to it contain — it does
not inherit the orchestrator's conversation history, other subagents' internal reasoning, or unrelated
workspace state (database schema discussions, CSS layouts, CLI flags from earlier in the session).

The orchestrator (`deliver-feature`) only ever consumes a subagent's final structured report (`analysis.md`,
`security-report.md`, etc.) — never its internal transcript. This is why every pipeline agent's contract
(`shared/contracts/`) matters: the report *is* the entire interface between agents, so it has to be complete
and self-contained on its own.

Practical implication: if an agent's output seems to be missing context "it should have known," the fix is
almost never "let it see more of the conversation" — it's "put that information in an upstream artifact it
actually reads" (the feature spec, `analysis.md`, `context-manifest.md`, etc.), or wire it in via
`context-manifest.md`'s Pinpoint Files.

See: `docs/runbooks/context-engineering.md`, section 1, "Isolation of Subagent Boundaries."
