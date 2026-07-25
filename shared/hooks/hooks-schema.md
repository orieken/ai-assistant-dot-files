# Hooks Configuration Schema

Hook files in `.claude/hooks/` or `shared/hooks/examples/` must conform to the following JSON/YAML schema:

```yaml
version: "1.0"
hooks:
  - id: "unique-hook-identifier"
    event: "on-event-name"
    enabled: true
    action:
      type: "agent" # Options: "agent", "skill", "script"
      target: "target-name"
      args:
        passContext: true
```

## Schema Specifications

- **`version`** (`string`, required): Schema version (e.g., `"1.0"`).
- **`hooks`** (`array`, required): List of hook definitions.
  - **`id`** (`string`, required): Unique slug.
  - **`event`** (`string`, required): Event name from catalog (`on-artifact-write`, `on-validation-pass`, `on-ki-created`, etc.).
  - **`enabled`** (`boolean`, optional): Defaults to `true`.
  - **`action`** (`object`, required):
    - **`type`** (`enum`: `["agent", "skill", "script"]`, required): Target type.
    - **`target`** (`string`, required): Target agent name (e.g., `knowledge-auditor`), skill name, or script path.
    - **`args`** (`object`, optional): Key-value parameters.
