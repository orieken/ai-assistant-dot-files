---
name: context-engineer
description: Optimizes the agent's context window by pruning unrelated files, identifying relevant Knowledge Items (KIs), and compiling a precise context manifest.
triggers:
  keywords: ["optimize-context", "prune-context", "context-engineering", "context", "manifest"]
  intentPatterns: ["/optimize-context", "Optimize my context", "Help me build context for *", "Clean up context"]
standalone: true
---

## When To Use
Use this skill when:
- You are about to start a new feature execution or a complex bug fix.
- The active context is crowded with unrelated files or terminal outputs.
- You want to ensure you are aligning with the latest Knowledge Items (KIs) and Architecture Decision Records (ADRs) without loading unnecessary code.

Do NOT use when:
- Performing simple, single-file edits that do not require architectural planning.

## Context To Load First
1. [context-engineering.md](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/docs/runbooks/context-engineering.md)
2. [ARCHITECTURE_RULES.md](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/ARCHITECTURE_RULES.md)
3. [DOMAIN_DICTIONARY.md](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/DOMAIN_DICTIONARY.md)

## Process

### 1. Identify Target Component Scope
Identify which layers, files, and domains are relevant to the requested task. Classify them using the Sunday/Saturday Clean Architecture boundary guidelines:
- Domain/Entity (inner-most)
- Application/Use Case
- Adapter/Presenter/Controller
- Infrastructure/Framework/UI (outer-most)

### 2. Audit Current Workspace Context
Analyze what context is currently loaded:
- List all open documents in the session.
- Identify files that are **out of scope** (e.g., from unrelated bounded contexts or different layers).
- Request the user or IDE to close/unload out-of-scope files.

### 3. Retrieve Domain Knowledge (Proactive RAG)
Invoke `search-ki` with the target component/domain, **before** doing independent analysis — it already
searches `shared/knowledge/`, `.claude/knowledge/`, and `docs/adrs/` and ranks results, so don't duplicate
that scan here. Separately, identify key interfaces/types that define the contract of the target component
(that's a codebase lookup, not a KI search).

### 4. Estimate the Token Budget
For each pinned file, estimate tokens (~line count × 8 chars/line ÷ 4 chars/token — a rough heuristic). Sum the total and compare against the consuming agent's tier budget (of a 200k-token context window): Analyst/Architect ≤60%, Developer ≤80%, Reviewer agents ≤40%. Flag `WARNING` if over budget and recommend specific cuts.

### 5. Compile a Context Manifest
Generate a concise `context-manifest.md` in the current feature workspace (e.g., `.claude/feature-workspace/context-manifest.md` or output directly). 

## Output Format

The skill produces a `context-manifest.md` matching this structure:

```markdown
# Context Manifest: [Task/Feature Name]

## 1. Scope and Boundaries
- **Target Component**: [e.g., billing-service, user-auth]
- **Relevant Layers**: [e.g., Domain Entities, Application Use Cases]
- **Bounded Context**: [e.g., Identity & Access]

## 2. Pinpoint Files (To Keep Open)
List specific files and line ranges that MUST be kept in active memory:
- [File Basename](file:///absolute/path/to/file#L1-L50) -- [Purpose/Contract]
- [File Basename](file:///absolute/path/to/file#L100-L150) -- [Implementation details]

## 3. Global Rules and Constraints
List reference files that establish the patterns:
- [ARCHITECTURE_RULES.md](file:///absolute/path/to/ARCHITECTURE_RULES.md)
- [DOMAIN_DICTIONARY.md](file:///absolute/path/to/DOMAIN_DICTIONARY.md)

## 4. Knowledge Items & ADRs (To Load)
- [KI Name](file:///absolute/path/to/ki/artifact) -- [Summary of relevance]
- [ADR Name](file:///absolute/path/to/adr) -- [Decision context]

## 5. Prune Recommendations (To Close)
List files currently open that should be closed immediately:
- `[ ]` [File Basename](file:///absolute/path/to/file)
- `[ ]` [File Basename](file:///absolute/path/to/file)

## 6. Token Budget
- **Estimated total tokens for pinned files**: ~<N>
- **Target agent tier**: [Analyst/Architect: ≤60% | Developer: ≤80% | Reviewer: ≤40%] of a 200k-token context window
- **Status**: OK | WARNING (exceeds tier budget — see cut recommendations below)
- **Cut recommendations (if WARNING)**: [file] -- [reason it's the lowest-signal pin]
```

## Guardrails
- **No directory dumps**: Do not include entire directories in the manifest. Specify files explicitly.
- **Limit line count**: Files over 500 lines must be referenced with specific line ranges.
- **Rule alignment**: Never recommend files that violate Clean Architecture dependency boundaries (e.g. loading infrastructure database models in the Domain context manifest).
- **Never** report a token budget as OK without having actually estimated it.
