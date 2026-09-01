# Token Optimization & Context Minimization Audit
**Date**: 2026-08-07  
**Scope**: Framework-wide Token Footprint, Multi-Platform System Prompt Analysis, and Context Minimization Protocols  
**Author**: Saturday Multi-Agent Feature Team  

---

## 1. Executive Summary

This audit evaluates the token usage, static prompt overhead, and context management protocols across the `ai-assistant-dot-files` Context Engineering Framework. 

While the framework provides advanced context management capabilities (`context-engineer`, `memory-compression`, `summarize-artifact`), the default platform generation scripts (`scripts/generate-configs.sh`) compile **monolithic rule files** (`AGENTS.md`, `.cursorrules`, `.openai.md`, `.windsurfrules`, `.roomodes`) that inject **19,000 to 58,000 tokens of static prompt overhead** into every single interaction.

By implementing **Stack-Scoped Rule Generation**, **Strict Always-Apply Rule Budgeting**, **Automated Token Auditing in `health-check.sh`**, and **Line-Bounded Tool Retrieval**, total baseline token consumption per turn can be reduced by **65% to 75%**.

---

## 2. Key Audit Findings & Context Bloat Sources

### 2.1 Monolithic Global Rule Bundles (~19,000–58,000 Input Tokens / Request)
- **`AGENTS.md` (Gemini / Antigravity)**: ~76.5 KB (~19,100 tokens). Concatenates Clean Architecture rules, approval gates, craftsmanship guidelines, persona roster, and **all 9 language conventions** (C#, Go, Java, Kotlin, Python, Rust, Swift, TypeScript, IaC).
- **`.cursorrules` (Cursor legacy)** & **`.openai.md` (Codex)**: ~76 KB (~19,000 tokens each).
- **`.roomodes` (Roo Code)**: ~232 KB (~58,000 tokens).
- **Impact**: Any ad-hoc chat, bug fix, or simple edit pays a multi-thousand token tax *before user input or code files are read*. Editing a single Python script still loads full C#, Java, Swift, Kotlin, and Go rule specifications.

### 2.2 Unbounded Tool Outputs & Context Pollution
- Uncapped shell outputs (`git log` without `-n`, `head`/`tail` omissions, full directory dumps) quickly fill context memory.
- File viewing without line-range constraints (`StartLine`/`EndLine`) loads 500+ line source files when only a 30-line function or interface signature is needed.

### 2.3 Subagent Context Leakage in Sequential Pipelines
- Multi-agent orchestrations (`deliver-feature`) can suffer context drift if intermediate subagents receive full prior conversation histories instead of strict isolated input contracts (`subagent-isolation-is-a-hard-boundary.md`).

---

## 3. Platform-Specific Optimization Protocols

### 3.1 Cursor Optimization Protocol
1. **Trim `alwaysApply` Rules**:
   - Limit `alwaysApply: true` strictly to lightweight safety/awareness rules: `approval-gates.mdc` and `agent-roster.mdc`.
   - Ensure `architecture.mdc`, `design-principles.mdc`, and language-specific rules use `alwaysApply: false` with targeted file globs (e.g. `globs: ["**/*.ts"]`) so they load only when editing matching files.
2. **Direct Subagent / Persona Tagging**:
   - In ad-hoc Cursor Chat (`Cmd+L`), directly `@`-tag the specific agent file (e.g. `#file:shared/agents/code-reviewer.md`) rather than invoking full multi-agent handoff contracts (`shared/contracts/`).

### 3.2 Gemini / Antigravity & Global (`AGENTS.md`) Protocol
1. **Stack-Scoped Config Generation**:
   - Update `scripts/generate-configs.sh` to accept `--stack <languages>` (e.g., `--stack go,typescript`).
   - Omit irrelevant language conventions from target project `AGENTS.md` files.
2. **Pre-flight Context Engineering**:
   - Always run `context-engineer` before multi-file feature tasks to generate a focused `context-manifest.md` with explicit line-range pins and prune out-of-context open tabs.

### 3.3 Memory & Knowledge Item (KI) Density
1. **Periodic Compression**:
   - Run `scripts/health-check.sh` and `install.sh --sync-memory` regularly to trigger `memory-compression` and consolidate stale Knowledge Items.
2. **High-Fidelity References**:
   - Pin high-fidelity code/schema/interface files instead of long prose summaries.

---

## 4. Proposed Action Plan

| Step | Action Item | Target File / Script | Expected Impact |
|---|---|---|---|
| 1 | Add `--stack` filter to config generator | [generate-configs.sh](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/scripts/generate-configs.sh) | Reduces `AGENTS.md` & `.cursorrules` from ~19k to ~4k tokens |
| 2 | Add token footprint audit pass | [health-check.sh](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/scripts/health-check.sh) | Warns when baseline static rules exceed 5k tokens (20 KB) |
| 3 | Enforce line-range bounds in tools | [AGENTS.md](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/AGENTS.md) / [context-engineer](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/shared/skills/context-engineer/SKILL.md) | Prevents reading whole files > 300 lines unnecessarily |
| 4 | Map utility agents to light models | [model-defaults.yaml](file:///Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files/shared/model-defaults.yaml) | Routes low-complexity tasks to high-speed / lower-cost tiers |

---

## 5. Verification & Health Check Integration

- Execute `bash scripts/health-check.sh` to verify zero drift and validate token budgets.
- Run `bash scripts/generate-configs.sh --dry-run` to inspect generated config sizes across platforms.
