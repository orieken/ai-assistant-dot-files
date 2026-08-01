# Adapter: Source Retrieval (DEFERRED)

**Corpus**: `project-source`
**Status**: Deferred — no implementation planned for v3.x

## Deferral Rationale

Retrieval over installed-project source code is deferred because:

1. **Modern LLM clients already do this well.** Claude Code's native Grep, Glob, and Read tools perform
   symbol-level, path-pattern, and content search over arbitrary codebases without any framework
   integration. Cursor does the same via its built-in codebase index. Building a parallel framework
   retrieval path over the same files adds operational cost (index maintenance, embedding drift, rebuild
   triggers) without clear uplift over what the client already provides.

2. **No telemetry signal yet.** ADR-002 commits to "add vector-backed code retrieval only if telemetry
   shows the built-in isn't sufficient." No telemetry data exists yet (the telemetry layer landed in
   Phase 1 but retrieval event logging isn't wired). Deferral is the conservative choice until a class of
   miss becomes chronic and visible in the log.

3. **Interface is ready.** The `Retrieve(query, "project-source")` call in `retriever.interface.md` is
   the integration point. When source retrieval is eventually built, it implements that interface without
   touching any other adapter or caller.

## When This Would Graduate Off Deferred

- Telemetry log shows repeated "miss" patterns — queries where Claude Code Grep returned nothing useful
  and the agent fell back to a wide read of the source tree, suggesting a semantic index would have
  helped.
- A class of useful query (e.g., "which service handles payment retries?") becomes common and Grep
  can't satisfy it without enumerating many files.
- The `saturday-mcp` vector adapter for `project-features` proves stable for > 3 months, validating
  the sqlite-vec embed → query pipeline, making the same approach for source low-risk.

## References

- Interface: [`../retriever.interface.md`](../retriever.interface.md)
- ADR: [`../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md`](../../../docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md)
- Migration plan: `docs/aos/migration-plan.md` Op 3.4
