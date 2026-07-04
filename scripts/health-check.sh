#!/usr/bin/env bash
set -euo pipefail

# Deterministic backbone for the `health-check` skill (shared/skills/health-check/SKILL.md).
# Covers everything that's actually scriptable; the skill itself adds AI judgment for anything
# that requires reading prose (e.g. whether an orphaned-looking domain term is really unused).
#
# Written for bash 3.2 (macOS default has no associative arrays) — see scripts/test-agents.sh for the
# same constraint, and every grep pipeline ends with `|| true` where a legitimate zero-match result
# would otherwise abort the script under `set -e` + `pipefail` (a real bug found and fixed twice
# already in this repo's other scripts).

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_DIR="$REPO_DIR/shared"

VERBOSE=false
FIX=false

for arg in "$@"; do
  case "$arg" in
    --verbose) VERBOSE=true ;;
    --fix) FIX=true ;;
  esac
done

PASS_COUNT=0
FAIL_COUNT=0
WARN_COUNT=0

pass() { if $VERBOSE; then echo "  PASS  $1"; fi; ((PASS_COUNT++)) || true; }
fail() { echo "  FAIL  $1"; ((FAIL_COUNT++)) || true; }
warn() { echo "  WARN  $1"; ((WARN_COUNT++)) || true; }

echo ""
echo "=== Framework Health Check ==="
echo "Repository: $REPO_DIR"
if $FIX; then echo "Mode: --fix (will attempt repairs)"; fi
echo ""

# --- 1. Symlinks resolve ---------------------------------------------------
echo "--- Symlinks ---"
for name in agents skills rules; do
  link=".claude/$name"
  expected="../shared/$name"
  if [[ -L "$link" ]]; then
    target=$(readlink "$link")
    if [[ "$target" == "$expected" && -e "$link" ]]; then
      pass "$link -> $target"
    else
      fail "$link -> $target (expected $expected, or target missing)"
      if $FIX; then
        rm -f "$link"
        ln -s "$expected" "$link"
        echo "        [fix] recreated $link -> $expected"
      fi
    fi
  else
    fail "$link is not a symlink"
    if $FIX; then
      rm -rf "$link"
      ln -s "$expected" "$link"
      echo "        [fix] created $link -> $expected"
    fi
  fi
done
echo ""

# --- 2. Agent frontmatter ---------------------------------------------------
echo "--- Agent Frontmatter (name, description, tools, model, version) ---"
for agent_file in "$SHARED_DIR/agents/"*.md; do
  base="$(basename "$agent_file")"
  [[ "$base" == "CHANGELOG.md" ]] && continue

  missing=""
  for field in name description tools model version; do
    if ! grep -q "^${field}:" "$agent_file"; then
      missing="$missing $field"
    fi
  done

  if [[ -z "$missing" ]]; then
    pass "$base"
  else
    fail "$base — missing frontmatter:$missing"
  fi
done
echo ""

# --- 3. Skill frontmatter ---------------------------------------------------
echo "--- Skill Frontmatter (name, description, triggers, standalone) ---"
for skill_dir in "$SHARED_DIR/skills/"*/; do
  skill_name="$(basename "$skill_dir")"
  skill_file="${skill_dir}SKILL.md"

  if [[ ! -f "$skill_file" ]]; then
    fail "$skill_name — no SKILL.md"
    continue
  fi

  missing=""
  for field in name description triggers standalone; do
    if ! grep -q "^${field}:" "$skill_file"; then
      missing="$missing $field"
    fi
  done

  if [[ -z "$missing" ]]; then
    pass "$skill_name"
  else
    fail "$skill_name — missing frontmatter:$missing"
  fi
done
echo ""

# --- 4. Platform config drift (delegates to check-parity.sh) ---------------
echo "--- Platform Config Drift ---"
if bash "$REPO_DIR/scripts/check-parity.sh" > /tmp/health-check-parity.$$ 2>&1; then
  pass "check-parity.sh — no drift"
else
  fail "check-parity.sh reported drift — see below"
  grep -E "DRIFT|MISS" /tmp/health-check-parity.$$ || true
  if $FIX; then
    echo "        [fix] running scripts/generate-configs.sh"
    bash "$REPO_DIR/scripts/generate-configs.sh" > /dev/null
    echo "        [fix] re-run scripts/check-parity.sh to confirm"
  fi
fi
rm -f /tmp/health-check-parity.$$
echo ""

# --- 5. Domain dictionary orphaned terms (best-effort) ----------------------
echo "--- Domain Dictionary Orphaned Terms (best-effort) ---"
DICT="$REPO_DIR/DOMAIN_DICTIONARY.md"
if [[ -f "$DICT" ]]; then
  terms=$(grep -oE '^\| \*\*[A-Za-z][^*]*\*\*' "$DICT" | sed 's/^| \*\*//; s/\*\*$//' || true)
  while IFS= read -r term; do
    [[ -z "$term" ]] && continue
    # Search everywhere except the dictionary file itself for a real usage.
    hits=$( (grep -rlF "$term" "$SHARED_DIR" "$REPO_DIR/docs" 2>/dev/null || true) | (grep -v "DOMAIN_DICTIONARY.md" || true) | wc -l | tr -d ' ')
    if [[ "$hits" -eq 0 ]]; then
      warn "\"$term\" — defined but not referenced anywhere in shared/ or docs/ (may be a framework-level term used only in generated project code, not this repo — verify before removing)"
    else
      pass "\"$term\" — referenced in $hits file(s)"
    fi
  done <<< "$terms"
