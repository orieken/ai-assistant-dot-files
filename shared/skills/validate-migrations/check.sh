#!/usr/bin/env bash

# check.sh

TARGET="${1:-.}"

echo "🔍 Validating migration files for destructive operations in: $TARGET"

# Check for common migration file extensions
FILES=$(find "$TARGET" -type f \( -name "*.sql" -o -name "*.up.sql" -o -name "*.ts" -o -name "*.py" \) | grep -i "migration")

if [ -z "$FILES" ]; then
  # Try just looking for SQL files if "migration" isn't in the path
  FILES=$(find "$TARGET" -type f -name "*.sql")
fi

if [ -z "$FILES" ]; then
  echo "⚠️  No obvious migration files found to check."
  exit 0
fi

# List of forbidden operations in a standard "Expand" phase
DESTRUCTIVE_PATTERN="DROP TABLE|DROP COLUMN|RENAME COLUMN|ALTER TABLE.*DROP|DELETE FROM"
NOT_NULL_WITHOUT_DEFAULT="NOT NULL(?!.*DEFAULT)" # Simplistic check, harder to do purely with grep

EXIT_CODE=0

for FILE in $FILES; do
  # Check if the file is explicitly marked as a "Contract" migration, which allows destructive ops
  if grep -qi "contract phase" "$FILE" || [[ "$FILE" == *"contract"* ]]; then
    echo "⚠️  $FILE is marked as a Contract migration. Destructive operations allowed."
    continue
  fi

  # Search for destructive patterns
  if grep -qiE "$DESTRUCTIVE_PATTERN" "$FILE"; then
    echo "❌ ERROR: Destructive database operation found in $FILE"
    echo "   Migrations MUST follow the Expand/Contract pattern."
    echo "   Found: $(grep -inE "$DESTRUCTIVE_PATTERN" "$FILE")"
    EXIT_CODE=1
  fi
done

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ Migration validation passed!"
else
  echo "❌ Migration validation failed. Please refactor to avoid downtime."
fi

exit $EXIT_CODE
