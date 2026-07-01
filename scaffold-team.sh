#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

TARGET_DIR="${1:-$(pwd)}"

exec "$REPO_DIR/install.sh" --project "$TARGET_DIR" --platform claude-code
