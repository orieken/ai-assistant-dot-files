#!/usr/bin/env bash
# check-parity.sh — Detects drift between CLAUDE.md and other platform config files.
# Run this after updating CLAUDE.md to see which platform configs may need syncing.

set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")" && pwd)"
REFERENCE_FILE="$REPO_DIR/CLAUDE.md"

PLATFORM_FILES=(
  ".cursor/rules/global.mdc"
  ".windsurfrules"
  ".openai.md"
  ".github/copilot-instructions.md"
  ".gemini/antigravity/instructions.md"
)

CORE_CONCEPTS=(
  "cyclomatic complexity"
  "< 7"
  "85%"
  "Clean Architecture"
  "SOLID"
  "TDD"
  "Boy Scout"
  "DOMAIN_DICTIONARY"
  "Expand/Contract"
  "Saturday Framework"
  "Sunday Framework"
  "OpenTelemetry"
  "fitness function"
)

echo "=== Cross-Platform Parity Check ==="
echo "Reference: CLAUDE.md ($(stat -f '%Sm' "$REFERENCE_FILE" 2>/dev/null || stat -c '%y' "$REFERENCE_FILE" 2>/dev/null | cut -d' ' -f1))"
echo ""

has_drift=0

for platform_file in "${PLATFORM_FILES[@]}"; do
  full_path="$REPO_DIR/$platform_file"

  if [ ! -f "$full_path" ]; then
    echo "MISSING: $platform_file"
    has_drift=1
    continue
  fi

  last_modified=$(stat -f '%Sm' "$full_path" 2>/dev/null || stat -c '%y' "$full_path" 2>/dev/null | cut -d' ' -f1)
  missing_concepts=()

  for concept in "${CORE_CONCEPTS[@]}"; do
    if ! grep -qi "$concept" "$full_path" 2>/dev/null; then
      missing_concepts+=("$concept")
    fi
  done

  if [ ${#missing_concepts[@]} -eq 0 ]; then
    echo "PASS: $platform_file (last modified: $last_modified)"
  else
    echo "DRIFT: $platform_file (last modified: $last_modified)"
    echo "  Missing concepts:"
    for concept in "${missing_concepts[@]}"; do
      echo "    - $concept"
    done
    has_drift=1
  fi
done

echo ""
if [ $has_drift -eq 0 ]; then
  echo "Result: All platform configs are in sync with core concepts from CLAUDE.md"
else
  echo "Result: Drift detected. Review the files above and sync missing concepts."
fi

exit $has_drift
