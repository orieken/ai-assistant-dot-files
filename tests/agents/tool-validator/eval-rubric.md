# Eval Rubric: tool-validator / input-skill-draft.md

- **Hidden MCP dependency is flagged**: `standalone: false` combined with `requires_mcp: ["memory-server"]` means the skill cannot run without the `memory-server` MCP — the auditor flags this as a hidden dependency that must be declared in the skill's documentation so callers know.
- **Invalid parameter name is flagged**: `42_invalid_param` starts with a digit, which is not a valid parameter name — the auditor identifies this specific parameter and explains the naming constraint violated.
- **`requires_mcp` frontmatter field is validated**: the auditor checks that `memory-server` is declared in the frontmatter `requires_mcp` array and is not just mentioned in the instructions body — in this case it IS in frontmatter, so the auditor reports this as compliant.
- **`standalone: false` declaration is accepted**: the skill explicitly declares `standalone: false` — the auditor does not flag this as a problem, only notes it means MCP is required.
- **No fabricated findings**: the auditor does not invent issues beyond the actual content — the `output_path` and `detail_level` parameters with valid names and types are not flagged.

## How to Grade
For each bullet, quote the specific line(s) of `actual-output.md` that satisfy it. If a bullet has no supporting quote, mark it FAIL and say what's missing.
