# Epic 53 — Deterministic Inventory Drift Check

Source: `docs/audits/framework-gap-audit-2026-07-25.md` § Dimension 5 (highest-priority Epic per the audit's ranking).

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`. Do NOT push.

## Prior context

Today's `scripts/health-check.sh` validates frontmatter presence but does NOT compare actual `shared/*/` inventories against the counts claimed in prose docs (README, `AGENT_REFERENCE.md`, generated platform configs, changelog). Drift already observed: `framework-gap-audit-2026-07-25.md` claims "12 contracts" but `ls shared/contracts/*.md | wc -l` returns 16.

## Scope

**One commit.** Create `scripts/check-inventory-drift.sh`:

- Counts actual files in `shared/agents/*.md` (excluding CHANGELOG.md), `shared/skills/*/SKILL.md`, `shared/rules/*.md`, `shared/contracts/*.md`, `shared/templates/*.md`, `shared/schemas/*.schema.json`
- Greps `README.md`, `docs/AGENT_REFERENCE.md` (if exists), and any generated platform configs under `.cursor/`, `.windsurfrules`, `.github/copilot-instructions.md`, `.openai.md`, `AGENTS.md` for stated counts (e.g., "36 agents", "N skills")
- Reports each drift: `EXPECTED: <n> per <file>:<line>, ACTUAL: <n>`
- Exit 0 if no drift; exit 1 with drift list on stderr

Wire into `scripts/health-check.sh` as a new check step (INFO level; upgrade to FAIL when count discrepancies are important). Also wire into `scripts/ci-check.sh` if it exists.

Update `docs/audits/framework-gap-audit-2026-07-25.md` if any counts in it are wrong post-fix.

Commit: `feat(ci): add inventory-drift check (Epic 53)`.

## Discipline

- Match commit + git-add discipline from other prompts in `docs/prompts/` (per-op commits, never `git add -A`, explicit paths).
- Verify the script runs green on the current repo state before commit.

## Escalation

- If the drift count exceeds ~5 across the repo, halt and report — that's audit-worthy on its own.
- If a prose doc has counts that are correct AT THE TIME OF WRITING but stale now (like an old blog draft citing "24 agents"), don't auto-fix — surface with `INFO` and let human decide.

## Report (under 100 words)

```
Commit: <sha>
Drift found (pre-fix): <n> discrepancies across <files>
Post-fix state: <n> drift items remaining (should be 0 or explained)
Wired into health-check.sh: yes | no
```

Go.
