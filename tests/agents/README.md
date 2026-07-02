# Agent Golden-File Tests

These are regression tests for the agent *prompts* in `shared/agents/` — not for application code. Each
subdirectory is a fixed scenario for one agent: a deliberately-flawed input, and a set of fuzzy patterns
that a correct agent output must contain.

## Why "fuzzy" and not exact-match
Agents are LLMs, not pure functions — the same prompt against the same input won't produce byte-identical
output twice. A traditional golden-file diff would fail on harmless wording changes constantly. Instead:
- **Structural check** (automatable, CI-safe): does the output contain the required section headings from
  its contract in `shared/contracts/`, and does it match a handful of scenario-specific keyword patterns
  that any competent run should surface (e.g. "sql injection" for a SQL-injection fixture)?
- **Quality check** (not automatable): is the *reasoning* actually good? That requires a human (or another
  LLM acting as judge) reading the output. This suite does not attempt that — see `docs/features/context-engineering-framework/TODO.md`
  Decision #3: structural checks are for CI, LLM-backed judgment is manual-only.

## Directory layout
```
tests/agents/
  <agent-name>/
    input-*.{md,ts}         <- the fixed fixture fed to the agent
    expected-patterns.txt   <- scenario-specific fuzzy patterns (one per line, extended regex, case-insensitive)
    actual-output.md        <- NOT checked in; you generate this locally when testing a prompt change
```

## How to run a test
Agent invocation itself is a manual step — it requires a live Claude Code session with these agents
registered (this repo's own scripts have no way to call an LLM). To test a prompt change:

1. Open a Claude Code session in this repo (or one with `shared/agents/` installed).
2. Invoke the agent against the fixture, e.g.: "Act as the `security-reviewer` agent (see
   `shared/agents/security-reviewer.md`) and review `tests/agents/security-reviewer/input-vulnerable-code.ts`."
3. Save its full markdown output to `tests/agents/security-reviewer/actual-output.md`.
4. Run `./scripts/test-agents.sh` — it checks `actual-output.md` against both `expected-patterns.txt` and
   the agent's contract in `shared/contracts/` (required section headings), and reports PASS/FAIL/SKIP.

Repeat for whichever agent(s) you changed. **Run this after editing any agent prompt** — it's the fastest
way to catch a prompt edit that accidentally dropped a required section or stopped surfacing an obvious
finding.

## Fixtures
| Agent | Fixture | Scenario |
|---|---|---|
| `security-reviewer` | `input-vulnerable-code.ts` | SQL injection, hardcoded secret, plaintext password logging, user enumeration |
| `code-reviewer` | `input-smelly-code.ts` | Mixed responsibilities, magic numbers, domain class importing infrastructure directly |
| `analyst` | `input-feature-spec.md` | Password reset feature spec with explicit non-enumeration and token-expiry constraints |
| `architect` | `input-analysis.md` | Analysis with a new cross-bounded-context call and a brand-new external dependency |
| `qa-engineer` | `input-implementation-notes.md` | Implementation notes calling out three specific edge cases QA must cover |

## Adding a new fixture
1. Create `tests/agents/<agent-name>/` if it doesn't exist.
2. Add `input-*` — the smallest input that reliably exercises the behavior you want to guard.
3. Add `expected-patterns.txt` — a handful of patterns any correct run should produce. Prefer 3-6 sharp,
   distinctive patterns over a long list of generic ones (generic patterns pass trivially and don't catch
   regressions).
4. If the agent has a contract in `shared/contracts/`, add it to the `CONTRACT_MAP` in
   `scripts/test-agents.sh` so section-heading checks run automatically.
