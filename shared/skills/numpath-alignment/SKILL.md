---
name: numpath-alignment
description: Checks that a newly built NumPath feature aligns with MacLellan's Intelligent Tutoring System framework and the KC/BKT research foundations. Use after shipping any feature that touches student modeling, problem selection, or feedback. Also useful before writing a paper section to verify the system's theoretical grounding holds.
triggers:
  keywords: ["alignment", "MacLellan", "ITS", "knowledge component", "tutoring framework", "research alignment", "theoretical grounding"]
  intentPatterns:
    - "check alignment with MacLellan"
    - "does this align with ITS research"
    - "is this theoretically grounded"
    - "check our research alignment"
    - "before I write the paper section"
standalone: true
---

## When To Use

- After shipping a feature that touches the adaptive engine, student modeler, or feedback loop
- Before writing a paper section that claims alignment with ITS/KC research
- When a design decision feels like it might deviate from established tutoring theory
- As a pre-flight check before Phase 4 experiment design

Do NOT use for: frontend styling, infrastructure changes, tooling decisions.

---

## MacLellan Alignment Checklist

Run through each principle and mark as ✅ aligned, ⚠️ partial, or ❌ not yet implemented.

### 1. Knowledge Component Decomposition
> Skills are atomic, prerequisite-linked KCs — not vague "topics"

Check:
- Are all skills defined in `DOMAIN_DICTIONARY.md` with a `code`, `name`, and `domain`?
- Does each problem tag exactly one `skill_id`?
- Are prerequisite relationships defined (even if empty in Phase 1)?

### 2. Mastery Before Advancement
> Students only advance when they demonstrate genuine understanding (p_mastery ≥ threshold)

Check:
- Does `AdaptiveEngine` respect `MASTERY_THRESHOLD = 0.80`?
- Is the next problem selection gated on KC mastery, not just session accuracy?
- Are students ever pushed to harder problems without meeting the BKT threshold?

### 3. Error as Diagnostic Signal
> Mistakes trigger targeted remediation, not generic retry

Check:
- Does `MistakeClassifier` tag errors with a structured `mistake_code`?
- Does a classified mistake influence the next problem selection (skill targeting)?
- Are `MistakeEvent` records persisted for longitudinal analysis?

### 4. Deliberate Practice at the Edge of Mastery
> Difficulty targets the current KC frontier — not too easy, not too hard

Check:
- Does `AdaptiveEngine` use `current_difficulty ± DIFFICULTY_STEP` logic?
- Is the frustration detection window (`FRUSTRATION_WINDOW = 3`) active?
- Is the mastery advancement window (`MASTERY_WINDOW = 3`) active?

### 5. Teacher-in-the-Loop
> No black-box AI — every insight must be explainable

Check:
- Does every `NextProblemResponse` include a `reason` field?
- Does the teacher dashboard show KC states, not just accuracy scores?
- Are AI-generated insights (LLM) traceable back to specific KC + mistake data?

### 6. Longitudinal Student Modeling
> The system learns from the full history of a student, not just the current session

Check:
- Are `kc_states` persisted across sessions (not reset on logout)?
- Are `MistakeEvent` records queryable by student over time?
- Is `opportunity_count` used in BKT updates?

---

## Context To Load First

1. `docs/phd-roadmap/MASTER_PROMPT.md` — current phase and stack
2. `DOMAIN_DICTIONARY.md` — canonical KC definitions
3. `numpath/ml/numpath_ml/bkt.py` — BKT implementation
4. `numpath/ml/numpath_ml/adaptive_engine.py` — engine logic
5. `numpath/backend/backend/use_cases/submit_attempt.py` — full attempt flow
6. The feature or component being reviewed

---

## Process

1. Read the six checklist items above
2. Read the relevant source files
3. For each item: mark ✅ / ⚠️ / ❌ with a one-sentence justification
4. List any gaps as concrete TODOs with file paths
5. If writing a paper section: confirm at least 5/6 are ✅ before drafting

---

## Output Format

```markdown
## MacLellan Alignment Check — [feature name] — [date]

| Principle | Status | Notes |
|-----------|--------|-------|
| KC Decomposition | ✅ | All skills in DOMAIN_DICTIONARY, each problem has skill_id |
| Mastery Before Advancement | ✅ | MASTERY_THRESHOLD = 0.80 enforced in AdaptiveEngine |
| Error as Diagnostic | ⚠️ | MistakeClassifier exists but doesn't influence problem selection yet |
| Deliberate Practice | ✅ | Frustration/mastery windows active |
| Teacher-in-the-Loop | ⚠️ | reason field present, but dashboard KC view not yet built |
| Longitudinal Modeling | ✅ | kc_states persisted, opportunity_count tracked |

### Gaps
- [ ] `numpath/ml/numpath_ml/adaptive_engine.py`: mistake_code should weight skill selection
- [ ] `numpath/frontend/src/views/TeacherView.vue`: add KC heatmap to dashboard

### Paper-ready?
3 of 6 principles fully implemented. Not yet ready for paper claims about ITS alignment.
```

Show to user for review. This output is informational — no file write unless user asks to save it.

---

## Guardrails

- Never claim full MacLellan alignment unless all 6 principles are ✅
- Flag ⚠️ honestly — "partial" is not the same as "done"
- If the paper section is being drafted: block on any ❌ items
- This skill reads code — it does not modify it

---

## Standalone Mode

Fully conversational. Reads source files, produces a checklist report. No external tools required.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
