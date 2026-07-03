# Platform Verification Protocols

Manual test protocols for confirming that Cursor, GitHub Copilot, and Gemini Antigravity actually load and
respect the configs this framework generates for them. This closes the one gap that can't be automated —
`scripts/check-parity.sh` and `scripts/health-check.sh` confirm the *files* are structurally correct, but
only opening the actual tool tells us whether it *reads* them as expected.

## Why this exists
Two real findings drove this:
1. Config format assumptions can go stale — GitHub Copilot's path-scoped `.github/instructions/` support
   didn't exist when this framework's Tier system was first designed.
2. Gemini Antigravity's actual config mechanism is genuinely uncertain (secondary sources, not confirmed
   against primary docs) — see `tests/platform-verification/antigravity.md`'s confidence table.

Neither gap can be closed from within this sandboxed environment — there's no Cursor/Copilot/Antigravity
installed here to test against.

## How to run
1. Pick a protocol: [cursor.md](cursor.md), [copilot.md](copilot.md), or [antigravity.md](antigravity.md).
2. Follow its setup step (Cursor/Copilot need nothing extra; Antigravity needs a fresh `install.sh --project`
   run first — see that file for why).
3. Run as many of the numbered tests as you can — partial results are useful, you don't need to finish all
   of them in one sitting.
4. Fill in that protocol's "Report back" checklist and paste it back into the conversation.

## Fixtures
`fixtures/` contains three deliberately flawed files (same pattern as `tests/agents/`'s golden-file
fixtures — intentional violations, not real framework code) used across all three protocols so results are
comparable tool-to-tool:
- `sample.go` — violates `go-backend` rules (untyped interface, SQL injection risk, no timeout, swallowed error)
- `sample.vue` — violates `vue-frontend` rules (Options API, custom CSS, business logic in component)
- `sample.test.ts` — violates `testing` rules (vague name, tests internal state, no mocking)

## After you report back
Whatever comes back gets used to:
- Check off (or leave open, with evidence) the remaining Epic 11 TODO item: "Test: verify each platform's
  AI tool acknowledges the persona/agent roster when prompted."
- Decide whether Gemini/Antigravity's hybrid approach in `shared/platform-registry.json` should be
  simplified, expanded, or left as-is.
- Fix anything that turns out to be structurally wrong, now backed by real tool behavior instead of docs.
