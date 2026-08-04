#!/usr/bin/env bash
set -euo pipefail

# Deterministic inventory drift check (Epic 53).
# Counts actual files in shared/ sub-directories and compares them against
# counts stated in prose docs. Exits 1 if any drift is found.
#
# Written for bash 3.2 (macOS default) — no associative arrays.

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_DIR="$REPO_DIR/shared"

DRIFT_COUNT=0
INFO_COUNT=0

drift() { echo "  DRIFT  $1"; ((DRIFT_COUNT++)) || true; }
info()  { echo "  INFO   $1"; ((INFO_COUNT++)) || true; }

echo ""
echo "=== Inventory Drift Check ==="
echo "Repository: $REPO_DIR"
echo ""

# --- Actual counts -----------------------------------------------------------

actual_agents=$(find "$SHARED_DIR/agents" -maxdepth 1 -name "*.md" ! -name "CHANGELOG.md" | wc -l | tr -d ' ')
actual_skills=$(find "$SHARED_DIR/skills" -mindepth 2 -maxdepth 2 -name "SKILL.md" | wc -l | tr -d ' ')
actual_rules=$(find "$SHARED_DIR/rules" -maxdepth 1 -name "*.md" | wc -l | tr -d ' ')
actual_contracts=$(find "$SHARED_DIR/contracts" -maxdepth 1 -name "*.md" | wc -l | tr -d ' ')
actual_templates=$(find "$SHARED_DIR/templates" -maxdepth 1 -name "*.md" | wc -l | tr -d ' ')
actual_schemas=$(find "$SHARED_DIR/schemas" -maxdepth 1 -name "*.schema.json" | wc -l | tr -d ' ')

echo "Actual counts:"
echo "  agents:    $actual_agents"
echo "  skills:    $actual_skills"
echo "  rules:     $actual_rules"
echo "  contracts: $actual_contracts"
echo "  templates: $actual_templates"
echo "  schemas:   $actual_schemas"
echo ""

# --- Helper: check a grep result for a specific claimed number ---------------
# check_claim FILE LINE_NUMBER LABEL CLAIMED_N ACTUAL_N
check_claim() {
  local file="$1" lineno="$2" label="$3" claimed="$4" actual="$5"
  if [[ "$claimed" != "$actual" ]]; then
    drift "$file:$lineno — '$label' claims $claimed, actual is $actual"
  fi
}

# --- README.md ---------------------------------------------------------------
readme="$REPO_DIR/README.md"

# Targeted pass — find any "N agents" or "N skills" claims in README
while IFS= read -r match; do
  lineno=$(echo "$match" | cut -d: -f1)
  text=$(echo "$match" | cut -d: -f2-)
  if echo "$text" | grep -qE '[0-9]+ agents'; then
    n=$(echo "$text" | grep -oE '[0-9]+ agents' | grep -oE '[0-9]+' | head -1)
    check_claim "README.md" "$lineno" "agents" "$n" "$actual_agents"
  fi
  if echo "$text" | grep -qE '[0-9]+ skills'; then
    n=$(echo "$text" | grep -oE '[0-9]+ skills' | grep -oE '[0-9]+' | head -1)
    check_claim "README.md" "$lineno" "skills" "$n" "$actual_skills"
  fi
done < <(grep -nE '[0-9]+ agents|[0-9]+ skills' "$readme" 2>/dev/null || true)

# Also check heading-format counts: "## Agent Roster (N)" and "## Skill Catalog (N)"
while IFS= read -r match; do
  lineno=$(echo "$match" | cut -d: -f1)
  text=$(echo "$match" | cut -d: -f2-)
  if echo "$text" | grep -qiE 'Agent Roster \([0-9]+\)'; then
    n=$(echo "$text" | grep -oE '\([0-9]+\)' | tr -d '()' | head -1)
    check_claim "README.md" "$lineno" "Agent Roster heading" "$n" "$actual_agents"
  fi
  if echo "$text" | grep -qiE 'Skill Catalog \([0-9]+\)'; then
    n=$(echo "$text" | grep -oE '\([0-9]+\)' | tr -d '()' | head -1)
    check_claim "README.md" "$lineno" "Skill Catalog heading" "$n" "$actual_skills"
  fi
