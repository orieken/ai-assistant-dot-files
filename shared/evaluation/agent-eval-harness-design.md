# Prompt-Regression Eval Harness — Design (Epic 61 Phase A)

Architecture decisions for `scripts/run-agent-evals.sh`. These rulings commit the implementation
shape; Phase B implements against them. A human must approve this document before Phase B begins.

---

## 1. Runner tech: `claude -p`

**Ruling: `claude -p` subshell per agent. No Agent SDK.**

Rationale: the Agent SDK adds a Python or Node.js dependency the framework currently has no other
reason to require. `claude -p` is already available wherever Claude Code is installed — the
population of users running this harness is exactly the population with `claude` in their PATH.

Fidelity assessment: agents declare `tools:` in frontmatter (Read, Write, Edit, etc.). In a real
interactive session, Claude Code enforces those tool grants. In a `claude -p` invocation, tool
grants follow the CLI's current permission model, not the agent's frontmatter. This is acceptable
for prompt regression testing because:

- Every current fixture (32 agents) is a **read-only evaluation**: the agent reads a provided
  input file and produces a markdown report. It does not need to actually execute writes, run
  bash, or call MCP tools to produce correct output — the fixture provides all the inputs it
  needs.
- The `agent-eval` skill already validates this assumption in practice: it "acts as the agent"
  without real tool grants and produces gradeable output. The harness automates what that skill
  does, using the same mechanism.
- If a future fixture genuinely requires tool execution to produce correct output, that fixture
  must document the dependency in its `eval-rubric.md` and the harness emits a SKIP for it.

The invocation pattern:

```bash
claude -p --model <resolved-model> \
  "Read shared/agents/<agent>.md in full. Act as that agent. Apply its complete Process and
   Output Format to the content of tests/agents/<agent>/<input-file>. Produce the full
   markdown output only — no preamble, no 'I will now...' framing." \
  > tests/agents/<agent>/actual-output.md
```

Model tier is resolved from `shared/model-defaults.yaml` at `claude_code.<tier>` using
`scripts/resolve-model-tier.py` (already exists). A `light` agent gets Haiku; `default` gets
`inherit` (whatever model the user has configured, typically Sonnet); `heavy` gets Opus.

---

## 2. Grading tech: two-pass (generate → judge)

**Ruling: generate with the agent's configured model, grade rubric with a second lighter judge.**

The two passes are independent `claude -p` invocations:

**Pass 1 — generation** (uses agent's model tier, as resolved above):
Produces `tests/agents/<agent>/actual-output.md`.

**Pass 2 — rubric judge** (uses `light` tier — Haiku):
```bash
claude -p --model claude-haiku-4-5-20251001 \
  "Read tests/agents/<agent>/eval-rubric.md and tests/agents/<agent>/actual-output.md.
   For EACH criterion in eval-rubric.md, output exactly: PASS <criterion-label> | <one-line quote>
   or FAIL <criterion-label> | <one-line explanation of what is missing>. Output nothing else." \
  > /tmp/rubric-grade-<agent>.txt
```

Structured output from the judge (one line per criterion) is parseable by bash without another
LLM call. The judge prompt is deliberately minimalist to reduce noise.

**Pattern checks** (pass 0 — deterministic, always first):
Same `grep -iE` logic as `scripts/test-agents.sh`. Run before either LLM pass. If the output
file is absent, emit SKIP rather than failing. Pattern checks never cost API tokens.

Why two-pass over one combined pass: separating generation from grading means the generation
step faithfully uses the agent's model tier (preserving what the harness is actually testing —
does THIS model produce correct output for THIS agent prompt?), while grading uses a cheap
judge. Combining them into a single prompt would force either the generation or the grading to
run on the wrong model.

---

## 3. Cost + cadence policy

**Full-sweep cost estimate (32 agents, Sonnet default + Haiku judge):**

| Component | Input tokens | Output tokens | Cost |
|---|---|---|---|
| Generation (32 × Sonnet) | 32 × ~8,000 = 256k | 32 × ~2,000 = 64k | ~$2.05 |
| Rubric judge (32 × Haiku) | 32 × ~3,000 = 96k | 32 × ~300 = 9.6k | ~$0.04 |
| **Total** | | | **~$2.10** |

Breakdown: input at $3/Mtok (Sonnet), $0.25/Mtok (Haiku); output at $20/Mtok (Sonnet),
$1.25/Mtok (Haiku). Assumes the user's current model is Sonnet for `default` tier.

