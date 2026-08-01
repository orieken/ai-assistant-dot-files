---
name: learning-engine
description: Analyzes past pipeline delivery retrospectives and extracts candidate lessons for docs/lessons-learned/. Opposing-force pair with forgetting-engine. In AOS Phase 3, wired as a hook on-retrospective-written (opt-in, disabled by default).
triggers:
  keywords: [learning engine, extract learning, process feedback, retrospective analysis]
  intentPatterns: ["extract lessons from retrospectives", "run learning engine"]
standalone: true
---

## When To Use

Use when:
- Running a feedback sweep across completed feature deliveries in `docs/features/`.
- Proposing new entries for `docs/lessons-learned/`.
- Invoked automatically by the `on-retrospective-written` hook (AOS Phase 3, opt-in).

Do NOT use when:
- Expiring or archiving old lessons (use `forgetting-engine`).

## Invocation Modes

| Mode | Trigger | Scope |
|---|---|---|
| Manual | Human runs `/learning-engine` | Full sweep of all retrospectives |
| Hook (scope=latest) | `on-retrospective-written` hook after a retrospective write | Latest retrospective only |
| Hook (scope=all) | `on-retrospective-written` hook with `scope: "all"` | Full sweep |

## Context To Load First

1. `docs/lessons-learned/` — existing lessons (to avoid proposing duplicates)
2. `docs/features/*/retrospective.md` — one file in hook/latest mode; all files in manual/all mode
3. `shared/knowledge/*.md` — check if a proposed lesson is already documented as a KI

## Process

1. **Check for duplicates**: load `docs/lessons-learned/` and `shared/knowledge/` to identify lessons already documented. A proposed lesson that duplicates an existing KI should note "already captured in KI: <file>" rather than proposing a new lessons-learned entry.

2. **Scan retrospectives** (per scope):
   - `latest`: read only the retrospective file passed via `passContext`
   - `all`: Glob `docs/features/*/retrospective.md`

3. **Identify promotable patterns**:
   - Recurring test breakages (same contract failing across ≥ 2 features)
   - Agent retry patterns that suggest a missing guardrail
   - Decisions made under ambiguity that resolved well and should be defaults
   - Anti-patterns discovered during QA or security review

4. **Draft proposal** in `.claude/feature-workspace/proposed-lessons.md`:

```markdown
# Proposed Lessons: YYYY-MM-DD

## Executive Summary
[1-2 sentences on the recurrent pattern]

## Proposed Entries

### 1. [Lesson title]
**Pattern**: [What recurred]
**Evidence**: [Feature(s) where this appeared]
**Recommendation**: [What to do differently]
**Promote to KI?**: [yes — suggest /create-ki after approval | no — lessons-learned is sufficient]

- [ ] Approve entry 1
```

5. **Pause for human confirmation**: Present the draft and ask explicitly:
   > "The learning engine found N promotable patterns. Reply 'approve all', 'approve N', or 'skip' to decide."
   
   Do NOT proceed to step 6 without explicit confirmation.

6. **Persist approved lessons**: for each approved entry, write to `docs/lessons-learned/lessons-YYYY-MM-DD.md`. If the entry is flagged "Promote to KI?": yes, remind the human to run `/create-ki` — do not invoke it automatically.

## Output Format

```markdown
# Proposed Lessons: YYYY-MM-DD

## Executive Summary
[Description of pattern]

## Proposed Entries
### 1. [Title]
**Pattern**: ...
**Evidence**: ...
**Recommendation**: ...
**Promote to KI?**: yes | no
- [ ] Approve entry 1
```

## Guardrails

- **Never write to `docs/lessons-learned/` without explicit human confirmation** ("approve" or "yes"). This is non-negotiable regardless of invocation mode or hook args.
- Never invoke `/create-ki` automatically — propose it, let the human trigger it.
- If invoked via hook and no promotable patterns are found, emit a brief "no lessons extracted" note and exit cleanly. Do not block the retrospective write.
- If `docs/lessons-learned/` does not exist, create it with a minimal `README.md` before writing the first lesson.

## Opposing Force

`forgetting-engine` is this skill's opposing force. It removes lessons that have become obsolete.
Together they maintain the quality of `docs/lessons-learned/` over time.

## Hook Configuration

See `shared/hooks/on-retrospective-written.yaml` for the opt-in hook definition. Copy to
`.claude/hooks/` and set `enabled: true` to activate automatic invocation.

## Standalone Mode

Operates purely using local file editing tools (Read, Write, Edit, Glob).

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) AOS Phase 3 Runtime layer. CC BY 4.0.*
