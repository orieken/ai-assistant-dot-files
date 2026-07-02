# Runbook: Adding a New Platform

Step-by-step guide for wiring a new AI coding tool into the `shared/` -> generated-config pipeline. See
[ARCHITECTURE.md](../ARCHITECTURE.md) section 2 for the tier system this fits into.

## 1. Decide the tier

| Question | Tier 1 (Full) | Tier 2 (Personas + Rules) | Tier 3 (System Prompt) |
|---|---|---|---|
| Native subagents with tool access? | Yes | No | No |
| Multiple rule files supported? | Yes | Yes | No — single file |
| Can it follow "read this file" references? | Yes | No — must inline | No — must inline |

If you're not sure, default to Tier 3 (single inlined instruction file) — it's the safest, most portable
option, and the tool almost certainly supports at least that much.

## 2. Register it in `shared/platform-registry.json`

Add an entry to the `platforms` array:
```json
{
  "name": "your-platform",
  "tier": 3,
  "tierLabel": "System Prompt",
  "capabilities": {
    "agents": false,
    "skills": false,
    "rules": false,
    "hooks": false,
    "subAgentOrchestration": false
  },
  "configPaths": {
    "systemPrompt": ".yourplatform.md"
  },
  "format": "markdown",
  "installStrategy": "generate-inline",
  "globalConfigDir": null,
  "notes": "Whatever's true and useful for the next person wiring this up."
}
```
Match the shape of an existing Tier 2 or Tier 3 entry as closely as your platform's real capabilities allow.

## 3. Add a generator function to `scripts/generate-configs.sh`

For Tier 3, this is usually one call to the existing `generate_tier3()` helper:
```bash
if should_generate "your-platform"; then
  echo ""
  echo "--- Your Platform (Tier 3: System Prompt) ---"
  generate_tier3 "your-platform" \
    "$OUTPUT_DIR/.yourplatform.md" \
    "# Your Platform Instructions (Saturday Framework)"
  ((GENERATED++))
fi
```
For Tier 2 with per-concern files, follow `generate_cursor()`'s pattern instead — one `generate_mdc()` call
per rule file, plus the persona roster via `collect_agent_roster()`.

**Never write ad-hoc content-collection logic** — reuse `collect_rules()`, `collect_craftsmanship_section()`,
and `collect_agent_roster()`. They're the single source of truth for what goes in every platform's output;
duplicating that logic per-platform is exactly the drift this framework exists to prevent.

## 4. Add a parity check to `scripts/check-parity.sh`

Add a section that verifies the generated file exists and contains the full persona roster:
```bash
if [[ -f "$OUTPUT_DIR/.yourplatform.md" ]]; then
  pass "Your Platform (.yourplatform.md)"
  check_agent_roster "$OUTPUT_DIR/.yourplatform.md" "Your Platform agent roster"
else
  miss "Your Platform (.yourplatform.md)"
fi
```

## 5. Wire it into `install.sh` / `uninstall.sh`

Add platform auto-detection (how do you tell the tool is actually installed on this machine?) and make sure
`--platform your-platform` is accepted as a filter value.

## 6. Test it end-to-end

```bash
scripts/generate-configs.sh --platform your-platform --dry-run
scripts/generate-configs.sh --platform your-platform
scripts/check-parity.sh
```
Then actually open the generated file in the target tool and confirm it loads without error — a `.mdc` file
with subtly invalid YAML frontmatter, for example, gets silently ignored by Cursor rather than erroring, so
"the script ran with no errors" is not the same as "the tool actually picked it up."

## 7. Document the terminology decision

If the platform is Tier 2/3, its generated persona roster must say "persona," never "agent" — this is
already handled by `collect_agent_roster()` for any platform using it, but if you write custom content
collection instead, make sure you match this rule (see `DOMAIN_DICTIONARY.md`'s Capability Tier definition
for why).
