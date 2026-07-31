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
    ((total_agents++)) || true

    if ! grep -q "$agent_name" "$file" 2>/dev/null; then
      ((missing_agents++)) || true
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
  for expected_mdc in architecture design-principles approval-gates agent-roster testing go-backend vue-frontend typescript-conventions python-conventions csharp-conventions java-conventions; do
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
  check_rule_content "$CURSOR_DIR/testing.mdc" "testing.mdc" "testing-conventions.md"
  check_rule_content "$CURSOR_DIR/go-backend.mdc" "go-backend.mdc" "go-conventions.md"
  check_rule_content "$CURSOR_DIR/typescript-conventions.mdc" "typescript-conventions.mdc" "typescript-conventions.md"
  check_rule_content "$CURSOR_DIR/python-conventions.mdc" "python-conventions.mdc" "python-conventions.md"
  check_rule_content "$CURSOR_DIR/csharp-conventions.mdc" "csharp-conventions.mdc" "csharp-conventions.md"
  check_rule_content "$CURSOR_DIR/java-conventions.mdc" "java-conventions.mdc" "java-conventions.md"
else
  miss ".cursor/rules/ directory"
fi

echo ""
echo "--- Cursor agents/skills symlinks ---"
for symlink in ".cursor/agents" ".cursor/skills"; do
  full_path="$REPO_DIR/$symlink"
  if [[ -L "$full_path" ]]; then
    target=$(readlink "$full_path")
    if [[ "$target" == *"shared"* ]]; then
      pass "$symlink -> $target"
    else
      fail "$symlink" "points to $target (expected shared/ directly, not a double-hop through .claude/)"
    fi
  elif [[ -d "$full_path" ]]; then
    fail "$symlink" "is a directory, not a symlink to shared/"
  else
    miss "$symlink"
  fi
done

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
for scoped_file in "testing" "go-backend" "vue-frontend" "typescript-conventions" "python-conventions" "csharp-conventions" "java-conventions"; do
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
echo "--- AGENTS.md (cross-tool convention, confirmed read by Gemini Antigravity 2026-07-02) ---"
if [[ -f "$REPO_DIR/AGENTS.md" ]]; then
  check_concept_coverage "$REPO_DIR/AGENTS.md" "AGENTS.md"
  check_agent_roster "$REPO_DIR/AGENTS.md" "AGENTS.md"
else
  miss "AGENTS.md"
fi

echo ""
echo "--- Roo Code ---"
roomodes_file="$REPO_DIR/.roomodes"
if [[ -f "$roomodes_file" ]]; then
  if python3 -c "
import sys, yaml
data = yaml.safe_load(open('$roomodes_file'))
modes = data.get('customModes', [])
sys.exit(0 if modes else 1)
" 2>/dev/null; then
    mode_count=$(python3 -c "import yaml; d=yaml.safe_load(open('$roomodes_file')); print(len(d.get('customModes', [])))" 2>/dev/null || echo "?")
    pass ".roomodes ($mode_count custom modes, valid YAML)"
  else
    fail ".roomodes" "file exists but customModes array is empty or YAML is invalid"
  fi

  # Verify each shared agent has a corresponding mode
  missing_modes=0
  for agent_file in "$SHARED_DIR/agents/"*.md; do
    [[ "$(basename "$agent_file")" == "CHANGELOG.md" ]] && continue
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    [[ -z "$agent_name" ]] && continue
    if ! python3 -c "
import yaml, sys
data = yaml.safe_load(open('$roomodes_file'))
slugs = [m['slug'] for m in data.get('customModes', [])]
sys.exit(0 if '$agent_name' in slugs else 1)
" 2>/dev/null; then
      ((missing_modes++)) || true
    fi
  done
  if [[ $missing_modes -eq 0 ]]; then
    pass ".roomodes agent roster complete"
  else
    fail ".roomodes agent roster" "$missing_modes agents missing — regenerate with: scripts/generate-configs.sh --platform roo-code"
  fi
else
  miss ".roomodes (run: scripts/generate-configs.sh --platform roo-code)"
fi

roo_rules_dir="$REPO_DIR/.roo/rules"
if [[ -d "$roo_rules_dir" ]]; then
  rule_count=$(find "$roo_rules_dir" -name "*.md" | wc -l | tr -d ' ')
  shared_rule_count=$(find "$SHARED_DIR/rules" -maxdepth 1 -name "*.md" | wc -l | tr -d ' ')
  if [[ "$rule_count" -eq "$shared_rule_count" ]]; then
    pass ".roo/rules/ ($rule_count rule files)"
  else
    fail ".roo/rules/" "$rule_count files present, expected $shared_rule_count — regenerate"
  fi
else
  miss ".roo/rules/ (run: scripts/generate-configs.sh --platform roo-code)"
fi

echo ""
echo "--- Cline ---"
clinerules_dir="$REPO_DIR/.clinerules"
expected_cline_files=("00-approval-gates.md" "01-design-principles.md" "02-architecture-guardrails.md" "03-agent-roster.md" "04-testing-conventions.md" "05-go-conventions.md" "06-typescript-conventions.md" "07-python-conventions.md" "08-csharp-conventions.md" "09-java-conventions.md")
if [[ -d "$clinerules_dir" ]]; then
  missing_cline=0
  for expected in "${expected_cline_files[@]}"; do
    if [[ ! -f "$clinerules_dir/$expected" ]]; then
      fail ".clinerules/$expected" "missing"
      ((missing_cline++)) || true
    fi
  done

  # Check path-scoped files have paths: frontmatter
  for scoped_file in "04-testing-conventions.md" "05-go-conventions.md" "06-typescript-conventions.md" "07-python-conventions.md" "08-csharp-conventions.md" "09-java-conventions.md"; do
    full_path="$clinerules_dir/$scoped_file"
    if [[ -f "$full_path" ]]; then
      if head -1 "$full_path" | grep -q '^---$' && grep -q '^paths:' "$full_path"; then
        true
      else
        fail ".clinerules/$scoped_file" "missing paths: frontmatter for scope restriction"
        ((missing_cline++)) || true
      fi
    fi
  done

  if [[ $missing_cline -eq 0 ]]; then
    pass ".clinerules/ (${#expected_cline_files[@]} files, path-scoped rules verified)"
  fi

  # Check agent roster is included
  if grep -q 'analyst\|developer\|qa-engineer' "$clinerules_dir/03-agent-roster.md" 2>/dev/null; then
    pass ".clinerules/03-agent-roster.md contains agent list"
  else
    fail ".clinerules/03-agent-roster.md" "agent names not found — regenerate"
  fi
else
  miss ".clinerules/ (run: scripts/generate-configs.sh --platform cline)"
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
