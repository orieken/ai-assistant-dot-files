#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SHARED_DIR="$REPO_DIR/shared"

MODE=""
TARGET_DIR=""
DRY_RUN=false
USE_COPY=false
PLATFORM_FILTER=""
SHOW_TOUR=false

usage() {
  cat <<'EOF'
Usage: install.sh [OPTIONS]

Install the Context Engineering Framework on this machine.

Modes (pick one):
  --global              Symlink configs to ~/  (always current)
  --project <path>      Copy configs to a target project directory
  --project             Copy configs to current directory

Options:
  --copy                Use copies instead of symlinks (Windows/WSL fallback)
  --platform <name>     Only install for a specific platform (claude-code, cursor, etc.)
  --dry-run             Show what would be installed without doing it
  --tour                Run the onboarding skill after install
  -h, --help            Show this help

Examples:
  ./install.sh --global
  ./install.sh --project /path/to/my-app
  ./install.sh --project --platform claude-code
  ./install.sh --global --dry-run
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
    --copy)      USE_COPY=true; shift ;;
    --platform)  PLATFORM_FILTER="$2"; shift 2 ;;
    --dry-run)   DRY_RUN=true; shift ;;
    --tour)      SHOW_TOUR=true; shift ;;
    -h|--help)   usage ;;
    *)           echo "Unknown option: $1"; usage ;;
  esac
done

if [[ -z "$MODE" ]]; then
  echo "Error: specify --global or --project <path>"
  echo "Run './install.sh --help' for usage."
  exit 1
fi

AGENT_COUNT=$(find "$SHARED_DIR/agents" -name "*.md" -not -name "CHANGELOG.md" | wc -l | tr -d ' ')
SKILL_COUNT=$(find "$SHARED_DIR/skills" -mindepth 1 -maxdepth 1 -type d | wc -l | tr -d ' ')
RULE_COUNT=$(find "$SHARED_DIR/rules" -name "*.md" | wc -l | tr -d ' ')

log()  { echo "  $1"; }
ok()   { echo "  [ok] $1"; }
skip() { echo "  [skip] $1"; }
dry()  { echo "  [dry-run] $1"; }

link_or_copy() {
  local src="$1"
  local dest="$2"

  if [[ -e "$dest" || -L "$dest" ]]; then
    if [[ -L "$dest" ]]; then
      local existing_target
      existing_target=$(readlink "$dest")
      if [[ "$existing_target" == "$src" ]]; then
        skip "$dest (already linked)"
        return
      fi
    fi
    local backup="${dest}.bak.$(date +%s)"
    if $DRY_RUN; then
      dry "would backup $dest -> $backup"
    else
      mv "$dest" "$backup"
      log "backed up $dest -> $backup"
    fi
  fi

  if $DRY_RUN; then
    if $USE_COPY; then
      dry "would copy $src -> $dest"
    else
      dry "would symlink $dest -> $src"
    fi
    return
  fi

  mkdir -p "$(dirname "$dest")"
  if $USE_COPY; then
    if [[ -d "$src" ]]; then
      cp -r "$src" "$dest"
    else
      cp "$src" "$dest"
    fi
    ok "copied $dest"
  else
    ln -s "$src" "$dest"
    ok "linked $dest -> $src"
  fi
}

detect_platforms() {
  local detected=()

  detected+=("claude-code")

  if command -v cursor &>/dev/null || [[ -d "$HOME/.cursor" ]]; then
    detected+=("cursor")
  fi

  if [[ -f "$HOME/.windsurfrules" ]] || command -v windsurf &>/dev/null; then
    detected+=("windsurf")
  fi

  if command -v gh &>/dev/null && gh extension list 2>/dev/null | grep -q copilot; then
    detected+=("github-copilot")
  elif [[ -d "$HOME/.config/github-copilot" ]]; then
    detected+=("github-copilot")
  fi

  if command -v gemini &>/dev/null || [[ -d "$HOME/.gemini" ]]; then
    detected+=("gemini")
  fi

  if [[ -f "$HOME/.openai.md" ]] || [[ -f ".openai.md" ]]; then
    detected+=("openai-codex")
  fi

  echo "${detected[@]}"
}

