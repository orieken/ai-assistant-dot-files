# Portable `model_tier` Abstraction Implementation Plan

- **Status**: Proposed (Phase A)
- **Target Repo**: `/Users/oscarrieken/Projects/Rieken/ai-assistant-dot-files`
- **Date**: 2026-07-29

---

## 1. Executive Summary

This plan introduces a portable, platform-agnostic way for agents and counter-agents to declare what class of model capability they require (`light`, `default`, `heavy`) in frontmatter rather than pinning a concrete vendor model ID (`claude-haiku-4-5-20251001` or `claude-opus-5`) or blindly defaulting to `inherit`.

Concrete model IDs are resolved at install time via `shared/model-defaults.yaml` for supported platforms (e.g., Claude Code). For platforms that do not support per-agent model selection (Cursor, Copilot, Gemini Antigravity, Roo Code, Cline), the installer strips the field from exported agent definitions and emits an explicit warning log with instructions for setting platform-wide model defaults. User override files (such as `.claude/model-overrides.yaml`) take precedence over `shared/model-defaults.yaml`.

---

## 2. Tier Taxonomy

The framework defines 3 abstract model tiers:

| Tier | Purpose & Characteristics | Typical Agent Roles |
|---|---|---|
| `light` | Pattern-matching, rubric-scoring, structural validation, regex/schema scanning. Fast & cost-effective. | Read-only counter-auditors, prompt evaluators, rule auditors, documentation managers. |
| `default` | General feature production, code generation, refactoring, test writing, standard code reviews. Session default model (`inherit` on Claude Code). | Analyst, Developer, Code-Reviewer, QA Engineer, Tech Writer, DevOps, Spec Writer, SRE. |
| `heavy` | Deep architectural reasoning, multi-file synthesis, complex threat modeling, critical security audits. Highest capability model. | Architect, Security-Reviewer. |

---

## 3. Frontmatter Contract Change & Precedence

1. **`model_tier`**: REQUIRED field in agent frontmatter going forward. Allowed enum values: `["light", "default", "heavy"]`.
2. **`model`**: OPTIONAL field in agent frontmatter.
3. **Precedence Rules**:
   - **Claude Code**: If both `model` and `model_tier` are present, `model` takes precedence as an explicit vendor override. Otherwise, `model_tier` resolves via `shared/model-defaults.yaml`.
   - **Other Platforms** (Cursor, Copilot, Gemini, Roo Code, Cline): `model` is ignored; `model_tier` is evaluated by the installer. If the platform mapping in `shared/model-defaults.yaml` is `null`, the installer strips the field and logs a `WARN`.
   - **User Overrides**: User override configs (`.claude/model-overrides.yaml` or `.cursor/model-overrides.yaml`) are checked *before* `shared/model-defaults.yaml`. User overrides always take highest precedence.

---

## 4. Per-Platform Resolution Matrix

| Platform | `light` Resolution | `default` Resolution | `heavy` Resolution | Installer Action |
|---|---|---|---|---|
| **Claude Code** | `claude-haiku-4-5-20251001` | `inherit` | `claude-opus-5` | Writes `model: <resolved-id>` into exported `.claude/agents/*.md`. |
| **Cursor** | `null` | `null` | `null` | Strips `model`/`model_tier` field; emits warning pointing to Cursor settings. |
| **GitHub Copilot** | `null` | `null` | `null` | Strips `model`/`model_tier` field; emits warning pointing to Copilot settings. |
| **Gemini Antigravity** | `null` | `null` | `null` | Strips `model`/`model_tier` field; emits warning pointing to AGY settings. |
| **Roo Code** | `null` | `null` | `null` | Strips `model`/`model_tier` field; emits warning pointing to Roo Code settings. |
| **Cline** | `null` | `null` | `null` | Strips `model`/`model_tier` field; emits warning pointing to Cline settings. |

---

## 5. Rollout Strategy

1. **Phase A**: Design proposal & plan committed (`docs(aos): draft model-tier abstraction plan`). Pause for user approval.
2. **Phase B**: Update frontmatter schema (`shared/schemas/agent-frontmatter.schema.json`), contract doc (`shared/contracts/agent-frontmatter-contract.md`), create `shared/model-defaults.yaml`, and update `shared/templates/agent.template.md`.
3. **Phase C**: Backfill all existing agents in `shared/agents/` across 3 atomic commits (`light` auditors, `default` producers, `heavy` architects), including single-line `#` rationale comments. Pause for user approval on tier assignments.
4. **Phase D**: Update `install.sh` to resolve `model_tier` against `shared/model-defaults.yaml` (with user override precedence) and emit warnings for `null`-mapped platforms.
5. **Phase E**: Extend `scripts/health-check.sh` to validate `model_tier` (emitting `WARN` for one release cycle before upgrading to `FAIL`), and introduce the `model-tier-auditor` counter-agent.

---

## 6. Open Questions for User Approval

- Confirm taxonomy: Are `light`, `default`, and `heavy` sufficient, or is a specialized 4th tier (e.g. `vision`) needed? (Recommendation: 3 tiers are sufficient for current agents).
- Confirm rollout cadence: Warn on missing `model_tier` in `health-check.sh` for 1 release cycle before enforcing failure?
