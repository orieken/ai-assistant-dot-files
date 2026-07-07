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
  --global              Symlink configs to ~/  (always current after a git pull)
  --project <path>      Symlink configs into a target project directory (also always current --
                         the target keeps depending on this repo's checkout staying where it is)
  --project             Symlink configs into the current directory

Options:
  --copy                Use real copies instead of symlinks -- required on Windows/WSL (symlinks need
                         elevated permissions there), and also the way to make a --project install
                         independent of this repo's checkout (a symlinked project install breaks if
                         this repo is later moved or deleted)
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
    elif $USE_COPY && [[ -e "$dest" ]] && diff -rq "$src" "$dest" > /dev/null 2>&1; then
      skip "$dest (already copied, content identical)"
      return
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

install_antigravity() {
  log ""
  log "--- Gemini Antigravity (confirmed 2026-07-02 via tests/platform-verification/antigravity.md) ---"

  if [[ "$MODE" == "global" ]]; then
    # Confirmed: Antigravity's global skills root is ~/.gemini/config/skills/ — NOT ~/.agents/skills/
    # (that was an earlier, unconfirmed guess). Global rules use a single ~/.gemini/GEMINI.md file per
    # secondary-source research, not yet confirmed the way the skills path and AGENTS.md have been —
    # left unhandled here rather than generating unconfirmed content.
    if ! $DRY_RUN; then mkdir -p "$HOME/.gemini/config" 2>/dev/null || true; fi
    link_or_copy "$SHARED_DIR/skills" "$HOME/.gemini/config/skills"
  else
    # Confirmed: project-level AGENTS.md is read for rules (generated separately, see
    # generate_agents_md in scripts/generate-configs.sh). Project-level .agents/skills/ and
    # .agents/rules/ are the documented convention but weren't directly exercised by the 2026-07-02
    # test (no .agents/ existed at session start, so it fell back to the global skills root) — kept
    # since they match the codelab's documented project-scope convention and don't contradict anything
    # confirmed so far.
    if ! $DRY_RUN; then mkdir -p "$TARGET_DIR/.agents" 2>/dev/null || true; fi
    link_or_copy "$SHARED_DIR/skills" "$TARGET_DIR/.agents/skills"
    link_or_copy "$SHARED_DIR/rules" "$TARGET_DIR/.agents/rules"
  fi
}

install_cursor() {
  log ""
  log "--- Cursor (native agents/skills confirmed 2026-07-06, rules still generated inline) ---"

  # Cursor reads .cursor/agents/*.md and .cursor/skills/*/SKILL.md using the same open standard
  # shared/agents/ and shared/skills/ already follow -- symlink directly instead of flattening into
  # .mdc personas (the retired Epic 11 workaround). Rules still go through generate-configs.sh
  # (install_generated_configs) since Cursor Rules can't follow file references, only agents/skills can.
  if [[ "$MODE" == "global" ]]; then
    local cursor_dir="$HOME/.cursor"
    if ! $DRY_RUN; then mkdir -p "$cursor_dir" 2>/dev/null || true; fi
    link_or_copy "$SHARED_DIR/agents" "$cursor_dir/agents"
    link_or_copy "$SHARED_DIR/skills" "$cursor_dir/skills"
  else
    local cursor_dir="$TARGET_DIR/.cursor"
    if ! $DRY_RUN; then mkdir -p "$cursor_dir" 2>/dev/null || true; fi
    link_or_copy "$SHARED_DIR/agents" "$cursor_dir/agents"
    link_or_copy "$SHARED_DIR/skills" "$cursor_dir/skills"
  fi
}

install_claude_code() {
  log ""
  log "--- Claude Code (Tier 1: Full) ---"

  if [[ "$MODE" == "global" ]]; then
    local claude_dir="$HOME/.claude"
    if ! $DRY_RUN; then mkdir -p "$claude_dir" 2>/dev/null || true; fi
    link_or_copy "$SHARED_DIR/agents" "$claude_dir/agents"
    link_or_copy "$SHARED_DIR/skills" "$claude_dir/skills"
    link_or_copy "$SHARED_DIR/rules" "$claude_dir/rules"
  else
    local claude_dir="$TARGET_DIR/.claude"
    if ! $DRY_RUN; then mkdir -p "$claude_dir" 2>/dev/null || true; fi
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

install_generated_configs() {
  local output_dir
  if [[ "$MODE" == "global" ]]; then
    output_dir="$HOME"
  else
    output_dir="$TARGET_DIR"
  fi

  local gen_args=("--output" "$output_dir")
  if $DRY_RUN; then
    gen_args+=("--dry-run")
  fi
  if [[ -n "$PLATFORM_FILTER" ]]; then
    gen_args+=("--platform" "$PLATFORM_FILTER")
  fi

  "$REPO_DIR/scripts/generate-configs.sh" "${gen_args[@]}"
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

if should_install "claude-code"; then
  install_claude_code
fi

if should_install "cursor"; then
  install_cursor
fi

if should_install "gemini"; then
  install_antigravity
fi

install_generated_configs

echo ""
echo "========================================"
echo "Installation summary"
echo "========================================"
echo "  Agents:    $AGENT_COUNT"
echo "  Skills:    $SKILL_COUNT"
echo "  Rules:     $RULE_COUNT"
echo "  Platforms: ${#DETECTED_PLATFORMS[@]} (${DETECTED_PLATFORMS[*]})"
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

if ! $DRY_RUN; then
  echo ""
  echo "========================================"
  echo "Framework health check"
  echo "========================================"
  bash "$REPO_DIR/scripts/health-check.sh" || echo "(health-check reported issues — see above; run 'bash scripts/health-check.sh --verbose --fix' for detail and auto-repair)"
fi

if $SHOW_TOUR && ! $DRY_RUN; then
  echo ""
  echo "========================================"
  echo "Onboarding tour"
  echo "========================================"
  echo "This script can't invoke an AI skill itself — it's plain bash. To get the tour, open your AI"
  echo "tool in this repo (or wherever you installed to) and ask it to run the 'onboard' skill:"
  echo ""
  echo "  > /onboard"
  echo ""
  echo "or in a tool without slash commands:"
  echo ""
  echo "  Act as the onboard skill described in shared/skills/onboard/SKILL.md and give me a tour."
  echo ""
  echo "It covers: the three context layers (rules/agents/skills), how to invoke an agent, how to"
  echo "trigger a skill, how to run a full pipeline, and the approval gates that pause for confirmation."
  echo "It ends by pointing you at shared/templates/my-first-feature.md — a complete, pre-written"
  echo "feature spec you can run through /deliver-feature immediately to see the whole pipeline for real."
fi
echo ""
