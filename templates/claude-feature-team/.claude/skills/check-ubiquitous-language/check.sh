#!/usr/bin/env bash

# check.sh

TARGET="${1:-.}"
DICTIONARY="DOMAIN_DICTIONARY.md"

echo "🔍 Analyzing Ubiquitous Language terminology in: $TARGET"

if [ ! -f "$DICTIONARY" ]; then
  echo "⚠️  $DICTIONARY not found in the root. Skipping ubiquitous language check."
  echo "✅  (Placeholder pass) The analyst should create DOMAIN_DICTIONARY.md."
  exit 0
fi

# Extract synonyms to avoid from the Entities table in DOMAIN_DICTIONARY.md
# The table format is: | **[Entity Name]** | `[Domain]` | [Description] | `[Synonym]`, `[Synonym]` |
# We'll parse out the last column and extract words wrapped in backticks.
SYNONYMS=$(awk -F'|' '/^\| \*\*.*\*\*/ {print $5}' "$DICTIONARY" | grep -o '`[^`]*`' | tr -d '`')

if [ -z "$SYNONYMS" ]; then
  echo "✅ No prohibited synonyms defined in $DICTIONARY. Check passed!"
  exit 0
fi

VIOLATIONS_FOUND=0

# Use find to get all source files (excluding node_modules, dist, etc.)
FILES=$(find "$TARGET" -type f -not -path "*/node_modules/*" -not -path "*/dist/*" -not -path "*/build/*" -not -name "*.md" 2>/dev/null)

if [ -z "$FILES" ]; then
  # TARGET might be a specific file
  if [ -f "$TARGET" ]; then
    FILES="$TARGET"
  else
    echo "⚠️  No source files found in target: $TARGET"
    exit 0
  fi
fi

for file in $FILES; do
  for term in $SYNONYMS; do
    # Perform case-insensitive search for the exact word boundary for the synonym
    # We use -i for case insensitivity so 'client', 'Client', 'CLIENT' are caught.
    if grep -iqE "\b$term\b" "$file"; then
      echo "❌ Violation in $file: Found prohibited synonym '$term'."
      VIOLATIONS_FOUND=1
    fi
  done
done

if [ $VIOLATIONS_FOUND -eq 1 ]; then
  echo "❌ Ubiquitous language limits exceeded! Please align your naming with $DICTIONARY."
  exit 1
fi

echo "✅ Ubiquitous Language checks passed!"
exit 0
