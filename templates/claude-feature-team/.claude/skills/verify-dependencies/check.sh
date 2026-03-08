#!/usr/bin/env bash

# check.sh

TARGET="${1:-.}"

echo "🔍 Verifying clean architecture dependency boundaries in: $TARGET"

if [ -f "dependency-cruiser.js" ] && [ -f "node_modules/.bin/depcruise" ]; then
    echo "🏃 Running dependency-cruiser validation..."
    node_modules/.bin/depcruise --validate dependency-cruiser.js "$TARGET" || exit 1
    echo "✅ Dependency logic validated successfully via dependency-cruiser."
    exit 0
fi

echo "⚠️  dependency-cruiser not found. Running naive checks."

EXIT_CODE=0

# Example target structures:
# src/domain/ (Entities)
# src/usecases/ (Application)
# src/adapters/ (Infrastructure/Web)

# Basic check: Ensuring Domain Layer Entities NEVER import external libraries (other than purely functional ones) or the DB layer
# Find imports matching "import ... from '...';" or "require('...');"
FILES=$(find "$TARGET" -type f \( -name "*.ts" -o -name "*.py" -o -name "*.java" -o -name "*.go" \))

if [ -z "$FILES" ]; then
  echo "⚠️  No source code files found to check."
  exit 0
fi

for FILE in $FILES; do
   # If the file is in the "domain" or "entity" path, we enforce strictness
   path=$(echo "$FILE" | tr '[:upper:]' '[:lower:]')
   if [[ "$path" == *"domain/"* ]] || [[ "$path" == *"entities/"* ]]; then
     
     # Look for obvious infrastructure imports
     if grep -nE "import.*from.*(express|axios|mysql|pg|typeorm|mongoose|sequelize|http|request|fetch)" "$FILE"; then
         echo "❌ ERROR: Found infrastructure import in the Domain layer ($FILE)"
         echo "   Rule: Domain Layer must have zero external dependencies."
         EXIT_CODE=1
     fi

     # Look for imports going 'up' the directory tree (often an indicator of importing an outer layer adapter)
     # This is naive but works well enough as an AI guardrail
     if grep -nE "import.*from.*(\.\./\.\./adapters|\.\./\.\./infrastructure|\.\./\.\./controllers)" "$FILE"; then
         echo "❌ ERROR: Found import from outer layers in the Domain layer ($FILE)"
         echo "   Rule: Dependencies ALWAYS point inward."
         EXIT_CODE=1
     fi
   fi
done

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ Architecture boundary checks passed!"
else
  echo "❌ Architecture rules violated! Dependencies must point inward."
fi

exit $EXIT_CODE
