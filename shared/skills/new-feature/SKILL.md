---
name: new-feature
description: Guided feature creation — interviews the user, produces a structured spec, creates the docs directory, and optionally kicks off the delivery pipeline.
triggers:
  keywords: ["new-feature", "new feature", "create feature", "start feature"]
  intentPatterns: ["Create a new feature for *", "I want to build *", "Start a new feature *", "/new-feature *"]
standalone: true
---

## When To Use
When the user wants to start a new feature from scratch. This skill creates the spec file, sets up the documentation directory, and bridges into the delivery pipeline.

Do NOT use when the user already has a feature spec written and wants to deliver it — use `/deliver-feature` instead.
Do NOT use when the user wants to review an existing spec — use `/spec-writer [path]` in review mode.

## Context To Load First
1. `features/TEMPLATE.md` — the spec structure to follow
2. `DOMAIN_DICTIONARY.md` — ubiquitous language reference
3. `CLAUDE.md` — project constraints
4. `docs/features/README.md` — artifact persistence conventions

## Process

1. **Ask for a working title** — one short phrase that describes the feature. Derive the kebab-case name from this (e.g., "User Authentication" becomes `user-auth`).

2. **Create the feature docs directory** — `docs/features/<kebab-name>/`. This is where pipeline artifacts will be persisted after delivery.

3. **Invoke the spec-writer agent in Write Mode** — hand off to the spec-writer to interview the user one question at a time. The spec-writer produces `features/<kebab-name>.md`.

4. **Run the spec-writer's Readiness Critique** — the spec-writer automatically enters Review Mode on the draft it just wrote. If the verdict is NEEDS WORK, offer to fix gaps and re-run.

5. **Confirm the spec is READY** — once the readiness critique passes, show the user:
   ```
   Feature spec written to: features/<kebab-name>.md
   Documentation directory created: docs/features/<kebab-name>/
   Readiness: READY

   Next steps:
   - Run /deliver-feature features/<kebab-name>.md to start the delivery pipeline
   - Run /spec-writer features/<kebab-name>.md to review or revise the spec
   - Edit features/<kebab-name>.md manually and re-run /spec-writer to re-critique
   ```

6. **Offer to start delivery** — ask: "Start the delivery pipeline now?" If the user confirms, invoke `/deliver-feature features/<kebab-name>.md`.

## Output Format

### Files Created
- `features/<kebab-name>.md` — the feature spec (produced by spec-writer)
- `docs/features/<kebab-name>/` — empty directory, ready for pipeline artifacts

### Confirmation Message
```
Feature: [Title]
Spec: features/<kebab-name>.md
Docs: docs/features/<kebab-name>/
Status: READY | NEEDS WORK
```

## Guardrails
- Never skip the spec-writer interview. Every feature must have a structured spec before entering the pipeline.
- Never create a feature spec that the spec-writer's readiness critique rates as NEEDS WORK without informing the user of the gaps.
- Never start the delivery pipeline without explicit user confirmation.
- If a feature with the same kebab-name already exists in `features/`, ask the user whether to overwrite or choose a different name.

## Standalone Mode
Works entirely offline. No external services required. The spec-writer interview and critique run locally.
