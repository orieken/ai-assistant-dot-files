My favorite kind of tooling bug: discovering the repo already had the fix, but nobody remembered why.

While adding native Cursor agent/skill support to `ai-assistant-dot-files`, we found `.cursor/agents` and
`.cursor/skills` already existed as symlinks to `shared/`.

They were committed on 2026-04-09.

The problem was not that the symlinks were missing.

The problem was that they were not documented, verified, or included in the parity checks. So the system
could not preserve the decision.

The fix became:

- document Cursor's mixed strategy
- update the platform registry
- wire install behavior
- add `.cursor/agents` and `.cursor/skills` to `scripts/check-parity.sh`

"It works" is not the same as "it is maintained."

Full post: https://dev.to/orieken/the-forgotten-symlink-why-it-works-is-not-the-same-as-it-is-maintained-17hh

