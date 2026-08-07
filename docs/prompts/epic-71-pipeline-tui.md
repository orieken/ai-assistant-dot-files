# Epic 71 — Pipeline TUI / Visual Delivery Dashboard

Source: `docs/audits/framework-audit-2026-08-07.md` §3 item 3.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` (this IS the git repo — commits
land here directly). Do NOT push.

## Prior context

`deliver-feature` and `deliver-bugfix` are text-based pipeline skills. Progress is visible
only through the scrolling LLM transcript. No visual summary of pipeline state exists.

Existing pipeline observability:
- `pipeline-trace.json` (written by `deliver-feature`): records phase names, timestamps,
  iteration count, duration per phase.
- `pipeline-state.json` (if it exists): holds the active phase name and the feature slug.
- `shared/telemetry/event-schema.md` v1.1.0: `gate_decision` events captured at gates.
- `.claude/feature-workspace/<feature-name>/` — per-feature workspace directory.

The audit asks for a TUI showing: stage progression, active agent, live duration, estimated
token spend, and pending gate approvals. This epic delivers that. The hardest design question
is whether the TUI runs as a live sidecar (requires IPC with the LLM process, which may not
be possible) or as a post-hoc reader polling `pipeline-trace.json`.

## Scope

**Phase A — Design (one commit, then PAUSE for user approval):**

Draft and commit as `docs(pipeline): design pipeline TUI (Epic 71 Phase A)`:

Produce `docs/patterns/pipeline-tui-design.md` answering these rulings:

1. **Runtime model**: live sidecar (reads state from a shared file the pipeline writes to
   on each phase transition) or post-hoc reader (polls the trace file after pipeline ends).
   Rationale: can `deliver-feature` be modified to write a `pipeline-state.json` update on
   each phase transition? If yes, the sidecar model is viable. If no, post-hoc only.

2. **Language + library**: Go (`bubbletea` — TUI framework, already used in the ecosystem)
   or Python (`rich` — already required by no framework dependency, lower overhead). Decision
   criteria: does the framework have an existing Go or Python toolchain dependency that
   justifies the choice?

3. **Data model**: what fields does the TUI display? Define the `pipeline-state.json` schema
   extension (or new file) the TUI reads. Minimum fields: `phase`, `agent`, `started_at`,
   `token_estimate`, `pending_gate` (nullable).

4. **MVP feature set**: progress bar only? Or full agent log tail? Define the 3 screens the
   MVP must show.

5. **Install path**: `install.sh --with-tui`? Or ship as a script in `scripts/`?

**Phase B — Implementation (after approval; one commit per op):**

Op 1 — `feat(pipeline): pipeline-state.json schema + deliver-feature write hook (Epic 71 Op 1)`:
- Define `pipeline-state.json` schema (add to `shared/telemetry/` or alongside `pipeline-trace.json`).
- Modify `shared/skills/deliver-feature/SKILL.md` to write a `pipeline-state.json` update
  instruction at each phase transition (the LLM executing the skill must write the file).
- Same update for `shared/skills/deliver-bugfix/SKILL.md`.
- `bash scripts/health-check.sh` green.

Op 2 — `feat(pipeline): TUI binary (Epic 71 Op 2)`:
- Implement the TUI binary per the Phase A language/library ruling.
- Binary reads `pipeline-state.json` from the path `./pipeline-state.json` (current dir)
  or `--workspace <path>` flag.
- Polls every 2 seconds (configurable via `--interval`).
- Displays: phase name, active agent, elapsed time, token estimate, gate status.
- Exits cleanly on Ctrl-C or when `pipeline-state.json` contains `"status": "complete"`.

Op 3 — `feat(scripts): install.sh --with-tui flag (Epic 71 Op 3)`:
- Add `--with-tui` to `scripts/install.sh`.
- Compiles and links the TUI binary to the install target.
- Adds detection to `scripts/health-check.sh` (WARN if `--with-mcp` was used but `--with-tui`
  was not — they are companion tools).

Op 4 — `docs(runbooks): pipeline TUI usage guide (Epic 71 Op 4)`:
- Add `docs/runbooks/pipeline-tui.md`: how to launch the TUI alongside a pipeline run,
  what each screen means, how to interpret token estimates, and known limitations
  (e.g., token spend is an estimate, not a live API metric).
- Update `docs/RUNBOOKS.md` index.

After every commit: `bash scripts/health-check.sh` green.

## Discipline

Standard — match other prompts in `docs/prompts/`: per-op commits, Conventional Commits,
explicit `git add` paths only, never push.

## Escalation

- If Phase A concludes that the pipeline skills cannot write `pipeline-state.json` updates
  (because the LLM executing them is not guaranteed to produce file writes in the right
  location), halt and propose an alternative: a post-hoc `pipeline-replay` command that
  renders a static TUI replay from an existing `pipeline-trace.json` file.
- If the chosen TUI library introduces a new language runtime dependency the framework has
  no other reason to include, halt and propose the alternative language.
- If the TUI binary compilation requires CI changes (new toolchain in `scripts/ci-check.sh`),
  document the CI change needed and halt for approval before adding it.
- If Phase B Op 1 requires changes to more than 3 skill files to add phase-write hooks,
  halt and reassess: the write hook may need to be a shared helper rather than per-skill.

## Report (under 250 words)

```
Phase A commit: <sha>
Phase A rulings:
  - Runtime model: <live sidecar | post-hoc reader>
  - Language + library: <Go bubbletea | Python rich | other>
  - pipeline-state.json schema: <field list>
  - MVP screens: <list>
  - Install path: <install.sh --with-tui | scripts/ only>

Phase B commits (if approved):
  <sha> <message>
  ...
Verification: health-check <pass>, TUI launches against fixture trace file <pass>.
```

Go.
