# Eval Rubric: context-engineer / input-task.md

- **File list is scoped to the task**: the manifest pins `auth.middleware.ts` and `redis-client.ts` as high-priority — does not include every file under `src/api/` or `src/cache/` as a blanket load.
- **ADR-007 is referenced, not just noted**: the manifest explicitly links ADR-007 as required reading for the Redis usage decision, not just as a "see also" footnote.
- **Previous delivery is pruned intelligently**: the auth-overhaul feature archive is included only if it contains rate-limiting notes — the manifest explains WHY it's included (prior art), not just that it exists.
- **Budget is calculated and within limit**: the context-manifest includes a `## Budget` section with an estimated token count for the pinned files, and the estimate does not exceed the defined budget ceiling.
- **No speculative files included**: the manifest does not include files that might be relevant (e.g., all middleware files, all Redis usages) — only files confirmed relevant to rate-limiting and the specific Redis/auth insertion point.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
