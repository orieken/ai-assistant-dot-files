#!/usr/bin/env bash
set -euo pipefail

# Enterprise KI sync for the Context Engineering Framework.
# See docs/adrs/ADR-003-enterprise-memory-sync.md for design rationale.
#
# Usage:
#   sync-memory.sh [pull|push] [--confirm] [--config <path>]
#
# pull  (default) — diff and optionally apply org KIs into shared/knowledge/
# push            — promote local .claude/knowledge/ KIs to the org repo via PR

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SHARED_DIR="$REPO_DIR/shared"
KI_SCHEMA="$SHARED_DIR/schemas/ki-frontmatter.schema.json"
VALIDATE_SCRIPT="$REPO_DIR/scripts/validate-frontmatter.py"

SUBCOMMAND="pull"
CONFIRM=false
CONFIG_PATH=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    pull|push) SUBCOMMAND="$1"; shift ;;
    --confirm) CONFIRM=true; shift ;;
    --config)  CONFIG_PATH="$2"; shift 2 ;;
    -h|--help)
      echo "Usage: sync-memory.sh [pull|push] [--confirm] [--config <path>]"
      exit 0
      ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done

if [[ -z "$CONFIG_PATH" ]]; then
  CONFIG_PATH="$REPO_DIR/.claude/sync-config.yaml"
fi

if [[ ! -f "$CONFIG_PATH" ]]; then
  cat <<'HELP'
Error: .claude/sync-config.yaml not found.

Create it at the root of your framework checkout:

  memory_sync:
    org_repo: git@github.com:<your-org>/knowledge-hub.git
    cache_dir: ~/.claude/sync-cache
    push_pr_base: main

See docs/adrs/ADR-003-enterprise-memory-sync.md for details.
HELP
  exit 1
fi

# Simple key reader for the minimal YAML structure (no pyyaml required).
read_config_key() {
  local key="$1"
  local default="${2:-}"
  local val
  val=$(grep "^\s*${key}:" "$CONFIG_PATH" 2>/dev/null | head -1 | sed "s/.*${key}:\s*//" | tr -d '"' || true)
  echo "${val:-$default}"
}

ORG_REPO=$(read_config_key "org_repo" "")
CACHE_BASE=$(read_config_key "cache_dir" "~/.claude/sync-cache")
PUSH_BASE=$(read_config_key "push_pr_base" "main")
CACHE_BASE="${CACHE_BASE/#\~/$HOME}"

if [[ -z "$ORG_REPO" ]]; then
  echo "Error: memory_sync.org_repo not set in $CONFIG_PATH"
  exit 1
fi

# Auth: MEMORY_SYNC_TOKEN env var converts SSH URL to HTTPS for enterprises
# where SSH is disabled.
if [[ -n "${MEMORY_SYNC_TOKEN:-}" && "$ORG_REPO" =~ ^git@github.com: ]]; then
  ORG_REPO=$(echo "$ORG_REPO" | sed "s|git@github.com:|https://$MEMORY_SYNC_TOKEN@github.com/|")
fi

SLUG=$(echo "$ORG_REPO" | sed 's|.*[:/]||; s|\.git$||')
CACHE_DIR="$CACHE_BASE/$SLUG"
LAST_SYNC_FILE="$CACHE_DIR/.last-sync"

log()  { echo "  $1"; }
ok()   { echo "  [ok] $1"; }
info() { echo "  [info] $1"; }
warn() { echo "  [warn] $1"; }

# --- Validate a KI file against the frontmatter schema (best-effort) --------
validate_ki() {
  local ki_file="$1"
  if command -v python3 &>/dev/null && [[ -f "$VALIDATE_SCRIPT" && -f "$KI_SCHEMA" ]]; then
    python3 "$VALIDATE_SCRIPT" "$KI_SCHEMA" "$ki_file" > /dev/null 2>&1
  fi
}

# --- Resolve target KI directory inside a repo root -------------------------
find_ki_dir() {
  local root="$1"
  for candidate in "$root/shared/knowledge" "$root/knowledge"; do
    if [[ -d "$candidate" ]]; then
      echo "$candidate"; return
    fi
  done
  echo ""
}

# --- Update or create the org cache (clone on first run, fetch on repeat) ---
refresh_cache() {
  if [[ -d "$CACHE_DIR/.git" ]]; then
    log "Updating cached org repo ($SLUG)..."
    git -C "$CACHE_DIR" fetch origin --quiet
    git -C "$CACHE_DIR" reset --hard "origin/$PUSH_BASE" --quiet 2>/dev/null || \
      git -C "$CACHE_DIR" reset --hard origin/HEAD --quiet
    ok "cache updated"
  else
    log "Cloning org repo (first run)..."
    mkdir -p "$CACHE_BASE"
    git clone --depth=1 "$ORG_REPO" "$CACHE_DIR" --quiet
    ok "cloned to $CACHE_DIR"
  fi
}

