#!/usr/bin/env bash

# run.sh

echo "🧪 Running Tests and Coverage..."

if [ ! -f "package.json" ]; then
  echo "❌ No package.json found. Cannot run standard tests."
  exit 1
fi

# Determine test command
TEST_CMD="npm test"

# Check if coverage script exists
if grep -q "\"coverage\":" "package.json"; then
  TEST_CMD="npm run coverage"
fi

# Check if it's a cucumber project
if grep -q "cucumber" "package.json"; then
  echo "🥒 Detected Cucumber tests."
fi

# Check if it's a playwright project
if grep -q "playwright" "package.json"; then
  echo "🎭 Detected Playwright tests."
fi

# Execute tests
echo "🏃 Executing: $TEST_CMD"
$TEST_CMD || {
  echo "❌ Tests failed!"
  exit 1
}

echo "✅ All tests passed successfully."

# Check coverage threshold if output is text format
# Note: This is a placeholder since Istanbul/nyc output format varies.
# The AI agent reading the output should verify the 85% rule from the stdout.
echo "📊 Please ensure coverage output above meets the >85% threshold rule."
exit 0
