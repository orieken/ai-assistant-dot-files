#!/usr/bin/env bash
set -euo pipefail

# Fitness function: "no agent exceeds its token budget tier" (docs/features/context-engineering-framework/TODO.md, Epic 16).
#
# Operationalized as: no persisted context-manifest.md may report a WARNING token-budget status without
# actionable cut recommendations. context-engineer's own guardrail (shared/agents/context-engineer.md) says
# a WARNING must come with specific files to cut — this script verifies that guardrail wasn't silently
# violated across every feature actually delivered, not just trusted as prose.
#
# This does not (and cannot) measure a live model's actual context usage — there's no hook into that at
# runtime. It checks the artifact context-engineer itself produces and persists.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FEATURES_DIR="$REPO_DIR/docs/features"

PASS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0

pass() { echo "  PASS  $1"; ((PASS_COUNT++)) || true; }
fail() { echo "  FAIL  $1"; ((FAIL_COUNT++)) || true; }
skip() { echo "  SKIP  $1"; ((SKIP_COUNT++)) || true; }

echo ""
echo "=== Context Budget Fitness Function ==="
echo "Scanning: $FEATURES_DIR/*/context-manifest.md"
echo ""

if [[ ! -d "$FEATURES_DIR" ]]; then
  echo "No docs/features/ directory found — nothing to check."
  exit 0
fi

shopt -s nullglob
manifests=("$FEATURES_DIR"/*/context-manifest.md)
shopt -u nullglob

if [[ ${#manifests[@]} -eq 0 ]]; then
  echo "No context-manifest.md files found in any delivered feature — nothing to check."
  exit 0
fi

for manifest in "${manifests[@]}"; do
  feature_name="$(basename "$(dirname "$manifest")")"

  status_line=$(grep -m1 '\*\*Status\*\*:' "$manifest" || true)

  if [[ -z "$status_line" ]]; then
    skip "$feature_name — no Token Budget Status line found in context-manifest.md"
    continue
  fi

  if echo "$status_line" | grep -qi "WARNING"; then
    cut_line=$(grep -m1 '\*\*Cut recommendations' "$manifest" || true)
    # A real cut recommendation names a file; "None" or an empty list after the colon is not actionable.
    cut_content=$(echo "$cut_line" | sed 's/.*Cut recommendations[^:]*: *//' | tr -d '[]')

    if [[ -z "$cut_content" || "$cut_content" =~ ^[Nn]one ]]; then
      fail "$feature_name — Status: WARNING with no actionable cut recommendations"
    else
      pass "$feature_name — Status: WARNING, cut recommendations present"
    fi
  else
    pass "$feature_name — Status: OK"
  fi
done

echo ""
echo "==========================================="
echo "Results: $PASS_COUNT passed, $FAIL_COUNT failed, $SKIP_COUNT skipped"
echo ""

if [[ $FAIL_COUNT -gt 0 ]]; then
  exit 1
fi
