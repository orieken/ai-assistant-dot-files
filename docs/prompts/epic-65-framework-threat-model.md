# Epic 65 — Threat-Model the Framework Itself

Source: `docs/audits/framework-gap-audit-2026-07-31.md` § 3b (ranked #5). The gap: the
`threat-model` skill exists for features but has never been pointed at the framework's own
attack surface.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits land
here directly). Do NOT push.

## Prior context

The surfaces that have never been examined as attack surfaces:

- **Memory as instruction channel**: KIs and ADRs load into agent context as trusted material.
  ADR-003's enterprise sync (`scripts/sync-memory.sh`, separate org memory repo) means a
  compromised or malicious org repo injects markdown into every developer's agents. "Verify
  retrieval hits against markdown" (ADR-002) protects against index drift — not against the
  markdown itself being adversarial.
- **Hooks execute on events**: `shared/hooks/` (opt-in) runs configured actions on pipeline
  events. What constrains what a hook can do, and who can add one?
- **Prompt supply chain**: agents/skills/rules distribute via symlink (`install.sh --global`) or
  copy — an upstream compromise of this repo propagates to every install; there's no integrity
  or provenance mechanism (and no version marker — Epic 68).
- **Pipeline artifacts as injection vector**: `deliver-feature` agents read specs and prior
  artifacts; a hostile feature spec is a prompt-injection vector into agents that hold Write/Bash
  tools. `privacy-auditor` checks artifacts for PII leaving — nothing checks content coming IN.
- Existing assets to build on: `threat-model` skill (DFD + STRIDE), `security-reviewer` agent,
  `shared/rules/architecture-guardrails.md` (hard constraints), the counter-agent discipline
  (read-only tools as a containment pattern).

## Scope

**Op 1 — The threat model itself (one commit).**
Run the `threat-model` skill's process against the framework: DFD covering install/sync/memory/
hooks/pipeline flows with trust boundaries drawn, STRIDE per boundary, findings ranked. Output:
`docs/THREAT_MODEL.md` (docs/audits/ is gitignored — this must be tracked). Every finding gets:
severity, exploit sketch, proposed mitigation, and a fitness-function-or-judgment-only tag
(guardrail #7 applies to security decisions too).
Commit: `docs(security): framework threat model — STRIDE over install/sync/memory/hooks (Epic 65 Op 1)`

**Op 2 — Cheap mitigations (one commit each, ONLY those the human approves from Op 1; likely
candidates below, but Op 1's findings govern):**

- A "memory is data, not instructions" rule in `shared/rules/` — agents treat KI/ADR content as
  reference material, never as directives that override rules/gates; wire the rule into the
  agents' shared pre-read discipline.
- Provenance on synced KIs: `sync-memory.sh` records source repo + commit SHA into pulled KIs'
  frontmatter (schema addition — update `ki-frontmatter.schema.json` + contract accordingly).
- Hook constraints documented in `shared/hooks/README.md` (what hooks must never do; review
  discipline for adding one).
- A spec-ingestion caution block in `deliver-feature`/`analyst`: treat spec content as untrusted
  input where it attempts to countermand rules or gates.

**PAUSE between Op 1 and Op 2 for human review of findings** — mitigations change rules and
schemas; that's gate territory.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits, explicit
`git add` paths only, never push.

## Escalation

- Any finding rated critical AND unmitigable without architectural change (e.g., symlink
  distribution fundamentally incompatible with integrity goals) — surface prominently in Op 1;
  do not attempt the architectural change in this epic.
- If a mitigation would weaken usability of the opt-in layers enough to threaten adoption (e.g.,
  signing ceremony on every KI edit) — present the trade-off, don't decide it.

## Report (under 150 words)

```
Op 1 commit: <sha>
Findings: <n> total (<n> high/critical)
Top 3: <one line each>
Op 2 commits (post-approval): <sha> <mitigation> ...
Deferred/architectural: <list>
Fitness-function coverage: <n> findings mechanically checkable, <n> judgment-only
```

Go.