should_install() {
  local platform="$1"
  if [[ -n "$PLATFORM_FILTER" ]]; then
    [[ "$platform" == "$PLATFORM_FILTER" ]]
  else
    return 0
  fi
}

install_claude_code() {
  log ""
  log "--- Claude Code (Tier 1: Full) ---"

  if [[ "$MODE" == "global" ]]; then
    local claude_dir="$HOME/.claude"
    mkdir -p "$claude_dir" 2>/dev/null || true
    link_or_copy "$SHARED_DIR/agents" "$claude_dir/agents"
    link_or_copy "$SHARED_DIR/skills" "$claude_dir/skills"
    link_or_copy "$SHARED_DIR/rules" "$claude_dir/rules"
  else
    local claude_dir="$TARGET_DIR/.claude"
    mkdir -p "$claude_dir" 2>/dev/null || true
    link_or_copy "$SHARED_DIR/agents" "$claude_dir/agents"
    link_or_copy "$SHARED_DIR/skills" "$claude_dir/skills"
    link_or_copy "$SHARED_DIR/rules" "$claude_dir/rules"

    link_or_copy "$SHARED_DIR/ARCHITECTURE_RULES.md" "$TARGET_DIR/ARCHITECTURE_RULES.md"
    link_or_copy "$SHARED_DIR/DOMAIN_DICTIONARY.md" "$TARGET_DIR/DOMAIN_DICTIONARY.md"

    if [[ ! -f "$TARGET_DIR/CLAUDE.md" ]]; then
      local template_claude="$REPO_DIR/templates/claude-feature-team/CLAUDE.md"
      if [[ -f "$template_claude" ]]; then
        if $DRY_RUN; then
          dry "would copy CLAUDE.md template to $TARGET_DIR/"
        else
          cp "$template_claude" "$TARGET_DIR/CLAUDE.md"
          ok "installed CLAUDE.md template"
        fi
      fi
    else
      skip "CLAUDE.md already exists"
    fi

    if [[ ! -d "$TARGET_DIR/features" ]]; then
      local template_features="$REPO_DIR/templates/claude-feature-team/features"
      if [[ -d "$template_features" ]]; then
        if $DRY_RUN; then
          dry "would copy features/ template"
        else
          mkdir -p "$TARGET_DIR/features"
          cp "$template_features"/*.md "$TARGET_DIR/features/" 2>/dev/null || true
          ok "installed features/ template"
        fi
      fi
    else
      skip "features/ already exists"
    fi
  fi
}

install_cursor() {
  log ""
  log "--- Cursor (Tier 2: Personas + Rules) ---"

  local rules_dir
  if [[ "$MODE" == "global" ]]; then
    rules_dir="$HOME/.cursor/rules"
  else
    rules_dir="$TARGET_DIR/.cursor/rules"
  fi

  mkdir -p "$rules_dir" 2>/dev/null || true

  local mdc_file="$rules_dir/global.mdc"
  if $DRY_RUN; then
    dry "would generate $mdc_file (inlined from shared/rules/)"
  else
    cat > "$mdc_file" <<MDCEOF
---
description: "Context Engineering Framework — architecture, design, and governance rules"
alwaysApply: true
---

MDCEOF

    for rule_file in "$SHARED_DIR/rules/"*.md; do
      echo "" >> "$mdc_file"
      cat "$rule_file" >> "$mdc_file"
      echo "" >> "$mdc_file"
    done

    cat >> "$mdc_file" <<'ROSTEREOF'

# Agent / Persona roster

The following specialized personas are available. Invoke them by name when you need domain-specific expertise.

ROSTEREOF

    for agent_file in "$SHARED_DIR/agents/"*.md; do
      local agent_name
      agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//')
      local agent_desc
      agent_desc=$(grep '^description:' "$agent_file" | head -1 | sed 's/description: *//')
      if [[ -n "$agent_name" ]]; then
        echo "- **$agent_name**: $agent_desc" >> "$mdc_file"
      fi
    done

    ok "generated $mdc_file"
  fi

  local cursorrules
  if [[ "$MODE" == "global" ]]; then
    cursorrules="$HOME/.cursorrules"
  else
    cursorrules="$TARGET_DIR/.cursorrules"
  fi

  if $DRY_RUN; then
    dry "would generate $cursorrules"
  else
    cp "$mdc_file" "$cursorrules" 2>/dev/null || true
    ok "generated $cursorrules"
  fi
}

