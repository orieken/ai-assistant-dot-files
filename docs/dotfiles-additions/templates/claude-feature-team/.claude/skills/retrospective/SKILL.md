---
name: retrospective
description: Reads all artifacts from a completed feature delivery and produces a retrospective — what went well, what the code-reviewer flagged, what security found, what QA had to fix, and patterns to watch for in future features. Invoke with /retrospective or /retrospective [feature-name]
---

# Feature Delivery Retrospective

You are synthesizing a retrospective from completed feature delivery artifacts.
Input: $ARGUMENTS

## Step 1 — Locate the artifacts

If `$ARGUMENTS` specifies a feature name, look for artifacts related to that feature.
Otherwise, use the most recently modified artifacts in `.claude/feature-workspace/`.

Check which of the following exist:
```
.claude/feature-workspace/analysis.md
.claude/feature-workspace/architecture-notes.md
.claude/feature-workspace/implementation-notes.md
.claude/feature-workspace/code-review-report.md
.claude/feature-workspace/security-report.md
.claude/feature-workspace/qa-report.md
.claude/feature-workspace/docs-report.md
.claude/feature-workspace/devops-report.md
.claude/feature-workspace/delivery-summary.md
```

If fewer than 3 of these exist, tell the user:
> "Not enough delivery artifacts found to run a meaningful retrospective. Complete a
> feature delivery first with `/deliver-feature features/my-feature.md`."
Then stop.

List which artifacts were found and which are missing before proceeding.

## Step 2 — Read all artifacts

Read every artifact that exists. Do not summarize yet — read fully first.

Pay particular attention to:
- **analysis.md**: What was the original scope? Were acceptance criteria specific or vague?
- **architecture-notes.md**: What structural decisions were made? What fitness functions defined?
- **implementation-notes.md**: What deviations from the analysis occurred? What Boy Scout
  cleanups were done? What did the developer flag for QA and DevOps?
- **code-review-report.md**: Was it APPROVED or CHANGES REQUESTED? How many rounds?
  What specific smells were found?
- **security-report.md**: What STRIDE threats were identified? What severity? Were any
  Critical or High findings?
- **qa-report.md**: How many bugs were found during testing? Were any acceptance criteria
  uncoverable? What did QA flag for the tech writer?
- **devops-report.md**: What manual steps were required? What env vars were added?
- **delivery-summary.md**: Final status. Any known issues or follow-ups?

## Step 3 — Pattern analysis

After reading all artifacts, identify:

### Code Quality Patterns
- Which Fowler smells appeared in the code review? (Feature Envy, Magic Numbers, DRY, etc.)
- Were the Sandi Metz rules violated? (class/method size, param count)
- Did the code reviewer need multiple rounds, or was it approved first pass?
- Were Boy Scout cleanups needed? What existing debt was found?

### Analysis Quality Patterns
- Did the analyst identify all the files that were actually affected?
- Did any deviations occur between analysis tasks and implementation? Were they avoidable?
- Were acceptance criteria specific enough that QA could verify them cleanly?
- Were open questions left that caused rework later?

### Security Patterns
- What STRIDE categories showed up? (Spoofing, Tampering, Information Disclosure, etc.)
- Were any findings Critical or High? Could they have been caught earlier (in analysis or architecture)?
- Did the security reviewer fix issues directly, or were they recommendations only?

### Testing Patterns
- How many bugs did QA find in implementation?
- Were any acceptance criteria untestable as written?
- Did characterization tests need to be written for legacy code?

### Process Patterns
- Which conditional agents were invoked (architect, security-reviewer)?
- Were the human checkpoints used effectively, or rubber-stamped?
- What manual steps were required that could potentially be automated?

## Step 4 — Write the retrospective

Write `.claude/feature-workspace/retrospective.md`:

