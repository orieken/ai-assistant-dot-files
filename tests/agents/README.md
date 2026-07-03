# Agent Golden-File Tests

These are regression tests for the agent *prompts* in `shared/agents/` — not for application code. Each
subdirectory is a fixed scenario for one agent: a deliberately-flawed input, and a set of fuzzy patterns
that a correct agent output must contain. Each subdirectory, taken together with its `eval-rubric.md`, is an
**Agent Eval Case** (see `DOMAIN_DICTIONARY.md`).

## Why "fuzzy" and not exact-match
Agents are LLMs, not pure functions — the same prompt against the same input won't produce byte-identical
output twice. A traditional golden-file diff would fail on harmless wording changes constantly. Instead:
- **Structural check** (automatable, CI-safe): does the output contain the required section headings from
  its contract in `shared/contracts/`, and does it match a handful of scenario-specific keyword patterns
  that any competent run should surface (e.g. "sql injection" for a SQL-injection fixture)?
- **Quality check** (LLM-as-judge, not CI-safe): is the *reasoning* actually good — correct severity, a
  specific fix rather than a generic one, no false positives? This used to be manual-only (a human reading
  the output). `eval-rubric.md`, run via the `agent-eval` skill (`shared/skills/agent-eval/SKILL.md`),
  automates it: the skill acts as the agent, then grades its own output against the rubric, citing the
  specific line that justifies each grade. Still not CI-safe — it's an LLM's judgment, not a deterministic
  assertion — but no longer requires a human to run it by hand.

## Directory layout
```
tests/agents/
  <agent-name>/
    input-*.{md,ts}         <- the fixed fixture fed to the agent
    expected-patterns.txt   <- scenario-specific fuzzy patterns (one per line, extended regex, case-insensitive)
    eval-rubric.md          <- qualitative reasoning criteria, graded by the agent-eval skill
    actual-output.md        <- NOT checked in; you generate this locally when testing a prompt change
```

## How to run a test

**Automated (recommended)**: invoke the `agent-eval` skill against the agent name — it acts as the agent
against its fixture, saves `actual-output.md`, and grades both `expected-patterns.txt` and `eval-rubric.md`
in one pass, flagging any regression against the last recorded eval in `docs/agent-metrics/evals/`.

**Manual (structural check only, no rubric)**:
1. Open a Claude Code session in this repo (or one with `shared/agents/` installed).
2. Invoke the agent against the fixture, e.g.: "Act as the `security-reviewer` agent (see
   `shared/agents/security-reviewer.md`) and review `tests/agents/security-reviewer/input-vulnerable-code.ts`."
3. Save its full markdown output to `tests/agents/security-reviewer/actual-output.md`.
4. Run `./scripts/test-agents.sh` — it checks `actual-output.md` against both `expected-patterns.txt` and
   the agent's contract in `shared/contracts/` (required section headings), and reports PASS/FAIL/SKIP.

Repeat for whichever agent(s) you changed. **Run this after editing any agent prompt** — it's the fastest
way to catch a prompt edit that accidentally dropped a required section, stopped surfacing an obvious
finding, or regressed the quality of its reasoning.

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
4. Add `eval-rubric.md` — 4-6 qualitative criteria about the *reasoning*, not just keyword presence (e.g.
   "correct severity assigned," "fix is specific, not generic," "no false positive on the safe code path").
   See any existing `eval-rubric.md` for the format.
5. If the agent has a contract in `shared/contracts/`, add it to the `CONTRACT_MAP` in
   `scripts/test-agents.sh` so section-heading checks run automatically.