install_windsurf() {
  log ""
  log "--- Windsurf (Tier 2: Personas + Rules) ---"

  local windsurfrules
  if [[ "$MODE" == "global" ]]; then
    windsurfrules="$HOME/.windsurfrules"
  else
    windsurfrules="$TARGET_DIR/.windsurfrules"
  fi

  local cursor_source
  if [[ "$MODE" == "global" ]]; then
    cursor_source="$HOME/.cursorrules"
  else
    cursor_source="$TARGET_DIR/.cursorrules"
  fi

  if [[ -f "$cursor_source" ]]; then
    if $DRY_RUN; then
      dry "would copy $cursor_source -> $windsurfrules"
    else
      cp "$cursor_source" "$windsurfrules"
      ok "generated $windsurfrules"
    fi
  else
    skip "no .cursorrules to copy for windsurf (install cursor first)"
  fi
}

generate_tier3_config() {
  local platform_name="$1"
  local dest_path="$2"
  local header="$3"

  if $DRY_RUN; then
    dry "would generate $dest_path"
    return
  fi

  mkdir -p "$(dirname "$dest_path")"

  cat > "$dest_path" <<HEADEREOF
$header

## AI Feature Team & Global Rules
You are part of the Saturday Multi-Agent Feature Team. Before beginning any complex task, architectural decision, or feature delivery, you MUST adhere to the rules below.

HEADEREOF

  for rule_file in "$SHARED_DIR/rules/"*.md; do
    echo "" >> "$dest_path"
    cat "$rule_file" >> "$dest_path"
    echo "" >> "$dest_path"
  done

  cat >> "$dest_path" <<'RULESEOF'

## Craftsmanship Rules
You must **strictly adhere** to the patterns defined in `ARCHITECTURE_RULES.md` (Clean Architecture, DDD, GoF patterns, and micro-rules).
- **TDD/BDD First**: Drive design through testing. Feature code is incomplete without tests. Practice Red-Green-Refactor.
- **Kent Beck (Simple Design)**: 1) Passes tests, 2) Reveals intention, 3) No duplication, 4) Fewest elements.
- **Martin Fowler (Refactoring)**: Use named refactoring operations (Extract Function, Inline Variable, etc.) instead of vague cleanups.
- **Architectural Constraints & Fitness Functions**: Enforce cyclomatic complexity `< 7` and functions `< 30` LOC.
- **The Boy Scout Rule**: Always leave the code cleaner than you found it.

## Tech Stack
- **Backend / MCP**: Go
- **Frontend**: Vue 3 + Tailwind CSS
- **Test Automation**: TypeScript, Playwright, Cucumber.js, k6

## Agent / Persona Roster

The following specialized personas are available. Invoke them by name when you need domain-specific expertise.

RULESEOF

  for agent_file in "$SHARED_DIR/agents/"*.md; do
    local agent_name
    agent_name=$(grep '^name:' "$agent_file" | head -1 | sed 's/name: *//')
    local agent_desc
    agent_desc=$(grep '^description:' "$agent_file" | head -1 | sed 's/description: *//')
    if [[ -n "$agent_name" ]]; then
      echo "- **$agent_name**: $agent_desc" >> "$dest_path"
    fi
  done

  ok "generated $dest_path"
}

