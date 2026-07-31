# Eval Rubric: release-manager / input-git-log.md

- **Major version bump is determined correctly**: the `BREAKING CHANGE` footer in the OAuth2 commit triggers a MAJOR version bump per semver — the agent determines `v2.0.0`, not `v1.8.0`. If it chooses v1.8.0, that's a FAIL.
- **Migration guide is called out in the release notes**: the BREAKING CHANGE message says "see migration guide" — the agent either links it or flags it as a required deliverable before releasing.
- **Stripe 13 major bump is flagged as a deployment risk**: `stripe` went from 12.x to 13.x — a major npm dependency bump. The agent flags this as needing a change-log review for any breaking SDK changes, not just listing it as a routine dep bump.
- **CHANGELOG is structured by type**: the generated CHANGELOG groups changes under `### Breaking Changes`, `### Features`, `### Bug Fixes`, `### Chores` — not a flat list.
- **Deployment checklist is concrete**: the release plan includes at minimum: confirm migration guide is live, confirm stripe@13 is tested in staging, confirm backup before deploy — not a generic "test before you ship."

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
