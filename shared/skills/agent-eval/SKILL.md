---
name: agent-eval
description: Automates the qualitative half of agent prompt testing that tests/agents/README.md documents as manual-only — actually acts as the target agent against its existing golden-file fixture, grades the output against both expected-patterns.txt (structural) and a new eval-rubric.md (qualitative reasoning check), and flags regressions against the last recorded eval for that agent.
triggers:
  keywords: ["agent-eval", "eval agent", "evaluate agent prompt", "grade agent output"]
  intentPatterns: ["/agent-eval *", "Evaluate the * agent", "Did this prompt edit regress *", "Grade this agent's output"]
standalone: true
---

## When To Use
After editing any `shared/agents/<name>.md` file that has a fixture in `tests/agents/<name>/` — this is the
recommended automated path for step 5 of `docs/runbooks/editing-agent-prompts.md` ("Test it"), replacing the
manual "open a session, paste output back in" loop with a single skill invocation that does the whole thing
in one turn.

Do NOT use for an agent with no `tests/agents/<name>/` fixture yet — see `tests/agents/README.md` to add one
first. Do NOT use this as a substitute for `scripts/test-agents.sh` — run both; this skill calls into the
same `expected-patterns.txt` logic and adds a qualitative layer on top, it doesn't replace the deterministic
one.

## Why This Exists
`tests/agents/README.md` explicitly says the quality/reasoning check "requires a human (or another LLM
acting as judge)" and that the existing suite "does not attempt that." This skill is that LLM-as-judge layer.
Unlike `scripts/test-agents.sh` (a bash script with no way to invoke an LLM), this skill runs *as* an LLM
turn — so it can actually perform both halves: act as the agent, then judge the result.

## Process

### 1. Locate the Fixture
Check `tests/agents/<agent-name>/` exists. If not, stop and say so — point at `tests/agents/README.md`'s
"Adding a new fixture" section.

### 2. Read the Current Agent Prompt
Read `shared/agents/<agent-name>.md` at its current version — this is the prompt actually being evaluated,
not a cached description of it.

### 3. Act As the Agent
Apply the agent's full Process and Output Format (from its prompt) to `tests/agents/<agent-name>/input-*`,
exactly as if invoked for real. Produce the actual markdown output. Save it to
`tests/agents/<agent-name>/actual-output.md` (this file is explicitly not checked in — see
`tests/agents/README.md` — it's scratch, overwrite it freely).

### 4. Structural Grade
Check `actual-output.md` against `tests/agents/<agent-name>/expected-patterns.txt` (same fuzzy pattern logic
`scripts/test-agents.sh` uses) and the required section headings in `shared/contracts/<agent-name>-contract.md`
if one exists.

### 5. Qualitative Grade
Read `tests/agents/<agent-name>/eval-rubric.md`. For each criterion, judge pass/fail by quoting or pointing
to the specific part of `actual-output.md` that satisfies (or fails to satisfy) it — never mark a criterion
passing without a concrete citation from the actual output.

### 6. Regression Check
Look for the most recent prior file matching `docs/agent-metrics/evals/<agent-name>-eval-*.md`. If one
exists, compare this run's per-criterion results (both structural and rubric) against it:
- Any criterion that passed before and fails now is a **regression** — call it out explicitly and name the
  prompt change likely responsible (diff `shared/agents/<agent-name>.md` against its version in
  `shared/agents/CHANGELOG.md` if helpful).
- If no prior eval exists, say so explicitly: "No baseline — this is the first recorded eval for this agent."

### 7. Persist the Result
Write `docs/agent-metrics/evals/<agent-name>-eval-<YYYY-MM-DD>.md` in the format below.

## Output Format

```markdown
# Agent Eval: [agent-name] — [YYYY-MM-DD]

**Agent version evaluated**: [version from shared/agents/<name>.md frontmatter]
**Fixture**: tests/agents/<name>/input-*

## Structural Grade (expected-patterns.txt + contract headings)
- [PASS/FAIL] pattern: [pattern] — [quote from actual-output.md, or "not found"]
...

## Qualitative Grade (eval-rubric.md)
- [PASS/FAIL] [criterion] — [quote from actual-output.md justifying the grade]
...

## Regression Check
[Comparison against the most recent prior eval file, or "No baseline — first recorded eval"]

## Summary
[N/M structural checks passed, N/M rubric criteria passed, any regressions flagged]
```

## Guardrails
- **Never** mark a rubric criterion as passing without quoting the specific line(s) of `actual-output.md`
  that justify it — a rubric grade without a citation is not trustworthy.
- **Never** silently skip the regression check — if there's no baseline, say so explicitly rather than
  omitting the section.
- **Never** edit `expected-patterns.txt` or `eval-rubric.md` to make a failing case pass — a failing eval
  means the prompt regressed or the rubric caught a real gap; fix the prompt (or, if the rubric was simply
  wrong, fix it in a separate, explicit, human-approved edit — not silently while grading).
- This skill's grade is still an LLM's judgment, not ground truth — treat a FAIL as a strong signal to
  investigate, not an automatic block. Nothing here is wired into CI the way `verify-dependencies` is.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