done < <(grep -nE 'Agent Roster \([0-9]+\)|Skill Catalog \([0-9]+\)' "$readme" 2>/dev/null || true)

# --- docs/AGENT_REFERENCE.md -------------------------------------------------
agent_ref="$REPO_DIR/docs/AGENT_REFERENCE.md"
if [[ -f "$agent_ref" ]]; then
  while IFS= read -r match; do
    lineno=$(echo "$match" | cut -d: -f1)
    text=$(echo "$match" | cut -d: -f2-)
    if echo "$text" | grep -qE '[0-9]+ agents'; then
      n=$(echo "$text" | grep -oE '[0-9]+ agents' | grep -oE '[0-9]+' | head -1)
      check_claim "docs/AGENT_REFERENCE.md" "$lineno" "agents" "$n" "$actual_agents"
    fi
  done < <(grep -nE '[0-9]+ agents' "$agent_ref" 2>/dev/null || true)
fi

# --- docs/audits/ — informational only (snapshots, not authoritative) --------
for audit_file in "$REPO_DIR/docs/audits/"*.md; do
  [[ -f "$audit_file" ]] || continue
  rel="${audit_file#"$REPO_DIR/"}"
  while IFS= read -r match; do
    lineno=$(echo "$match" | cut -d: -f1)
    text=$(echo "$match" | cut -d: -f2-)
    for pattern in "agents" "skills" "contracts"; do
      if echo "$text" | grep -qE "[0-9]+ ${pattern}"; then
        n=$(echo "$text" | grep -oE "[0-9]+ ${pattern}" | grep -oE '[0-9]+' | head -1)
        case "$pattern" in
          agents)    actual_val="$actual_agents" ;;
          skills)    actual_val="$actual_skills" ;;
          contracts) actual_val="$actual_contracts" ;;
        esac
        if [[ "$n" != "$actual_val" ]]; then
          info "$rel:$lineno — '$pattern' snapshot says $n, actual is $actual_val (audit docs are snapshots — update manually if needed)"
        fi
      fi
    done
  done < <(grep -nE '[0-9]+ agents|[0-9]+ skills|[0-9]+ contracts' "$audit_file" 2>/dev/null || true)
done

# --- Platform config files — informational only ------------------------------
for cfg in \
  "$REPO_DIR/.github/copilot-instructions.md" \
  "$REPO_DIR/.openai.md" \
  "$REPO_DIR/AGENTS.md" \
  "$REPO_DIR/.windsurfrules"; do
  [[ -f "$cfg" ]] || continue
  rel="${cfg#"$REPO_DIR/"}"
  while IFS= read -r match; do
    lineno=$(echo "$match" | cut -d: -f1)
    text=$(echo "$match" | cut -d: -f2-)
    for pattern in "agents" "skills"; do
      if echo "$text" | grep -qE "[0-9]+ ${pattern}"; then
        n=$(echo "$text" | grep -oE "[0-9]+ ${pattern}" | grep -oE '[0-9]+' | head -1)
        case "$pattern" in
          agents) actual_val="$actual_agents" ;;
          skills) actual_val="$actual_skills" ;;
        esac
        if [[ "$n" != "$actual_val" ]]; then
          info "$rel:$lineno — '$pattern' claims $n, actual is $actual_val (generated config — re-run platform exporter to sync)"
        fi
      fi
    done
  done < <(grep -nE '[0-9]+ agents|[0-9]+ skills' "$cfg" 2>/dev/null || true)
done

# --- Summary -----------------------------------------------------------------
echo ""
echo "Results: $DRIFT_COUNT drift(s) found, $INFO_COUNT info notice(s)"
echo ""

if [[ $DRIFT_COUNT -gt 0 ]]; then
  echo "Fix: update the prose count(s) above to match actual values."
  exit 1
fi
