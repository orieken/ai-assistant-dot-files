# Runbook: Cross-Agent Self-Audit Prompt

## When to use
Periodically, or after a burst of changes to `shared/` (new agents, skills, or contracts) — hand this
prompt to a *different* agent than whichever one did the work (a different AI tool, or a fresh session with
no memory of the changes) to get an independent check. An agent auditing its own recent work tends to
confirm its own assumptions; a different agent reading only the files has no such bias.

This is not a substitute for `scripts/check-parity.sh` or `scripts/health-check.sh` — those catch what's
scriptable (symlink drift, frontmatter presence, changelog/version pairing). This prompt is for everything
those scripts can't check: whether two files that are supposed to agree with each other actually do,
whether a doc still describes the repo as it used to be, whether a "do NOT use X, use Y" disambiguation is
actually correct.

## Prerequisites
- The target agent needs read/search access to the repo (any tool — this prompt makes no assumptions beyond
  that).
- No special setup. This can be run against any checkout of the repo, including a fresh clone.

## The Prompt

Copy everything below the line into a new session with a different agent than the one that made recent
changes:

---

```
You are auditing the "ai-assistant-dot-files" Context Engineering Framework — a repo that defines a
canonical set of AI agents, skills, and rules once in shared/, then generates or symlinks them into six
different AI coding tools (Claude Code, Cursor, Windsurf, GitHub Copilot, Gemini/Antigravity, OpenAI/Codex).
It also ships a 14-agent feature-delivery pipeline (spec -> analysis -> architecture -> implementation ->
review -> security -> QA -> docs -> deploy) orchestrated by shared/skills/deliver-feature/SKILL.md.

Read first, in this order:
1. README.md — architecture, agent roster, skill catalog, pipeline diagram
2. docs/ARCHITECTURE.md — the shared/ canonical layer and tier system
3. docs/CONTRIBUTING.md — the conventions every addition is supposed to follow
4. docs/runbooks/context-engineering.md — how context, memory, and learning are supposed to work together
5. shared/DOMAIN_DICTIONARY.md — the ubiquitous language every file is supposed to match

Your job is not to summarize what the framework claims to do. It's to verify whether it actually does it,
by checking the source files against each other. Do not accept a skill's or agent's own description as
proof it's correct — read the file it describes, and read the files it references, and confirm they agree.

Do not modify files during this audit unless the user explicitly asks for fixes. Produce findings only.

Audit these specific dimensions:

1. **shared/ is the single source of truth.** Every generated platform config (.cursor/, .github/,
   .windsurfrules, .openai.md, AGENTS.md, .claude/ symlinks) should be derivable from shared/ with no
   drift. If scripts/check-parity.sh exists, run it — but also spot-check by hand, since an automated
   check can itself have a blind spot: at minimum, check one rule file, one generated persona file, one flat
   prompt file, and one symlink target against shared/ directly, rather than trusting a green run alone.
   When any validation script passes, inspect what it actually checks and name at least one class of drift
   it does not cover — a passing script proves the checks it runs are clean, not that nothing is wrong.
   Concretely: a plain local run of these scripts is not proof they pass in CI — if Docker is available, run
   `scripts/ci-check.sh`, which executes them inside the same OS/bash version the CI runner uses (this repo
   already hit a bash-version-specific bug that a local macOS run couldn't detect at all).

2. **"Twin" files must actually match.** Where an agent has both a shared/agents/<name>.md (native agent)
   and a shared/skills/<name>/SKILL.md (standalone version), their Output Format sections — exact headings,
   not just intent — must be identical. A mismatch here means a contract written against one will silently
   fail against the other.

3. **Every pipeline handoff that has a contract is actually validated, and every artifact that should have
   a contract does.** Before reporting findings for this dimension, build and show the inventory you worked
   from — a table with one row per artifact deliver-feature's agents produce:
   `Artifact | Producer | Produced at step | Contract exists? | validate-artifact mapping? | deliver-feature invocation? | Status`.
   Building this table first (rather than jumping straight to findings) is what catches an artifact that
   isn't obviously broken but is still missing coverage — and it lets a human check the audit itself quickly.
   Flag any artifact missing any of the three columns — especially early-pipeline artifacts everything
   downstream depends on.

4. **Numbered steps stay internally consistent.** If deliver-feature (or any other skill) has numbered
   process steps, confirm every cross-reference to "step N" elsewhere in the same file (and in other files
   that cite it) still points at the right step after any edits.

5. **Documentation matches current reality, not a past version of the repo.** For every runbook/doc, check:
   does it reference files, scripts, or paths that still exist? Does it describe the current agent/skill
   count, current pipeline shape, current platform list? A doc that was accurate when written but never
   updated as the repo grew is worse than no doc, because it's trusted by default. Separate active
   docs/runbooks from clearly archived or historical planning material first — this repo has both (e.g. old
   prompt-engineering templates, a pre-restructure legacy tree kept as a historical snapshot). Report stale
   archived material only if it's still linked from somewhere as current operational guidance; otherwise it's
   intentionally preserved history, not a bug.

6. **Skills with overlapping-sounding purposes actually disambiguate correctly.** Where multiple skills
   sound similar (e.g. anything with "retrospective," "score," "audit," or "eval" in the name), check that
   each one's "do NOT use this for X, use Y instead" language is accurate given what X and Y actually do —
   not just present, but correct.

7. **Versioning integrity.** Every shared/agents/*.md version bump has a matching shared/agents/CHANGELOG.md
   entry, and vice versa. No agent's version implies a behavior change that isn't reflected in its CHANGELOG
   description.

8. **No dead references.** Grep for every internal file path, skill name, and script name mentioned in
   shared/ and docs/ and confirm the target still exists. Flag anything pointing at a renamed, moved, or
   deleted file.

9. **Ubiquitous language matches DOMAIN_DICTIONARY.md.** Spot-check whether terms defined there are used
   consistently, and whether any term used repeatedly across the repo is missing from the dictionary.

10. **Irreversible actions are actually gated.** Confirm every commit, deploy, migration, and external API
    call described in the pipeline routes through shared/rules/approval-gates.md's checkpoints, not just in
    prose but in the actual invoking skill's steps.

Report findings ranked most-severe first. For each finding, state:
- The specific file(s) and line(s) involved
- What's actually there, quoted or closely paraphrased — not a summary
- The concrete failure scenario this causes (not "this is inconsistent" — "this causes X to happen when Y")
- Whether you verified it by reading both sides yourself (CONFIRMED) or you're inferring from partial
  evidence (PLAUSIBLE — flag it as such, don't present a guess as a finding)

Do not report a finding you haven't personally traced through the actual files. A plausible-sounding
inconsistency that turns out to be correct on inspection is worse than no finding at all.
```

