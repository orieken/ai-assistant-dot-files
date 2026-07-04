---
name: context-engineer-must-be-wired-into-pipeline
tags: [context-engineering, deliver-feature, pipeline, orchestration]
domain: framework-internals
created: 2026-07-02
---

The `context-engineer` agent existed for a long time as a fully-specified pre-flight optimizer (context
manifest, KI/ADR lookup, pruning checklist) but `deliver-feature/SKILL.md` never actually invoked it — it
was a dead capability. An agent or skill being well-specified in isolation is not the same as it being wired
into the pipeline that's supposed to use it.

When adding a new agent to this framework, the definition alone is not enough: grep the orchestrating
skill (`deliver-feature/SKILL.md`) to confirm it's actually invoked at the right phase, and confirm the
consuming agents (here: `analyst`, `developer`) have an explicit step to read that agent's output. Otherwise
the capability silently rots.

Fixed by: adding context-engineer as step 6 of `deliver-feature` Phase 1 (before analyst), and adding a
manifest-check step to `analyst.md` and `developer.md`'s process before their "explore the codebase" step.
(Exact step/phase numbers drift as the pipeline grows — check `deliver-feature/SKILL.md` directly for the
current numbers rather than trusting this note long-term; the lesson, not the number, is what matters.)
