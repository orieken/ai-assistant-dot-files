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

# 3. Copy the Domain Dictionary Template
cp "$TEMPLATE_DIR/DOMAIN_DICTIONARY.template.md" "$TARGET_DIR/"
echo "✅ Domain Dictionary template installed."

# 4. Copy the Subagents, Skills & Rules
mkdir -p "$TARGET_DIR/.claude"
cp -r "$TEMPLATE_DIR/.claude/agents" "$TARGET_DIR/.claude/"
cp -r "$TEMPLATE_DIR/.claude/skills" "$TARGET_DIR/.claude/"
cp -r "$TEMPLATE_DIR/.claude/rules" "$TARGET_DIR/.claude/"
echo "✅ Subagents (Spec Writer, Analyst, Architect, Performance, Data, Developer, Code Reviewer, A11y, Security, QA, SRE, Tech Writer, DevOps) installed."
echo "✅ Skills installed."
echo "✅ Rules installed."

echo
echo "🎉 AI Feature Team successfully scaffolded!"
echo "Next Steps:"
echo "1. Cd into your project:  cd $TARGET_DIR"
echo "2. Edit CLAUDE.md:        Update the [SET IN PROJECT CLAUDE.md] placeholders with your actual stack."
echo "3. Write a feature:       Copy features/TEMPLATE.md to features/my-feature.md"
echo "4. Launch the team:       Run 'claude' and type '/deliver-feature features/my-feature.md'"

echo
echo "--- Deployment Verification ---"
echo "✅ Agents deployed:  $(ls -1 "$TARGET_DIR/.claude/agents/" | grep -c "\.md$")"
echo "✅ Skills deployed:  $(find "$TARGET_DIR/.claude/skills/" -mindepth 1 -maxdepth 1 -type d | wc -l | tr -d ' ')"
echo "✅ Rules deployed:   $(ls -1 "$TARGET_DIR/.claude/rules/" | grep -c "\.md$")"
echo "✅ Templates:        $(ls -1d "$TARGET_DIR/features/TEMPLATE.md" "$TARGET_DIR/DOMAIN_DICTIONARY.template.md" 2>/dev/null | wc -l | tr -d ' ')"
