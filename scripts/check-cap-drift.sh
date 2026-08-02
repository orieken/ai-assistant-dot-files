#!/usr/bin/env bash
set -euo pipefail

# Cap-drift check (Epic 64 Op 3).
# Verifies that each config in shared/configs/ enforces the same cyclomatic
# complexity cap prescribed in its source convention file.
#
# Strategy:
#   - Each config carries a "# CAP: N" marker line in its header.
#   - Each convention file has a parseable cap claim (grep pattern below).
#   - This script extracts both and compares them.
#
# FAILs on mismatch; WARNs when a config or convention cannot be parsed
# (never false-FAILs on missing tooling — see the 6c422cb lesson).
#
# Written for bash 3.2 (macOS default) — no associative arrays, no pipefail in subshells.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIGS_DIR="$REPO_DIR/shared/configs"
RULES_DIR="$REPO_DIR/shared/rules"

DRIFT_COUNT=0
WARN_COUNT=0

drift() { echo "  FAIL   $1"; ((DRIFT_COUNT++)) || true; }
warn()  { echo "  WARN   $1"; ((WARN_COUNT++)) || true; }
pass()  { echo "  PASS   $1"; }

echo ""
echo "=== Cap Drift Check (shared/configs/ vs. shared/rules/) ==="
echo "Repository: $REPO_DIR"
echo ""

# ---------------------------------------------------------------------------
# extract_config_cap FILE
#   Greps for "# CAP: N" in the first 10 lines of FILE.
#   Prints the number N, or empty string if not found.
# ---------------------------------------------------------------------------
extract_config_cap() {
  local file="$1"
  head -10 "$file" | grep -oE 'CAP: [0-9]+' | grep -oE '[0-9]+' | head -1 || true
}

# ---------------------------------------------------------------------------
# extract_convention_cap FILE GREP_PATTERN [lt]
#   Greps FILE for GREP_PATTERN and extracts the first digit group.
#   If the third argument is "lt" (less-than), the number is treated as an
#   exclusive upper bound and returned as N-1 (e.g. "< 7" → cap 6).
#   Prints the cap number, or empty string if not found.
# ---------------------------------------------------------------------------
extract_convention_cap() {
  local file="$1"
  local pattern="$2"
  local mode="${3:-}"
  local raw
  raw=$(grep -iE "$pattern" "$file" 2>/dev/null | grep -oE '[0-9]+' | head -1 || true)
  if [[ -z "$raw" ]]; then
    echo ""
    return
  fi
  if [[ "$mode" == "lt" ]]; then
    echo $((raw - 1))
  else
    echo "$raw"
  fi
}

# ---------------------------------------------------------------------------
# check_cap CONFIG_FILE CONVENTION_FILE CONVENTION_GREP_PATTERN LABEL [lt]
#   Pass "lt" as fifth arg when the convention states an exclusive upper bound
#   (e.g. "< 7") rather than a direct cap (e.g. "capped at 6").
# ---------------------------------------------------------------------------
check_cap() {
  local config_file="$1"
  local convention_file="$2"
  local convention_pattern="$3"
  local label="$4"
  local mode="${5:-}"

  if [[ ! -f "$config_file" ]]; then
    warn "$label — config file missing: $config_file"
    return
  fi
  if [[ ! -f "$convention_file" ]]; then
    warn "$label — convention file missing: $convention_file"
    return
  fi

  local config_cap
  config_cap=$(extract_config_cap "$config_file")

  local convention_cap
  convention_cap=$(extract_convention_cap "$convention_file" "$convention_pattern" "$mode")

  if [[ -z "$config_cap" ]]; then
    warn "$label — cannot parse CAP from config $(basename "$config_file") (add '# CAP: N' header)"
    return
  fi
  if [[ -z "$convention_cap" ]]; then
    warn "$label — cannot parse cap from convention $(basename "$convention_file") (pattern: '$convention_pattern')"
    return
  fi

  if [[ "$config_cap" == "$convention_cap" ]]; then
    pass "$label — cap: $config_cap (config) == $convention_cap (convention)"
  else
    drift "$label — cap mismatch: config says $config_cap, convention says $convention_cap"
  fi
}

# ---------------------------------------------------------------------------
# Checks: one per config file
# ---------------------------------------------------------------------------

# TypeScript — ESLint complexity max 6
# Convention: "capped at 6" (typescript-conventions.md line 11-12)
check_cap \
  "$CONFIGS_DIR/eslint.framework.config.js" \
  "$RULES_DIR/typescript-conventions.md" \
  "capped at [0-9]" \
  "TypeScript/ESLint"

# Go — golangci-lint gocyclo min-complexity 7 (= cap 6)
# Convention: CLAUDE.md line 10 "Cyclomatic complexity < 7" → cap = 7-1 = 6.
# go-conventions.md names gocyclo/revive; CLAUDE.md is the authoritative cap source.
check_cap \
  "$CONFIGS_DIR/.golangci.yml" \
  "$REPO_DIR/CLAUDE.md" \
  "Cyclomatic complexity < [0-9]" \
  "Go/golangci-lint" \
  "lt"

# Python — ruff max-complexity 6
# Convention: CLAUDE.md line 10 "Cyclomatic complexity < 7" → cap = 7-1 = 6.
check_cap \
  "$CONFIGS_DIR/ruff.toml" \
  "$REPO_DIR/CLAUDE.md" \
  "Cyclomatic complexity < [0-9]" \
  "Python/ruff" \
  "lt"

# Kotlin — detekt ComplexMethod threshold 6
# Convention: kotlin-conventions.md "ComplexMethod rule capped at `6`"
check_cap \
  "$CONFIGS_DIR/detekt.yml" \
  "$RULES_DIR/kotlin-conventions.md" \
  "capped at .?[0-9]" \
  "Kotlin/detekt"

# Swift — SwiftLint cyclomatic_complexity warning 6
# Convention: swift-conventions.md "capped at `6`"
check_cap \
  "$CONFIGS_DIR/.swiftlint.yml" \
  "$RULES_DIR/swift-conventions.md" \
  "capped at .?[0-9]" \
  "Swift/SwiftLint"

# Rust — clippy cognitive-complexity-threshold 6
# Convention: rust-conventions.md "cap at `6`"
check_cap \
  "$CONFIGS_DIR/clippy.toml" \
  "$RULES_DIR/rust-conventions.md" \
  "cap at .?[0-9]" \
  "Rust/clippy"

# Java — Checkstyle CyclomaticComplexity max 6
# Convention: CLAUDE.md line 10 "Cyclomatic complexity < 7" → cap = 7-1 = 6.
check_cap \
  "$CONFIGS_DIR/checkstyle.xml" \
  "$REPO_DIR/CLAUDE.md" \
  "Cyclomatic complexity < [0-9]" \
  "Java/Checkstyle" \
  "lt"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "Results: $DRIFT_COUNT cap mismatch(es), $WARN_COUNT warning(s)"
echo ""

if [[ $DRIFT_COUNT -gt 0 ]]; then
  echo "Fix: update the config CAP marker and/or the enforced value to match the convention."
  exit 1
fi