else
  fail "DOMAIN_DICTIONARY.md not found"
fi
echo ""

# --- 6. Inter-agent contracts exist for pipeline agents ---------------------
echo "--- Inter-Agent Contracts ---"
for pair in "analyst:analysis-contract.md" "architect:architecture-contract.md" \
            "developer:implementation-contract.md" "code-reviewer:review-contract.md" \
            "security-reviewer:security-contract.md" "qa-engineer:qa-contract.md" \
            "sre-engineer:observability-contract.md"; do
  agent="${pair%%:*}"
  contract="${pair##*:}"
  if [[ -f "$SHARED_DIR/contracts/$contract" ]]; then
    pass "$agent -> $contract"
  else
    fail "$agent has no contract at shared/contracts/$contract"
  fi
done
echo ""

# --- 7. Agent changelog up to date (no version mismatches) -----------------
echo "--- Changelog / Version Consistency ---"
CHANGELOG="$SHARED_DIR/agents/CHANGELOG.md"
if [[ -f "$CHANGELOG" ]]; then
  for agent_file in "$SHARED_DIR/agents/"*.md; do
    base="$(basename "$agent_file")"
    [[ "$base" == "CHANGELOG.md" ]] && continue
    name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//' || true)
    version=$(grep '^version:' "$agent_file" | head -1 | sed 's/version: *//' || true)
    [[ -z "$name" || -z "$version" ]] && continue

    if grep -qF "$version" "$CHANGELOG" && grep -qF "$name" "$CHANGELOG"; then
      pass "$name $version — mentioned in CHANGELOG.md"
    else
      warn "$name $version — not found together in CHANGELOG.md (current version may be undocumented)"
    fi
  done
else
  fail "shared/agents/CHANGELOG.md not found"
fi
echo ""

# --- 8. Knowledge Item frontmatter valid ------------------------------------
echo "--- Knowledge Item Frontmatter (name, tags, domain, created) ---"
for ki_dir in "$SHARED_DIR/knowledge" "$REPO_DIR/.claude/knowledge"; do
  [[ -d "$ki_dir" ]] || continue
  for ki_file in "$ki_dir"/*.md; do
    [[ -f "$ki_file" ]] || continue
    base="$(basename "$ki_file")"
    [[ "$base" == "README.md" ]] && continue

    missing=""
    for field in name tags domain created; do
      if ! grep -q "^${field}:" "$ki_file"; then
        missing="$missing $field"
      fi
    done

    if [[ -z "$missing" ]]; then
      pass "$base"
    else
      fail "$base — missing frontmatter:$missing"
    fi
  done
done
echo ""

echo "--- Memory Registry (shared/memory-registry.json) ---"
REGISTRY="$SHARED_DIR/memory-registry.json"
if [[ -f "$REGISTRY" ]]; then
  if python3 -c "import json; json.load(open('$REGISTRY'))" 2>/dev/null; then
    pass "memory-registry.json is valid JSON"
  else
    fail "memory-registry.json is not valid JSON"
  fi

  # Every path each source declares must actually exist.
  registry_paths=$(python3 -c "
import json
data = json.load(open('$REGISTRY'))
for s in data.get('sources', []):
    for p in s.get('paths', []):
        print(p)
" 2>/dev/null || true)
  optional_paths=$(python3 -c "
import json
data = json.load(open('$REGISTRY'))
for p in data.get('optionalPaths', []):
    print(p)
" 2>/dev/null || true)
  while IFS= read -r rpath; do
    [[ -z "$rpath" ]] && continue
    full_path="$REPO_DIR/$rpath"
    if [[ -e "$full_path" ]]; then
      pass "registry path exists: $rpath"
    elif echo "$optional_paths" | grep -qxF "$rpath"; then
      warn "registry path missing (marked optional): $rpath"
    else
      fail "registry path missing: $rpath"
    fi
  done <<< "$registry_paths"

  # No two KIs should share an exact frontmatter name: — a real duplicate, not just an overlap
  # memory-engineer would judge more subtly; this is the cheap, deterministic half of that check.
  ki_names=$( (grep -h '^name:' "$SHARED_DIR"/knowledge/*.md "$REPO_DIR"/.claude/knowledge/*.md 2>/dev/null || true) | sed 's/^name: *//' | sort)
  dupe_names=$(echo "$ki_names" | uniq -d || true)
  if [[ -z "$dupe_names" ]]; then
    pass "no duplicate KI frontmatter names"
  else
    while IFS= read -r dname; do
      [[ -z "$dname" ]] && continue
      fail "duplicate KI frontmatter name: $dname — memory-engineer should audit these for a merge"
    done <<< "$dupe_names"
  fi
else
  warn "shared/memory-registry.json not found — skipping Memory Registry checks"
fi
echo ""

echo "==========================================="
echo "Results: $PASS_COUNT passed, $WARN_COUNT warned, $FAIL_COUNT failed"
if ! $VERBOSE; then
  echo "(pass details hidden — re-run with --verbose to see them)"
fi
echo ""

if [[ $FAIL_COUNT -gt 0 ]]; then
  exit 1
fi
