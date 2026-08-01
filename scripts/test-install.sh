#!/usr/bin/env bash
set -uo pipefail

# Install verification matrix: runs a real `install.sh --project <scratch-dir> --platform <name>`
# for every platform in shared/platform-registry.json and asserts the expected symlinks/files exist
# and resolve correctly. Complements check-parity.sh (which validates this repo's own generated
# output) by validating what a FRESH install produces in a target project -- the two are not
# duplicates, they cover different failure modes (drift in this repo vs. a broken install path).
#
# Scoped to --project mode only, not --global -- --global would touch $HOME, which is fine inside
# ci-check.sh's ephemeral Docker container but not appropriate to also run against a developer's
# real home directory from a plain local run of this script.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRATCH_ROOT="$(mktemp -d)"

PASS_COUNT=0
FAIL_COUNT=0

pass() { echo "  PASS  $1"; PASS_COUNT=$((PASS_COUNT + 1)); }
fail() { echo "  FAIL  $1"; FAIL_COUNT=$((FAIL_COUNT + 1)); }

cleanup() { rm -rf "$SCRATCH_ROOT"; }
trap cleanup EXIT

assert_symlink() {
  local path="$1"
  local expected_substring="$2"
  local label="$3"

  if [[ ! -L "$path" ]]; then
    fail "$label -- not a symlink ($path)"
    return
  fi
  local target
  target=$(readlink "$path")
  if [[ "$target" == *"$expected_substring"* ]]; then
    pass "$label -> $target"
  else
    fail "$label -- points to $target, expected something containing $expected_substring"
  fi
}

assert_file() {
  local path="$1"
  local label="$2"

  if [[ -f "$path" ]]; then
    pass "$label"
  else
    fail "$label -- missing ($path)"
  fi
}

assert_install_marker() {
  local path="$1"
  local label="$2"
  if [[ ! -f "$path" ]]; then
    fail "$label -- marker missing ($path)"
    return
  fi
  if ! command -v python3 &>/dev/null; then
    pass "$label (exists; content validation skipped — no python3)"
    return
  fi
  if python3 -c "
import json, sys
d = json.load(open('$path'))
required = ['source_repo', 'git_tag', 'commit_sha', 'installed_at', 'mode', 'framework_level', 'platforms']
missing = [k for k in required if k not in d]
sys.exit(1 if missing else 0)
" 2>/dev/null; then
    pass "$label (valid JSON with required fields)"
  else
    fail "$label -- invalid JSON or missing required fields"
  fi
}

assert_frontmatter() {
  local path="$1"
  local label="$2"
  local required_field="$3"

  if [[ ! -f "$path" ]]; then
    fail "$label -- missing ($path)"
    return
  fi
  if head -1 "$path" | grep -q '^---$' && grep -q "^${required_field}:" "$path"; then
    pass "$label"
  else
    fail "$label -- missing valid frontmatter or $required_field field"
  fi
}

echo ""
echo "=== Install Verification Matrix ==="
echo "Scratch root: $SCRATCH_ROOT"
echo ""

echo "--- claude-code ---"
target="$SCRATCH_ROOT/claude-code"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform claude-code >/dev/null 2>&1
assert_symlink "$target/.claude/agents" "shared/agents" "claude-code .claude/agents"
assert_symlink "$target/.claude/skills" "shared/skills" "claude-code .claude/skills"
assert_symlink "$target/.claude/rules" "shared/rules" "claude-code .claude/rules"
assert_file "$target/ARCHITECTURE_RULES.md" "claude-code ARCHITECTURE_RULES.md"
assert_file "$target/DOMAIN_DICTIONARY.md" "claude-code DOMAIN_DICTIONARY.md"
assert_file "$target/CLAUDE.md" "claude-code CLAUDE.md template"
assert_install_marker "$target/.claude/framework-install.json" "claude-code framework-install.json"
echo ""

echo "--- cursor ---"
target="$SCRATCH_ROOT/cursor"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform cursor >/dev/null 2>&1
assert_symlink "$target/.cursor/agents" "shared/agents" "cursor .cursor/agents"
assert_symlink "$target/.cursor/skills" "shared/skills" "cursor .cursor/skills"
for mdc in architecture design-principles approval-gates agent-roster testing go-backend vue-frontend \
  typescript-conventions python-conventions csharp-conventions java-conventions; do
  assert_frontmatter "$target/.cursor/rules/${mdc}.mdc" "cursor ${mdc}.mdc" "description"
done
echo ""

echo "--- windsurf ---"
target="$SCRATCH_ROOT/windsurf"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform windsurf >/dev/null 2>&1
assert_file "$target/.windsurfrules" "windsurf .windsurfrules"
echo ""

echo "--- github-copilot ---"
target="$SCRATCH_ROOT/github-copilot"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform github-copilot >/dev/null 2>&1
assert_file "$target/.github/copilot-instructions.md" "copilot copilot-instructions.md"
for scoped in testing go-backend vue-frontend typescript-conventions python-conventions csharp-conventions java-conventions; do
  assert_frontmatter "$target/.github/instructions/${scoped}.instructions.md" "copilot ${scoped}.instructions.md" "applyTo"
done
echo ""

echo "--- gemini ---"
target="$SCRATCH_ROOT/gemini"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform gemini >/dev/null 2>&1
assert_symlink "$target/.agents/skills" "shared/skills" "gemini .agents/skills"
assert_symlink "$target/.agents/rules" "shared/rules" "gemini .agents/rules"
assert_file "$target/AGENTS.md" "gemini AGENTS.md"
echo ""

echo "--- openai-codex ---"
target="$SCRATCH_ROOT/openai-codex"
mkdir -p "$target"
"$REPO_DIR/install.sh" --project "$target" --platform openai-codex >/dev/null 2>&1
assert_file "$target/.openai.md" "openai-codex .openai.md"
echo ""

echo "==========================================="
echo "Results: $PASS_COUNT passed, $FAIL_COUNT failed"
if [[ $FAIL_COUNT -gt 0 ]]; then
  exit 1
fi
exit 0
