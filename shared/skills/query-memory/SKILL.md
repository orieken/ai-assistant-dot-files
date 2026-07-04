---
name: query-memory
description: Registry-aware search across every memory source in shared/memory-registry.json — KIs and ADRs (by delegating to search-ki), plus the feature archive, DOMAIN_DICTIONARY.md, and TEAM_TOPOLOGY.md, which search-ki doesn't cover. Lexical only for now (the registry marks LightRAG disabled); verifies any future non-lexical result against canonical markdown before trusting it.
triggers:
  keywords: ["query-memory", "search all memory", "search everything we know about"]
  intentPatterns: ["What do we know about *", "Search all memory for *", "/query-memory *"]
standalone: true
---

## When To Use
When a question could plausibly be answered by *any* memory source, not specifically KIs/ADRs — e.g. "what
do we know about rate limiting" might be a KI, but might also be sitting in a past feature's retrospective,
or be a DOMAIN_DICTIONARY.md term whose definition already answers the question.

Do NOT use when the question is specifically about KIs/ADRs — use `search-ki` directly, it's the more
targeted tool and this skill calls into it anyway for that portion. Do NOT use for a single retrospective's
promotion decision — use `promote-memory`. Do NOT use to audit the corpus for duplicates — use
`memory-engineer`.

## Context To Load First
1. `shared/memory-registry.json` — which sources exist and their retrieval backend
2. Per-source files, only after step 1 tells you which ones are relevant to the query (see Process)

## Process

### 1. Read the Registry First
Load `shared/memory-registry.json`. This tells you which sources exist and whether any use a backend other
than `lexical` — currently none do (`lightrag` is `"status": "disabled"`), so every search in practice is a
markdown read, but check the registry rather than hardcoding that assumption, since it's meant to change
without this skill needing a rewrite.

### 2. Delegate the KI/ADR Portion
Invoke `search-ki` for the `knowledge-items` and `adrs` registry sources — don't reimplement its pre-filter
and judgment-based ranking here, that would be exactly the duplicated-logic problem `memory-engineer` exists
to catch elsewhere.

### 3. Search the Remaining Sources Directly
For `feature-archive`, `domain-dictionary`, and `team-topology` (the sources `search-ki` doesn't cover):
- **feature-archive**: grep `docs/features/*/retrospective.md` and `docs/features/*/analysis.md` for
  relevant content — same judgment-based reading `search-ki` applies, not a keyword-only match.
- **domain-dictionary**: check if the query matches a defined term directly — if so, that's usually the
  most authoritative answer available.
- **team-topology**: only relevant for ownership/crossing questions, not general knowledge queries.

### 4. If a Non-Lexical Backend Is Ever Active
Per the registry's `lightrag` entry: any result retrieved through it MUST be verified against the source
markdown file before being presented as fact — never surface a LightRAG result standalone. This step is
inert today (nothing uses that backend yet) but is not optional once something does.

### 5. Rank and Cap
Combine results from all sources, rank by actual relevance to the question (not by which source it came
from), and cap at 5 total — same discipline as `search-ki`'s own cap, for the same reason.

## Output Format

```markdown
## Memory Query: "[query]"

### Matches
- [Result title] ([source: KI | ADR | Feature Archive | Domain Dictionary]) -- [one-line why it matches] -- [file link]

### No Match
[If nothing matches: "No memory source covers this — consider `create-ki` or `promote-memory` after this
task if the answer turns out to be reusable."]
```

## Guardrails
- **Never** duplicate `search-ki`'s KI/ADR logic — delegate to it.
- **Never** fabricate a match — say "no match" plainly, same as `search-ki`.
- **Never** present a non-lexical-backend result without markdown verification, if/when that backend is ever active.
- Cap at 5 results total across all sources combined, not 5 per source.

## Standalone Mode
Pure local file reads plus a delegated call to `search-ki` (also pure local reads). No external calls, no
embeddings, today. If `lightrag` is ever enabled per the registry, this skill degrades gracefully back to
lexical-only if that backend is unavailable — markdown is always the fallback, never the other way around.

---
*Part of the [ai-assistant-dot-files](https://github.com/orieken/ai-assistant-dot-files) Context Engineering Framework by Oscar Rieken — licensed under [CC BY 4.0](https://github.com/orieken/ai-assistant-dot-files/blob/main/LICENSE-CONTENT.md). If you copy or adapt this file, please keep this attribution.*
