Long-running AI coding sessions inevitably degrade as context windows get heavy with stale terminal output and old stack traces.

Instead of fighting prompt decay in a heavy chat, we check the next unit of work into git as a **First-Class Handoff Prompt** in `docs/prompts/`.

Every handoff prompt is a self-contained markdown file containing:
1. Target repo path
2. Context files to read first
3. Explicit scope of tasks
4. Guardrails & Conventional Commit rules
5. Escalation criteria
6. Standardized report format

A fresh AI agent (or developer) can pick up the prompt in a clean context window and execute with 100% precision.

When you treat prompts like version-controlled source code, context exhaustion disappears.

Full breakdown: TODO_DEVTO_URL
