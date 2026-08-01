#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

MODE=""
TARGET_DIR=""
DRY_RUN=false

usage() {
  cat <<'EOF'
Usage: uninstall.sh [OPTIONS]

Remove the Context Engineering Framework from this machine.

Modes (pick one):
  --global              Remove global config symlinks from ~/
  --project <path>      Remove framework files from a target project
  --project             Remove framework files from current directory

Options:
  --dry-run             Show what would be removed without doing it
  -h, --help            Show this help
EOF
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --global)    MODE="global"; shift ;;
    --project)
      MODE="project"
      if [[ $# -gt 1 && ! "$2" =~ ^-- ]]; then
        TARGET_DIR="$2"; shift 2
      else
        TARGET_DIR="$(pwd)"; shift
      fi
      ;;
    --dry-run)   DRY_RUN=true; shift ;;
    -h|--help)   usage ;;
    *)           echo "Unknown option: $1"; usage ;;
  esac
done

if [[ -z "$MODE" ]]; then
  echo "Error: specify --global or --project <path>"
  exit 1
fi

log()  { echo "  $1"; }
ok()   { echo "  [removed] $1"; }
skip() { echo "  [skip] $1"; }
dry()  { echo "  [dry-run] $1"; }

remove_marker() {
  local path="$1"
  if [[ ! -f "$path" ]]; then
    skip "$path (not found)"
    return
  fi
  if $DRY_RUN; then
    dry "would remove $path"
  else
    rm -f "$path"
    ok "$path"
  fi
}

remove_if_framework() {
  local path="$1"

  if [[ ! -e "$path" && ! -L "$path" ]]; then
    skip "$path (not found)"
    return
  fi

  if [[ -L "$path" ]]; then
    local target
    target=$(readlink "$path")
    if [[ "$target" == *"shared"* || "$target" == *"ai-assistant-dot-files"* ]]; then
      if $DRY_RUN; then
        dry "would remove symlink $path -> $target"
      else
        rm "$path"
        ok "$path (symlink)"
      fi
      restore_backup "$path"
      return
    else
      skip "$path (symlink to $target — not ours)"
      return
    fi
  fi

  if $DRY_RUN; then
    dry "would remove $path"
  else
    rm -rf "$path"
    ok "$path"
  fi
  restore_backup "$path"
}

restore_backup() {
  local path="$1"
  local latest_backup
  latest_backup=$(ls -t "${path}.bak."* 2>/dev/null | head -1 || true)

  if [[ -n "$latest_backup" ]]; then
    if $DRY_RUN; then
      dry "would restore $latest_backup -> $path"
    else
      mv "$latest_backup" "$path"
      log "restored $latest_backup -> $path"
    fi
  fi
}

echo ""
echo "Context Engineering Framework Uninstaller"
echo "=========================================="
echo ""
echo "Mode:    $MODE"
if [[ "$MODE" == "project" ]]; then
  echo "Target:  $TARGET_DIR"
fi
echo "Dry run: $DRY_RUN"
echo ""

if [[ "$MODE" == "global" ]]; then
  log "--- Claude Code ---"
  remove_if_framework "$HOME/.claude/agents"
  remove_if_framework "$HOME/.claude/skills"
  remove_if_framework "$HOME/.claude/rules"

  log ""
  log "--- Cursor ---"
  remove_if_framework "$HOME/.cursor/agents"
  remove_if_framework "$HOME/.cursor/skills"
  remove_if_framework "$HOME/.cursor/rules/global.mdc"
  for mdc in architecture design-principles approval-gates agent-roster testing go-backend vue-frontend; do
    remove_if_framework "$HOME/.cursor/rules/${mdc}.mdc"
  done
  remove_if_framework "$HOME/.cursorrules"

  log ""
  log "--- Windsurf ---"
  remove_if_framework "$HOME/.windsurfrules"

  log ""
  log "--- GitHub Copilot ---"
  remove_if_framework "$HOME/.github/copilot-instructions.md"

  log ""
  log "--- Gemini ---"
  remove_if_framework "$HOME/.gemini/antigravity/instructions.md"

  log ""
  log "--- OpenAI ---"
  remove_if_framework "$HOME/.openai.md"

  log ""
  log "--- Framework Version Marker ---"
  remove_marker "$HOME/.claude/framework-install.json"
else
  log "--- Claude Code ---"
  remove_if_framework "$TARGET_DIR/.claude/agents"
  remove_if_framework "$TARGET_DIR/.claude/skills"
  remove_if_framework "$TARGET_DIR/.claude/rules"
  remove_if_framework "$TARGET_DIR/ARCHITECTURE_RULES.md"
  remove_if_framework "$TARGET_DIR/DOMAIN_DICTIONARY.md"

  log ""
  log "--- Cursor ---"
  remove_if_framework "$TARGET_DIR/.cursor/agents"
  remove_if_framework "$TARGET_DIR/.cursor/skills"
  remove_if_framework "$TARGET_DIR/.cursor/rules/global.mdc"
  for mdc in architecture design-principles approval-gates agent-roster testing go-backend vue-frontend; do
    remove_if_framework "$TARGET_DIR/.cursor/rules/${mdc}.mdc"
  done
  remove_if_framework "$TARGET_DIR/.cursorrules"

  log ""
  log "--- Windsurf ---"
  remove_if_framework "$TARGET_DIR/.windsurfrules"

  log ""
  log "--- GitHub Copilot ---"
  remove_if_framework "$TARGET_DIR/.github/copilot-instructions.md"

  log ""
  log "--- Gemini ---"
  remove_if_framework "$TARGET_DIR/.gemini/antigravity/instructions.md"

  log ""
  log "--- OpenAI ---"
  remove_if_framework "$TARGET_DIR/.openai.md"

  log ""
  log "--- Framework Version Marker ---"
  remove_marker "$TARGET_DIR/.claude/framework-install.json"
fi

echo ""
if $DRY_RUN; then
  echo "This was a dry run. No files were modified."
else
  echo "Uninstall complete. Backups were restored where available."
fi
echo ""