---

## Verification
A well-run audit produces a short list (typically under 10) of findings that are each independently
checkable in under a minute by opening the two files it names. If the report reads as a restatement of each
file's own description rather than a cross-check between files, the audit didn't actually do its job —
re-run with an explicit reminder to quote actual file content, not summarize intent. Dimension 3's artifact
matrix should be present in the report itself, not just implied by the findings — its absence is a sign the
auditor skipped straight to conclusions.

## Known Judgment Calls (Reviewed — Do Not Re-Flag As New Findings)
Some things an audit surfaces are real observations but not gaps — they're already-documented, deliberate
tradeoffs. Before treating one of these as a new finding, check whether it's already covered here:

- **Standalone agents have uneven governance depth** (some, like `finops-engineer`/`dx-engineer`/
  `dependency-auditor`, have no check beyond a human reading the report). Reviewed 2026-07-05: this is
  already called out explicitly in `docs/AGENT_REFERENCE.md`'s "What this survey actually shows" section
  (#2) as a conscious judgment call tied to each agent's blast radius, not an oversight. Only re-flag this
  if a *specific* standalone agent's risk profile has changed (e.g. it gained the ability to take an
  irreversible action) since that entry was written.
- **Some platform behaviors in `docs/ARCHITECTURE.md` are marked unconfirmed rather than confirmed.**
  Reviewed 2026-07-05: this is intentional disclosure of actual verification status (confirmed = live-tested
  against that platform; unconfirmed = documented honestly rather than assumed), not an unaddressed gap.
  Building automated fixtures for every claimed platform path isn't proportionate — most of these require a
  live account/session on a third-party tool to verify. Only re-flag a *specific* line if the platform's
  actual behavior has since been confirmed to contradict what's written.

If a future audit re-raises one of these unchanged, treat it as confirmation the existing documentation is
correct, not a new action item — cite this section back rather than re-litigating from scratch.

## Escalation / Acting on Results
- Treat every CONFIRMED finding as real; verify PLAUSIBLE findings yourself before acting (read both files
  the same way the audit was supposed to). Check "Known Judgment Calls" above first — a finding that matches
  one exactly isn't new.
- Fix findings in the same spirit as the framework's own conventions: bump affected agent versions and add
  a `shared/agents/CHANGELOG.md` entry (see [editing-agent-prompts.md](editing-agent-prompts.md)), re-run
  `scripts/check-parity.sh` and `scripts/health-check.sh` after any fix, and commit per logical unit of work
  rather than batching unrelated fixes into one commit.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/loom) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/loom/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
