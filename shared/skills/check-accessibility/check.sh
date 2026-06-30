#!/usr/bin/env bash

# check.sh

TARGET="${1:-.}"

echo "🔍 Scanning frontend files for Accessibility & Semantic HTML violations in: $TARGET"

# Check if eslint with jsx-a11y is available
if [ -f "node_modules/.bin/eslint" ] && grep -q "jsx-a11y" package.json 2>/dev/null; then
   echo "🏃 Project uses jsx-a11y. Running native linter..."
   # This is a best effort to run the project's linter
   npm run lint || exit 1
   echo "✅ Native linter passed."
   exit 0
fi

echo "⚠️  ESLint jsx-a11y not detected. Running basic naive checks via grep..."

EXIT_CODE=0

# Find typical frontend files
FILES=$(find "$TARGET" -type f \( -name "*.tsx" -o -name "*.jsx" -o -name "*.vue" -o -name "*.svelte" -o -name "*.html" \))

if [ -z "$FILES" ]; then
  echo "⚠️  No frontend UI files found to check."
  exit 0
fi

for FILE in $FILES; do
  # Check for interaction on generic elements
  if grep -nE "<(div|span)[^>]*(onClick|onkeyup|onkeydown|@click)" "$FILE"; then
    echo "❌ ERROR: Found interaction handler on generic <div> or <span> in $FILE"
    echo "   Rule: Semantic HTML over div-soup. Use <button> or <a> for interactive elements."
    EXIT_CODE=1
  fi

  # Rough check for missing alt on img
  if grep -nE "<img[^>]*(?!alt)[^>]*>" "$FILE" | grep -v "alt="; then
    echo "⚠️  WARNING: Found <img> tag potentially missing 'alt' attribute in $FILE"
    # We might not fail the build for a rough grep, but warn
  fi
done

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ Accessibility checks passed!"
else
  echo "❌ Accessibility violations found. Please use semantic HTML."
fi

exit $EXIT_CODE