# === PULL ===================================================================
pull_kis() {
  echo ""
  echo "=== Memory Sync: Pull ==="
  echo "Org repo: $ORG_REPO"
  echo "Cache:    $CACHE_DIR"
  echo ""

  refresh_cache

  local org_ki_dir
  org_ki_dir=$(find_ki_dir "$CACHE_DIR")
  if [[ -z "$org_ki_dir" ]]; then
    warn "No knowledge/ or shared/knowledge/ directory found in org repo."
    warn "The org repo should contain KIs under knowledge/ or shared/knowledge/."
    exit 1
  fi

  local local_ki_dir="$SHARED_DIR/knowledge"
  local project_ki_dir="$REPO_DIR/.claude/knowledge"

  # Phase 1: collision detection (org name slug vs project KI slug)
  local collisions=0
  for org_ki in "$org_ki_dir"/*.md; do
    [[ -f "$org_ki" ]] || continue
    base="$(basename "$org_ki")"
    [[ "$base" == "README.md" ]] && continue
    if [[ -f "$project_ki_dir/$base" ]]; then
      org_name=$(grep '^name:' "$org_ki" | head -1 | sed 's/name: *//' || true)
      local_name=$(grep '^name:' "$project_ki_dir/$base" | head -1 | sed 's/name: *//' || true)
      if [[ "$org_name" == "$local_name" ]]; then
        echo "  COLLISION  $base (org name '$org_name' matches .claude/knowledge/$base)"
        ((collisions++)) || true
      fi
    fi
  done

  if [[ "$collisions" -gt 0 ]]; then
    echo ""
    echo "ERROR: $collisions slug collision(s) between org KIs and .claude/knowledge/."
    echo "Rename the conflicting project or org KI before running --confirm."
    exit 1
  fi

  # Phase 2: diff
  local to_add=()
  local to_update=()

  echo "Diffing org KIs against local shared/knowledge/..."
  echo ""

  for org_ki in "$org_ki_dir"/*.md; do
    [[ -f "$org_ki" ]] || continue
    base="$(basename "$org_ki")"
    [[ "$base" == "README.md" ]] && continue

    local_ki="$local_ki_dir/$base"
    if [[ ! -f "$local_ki" ]]; then
      echo "  + ADD     $base"
      to_add+=("$base")
    elif ! diff -q "$org_ki" "$local_ki" > /dev/null 2>&1; then
      echo "  ~ UPDATE  $base"
      diff "$local_ki" "$org_ki" | grep -E '^[<>]' | head -4 | sed 's/^/             /' || true
      to_update+=("$base")
    fi
  done

  local changes=$(( ${#to_add[@]} + ${#to_update[@]} ))

  if [[ "$changes" -eq 0 ]]; then
    echo "  (no changes — shared/knowledge/ is up to date with org repo)"
    date -u "+%Y-%m-%dT%H:%M:%SZ" > "$LAST_SYNC_FILE"
    echo ""
    echo "Last-sync timestamp updated."
    return
  fi

  echo ""
  echo "$changes change(s) pending."

  if ! $CONFIRM; then
    echo ""
    echo "Dry run — no files changed. Re-run with --confirm to apply:"
    echo "  bash scripts/sync-memory.sh pull --confirm"
    return
  fi

  # Phase 3: apply
  echo ""
  echo "Applying changes..."

  for base in "${to_add[@]}" "${to_update[@]}"; do
    org_ki="$org_ki_dir/$base"
    local_ki="$local_ki_dir/$base"

    if ! validate_ki "$org_ki" 2>/dev/null; then
      warn "SKIP $base — failed KI frontmatter validation (malformed org KI)"
      continue
    fi

    # Stamp sync_source + sync_pulled + sync_commit_sha if not already present.
    # sync_commit_sha records the org repo's HEAD SHA at pull time, enabling auditors
    # to trace which org-repo state introduced this KI (THREAT_MODEL.md F-04).
    local org_head_sha
    org_head_sha=$(git -C "$CACHE_DIR" rev-parse HEAD 2>/dev/null || echo "unknown")
    if ! grep -q '^sync_source:' "$org_ki"; then
      local today
      today=$(date +%Y-%m-%d)
      python3 - "$org_ki" "$local_ki" "$ORG_REPO" "$today" "$org_head_sha" <<'PYEOF'
import sys
src, dst, repo, today, sha = sys.argv[1:]
content = open(src).read()
parts = content.split('---', 2)
if len(parts) >= 3:
    parts[1] = parts[1].rstrip('\n') + (
        f'\nsync_source: {repo}\nsync_pulled: {today}\nsync_commit_sha: {sha}\n'
    )
    content = '---'.join(parts)
open(dst, 'w').write(content)
PYEOF
    else
      cp "$org_ki" "$local_ki"
    fi
    ok "applied $base"
  done

  date -u "+%Y-%m-%dT%H:%M:%SZ" > "$LAST_SYNC_FILE"

  echo ""
  echo "Sync complete. $changes KI(s) written to shared/knowledge/."
  echo "Review, then commit:"
  echo "  git add shared/knowledge/ && git commit -m 'chore(knowledge): sync KIs from org repo'"
}

# === PUSH ===================================================================
push_kis() {
  echo ""
  echo "=== Memory Sync: Push ==="
  echo "Org repo: $ORG_REPO"
  echo ""

  local project_ki_dir="$REPO_DIR/.claude/knowledge"
  if [[ ! -d "$project_ki_dir" ]]; then
    info "No .claude/knowledge/ directory — nothing to push."
    exit 0
  fi

  # Promotion heuristic: KIs created more than 30 days ago
  local thirty_days_ago
  thirty_days_ago=$(python3 -c "
import datetime
print((datetime.date.today() - datetime.timedelta(days=30)).isoformat())
" 2>/dev/null || date -v-30d +%Y-%m-%d 2>/dev/null || echo "1970-01-01")

  local candidates=()
  echo "Candidate KIs (created > 30 days ago, excluding already-synced):"

  for ki in "$project_ki_dir"/*.md; do
    [[ -f "$ki" ]] || continue
    base="$(basename "$ki")"
    [[ "$base" == "README.md" ]] && continue

    # Skip if already came from sync
    grep -q '^sync_source:' "$ki" && continue

    created=$(grep '^created:' "$ki" | head -1 | sed 's/created: *//' | tr -d '"' || true)
    if [[ -z "$created" || "$created" < "$thirty_days_ago" ]]; then
      name=$(grep '^name:' "$ki" | head -1 | sed 's/name: *//' || true)
      echo "  + $base  (name: $name, created: ${created:-unknown})"
      candidates+=("$base")
    fi
  done

  if [[ "${#candidates[@]}" -eq 0 ]]; then
    info "No eligible KIs to push (none older than 30 days, or all already synced from org)."
    exit 0
  fi

  echo ""
  if ! $CONFIRM; then
    echo "Dry run — re-run with --confirm to create a PR:"
    echo "  bash scripts/sync-memory.sh push --confirm"
    return
  fi

  local tmp_dir
  tmp_dir=$(mktemp -d)
  # shellcheck disable=SC2064
  trap "rm -rf '$tmp_dir'" EXIT

  log "Cloning org repo to temp dir..."
  git clone --depth=1 "$ORG_REPO" "$tmp_dir" --quiet
  ok "cloned"

  local branch_name="ki-sync/$(basename "$REPO_DIR")-$(date +%Y%m%d)"
  git -C "$tmp_dir" checkout -b "$branch_name" --quiet

  local org_ki_target
  org_ki_target=$(find_ki_dir "$tmp_dir")
  if [[ -z "$org_ki_target" ]]; then
    mkdir -p "$tmp_dir/shared/knowledge"
    org_ki_target="$tmp_dir/shared/knowledge"
  fi

  for base in "${candidates[@]}"; do
    cp "$project_ki_dir/$base" "$org_ki_target/$base"
    ok "staged $base"
  done

  git -C "$tmp_dir" add "$org_ki_target" --quiet
  git -C "$tmp_dir" commit -m "feat(knowledge): promote KIs from $(basename "$REPO_DIR")" --quiet

  git -C "$tmp_dir" push origin "$branch_name" --quiet

  # Open PR via gh CLI when available, otherwise print instructions
  if command -v gh &>/dev/null && [[ "$ORG_REPO" =~ github.com ]]; then
    local org_slug
    org_slug=$(echo "$ORG_REPO" | sed 's|.*github.com[:/]||; s|\.git$||')
    local pr_body
    pr_body="## KI Promotion from \`$(basename "$REPO_DIR")\`

$(for base in "${candidates[@]}"; do
  name=$(grep '^name:' "$project_ki_dir/$base" | head -1 | sed 's/name: *//' || true)
  echo "- \`$base\` — $name"
done)

Promoted via \`sync-memory.sh push\` on $(date +%Y-%m-%d)."

    gh pr create \
      --repo "$org_slug" \
      --head "$branch_name" \
      --base "$PUSH_BASE" \
      --title "feat(knowledge): promote KIs from $(basename "$REPO_DIR")" \
      --body "$pr_body"
    ok "PR opened in $org_slug"
  else
    echo ""
    echo "Branch pushed. Open a PR manually:"
    echo "  Org repo: $ORG_REPO"
    echo "  Branch:   $branch_name"
    echo "  Base:     $PUSH_BASE"
  fi
}

# === Dispatch ===============================================================
case "$SUBCOMMAND" in
  pull) pull_kis ;;
  push) push_kis ;;
  *)    echo "Unknown subcommand: '$SUBCOMMAND' — use 'pull' or 'push'"; exit 1 ;;
esac
