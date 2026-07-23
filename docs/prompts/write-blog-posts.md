# Write blog posts about recent framework + agent updates

Author dev.to blog posts (with LinkedIn companions) covering material framework and agent developments. Not marketing — technical narrative with real numbers.

## Target repo

`/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files` — this IS the git repo. Drafts land in `docs/blog-posts/`.

## Existing convention (read first)

- **Index**: `docs/blog-posts/README.md` — lists prior posts (01 through 04) and the "Before Publishing" checklist
- **Publishing workflow**: `docs/blog-posts/publishing.md` — `scripts/publish-blog-posts.py` is the tool; dry-run friendly; explicit `--confirm send` required for API calls
- **Sample posts to match style/voice**:
  - `01-context-window-budget-devto.md` + `01-context-window-budget-linkedin.md` — "Treat the Context Window Like a Budget"
  - `02-memory-promotion-pipeline-devto.md` + companion — "Memory Engineering Is a Promotion Pipeline"
  - `03-forgotten-symlink-devto.md` + companion — "The Forgotten Symlink" (an incident story)
  - `04-aos-context-engineering-devto.md` + companion — "Toward an AI Operating System"

Read at least ONE full dev.to + linkedin pair before drafting, to match voice + structure.

## Convention

- **Numbering**: continue the sequence. Next post is `05-`. If drafting multiple topics in one session, use 05, 06, 07, etc.
- **Dual variant**: every topic has TWO files — one for dev.to (long-form, 800-1500 words) + one for LinkedIn (300-500 words, teases the dev.to post via `TODO_DEVTO_URL` placeholder for later replacement).
- **File naming**: `NN-topic-slug-devto.md` and `NN-topic-slug-linkedin.md`. Slug is short kebab-case (2-5 words).
- **Dev.to frontmatter** (copy from sample):
  ```yaml
  ---
  title: "<Title in title case, quoted>"
  published: false
  description: "<One-sentence hook, quoted>"
  tags: <3-5 comma-separated, no quotes — e.g., ai, productivity, architecture, devtools>
  canonical_url:
  cover_image:
  ---
  ```
- **LinkedIn variant**: shorter, first-person, ends with a link. Match sample 01's `01-context-window-budget-linkedin.md` structure.

## Voice and tone

- **Technical narrative, not marketing.** Grounded in real commits, real numbers, real files. Prior post 01 says "24 agents, 53 skills, 13 inter-agent contracts, 6 platform targets" — that specificity is the bar.
- **First-person Oscar.** "I wanted to treat the context window as a budget." "The framework has a dedicated context-engineer agent."
- **Show mechanics, not just outcomes.** Explain the code, the pattern, the tradeoff — not just "we improved things."
- **Craftsmanship framing.** Reference Kent Beck, Martin Fowler, Sandi Metz, Uncle Bob, Neal Ford where they naturally fit. Not name-drops — real applications of their disciplines.
- **Every post should end with a "what did we learn / what's next" beat**, matching the sample posts.

## Source material for new topics

Recent commits contain the raw material. Skim these before picking topics:

```bash
git log --oneline --since="2026-07-15" -- shared/ docs/
```

And in the neighboring saturday-monorepo:

```bash
cd /Users/oscarrieken/Projects/Rieken/saturday-monorepo && git log --oneline --since="2026-07-15" -- saturday-mcp/
```

Key artifacts to draw from:

- `docs/aos/migration-plan.md` — 4-phase AOS evolution with real checklists
- `docs/aos/AOS_Governance_Design_Pack/` — the vision docs (6 principles, 15 governance pairs)
- `docs/adrs/ADR-001-adopt-rag-friendly-docs-structure.md` (April 2026) + `ADR-002-corpus-aware-retrieval-strategy.md` (2026-07-22) — the two framework ADRs
- `docs/patterns/frontmatter-conventions.md` — the frontmatter conventions doc
- `saturday-monorepo/saturday-mcp/mcp-add-plan.md` — the 63-commit retrofit historical record
- `saturday-monorepo/saturday-mcp/mcp-expand-plan.md` — the framework-MCP expansion plan

## Suggested topics (pick 1-3 per session — don't try all at once)

Each is a topic that maps cleanly to a numbered draft.

### 05: The 63-Commit Retrofit — Trinity in a Real MCP Server
- **Story**: saturday-mcp's Handler was an 880-LOC God Object registering 17 tools inline as `map[string]interface{}`. 63 commits later it's a 93-LOC composite; every tool is a `domain.Tool` type; workflows compose via `WorkflowTool` adapter; adding a tool is now "drop a file in `internal/tools/`."
- **Real numbers**: 63 commits, 880 → 93 LOC (89% reduction), 26 files touched, 86.7% test coverage on tools package.
- **Mechanics**: the mcp-add skill's approval-gated plan artifact; subagent orchestration for mechanical Extract-Class ops; per-op commit discipline through session limits.
- **Source**: `saturday-monorepo/saturday-mcp/mcp-add-plan.md`, especially the Retrofit Complete section.

