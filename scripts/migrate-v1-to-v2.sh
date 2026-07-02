#!/usr/bin/env bash
set -euo pipefail

# Migrates a pre-restructure ("v1") checkout of this repo to the "v2" canonical shared/ layer.
#
# v1 signature: .claude/agents/, .claude/skills/, .claude/rules/ are real directories with content,
# and ARCHITECTURE_RULES.md / DOMAIN_DICTIONARY.md live at the repo root (not under shared/).
# v2 signature: shared/{agents,skills,rules} hold the real content; .claude/{agents,skills,rules} are
# symlinks to them; ARCHITECTURE_RULES.md / DOMAIN_DICTIONARY.md live under shared/ with root-level
# symlinks pointing at them.
#
# This script only MOVES content (git mv where possible, falling back to mv) and creates symlinks —
# it never deletes anything, and if a v2 structure already exists for a given item, that item is
# skipped rather than overwritten. See docs/MIGRATION.md for the full v1 -> v2 breaking-changes list.
#
# Written for bash 3.2 (macOS default) — no associative arrays, and grep/test results that might
# legitimately be empty are guarded with `|| true` so `set -e` + `pipefail` doesn't abort mid-script
# (the same class of bug found and fixed multiple times elsewhere in this repo).

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_DIR"

DRY_RUN=false
for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=true ;;
  esac
done

log() { echo "  $1"; }
did() { echo "  [done] $1"; }
dry() { echo "  [dry-run] $1"; }
skip() { echo "  [skip] $1"; }

move_dir() {
  local old="$1" new="$2"

  if [[ -L "$old" ]]; then
    skip "$old is already a symlink — v2 structure already in place"
    return
  fi

  if [[ ! -d "$old" ]]; then
    skip "$old doesn't exist — nothing to migrate"
    return
  fi

  if [[ -e "$new" ]]; then
    echo "  [WARN] $new already exists and $old is a real directory — not overwriting. Resolve manually."
    return
  fi

  if $DRY_RUN; then
    dry "would move $old -> $new, then symlink $old -> $(basename "$new" | sed "s|.*|../$(dirname "$new" | xargs basename)/&|")"
    return
  fi

  mkdir -p "$(dirname "$new")"
  if git -C "$REPO_DIR" rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    git mv "$old" "$new" 2>/dev/null || mv "$old" "$new"
  else
    mv "$old" "$new"
  fi
  ln -s "../$(basename "$(dirname "$new")")/$(basename "$new")" "$old"
  did "moved $old -> $new, symlinked $old -> $new"
}

move_file_to_shared_with_root_symlink() {
  local old="$1"
  local new="shared/$(basename "$old")"

  if [[ -L "$old" ]]; then
    skip "$old is already a symlink — v2 structure already in place"
    return
  fi

  if [[ ! -f "$old" ]]; then
    skip "$old doesn't exist — nothing to migrate"
    return
  fi

  if [[ -e "$new" ]]; then
    echo "  [WARN] $new already exists and $old is a real file — not overwriting. Resolve manually."
    return
  fi

  if $DRY_RUN; then
    dry "would move $old -> $new, then symlink $old -> $new"
    return
  fi

  mkdir -p shared
  if git -C "$REPO_DIR" rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    git mv "$old" "$new" 2>/dev/null || mv "$old" "$new"
  else
    mv "$old" "$new"
  fi
  ln -s "$new" "$old"
  did "moved $old -> $new, symlinked $old -> $new"
}

echo ""
echo "=== v1 -> v2 Migration ==="
if $DRY_RUN; then echo "(dry run — nothing will be changed)"; fi
echo ""

echo "--- Agents, Skills, Rules ---"
move_dir ".claude/agents" "shared/agents"
move_dir ".claude/skills" "shared/skills"
move_dir ".claude/rules" "shared/rules"
echo ""

echo "--- Root Reference Files ---"
move_file_to_shared_with_root_symlink "ARCHITECTURE_RULES.md"
move_file_to_shared_with_root_symlink "DOMAIN_DICTIONARY.md"
echo ""

echo "==========================================="
if $DRY_RUN; then
  echo "Dry run complete. Re-run without --dry-run to apply."
else
  echo "Migration complete. Run 'bash scripts/health-check.sh --verbose' to verify the result."
fi
echo ""