```markdown
# Feature Retrospective: [Feature Name from delivery-summary or analysis]

**Date**: [today]
**Pipeline status**: [from delivery-summary — Complete / Complete with notes / Blocked]
**Agents invoked**: [list which ran, which were skipped]

---

## What Went Well 🟢

[3-5 specific things that worked — be concrete, not generic.
"The analyst identified all 7 affected files correctly" not "Analysis was good."
"Security review found no issues — the developer used auth providers correctly from the start"
not "No security problems."]

---

## Code Quality Findings 🔵

**Code review result**: APPROVED first pass / CHANGES REQUESTED ([N] rounds)

**Smells found**:
- [Smell]: [Where it appeared] → [Refactoring applied]
— or "None identified — clean implementation"

**Complexity debt discovered** (Boy Scout cleanups):
- [File]: [What was cleaned]
— or "None — touched files were already clean"

**Pattern to watch**: [If the same smell appeared, note it as a recurring pattern]

---

## Security Findings 🔴

**Threats identified**: [N total — breakdown by STRIDE category]
- Critical: [N] (all fixed before QA)
- High: [N]
- Medium: [N]
- Low: [N]
— or "No security surface / No findings"

**Could any have been caught earlier?**
[Was there a structural decision in architecture-notes that created the vulnerability?
Was there an acceptance criterion that should have been a security requirement?]

---

## Testing Findings 🟡

**Bugs found by QA**: [N]
- [Bug]: [Root cause — was it an edge case, missing validation, wrong assumption?]
— or "None — implementation matched spec"

**Untestable criteria**: [N]
- [Criterion]: [Why it couldn't be verified]
— or "All criteria verified"

**Legacy code encountered**: [Did characterization tests need to be written? What was found?]
— or "No legacy code touched"

---

## Process Observations ⚙️

**Analysis accuracy**:
- Files correctly predicted: [N/total]
- Tasks that deviated from analysis: [list or "None"]
- Acceptance criteria quality: [Were they specific enough? Any that caused ambiguity?]

**Human checkpoints**:
- Analyst approval: [Was the analysis read carefully or rubber-stamped?]
- Architect approval (if invoked): [Were structural decisions reviewed?]

**Manual steps required**: [N]
- [Step]: [Could this be automated in a future DevOps task?]
— or "None"

---

## Recommendations for Next Time 📋

*Specific, actionable — not generic advice.*

1. [Concrete recommendation]: [Why, based on what was observed in this delivery]

Example format:
- "Add a performance SLA to acceptance criteria for any feature touching the auth endpoint —
  the analyst didn't include one and the devops engineer had to add a timeout retrospectively."
- "The `UserSession` entity is approaching 100 lines — next time we touch it, schedule an
  Extract Class refactor before adding features."
- "Rate limiting features should always trigger the security-reviewer — add it to the
  conditional rules in CLAUDE.md."

---

## Fitness Functions to Add

[Based on what was found — are there architectural constraints that should now be enforced in CI?]
- [New fitness function]: [What it guards, how to implement it]
— or "No new fitness functions identified"

---

## DOMAIN_DICTIONARY Updates

[Were any new terms introduced? Any ubiquitous language violations found by the code-reviewer?]
- Added: [term] — [definition]
- Corrected: [wrong term used] → [correct term]
— or "No dictionary changes needed"
```

## Step 5 — Present to the user

Show the retrospective summary and offer:
> "Retrospective complete. Saved to `.claude/feature-workspace/retrospective.md`.
>
> Top recommendation: [first and most important recommendation]
>
> Would you like me to:
> - [A] Apply any of the recommendations now (update CLAUDE.md rules, add fitness functions, etc.)
> - [B] Update DOMAIN_DICTIONARY.md with any new terms
> - [C] Just save it for reference"

## Rules

- Be specific in every finding — "the developer used manual header injection instead of
  BearerAuthProvider in auth.service.ts line 47" not "there was a security issue"
- Always include "What Went Well" — retrospectives that only surface problems create
  anxiety, not learning
- Never fabricate findings — if an artifact doesn't mention something, don't invent it
- Recommendations must be actionable — "be more careful" is not a recommendation
- If this is the first delivery in the project, note it:
  > "First delivery — no historical patterns to compare against yet. Baseline established."