**Ruling: on-demand + pre-release. NOT nightly by default.**

A $2.10 sweep is affordable but not free. Nightly runs would cost ~$63/month for no signal gain
(agent prompts don't change nightly). The value is in detecting regressions from model upgrades
or prompt edits — both of which happen infrequently and intentionally.

**Flags the runner must support:**

- `--agents <comma-separated-names>`: run only the specified subset
  (e.g., `--agents analyst,code-reviewer` for targeted post-edit regression)
- `--pattern-only`: skip LLM passes entirely; deterministic check only (zero cost)
- `--no-judge`: run generation but skip rubric grading (cheap spot-check: ~$2.00)

**No API key → SKIP, not FAIL.** The script checks for `ANTHROPIC_API_KEY` before any LLM
invocation. If absent, all agents emit SKIP (exit 0). Pattern checks run regardless (they need
no key). This is the same lesson as commit `6c422cb` — a harness that fails because of missing
infra is worse than one that gracefully degrades.

---

## 4. Regression record format

**Location: `shared/evaluation/agent-evals/`** (new subdirectory of `shared/evaluation/`).

This aligns with the `shared/evaluation/` home the epic specifies and keeps evaluation artifacts
away from the documentation tree (`docs/agent-metrics/evals/`). The agent-eval skill will be
updated to write here too (Op 4 — docs update).

**Schema (one Markdown file per agent per run):**

```
shared/evaluation/agent-evals/<agent>-eval-<YYYY-MM-DD>.md
```

```markdown
# Agent Eval: <agent> — <YYYY-MM-DD>

**Agent version**: <from frontmatter>
**Model used**: <concrete model id, e.g. claude-sonnet-4-6>
**Fixture**: tests/agents/<agent>/<input-file>
**Run mode**: full | pattern-only | no-judge

## Pattern Grade
- [PASS|FAIL|SKIP] `<pattern>` — <"matched" or "not found in output">
**Pattern overall**: PASS | FAIL

## Rubric Grade
- [PASS|FAIL] <criterion label> — <one-line evidence or gap>
**Rubric overall**: PASS | FAIL | SKIP (if --pattern-only or no key)

## Regression Delta
Compared against: <previous file name, or "no baseline — first recorded eval">
- [REGRESSION|STABLE|IMPROVED] Pattern: <was X, now Y>
- [REGRESSION|STABLE|IMPROVED] Rubric: <was X, now Y>
**Overall delta**: REGRESSION | STABLE | IMPROVED | BASELINE
```

**Comparison logic:**
1. Find the most recent `<agent>-eval-*.md` in `shared/evaluation/agent-evals/` (sorted by date,
   pick latest).
2. Parse the previous file's "Pattern overall" and "Rubric overall" fields.
3. A REGRESSION is: previous=PASS + current=FAIL for either dimension. This fails the harness
   run (exit 1) and emits the agent name to stderr.
4. On first run (no prior file), emit "BASELINE" — not a failure.

Machine-readable fields use a consistent `**Field**: value` format so bash `grep` can extract
them without a full Markdown parser.

---

## Relationship to existing tooling

- `scripts/test-agents.sh`: structural check only (deterministic), no LLM, always CI-safe.
  The harness is additive — it never replaces this script; it adds the LLM-graded layer.
- `agent-eval` skill: interactive, single-agent. The harness automates the same logic in batch.
- `shared/evaluation/agent-evals/`: the harness writes here; the `agent-eval` skill will be
  updated to also write here (today it writes to `docs/agent-metrics/evals/`), so both paths
  contribute to the same baseline pool.

---

## Implementation plan (Phase B)

| Op | File | Commit |
|---|---|---|
| 1 | `scripts/run-agent-evals.sh` | `feat(evaluation): headless agent-eval runner (Epic 61 Op 1)` |
| 2 | `shared/evaluation/agent-evals/` schema doc + comparison logic | `feat(evaluation): regression record schema + diffing (Epic 61 Op 2)` |
| 3 | `.github/workflows/framework-ci.yml` opt-in CI job | `feat(ci): opt-in agent-eval regression job (Epic 61 Op 3)` |
| 4 | `tests/agents/README.md` harness documentation | `docs(tests): document eval harness (Epic 61 Op 4)` |

---

*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md).*
