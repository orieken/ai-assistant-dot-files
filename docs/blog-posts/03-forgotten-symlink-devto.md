---
title: "The Forgotten Symlink: Why 'It Works' Is Not the Same as 'It Is Maintained'"
published: false
description: "A small Cursor integration story about drift, rediscovery, and turning accidents into checks."
tags: ai, devtools, maintenance, architecture
canonical_url:
cover_image:
---

The best bug in a tooling repo is the one where you discover someone already solved your problem.

The worst version is realizing nobody remembered.

While adding native Cursor agent and skill support to `ai-assistant-dot-files`, we found that `.cursor/agents`
and `.cursor/skills` already existed in the repo as symlinks.

They pointed to `../shared/agents` and `../shared/skills`.

That was exactly the design we needed.

They had been committed on 2026-04-09 in commit `d0b54d3`, with the message "expanded to work for all
platforms."

Then they faded out of the system's working memory.

## The context

The framework has one canonical source of truth: `shared/`.

Agents, skills, rules, contracts, Knowledge Items, the domain dictionary, team topology, and platform
registry all live there. Platform-specific files are either symlinked or generated from that canonical
layer.

That is how one repo can project the same rules into Claude Code, Cursor, Windsurf, GitHub Copilot,
Gemini/Antigravity, and OpenAI Codex.

The repo has a capability tier system because those tools do not all support the same primitives.

Claude Code has full native agent orchestration. Cursor used to be treated mainly as a rules/persona target,
with generated `.mdc` files. Then Cursor shipped native Agent Skills and subagents:

- `.cursor/skills/*/SKILL.md`
- `.cursor/agents/*.md`

That changed the design.

The right move was not to generate Cursor-specific copies of the agents and skills. The right move was to
symlink Cursor directly to `shared/`, just like Claude Code already did.

And then we found out the repo already had those symlinks.

## The missing part was not the symlink

The symlinks existed.

What did not exist was a system that made them visible, verified, and meaningful.

They were created before the capability tier system explained why they mattered. They were not part of the
parity check. They were not clearly represented in the platform registry. They were not woven into the
install flow as a first-class design choice.

So the fix was not "create symlinks."

The fix was:

- document Cursor's mixed strategy in `docs/ARCHITECTURE.md`
- update `shared/platform-registry.json`
- teach `install.sh` to symlink `.cursor/agents` and `.cursor/skills`
- update `scripts/check-parity.sh` so those symlinks cannot silently disappear
- update the README capability matrix with the real Cursor behavior

That last part matters most.

If a behavior is important but not checked, it is folklore.

Folklore decays.

## Drift is not only duplication

When people talk about configuration drift, they usually mean duplicated text getting out of sync.

That is one kind.

This was another: a correct structural decision existed, but the rest of the system did not know how to
protect it.

The repo already had a parity script because earlier versions had copied instructions into `.cursorrules`,
`copilot-instructions.md`, and `CLAUDE.md` independently. That drift was easy to see once you knew to look
for it.

The forgotten symlink was quieter.

It did not fail loudly. It just stopped influencing future design.

## "It works" is a weak invariant

There is a tempting maintenance posture that says: if the file exists and the app loads it, we are done.

This story made the opposite case.

For cross-tool AI configuration, "it works" is not enough. The stronger invariant is:

- the source of truth is named
- each projection mechanism is documented
- platform capability differences are explicit
- parity checks fail when the projection drifts
- install behavior recreates the intended shape

That turns a one-off fix into a maintained feature.

## The useful lesson

The lesson is not "use symlinks."

Sometimes copying is right. Sometimes generation is right. Sometimes a platform cannot follow references,
so inlining is the only honest option. Cursor rules still need generated `.mdc` files with inlined content,
even though Cursor agents and skills can symlink directly to `shared/`.

The lesson is: make the maintenance contract match the platform's actual capabilities.

For Cursor, the result is now mixed:

- agents and skills: symlink to `shared/`
- rules: generate fully inlined `.mdc` files

That mixed strategy is less elegant than pretending everything works the same way.

It is also truer.

And in tooling, true ages better than elegant.

## Source trail

- `docs/features/context-engineering-framework/TODO.md` — Epic 30 Cursor native skills/agents parity and the 2026-04-09 symlink rediscovery.
- `docs/ARCHITECTURE.md` — Cursor mixed strategy and capability tier update.
- `shared/platform-registry.json` — Cursor platform notes and install strategy.
- `scripts/check-parity.sh` — explicit `.cursor/agents` and `.cursor/skills` symlink checks.
- `README.md` — updated platform capability matrix.

## Image prompt

Create a polished editorial hero image for a technical blog post titled "The Forgotten Symlink." Show an
archaeological dig inside a modern code repository: a small glowing chain-link/symlink artifact is being
uncovered from layers labeled April, May, June, July, while clean verification beams connect it to a
`shared/` source of truth and a Cursor folder. Style: cinematic but minimal vector illustration, dark navy
background, teal and amber highlights, no real logos, no readable code except tiny labels `shared/` and
`.cursor/`, wide 16:9 composition.