### 06: Templates Beat Prompts — Fixing Agent Output Structure Drift
- **Story**: every pipeline agent had its output format inlined in its own prompt as example markdown. Contracts in `shared/contracts/` checked structure separately. When LLM drift caused headings to shift ("Feature Summary" instead of "Summary"), audits started failing.
- **Fix**: extract templates as the single source of truth. 13 templates created, 13 agent prompts updated to reference them, all with structural verification against contracts (zero divergences found — the drift risk was real but hadn't hit yet).
- **Mechanics**: the extraction pattern applied uniformly; the "escape hatch" of putting AGENT_TEMPLATE outside `shared/agents/` so the loader doesn't register it as a real agent.
- **Source**: commit `6973541 feat(agents): extract output templates` + the follow-up `ca827dc fix(agents): move AGENT_TEMPLATE out of shared/agents/ loader path`.

### 07: Corpus-Aware RAG — Why One-Size-Fits-All Retrieval Is Wrong
- **Story**: adding RAG is the fashionable answer. But a framework's own KI corpus (30-200 items) doesn't need embeddings — LLM-as-retriever fits in a context window. A project's docs (50-500 files) is BM25 territory. A project's feature archive is vector territory. A project's source code is what Claude Code's Grep/Glob already do well.
- **Mechanics**: ADR-002's decision matrix; the graduated approach; MCP-native packaging so every client benefits equally; the "GEO the code files" prescriptive side (structured docs as the retrieval-quality multiplier).
- **Source**: `docs/adrs/ADR-002-corpus-aware-retrieval-strategy.md` + the sibling `ADR-001-adopt-rag-friendly-docs-structure.md` from April.

### 08: The Framework as an OS — Stage 3 Excellence Over Stage 4 Fantasies
- **Story**: adoption stages for AI coding — individual assistance → task augmentation → workflow orchestration → semi-autonomous → AI-native org. Most tools sit at stage 2. Most vendors chase stage 4. Deliberately capping at stage 3 excellence — human gates by policy, not by convention.
- **Mechanics**: AOS Vision's six principles (Capability, Governance, Learning, Memory Engineering, Context Engineering, Continuous Improvement); 15 governance pairs (11 audit + 4 opposing-force); the migration path (v3.0 → v3.3) that never breaks a base install.
- **Source**: `docs/aos/AOS_Governance_Design_Pack/00-AOS-Vision.md`, `01-Governance-Checks-and-Balances.md`, `docs/aos/migration-plan.md`.

### 09: Handoff Prompts as First-Class Artifacts
- **Story**: session limits, context bloat, and orchestration exhaustion are real. The pattern that survives them: check the next unit of work in as a markdown file under `docs/prompts/`. Fresh chat, fresh subagent, fresh install — anyone can pick it up cold.
- **Mechanics**: how a handoff prompt is structured (self-contained target, plan reference, scope, escalation criteria, report format); the discipline of drafting the prompt BEFORE starting the work; how prompts unlock async multi-agent orchestration.
- **Source**: `docs/prompts/README.md`, `docs/aos/prompts/README.md`, `saturday-monorepo/saturday-mcp/docs/prompts/README.md` — three parallel prompt directories all following the same convention.

### 10: Frontmatter Conventions Nobody Agrees On
- **Story**: YAML frontmatter in agent/skill/KI markdown is a de facto standard, but every vendor's shape is different. Claude Code uses one shape. Cursor's `.mdc` uses another. Windsurf and Copilot use flat markdown. There is no cross-tool standard, so the framework has to be its own source of truth and project into each platform.
- **Mechanics**: the frontmatter contracts, the JSON schemas, the IDE integration templates, the health-check validation, the platform generation pipeline (`scripts/generate-configs.sh`).
- **Source**: `docs/patterns/frontmatter-conventions.md`, `shared/schemas/*.schema.json`, `shared/platform-registry.json`.

## Discipline (non-negotiable)

- **One commit per topic** — a topic = one dev.to draft + one LinkedIn draft. Commit both files together with a message like `docs(blog): draft 05 — trinity in a real mcp server`.
- **Conventional commits** under `docs(blog):`.
- **NEVER `git add -A`** — this repo has known untracked directories (`docs/audits/`, `docs/aos/AOS_Governance_Design_Pack.zip`, `.gitignore` M) that must NOT be swept in. Always stage explicit paths: `git add docs/blog-posts/NN-topic-devto.md docs/blog-posts/NN-topic-linkedin.md`.
- **Update the README.md index** after each draft — add a new row to the drafts table.
- **Do NOT push.** Human reviews drafts before publishing.
- **Do NOT run `scripts/publish-blog-posts.py`.** That's human-triggered after review.
- **Every dev.to draft must include an image prompt** at the end (see sample 01 for format) — even if you generate no image, the prompt is needed for downstream commissioning.
- **LinkedIn drafts use `TODO_DEVTO_URL`** as placeholder — the human replaces it after the dev.to post goes live.

## Fact-checking (per the Before-Publishing checklist)

Before committing a draft, cross-check any specific claim against:
- The commit history (`git log --oneline -- <path>` for anything you cite)
- The referenced plan/ADR files
- Current file sizes / LOC counts (`wc -l <file>`)
- Current test coverage (`go test -cover ./...` if a Go claim)

If a claim can't be verified, either cite the source explicitly or hedge ("approximately", "around") — never invent specifics.

## Escalation criteria

Stop and report if:
- A topic doesn't have enough source material to write 800+ words honestly — skip it or ask for more context
- The convention conflicts (e.g., README.md format shifts between sample 01 and 04) — flag which sample you followed
- A number you'd cite doesn't verify against the repo — halt with the specific discrepancy

## Report (under 300 words)

```
STATUS: complete | stopped-at-<reason>

Topics drafted (list each with slug):
  - 05: <topic slug>
  - 06: <topic slug>
  ...

Commits landed:
  <sha> <message>
  ...

README.md index updated: yes | no
Image prompts included in every dev.to draft: yes | no
LinkedIn TODO_DEVTO_URL placeholders present: yes | no

Fact-check notes:
  - <any claim you had to hedge or verify from a non-obvious source>
```

Go.
