#!/usr/bin/env bash
set -euo pipefail

# Golden-file structural tests for pipeline agent prompts (shared/agents/).
#
# This script does NOT invoke any LLM agent itself — generating actual-output.md
# is a manual step (see tests/agents/README.md). It only validates a
# previously-generated actual-output.md against:
#   1. expected-patterns.txt — fuzzy, scenario-specific patterns
#   2. the agent's contract in shared/contracts/ — required section headings, if one exists
#
# Run this after editing any agent prompt in shared/agents/.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TESTS_DIR="$REPO_DIR/tests/agents"
CONTRACTS_DIR="$REPO_DIR/shared/contracts"

# bash 3.2 (macOS default) has no associative arrays — use a function instead.
contract_for_agent() {
  case "$1" in
    analyst) echo "analysis-contract.md" ;;
    architect) echo "architecture-contract.md" ;;
    code-reviewer) echo "review-contract.md" ;;
    security-reviewer) echo "security-contract.md" ;;
    qa-engineer) echo "qa-contract.md" ;;
    *) echo "" ;;
  esac
}

PASS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0

pass() { echo "  PASS  $1"; ((PASS_COUNT++)) || true; }
fail() { echo "  FAIL  $1"; ((FAIL_COUNT++)) || true; }
skip() { echo "  SKIP  $1"; ((SKIP_COUNT++)) || true; }

echo ""
echo "=== Agent Golden-File Tests ==="
echo "Fixtures: $TESTS_DIR"
echo ""

if [[ ! -d "$TESTS_DIR" ]]; then
  echo "No tests/agents/ directory found — nothing to test."
  exit 0
fi

for agent_dir in "$TESTS_DIR"/*/; do
  agent_name="$(basename "$agent_dir")"
  actual="$agent_dir/actual-output.md"
  expected="$agent_dir/expected-patterns.txt"

  echo "--- $agent_name ---"

  if [[ ! -f "$actual" ]]; then
    skip "$agent_name — no actual-output.md. Run this agent against its input-* fixture in a live Claude Code session, save the output here, then re-run this script."
    echo ""
    continue
  fi

  if [[ -f "$expected" ]]; then
    while IFS= read -r pattern || [[ -n "$pattern" ]]; do
      [[ -z "$pattern" || "$pattern" == \#* ]] && continue
      if grep -qEi -- "$pattern" "$actual"; then
        pass "pattern: $pattern"
      else
        fail "pattern: $pattern — not found"
      fi
    done < "$expected"
  else
    echo "  (no expected-patterns.txt — skipping pattern checks)"
  fi

  contract_file="$(contract_for_agent "$agent_name")"
  if [[ -n "$contract_file" && -f "$CONTRACTS_DIR/$contract_file" ]]; then
    while IFS= read -r line; do
      if [[ "$line" =~ ^-\ \`(.+)\`$ ]]; then
        section="${BASH_REMATCH[1]}"
        if grep -qF -- "$section" "$actual"; then
          pass "contract section: $section"
        else
          fail "contract section: $section — missing (see $contract_file)"
        fi
      fi
    done < "$CONTRACTS_DIR/$contract_file"
  fi

  echo ""
done

echo "==========================================="
echo "Results: $PASS_COUNT passed, $FAIL_COUNT failed, $SKIP_COUNT skipped"
echo ""

if [[ $FAIL_COUNT -gt 0 ]]; then
  exit 1
fi
