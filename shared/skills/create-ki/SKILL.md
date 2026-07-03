---
name: create-ki
description: Captures a pattern, bug fix, or architectural decision as a searchable Knowledge Item in shared/knowledge/ (portable, cross-project) or .claude/knowledge/ (project-specific), following the KI format defined in shared/knowledge/README.md.
triggers:
  keywords: ["create-ki", "save as knowledge item", "capture this pattern", "document this fix"]
  intentPatterns: ["Save this as a KI", "Document this pattern", "Capture this fix as a knowledge item", "/create-ki *"]
standalone: true
---

## When To Use
- After solving something non-obvious that will plausibly recur: a bug fix with a subtle root cause, a
  pattern adopted for a specific reason, a decision that isn't already an ADR.
- At the end of a `deliver-feature` run, if `retrospective.md`'s "What To Improve" surfaces something worth
  making reusable (this overlaps with `extract-lessons`, Epic 15 — `extract-lessons` calls this skill rather
  than duplicating the write logic).

Do NOT use for:
- Anything already covered by an ADR (`docs/adrs/`) — link to the ADR instead of duplicating it as a KI.
- Project structure, file paths, or conventions derivable by reading the code — KIs are for the *why*, not
  the *what* (same principle as this framework's own memory system).
- A one-off fix with no plausible recurrence — not everything needs to be preserved.

## Context To Load First
1. Run `search-ki` first for the same topic — if a KI already covers this, update it instead of creating a
   duplicate.
2. `shared/knowledge/README.md` — the exact frontmatter format required

## Process
1. **Search first** (via `search-ki`). If an existing KI covers the same ground, propose updating it instead
   and stop here unless the user explicitly wants a separate new KI.
2. **Decide location**: `shared/knowledge/` if the pattern applies to any project using this framework;
   `.claude/knowledge/` if it's specific to this codebase only (create the directory if it doesn't exist).
3. **Name the file**: kebab-case, descriptive, not generic (`retry-storm-from-missing-jitter.md`, not
   `bug-fix-1.md`).
4. **Write the KI** with frontmatter (`name`, `tags`, `domain`, `created`) and a body that states: what the
   pattern/issue is, why it exists or why it happened, and when it applies — not a step-by-step "what we did"
   narrative. Keep it under ~40 lines; a KI is a pointer to a lesson, not a full writeup.
5. **Confirm** the file path back to the caller.

## Output Format
The KI file itself:
```markdown
---
name: kebab-case-slug
tags: [tag-one, tag-two]
domain: bounded-context-or-area
created: YYYY-MM-DD
---

Body: the pattern, decision, or fix — what it is, why it exists, when it applies.
```

## Guardrails
- **Always** search first — a duplicate KI with slightly different wording is worse than no KI, since
  `search-ki` can't tell which one is authoritative.
- **Never** write a KI that just restates something in `ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, or
  an existing ADR — link to those instead.
- **Never** include secrets, credentials, or customer data in a KI body — same rule as everywhere else in
  this framework.
- Keep it short. A KI that requires its own table of contents has become documentation, not a knowledge item
  — split it or move it to `docs/`.

## Standalone Mode
Pure local file writes. No external calls.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
