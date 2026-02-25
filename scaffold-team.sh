#!/usr/bin/env bash

# scaffold-team.sh
# Quickly deploy the Claude Code Multi-Agent "Feature Team" to any workspace.

set -e

# Determine the absolute path to this repository
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$REPO_DIR/templates/claude-feature-team"

echo "🤖 Scaffolding Claude Feature Team..."

# Check if a target directory was provided
if [ -z "$1" ]; then
    TARGET_DIR="$(pwd)"
    echo "No target directory specified. Installing to current directory: $TARGET_DIR"
else
    TARGET_DIR="$1"
    mkdir -p "$TARGET_DIR"
    echo "Installing to specified directory: $TARGET_DIR"
fi

# 1. Copy the Orchestrator (CLAUDE.md)
if [ -f "$TARGET_DIR/CLAUDE.md" ]; then
    echo "⚠️  CLAUDE.md already exists in target. Backing up to CLAUDE.md.bak"
    mv "$TARGET_DIR/CLAUDE.md" "$TARGET_DIR/CLAUDE.md.bak"
fi
cp "$TEMPLATE_DIR/CLAUDE.md" "$TARGET_DIR/"
echo "✅ Orchestrator (CLAUDE.md) installed."

# 2. Copy the Feature Templates
mkdir -p "$TARGET_DIR/features"
cp "$TEMPLATE_DIR/features/TEMPLATE.md" "$TARGET_DIR/features/"
echo "✅ Feature templates installed."

# 3. Copy the Subagents & Skills
mkdir -p "$TARGET_DIR/.claude"
cp -r "$TEMPLATE_DIR/.claude/agents" "$TARGET_DIR/.claude/"
cp -r "$TEMPLATE_DIR/.claude/skills" "$TARGET_DIR/.claude/"
echo "✅ Subagents (Analyst, Dev, QA, Tech Writer, DevOps) installed."
echo "✅ '/deliver-feature' skill installed."

echo
echo "🎉 AI Feature Team successfully scaffolded!"
echo "Next Steps:"
echo "1. Cd into your project:  cd $TARGET_DIR"
echo "2. Edit CLAUDE.md:        Update the [SET IN PROJECT CLAUDE.md] placeholders with your actual stack."
echo "3. Write a feature:       Copy features/TEMPLATE.md to features/my-feature.md"
echo "4. Launch the team:       Run 'claude' and type '/deliver-feature features/my-feature.md'"
