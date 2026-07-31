# SKILL.md Draft: summarize-session (under audit)

---
name: summarize-session
description: Summarizes the current session into a markdown handoff document
standalone: false
requires_mcp: ["memory-server"]
parameters:
  - name: output_path
    type: string
    required: false
    default: ".claude/session-summary.md"
  - name: detail_level
    type: string
    required: false
    valid_values: ["brief", "full"]
  - name: 42_invalid_param
    type: integer
    required: false
---

## Instructions

You are summarizing the current session. Use the `memory-server` MCP tool to retrieve
recent context items. Write the summary to `{{ output_path }}`.

If `detail_level` is "full", include all tool calls. If "brief", include only key decisions.
