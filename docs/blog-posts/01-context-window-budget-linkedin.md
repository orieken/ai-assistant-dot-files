AI coding tools make it very easy to confuse "more context" with "better context."

While building `ai-assistant-dot-files`, I ended up treating the context window like a budget instead of a
junk drawer.

The framework now uses a dedicated `context-engineer` agent to produce a `context-manifest.md` before the
rest of the delivery pipeline starts. That manifest scopes relevant files, prior deliveries, Knowledge
Items, ADRs, and token pressure before implementation begins.

The useful distinction:

- Context = what is loaded right now
- Memory = durable knowledge that outlives the run
- Learning = feedback loops that change future behavior

That distinction alone prevents a lot of accidental "RAG will fix it" thinking.

Full post: https://dev.to/orieken/treat-the-context-window-like-a-budget-not-a-junk-drawer-5602

