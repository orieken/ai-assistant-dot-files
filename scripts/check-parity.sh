#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_DIR="$REPO_DIR/shared"

has_drift=0

echo ""
echo "=== Cross-Platform Parity Check ==="
echo "Canonical source: $SHARED_DIR"
echo ""

pass() { echo "  PASS  $1"; }
fail() { echo "  DRIFT $1"; echo "        $2"; has_drift=1; }
miss() { echo "  MISS  $1"; has_drift=1; }

CORE_CONCEPTS=(
  "cyclomatic complexity"
  "< 7"
  "Clean Architecture"
  "SOLID"
  "TDD"
  "Boy Scout"
  "Expand/Contract"
  "Saturday Framework"
  "Sunday Framework"
  "fitness function"
)

check_concept_coverage() {
  local file="$1"
  local label="$2"
  local missing=()

  for concept in "${CORE_CONCEPTS[@]}"; do
    if ! grep -qi "$concept" "$file" 2>/dev/null; then
      missing+=("$concept")
    fi
  done

  if [[ ${#missing[@]} -eq 0 ]]; then
    pass "$label"
  else
    fail "$label" "missing: ${missing[*]}"
  fi
}

check_agent_roster() {
  local file="$1"
  local label="$2"
  local missing_agents=0
  local total_agents=0

  for agent_file in "$SHARED_DIR/agents/"*.md; do
    local agent_name
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    if [[ -z "$agent_name" ]]; then continue; fi
    ((total_agents++))

    if ! grep -q "$agent_name" "$file" 2>/dev/null; then
      ((missing_agents++))
    fi
  done

  if [[ $missing_agents -eq 0 ]]; then
    pass "$label agent roster ($total_agents agents)"
  else
    fail "$label agent roster" "$missing_agents of $total_agents agents missing"
  fi
}

check_rule_content() {
  local file="$1"
  local label="$2"
  local rule_name="$3"
  local rule_file="$SHARED_DIR/rules/$rule_name"

  if [[ ! -f "$rule_file" ]]; then return; fi

  local rule_heading
  rule_heading=$(head -1 "$rule_file")

  if grep -qF "$rule_heading" "$file" 2>/dev/null; then
    pass "$label contains $rule_name"
  else
    fail "$label" "missing content from $rule_name"
  fi
}

echo "--- Cursor .mdc files ---"
CURSOR_DIR="$REPO_DIR/.cursor/rules"
if [[ -d "$CURSOR_DIR" ]]; then
  for expected_mdc in architecture design-principles approval-gates agent-roster testing go-backend vue-frontend; do
    local_file="$CURSOR_DIR/${expected_mdc}.mdc"
    if [[ -f "$local_file" ]]; then
      if ! head -1 "$local_file" | grep -q '^---$'; then
        fail "$expected_mdc.mdc" "missing YAML frontmatter"
      else
        pass "$expected_mdc.mdc"
      fi
    else
      miss "$expected_mdc.mdc"
    fi
  done

  check_rule_content "$CURSOR_DIR/architecture.mdc" "architecture.mdc" "architecture-guardrails.md"
  check_rule_content "$CURSOR_DIR/design-principles.mdc" "design-principles.mdc" "design-principles.md"
  check_rule_content "$CURSOR_DIR/approval-gates.mdc" "approval-gates.mdc" "approval-gates.md"
  check_agent_roster "$CURSOR_DIR/agent-roster.mdc" "agent-roster.mdc"

  persona_missing=0
  persona_total=0
  for agent_file in "$SHARED_DIR/agents/"*.md; do
    agent_base="$(basename "$agent_file")"
    [[ "$agent_base" == "CHANGELOG.md" ]] && continue
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    [[ -z "$agent_name" ]] && continue
    ((persona_total++)) || true
    if [[ ! -f "$CURSOR_DIR/${agent_name}.mdc" ]]; then
      ((persona_missing++)) || true
    fi
  done
  if [[ $persona_missing -eq 0 ]]; then
    pass "per-agent persona files ($persona_total/$persona_total present)"
  else
    fail "per-agent persona files" "$persona_missing of $persona_total missing"
  fi
else
  miss ".cursor/rules/ directory"
fi

echo ""
echo "--- Flat rule files ---"
for flat_file in ".cursorrules" ".windsurfrules"; do
  full_path="$REPO_DIR/$flat_file"
  if [[ -f "$full_path" ]]; then
    check_concept_coverage "$full_path" "$flat_file"
    check_agent_roster "$full_path" "$flat_file"
  else
    miss "$flat_file"
  fi
done

echo ""
echo "--- Tier 3 system prompts ---"
TIER3_FILES=(
  ".github/copilot-instructions.md:Copilot"
  ".gemini/antigravity/instructions.md:Gemini"
  ".openai.md:OpenAI"
)

for entry in "${TIER3_FILES[@]}"; do
  file_path="${entry%%:*}"
  label="${entry##*:}"
  full_path="$REPO_DIR/$file_path"

  if [[ -f "$full_path" ]]; then
    check_concept_coverage "$full_path" "$label ($file_path)"
    check_agent_roster "$full_path" "$label"
  else
    miss "$label ($file_path)"
  fi
done

echo ""
echo "--- GitHub Copilot scoped instructions ---"
for scoped_file in "testing" "go-backend" "vue-frontend"; do
  full_path="$REPO_DIR/.github/instructions/${scoped_file}.instructions.md"
  if [[ -f "$full_path" ]]; then
    if head -1 "$full_path" | grep -q '^---$' && grep -q '^applyTo:' "$full_path"; then
      pass "${scoped_file}.instructions.md"
    else
      fail "${scoped_file}.instructions.md" "missing applyTo frontmatter"
    fi
  else
    miss "${scoped_file}.instructions.md"
  fi
done

echo ""
echo "--- AGENTS.md (cross-tool convention, also read by Gemini Antigravity per current research) ---"
if [[ -f "$REPO_DIR/AGENTS.md" ]]; then
  check_concept_coverage "$REPO_DIR/AGENTS.md" "AGENTS.md"
  check_agent_roster "$REPO_DIR/AGENTS.md" "AGENTS.md"
else
  miss "AGENTS.md"
fi

echo ""
echo "--- Claude Code symlinks ---"
for symlink in ".claude/agents" ".claude/skills" ".claude/rules"; do
  full_path="$REPO_DIR/$symlink"
  if [[ -L "$full_path" ]]; then
    target=$(readlink "$full_path")
    if [[ "$target" == *"shared"* ]]; then
      pass "$symlink -> $target"
    else
      fail "$symlink" "points to $target (expected shared/)"
    fi
  elif [[ -d "$full_path" ]]; then
    fail "$symlink" "is a directory, not a symlink to shared/"
  else
    miss "$symlink"
  fi
done

echo ""
echo "==========================================="
if [[ $has_drift -eq 0 ]]; then
  echo "Result: All platform configs in sync with shared/ canonical source"
else
  echo "Result: DRIFT detected. Run './scripts/generate-configs.sh' to regenerate."
fi
echo ""

exit $has_drift
