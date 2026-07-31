#!/usr/bin/env bash
set -euo pipefail

# Fitness function: "no agent exceeds its token budget tier" (docs/features/context-engineering-framework/TODO.md, Epic 16).
#
# Validates every context-manifest.md it finds against two rules derived from
# context-engineer.md's own guardrails:
#   1. The manifest MUST contain a "**Status**:" line in section 7 (Token Budget).
#      A missing Status line is treated as FAIL — an omitted estimate is a missing guardrail,
#      not a passing one (per the context-engineer prompt).
#   2. A Status of WARNING MUST be accompanied by actionable cut recommendations naming at least
#      one file. "None" or an empty list is FAIL.
#
# Scans two locations (both always checked):
#   a. docs/features/*/context-manifest.md  — real manifests persisted after delivery
#   b. tests/fixtures/context-manifests/*.md — hand-authored regression fixtures (Epic 57)
#      These fixtures exist because no real manifests were present at time of writing; they are
#      kept permanently so the script is never in the trivially-passing "nothing to check" state.
#      Expected: passing-manifest.md → PASS, warning-no-cuts-manifest.md → FAIL,
#                missing-status-manifest.md → FAIL
#
# This does not (and cannot) measure a live model's actual context usage — there is no runtime
# hook into that. It checks the artifact context-engineer itself produces and persists.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FEATURES_DIR="$REPO_DIR/docs/features"
FIXTURES_DIR="$REPO_DIR/tests/fixtures/context-manifests"

PASS_COUNT=0
FAIL_COUNT=0

pass() { echo "  PASS  $1"; ((PASS_COUNT++)) || true; }
fail() { echo "  FAIL  $1"; ((FAIL_COUNT++)) || true; }

echo ""
echo "=== Context Budget Fitness Function ==="
echo ""

validate_manifest() {
  local manifest="$1"
  local label="$2"

  status_line=$(grep -m1 '\*\*Status\*\*:' "$manifest" || true)

  if [[ -z "$status_line" ]]; then
    fail "$label — missing **Status**: line in Token Budget section (required guardrail)"
    return
  fi

  if echo "$status_line" | grep -qi "WARNING"; then
    cut_line=$(grep -m1 '\*\*Cut recommendations' "$manifest" || true)
    # A real cut recommendation names a file; "None" or an empty list is not actionable.
    cut_content=$(echo "$cut_line" | sed 's/.*Cut recommendations[^:]*: *//' | tr -d '[]')

    if [[ -z "$cut_content" || "$cut_content" =~ ^[Nn]one ]]; then
      fail "$label — Status: WARNING with no actionable cut recommendations"
    else
      pass "$label — Status: WARNING, cut recommendations present"
    fi
  else
    pass "$label — Status: OK"
  fi
}

# --- Real delivered feature manifests -----------------------------------------
echo "--- Delivered features (docs/features/*/context-manifest.md) ---"
shopt -s nullglob
real_manifests=("$FEATURES_DIR"/*/context-manifest.md)
shopt -u nullglob

if [[ ${#real_manifests[@]} -eq 0 ]]; then
  echo "  (none — no context-manifest.md found in docs/features/; fixtures carry coverage)"
else
  for manifest in "${real_manifests[@]}"; do
    feature_name="$(basename "$(dirname "$manifest")")"
    validate_manifest "$manifest" "$feature_name"
  done
fi
echo ""

# --- Regression fixtures -------------------------------------------------------
echo "--- Fixtures (tests/fixtures/context-manifests/) ---"
shopt -s nullglob
fixture_manifests=("$FIXTURES_DIR"/*.md)
shopt -u nullglob

if [[ ${#fixture_manifests[@]} -eq 0 ]]; then
  fail "tests/fixtures/context-manifests/ is empty — at least the 3 hand-authored fixtures from Epic 57 must be present"
else
  for manifest in "${fixture_manifests[@]}"; do
    fixture_name="$(basename "$manifest" .md)"
    validate_manifest "$manifest" "fixture:$fixture_name"
  done
fi
echo ""

echo "==========================================="
echo "Results: $PASS_COUNT passed, $FAIL_COUNT failed"
echo "  (expected: fixture:passing-manifest PASS; fixture:warning-no-cuts-manifest FAIL; fixture:missing-status-manifest FAIL)"
echo ""

# Exit 1 only when unexpected failures occur: the two fixture FAILs above are by design.
# Count FAILs that are NOT from the known-bad fixtures.
known_bad_fixtures=2
unexpected_fail_count=$(( FAIL_COUNT - known_bad_fixtures ))
if [[ $unexpected_fail_count -gt 0 ]]; then
  echo "  $unexpected_fail_count unexpected failure(s) — investigate"
  exit 1
fi
