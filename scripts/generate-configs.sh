#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_DIR="$REPO_DIR/shared"

OUTPUT_DIR=""
DRY_RUN=false
PLATFORM_FILTER=""

usage() {
  cat <<'EOF'
Usage: generate-configs.sh [OPTIONS]

Generate platform-specific config files from the canonical shared/ layer.

Options:
  --output <dir>    Write generated files to <dir> (default: repo root)
  --platform <name> Only generate for a specific platform
  --dry-run         Show what would be generated without writing
  -h, --help        Show this help

Platforms: claude-code, cursor, windsurf, github-copilot, gemini, openai-codex
EOF
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --output)    OUTPUT_DIR="$2"; shift 2 ;;
    --platform)  PLATFORM_FILTER="$2"; shift 2 ;;
    --dry-run)   DRY_RUN=true; shift ;;
    -h|--help)   usage ;;
    *)           echo "Unknown option: $1"; usage ;;
  esac
done

OUTPUT_DIR="${OUTPUT_DIR:-$REPO_DIR}"

ok()   { echo "  [ok] $1"; }
dry()  { echo "  [dry-run] would generate $1"; }
skip() { echo "  [skip] $1"; }

should_generate() {
  [[ -z "$PLATFORM_FILTER" || "$1" == "$PLATFORM_FILTER" ]]
}

write_file() {
  local dest="$1"
  local content="$2"

  if $DRY_RUN; then
    dry "$dest"
    return
  fi

  mkdir -p "$(dirname "$dest")"
  echo "$content" > "$dest"
  ok "$dest"
}

collect_rules() {
  local result=""
  for rule_file in "$SHARED_DIR/rules/"*.md; do
    result+=$'\n'
    result+="$(cat "$rule_file")"
    result+=$'\n'
  done
  echo "$result"
}

collect_agent_roster() {
  local result=""
  result+=$'\n'"# Persona Roster"$'\n'
  result+=$'\n'"The following specialized personas are available. Invoke them by name when you need domain-specific expertise. Note: on this platform these are personas — context frames with no tool access or autonomous pipeline participation, per \`DOMAIN_DICTIONARY.md\`. Full multi-step agent orchestration is only available on Tier 1 (Claude Code)."$'\n'

  for agent_file in "$SHARED_DIR/agents/"*.md; do
    local agent_name agent_desc
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    agent_desc=$(grep '^description:' "$agent_file" | head -1 | sed 's/description: *//' || true)
    if [[ -n "$agent_name" ]]; then
      result+=$'\n'"- **$agent_name**: $agent_desc"
    fi
  done
  echo "$result"
}

extract_rule_content() {
  local file="$1"
  cat "$file"
}

extract_agent_body() {
  local file="$1"
  # Agent files have a preamble line (referencing .claude/rules/*.md — a file reference Cursor can't
  # follow) before the frontmatter, then the frontmatter itself, then the actual persona body. Skip
  # everything through the second '---' delimiter; the always-apply rule .mdc files already cover the
  # preamble's content, so it isn't lost, just not redundantly repeated per persona.
  awk '/^---$/{delim++; next} delim==2{print}' "$file"
}