install_copilot() {
  log ""
  log "--- GitHub Copilot (Tier 3: System Prompt) ---"

  local dest
  if [[ "$MODE" == "global" ]]; then
    dest="$HOME/.github/copilot-instructions.md"
  else
    dest="$TARGET_DIR/.github/copilot-instructions.md"
  fi

  generate_tier3_config "github-copilot" "$dest" \
    "# Copilot Instructions (Saturday Framework)"
}

install_gemini() {
  log ""
  log "--- Gemini / Antigravity (Tier 3: System Prompt) ---"

  local dest
  if [[ "$MODE" == "global" ]]; then
    dest="$HOME/.gemini/antigravity/instructions.md"
  else
    dest="$TARGET_DIR/.gemini/antigravity/instructions.md"
  fi

  generate_tier3_config "gemini" "$dest" \
    "# Antigravity Instructions (Saturday Framework)"
}

install_openai() {
  log ""
  log "--- OpenAI / Codex (Tier 3: System Prompt) ---"

  local dest
  if [[ "$MODE" == "global" ]]; then
    dest="$HOME/.openai.md"
  else
    dest="$TARGET_DIR/.openai.md"
  fi

  generate_tier3_config "openai-codex" "$dest" \
    "# OpenAI / Codex Instructions (Saturday Framework)"
}

echo ""
echo "Context Engineering Framework Installer"
echo "========================================"
echo ""
echo "Mode:     $MODE"
if [[ "$MODE" == "project" ]]; then
  echo "Target:   $TARGET_DIR"
fi
echo "Strategy: $(if $USE_COPY; then echo 'copy'; else echo 'symlink'; fi)"
echo "Dry run:  $DRY_RUN"
if [[ -n "$PLATFORM_FILTER" ]]; then
  echo "Platform: $PLATFORM_FILTER"
fi
echo ""

echo "Detecting installed platforms..."
DETECTED_PLATFORMS=($(detect_platforms))
echo "Found: ${DETECTED_PLATFORMS[*]}"
echo ""

INSTALLED_COUNT=0

if should_install "claude-code"; then
  install_claude_code
  ((INSTALLED_COUNT++))
fi

if should_install "cursor"; then
  install_cursor
  ((INSTALLED_COUNT++))
fi

if should_install "windsurf"; then
  install_windsurf
  ((INSTALLED_COUNT++))
fi

if should_install "github-copilot"; then
  install_copilot
  ((INSTALLED_COUNT++))
fi

if should_install "gemini"; then
  install_gemini
  ((INSTALLED_COUNT++))
fi

if should_install "openai-codex"; then
  install_openai
  ((INSTALLED_COUNT++))
fi

echo ""
echo "========================================"
echo "Installation summary"
echo "========================================"
echo "  Agents:    $AGENT_COUNT"
echo "  Skills:    $SKILL_COUNT"
echo "  Rules:     $RULE_COUNT"
echo "  Platforms: $INSTALLED_COUNT configured"
echo ""

if $DRY_RUN; then
  echo "This was a dry run. No files were modified."
  echo "Remove --dry-run to install for real."
else
  echo "Installation complete."
  echo ""
  echo "Next steps:"
  if [[ "$MODE" == "project" ]]; then
    echo "  1. cd $TARGET_DIR"
    echo "  2. Edit CLAUDE.md — update the stack placeholders"
    echo "  3. Write a feature: cp features/TEMPLATE.md features/my-feature.md"
    echo "  4. Launch: run 'claude' and type '/deliver-feature features/my-feature.md'"
  else
    echo "  1. Open any project with your AI tool of choice"
    echo "  2. The framework rules and personas are active globally"
    echo "  3. For project-specific setup: ./install.sh --project /path/to/project"
  fi
fi
echo ""