testing_rules_body() {
  echo "# Testing Rules

## Saturday Framework (E2E / UI Testing)
ALWAYS use the Site-Centric pattern: \`BaseSite\`, \`BasePage\`, \`BaseElement\`, \`BaseFlow\`.
NEVER use traditional Page Object Model (POM).
ALWAYS use Playwright driven by Cucumber.js for UI automation.
ALWAYS include OpenTelemetry instrumentation for every BDD scenario.

## Sunday Framework (API Testing)
ALWAYS use Vitest for unit tests and Playwright for integration/E2E API tests.
ALWAYS use the custom \`api\` fixture and fluent matchers (\`toHaveStatus\`, \`toBeSuccessful\`, \`toRespondWithin\`).
ALWAYS extend \`BaseApiClient\` for domain-specific API clients.
ALWAYS validate schemas with Zod (\`validateSchema()\`).
NEVER use custom retry loops — use \`CircuitBreaker\` or \`ExponentialBackoffStrategy\`.

## Test Quality
CRITICAL: Test coverage MUST be >= 85%.
CRITICAL: Cyclomatic complexity per function MUST be < 7.
ALWAYS practice TDD/BDD — Red-Green-Refactor.
NEVER write feature code without tests."
}

go_backend_rules_body() {
  echo "# Go Backend Conventions

ALWAYS follow Clean Architecture layers: Entities → Use Cases → Adapters → Frameworks.
NEVER let domain entities import adapter or framework packages.
ALWAYS define interfaces in the use-case layer, implement in adapters.
ALWAYS use structured logging with low-cardinality message strings.
NEVER use \`any\` or \`interface{}\` — use typed interfaces.
ALWAYS handle errors explicitly — no silent swallows.
ALWAYS set explicit timeouts on network calls.
NEVER use raw SQL without parameterized queries.
ALWAYS use the expand/contract pattern for database migrations."
}

vue_frontend_rules_body() {
  echo "# Vue 3 + Tailwind Frontend Conventions

ALWAYS use Vue 3 Composition API with \`<script setup>\`.
NEVER use Options API in new components.
ALWAYS use Tailwind CSS utility classes — no custom CSS unless absolutely necessary.
ALWAYS extract reusable UI into composables (\`use*.ts\`) or components.
NEVER put business logic in components — extract to composables or services.
ALWAYS use TypeScript with strict mode.
NEVER use \`any\` types — use \`unknown\` with Zod validation at boundaries.
ALWAYS co-locate component tests alongside components.
CRITICAL: Components MUST be < 100 lines. Extract when larger."
}

# Broad glob covering common source file extensions across this framework's supported stacks
# (TypeScript, Go, Python, Java) — used to Auto Attach rules that should apply whenever code is
# being edited, without paying their token cost on every single request (Cursor's own guidance:
# combined alwaysApply rules should stay under ~2,000 tokens; architecture.mdc alone exceeds that).
CODE_FILE_GLOBS='["**/*.ts", "**/*.tsx", "**/*.js", "**/*.jsx", "**/*.go", "**/*.py", "**/*.java"]'

generate_mdc() {
  local dest="$1"
  local description="$2"
  local always_apply="$3"
  local globs="$4"
  local body="$5"

  local frontmatter="---"$'\n'
  frontmatter+="description: \"$description\""$'\n'
  frontmatter+="alwaysApply: $always_apply"$'\n'
  if [[ -n "$globs" ]]; then
    frontmatter+="globs: $globs"$'\n'
  fi
  frontmatter+="---"$'\n'

  write_file "$dest" "${frontmatter}${body}"
}

generate_cursor_personas() {
  local rules_dir="$1"
  local persona_count=0

  for agent_file in "$SHARED_DIR/agents/"*.md; do
    local base
    base="$(basename "$agent_file")"
    [[ "$base" == "CHANGELOG.md" ]] && continue

    local agent_name agent_desc agent_desc_safe
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    agent_desc=$(grep '^description:' "$agent_file" | head -1 | sed 's/description: *//' || true)
    [[ -z "$agent_name" ]] && continue

    # Escape embedded double quotes — several agent descriptions quote example user phrases
    # (e.g. dependency-auditor: `"audit dependencies"`), which would otherwise terminate the YAML
    # frontmatter's description string early and produce an invalid .mdc file.
    agent_desc_safe=$(printf '%s' "$agent_desc" | sed 's/"/\\"/g')

    generate_mdc "$rules_dir/${agent_name}.mdc" \
      "Persona: $agent_name — $agent_desc_safe" \
      "false" "" \
      "$(extract_agent_body "$agent_file")"
    ((persona_count++)) || true
  done

  echo "  ($persona_count persona files generated)"
}

generate_cursor() {
  echo ""
  echo "--- Cursor (Tier 2: Personas + Rules) ---"

  local rules_dir="$OUTPUT_DIR/.cursor/rules"

  # Only approval-gates and agent-roster stay alwaysApply: true — both are cheap (a few KB) and
  # safety/awareness-critical regardless of what the user is doing (chatting, planning, or coding).
  # architecture.mdc and design-principles.mdc are large (architecture.mdc alone inlines the full
  # ARCHITECTURE_RULES.md) and code-editing-specific, so they're Auto Attached via CODE_FILE_GLOBS
  # instead — loaded when actually touching source, not on every single request. Cursor's own
  # guidance: combined alwaysApply rules should stay under ~2,000 tokens; the previous all-four-
  # always-apply setup was 2-3x over that.
  generate_mdc "$rules_dir/architecture.mdc" \
    "Architecture guardrails — Clean Architecture, SOLID, no hardcoded secrets, expand/contract migrations" \
    "false" "$CODE_FILE_GLOBS" \
    "$(extract_rule_content "$SHARED_DIR/rules/architecture-guardrails.md")

$(extract_rule_content "$SHARED_DIR/ARCHITECTURE_RULES.md")"

  generate_mdc "$rules_dir/design-principles.mdc" \
    "Design principles — Kent Beck simple design, Fowler refactoring, Sandi Metz limits, Boy Scout Rule" \
    "false" "$CODE_FILE_GLOBS" \
    "$(extract_rule_content "$SHARED_DIR/rules/design-principles.md")"

  generate_mdc "$rules_dir/approval-gates.mdc" \
    "Approval gates — NEVER skip human checkpoints for commits, deploys, migrations, external API calls" \
    "true" "" \
    "$(extract_rule_content "$SHARED_DIR/rules/approval-gates.md")"

  generate_mdc "$rules_dir/agent-roster.mdc" \
    "Persona roster — ALWAYS check this list before beginning specialized work" \
    "true" "" \
    "$(collect_agent_roster)"

  generate_mdc "$rules_dir/testing.mdc" \
    "Testing framework rules — Saturday (E2E) and Sunday (API) conventions" \
    "false" '["**/*.spec.*", "**/*.test.*", "**/*.feature", "**/steps/**"]' \
    "$(testing_rules_body)"

  generate_mdc "$rules_dir/go-backend.mdc" \
    "Go backend conventions — Clean Architecture, error handling, migrations" \
    "false" '["**/*.go", "**/go.mod", "**/go.sum"]' \
    "$(go_backend_rules_body)"

  generate_mdc "$rules_dir/vue-frontend.mdc" \
    "Vue 3 + Tailwind conventions — Composition API, strict typing, component limits" \
    "false" '["**/*.vue", "**/*.tsx", "**/*.jsx", "**/components/**"]' \
    "$(vue_frontend_rules_body)"

  generate_cursor_personas "$rules_dir"
}

collect_craftsmanship_section() {
  local result=""
  result+=$'\n'"## Craftsmanship Rules"$'\n'
  result+='You must **strictly adhere** to the patterns defined in `ARCHITECTURE_RULES.md` (Clean Architecture, DDD, GoF patterns, and micro-rules).'$'\n'
  result+='- **TDD/BDD First**: Drive design through testing. Feature code is incomplete without tests. Practice Red-Green-Refactor.'$'\n'
  result+='- **Kent Beck (Simple Design)**: 1) Passes tests, 2) Reveals intention, 3) No duplication, 4) Fewest elements.'$'\n'
  result+='- **Martin Fowler (Refactoring)**: Use named refactoring operations (Extract Function, Inline Variable, etc.) instead of vague cleanups.'$'\n'
  result+='- **Architectural Constraints & Fitness Functions**: Enforce cyclomatic complexity `< 7` and functions `< 30` LOC.'$'\n'
  result+='- **The Boy Scout Rule**: Always leave the code cleaner than you found it.'$'\n'
  result+=$'\n'
  result+="## Tech Stack"$'\n'
  result+="- **Backend / MCP**: Go"$'\n'
  result+="- **Frontend**: Vue 3 + Tailwind CSS"$'\n'
  result+="- **Test Automation (Saturday Framework)**: TypeScript, Playwright, Cucumber.js, k6"$'\n'
  result+="- **API Testing (Sunday Framework)**: Vitest, Playwright, Zod, CircuitBreaker"$'\n'
  echo "$result"
}

generate_cursorrules() {
  local dest="$OUTPUT_DIR/.cursorrules"
  local content=""

  content+="# Context Engineering Framework — All Rules"$'\n'
  content+=$'\n'
  content+="$(collect_rules)"
  content+=$'\n'
  content+="$(collect_craftsmanship_section)"
  content+=$'\n'
  content+="$(collect_agent_roster)"

  write_file "$dest" "$content"
}

generate_windsurfrules() {
  local dest="$OUTPUT_DIR/.windsurfrules"
  local content=""

  content+="# Context Engineering Framework — All Rules"$'\n'
  content+=$'\n'
  content+="$(collect_rules)"
  content+=$'\n'
  content+="$(collect_craftsmanship_section)"
  content+=$'\n'
  content+="$(collect_agent_roster)"

  write_file "$dest" "$content"
}

generate_tier3() {
  local platform_name="$1"
  local dest_path="$2"
  local header="$3"

  local content=""
  content+="$header"$'\n'
  content+=$'\n'
  content+="## AI Feature Team & Global Rules"$'\n'
  content+="You are part of the Saturday Multi-Agent Feature Team. Before beginning any complex task, architectural decision, or feature delivery, you MUST adhere to the rules below."$'\n'
  content+=$'\n'
  content+="$(collect_rules)"
  content+=$'\n'
  content+="## Craftsmanship Rules"$'\n'
  content+='You must **strictly adhere** to the patterns defined in `ARCHITECTURE_RULES.md` (Clean Architecture, DDD, GoF patterns, and micro-rules).'$'\n'
  content+='- **TDD/BDD First**: Drive design through testing. Feature code is incomplete without tests. Practice Red-Green-Refactor.'$'\n'
  content+='- **Kent Beck (Simple Design)**: 1) Passes tests, 2) Reveals intention, 3) No duplication, 4) Fewest elements.'$'\n'
  content+='- **Martin Fowler (Refactoring)**: Use named refactoring operations (Extract Function, Inline Variable, etc.) instead of vague cleanups.'$'\n'
  content+='- **Architectural Constraints & Fitness Functions**: Enforce cyclomatic complexity `< 7` and functions `< 30` LOC.'$'\n'
  content+='- **The Boy Scout Rule**: Always leave the code cleaner than you found it.'$'\n'
  content+=$'\n'
  content+="## Tech Stack"$'\n'
  content+="- **Backend / MCP**: Go"$'\n'
  content+="- **Frontend**: Vue 3 + Tailwind CSS"$'\n'
  content+="- **Test Automation**: TypeScript, Playwright, Cucumber.js, k6"$'\n'
  content+=$'\n'
  content+="$(collect_agent_roster)"

  write_file "$dest_path" "$content"
}

generate_instructions_md() {
  local dest="$1"
  local apply_to="$2"
  local body="$3"

  local frontmatter="---"$'\n'
  frontmatter+="applyTo: \"$apply_to\""$'\n'
  frontmatter+="---"$'\n'

  write_file "$dest" "${frontmatter}${body}"
}

generate_copilot_scoped_instructions() {
  local dest_dir="$OUTPUT_DIR/.github/instructions"

  # Path-scoped instructions coexist with and combine with copilot-instructions.md (the Tier 3
  # style file generated separately) — per GitHub's docs, both apply when a file matches. Mirrors
  # Cursor's testing/go-backend/vue-frontend .mdc files; applyTo uses comma-separated glob patterns
  # in a single quoted string, not an array like Cursor's `globs` field.
  generate_instructions_md "$dest_dir/testing.instructions.md" \
    '**/*.spec.*,**/*.test.*,**/*.feature' \
    "$(testing_rules_body)"

  generate_instructions_md "$dest_dir/go-backend.instructions.md" \
    '**/*.go,**/go.mod,**/go.sum' \
    "$(go_backend_rules_body)"

  generate_instructions_md "$dest_dir/vue-frontend.instructions.md" \
    '**/*.vue,**/*.tsx,**/*.jsx' \
    "$(vue_frontend_rules_body)"
}

generate_agents_md() {
  local dest_path="$OUTPUT_DIR/AGENTS.md"

  local content=""
  content+="# AGENTS.md"$'\n'
  content+=$'\n'
  content+="Cross-tool agent instructions, following the https://agents.md convention. Also the file Gemini"$'\n'
  content+="Antigravity is believed to read as its project-level rules (medium confidence — see the"$'\n'
  content+="\`gemini\` entry's notes in \`shared/platform-registry.json\` for the caveat)."$'\n'
  content+=$'\n'
  content+="## AI Feature Team & Global Rules"$'\n'
  content+="You are part of the Saturday Multi-Agent Feature Team. Before beginning any complex task, architectural decision, or feature delivery, you MUST adhere to the rules below."$'\n'
  content+=$'\n'
  content+="$(collect_rules)"
  content+=$'\n'
  content+="$(collect_craftsmanship_section)"
  content+=$'\n'
  content+="$(collect_agent_roster)"

  write_file "$dest_path" "$content"
}

echo ""
echo "Context Engineering Framework — Config Generator"
echo "================================================="
echo ""
echo "Source:  $SHARED_DIR"
echo "Output:  $OUTPUT_DIR"
echo "Dry run: $DRY_RUN"
if [[ -n "$PLATFORM_FILTER" ]]; then
  echo "Platform: $PLATFORM_FILTER"
fi
echo ""

GENERATED=0

if should_generate "cursor"; then
  generate_cursor
  generate_cursorrules
  agent_persona_count=$(find "$SHARED_DIR/agents" -name "*.md" -not -name "CHANGELOG.md" | wc -l | tr -d ' ')
  ((GENERATED += 8 + agent_persona_count))
fi

if should_generate "windsurf"; then
  generate_windsurfrules
  ((GENERATED++))
fi

if should_generate "github-copilot"; then
  echo ""
  echo "--- GitHub Copilot (Tier 2: Personas + Rules) ---"
  generate_tier3 "github-copilot" \
    "$OUTPUT_DIR/.github/copilot-instructions.md" \
    "# Copilot Instructions (Saturday Framework)"
  generate_copilot_scoped_instructions
  ((GENERATED += 4))
fi

if should_generate "gemini"; then
  echo ""
  echo "--- Gemini / Antigravity (Tier 3, hybrid — see platform-registry.json notes) ---"
  generate_tier3 "gemini" \
    "$OUTPUT_DIR/.gemini/antigravity/instructions.md" \
    "# Antigravity Instructions (Saturday Framework)"
  generate_agents_md
  ((GENERATED += 2))
fi

if should_generate "openai-codex"; then
  echo ""
  echo "--- OpenAI / Codex (Tier 3: System Prompt) ---"
  generate_tier3 "openai-codex" \
    "$OUTPUT_DIR/.openai.md" \
    "# OpenAI / Codex Instructions (Saturday Framework)"
  ((GENERATED++))
fi

echo ""
echo "================================================="
echo "Generated $GENERATED config file(s)"
if $DRY_RUN; then
  echo "This was a dry run. No files were written."
fi
echo ""
