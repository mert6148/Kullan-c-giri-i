#!/usr/bin/env bash
set -euo pipefail

REMOTE_URL=${1:-}
BRANCH=${2:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo main)}
MSG=${3:-"Update project files"}
FORCE=${FORCE:-false}

if ! command -v git >/dev/null 2>&1; then
  echo "Git is not installed. Aborting." >&2
  exit 1
fi

if [ -z "$REMOTE_URL" ]; then
  if git remote get-url origin >/dev/null 2>&1; then
    REMOTE_URL=$(git remote get-url origin)
  else
    echo "No remote specified and origin does not exist. Pass a remote url as first arg." >&2
    exit 1
  fi
fi

git add -A
if git diff --staged --quiet; then
  echo "No changes to commit."
  [ "$FORCE" = true ] || exit 0
fi

if ! git rev-parse --verify HEAD >/dev/null 2>&1; then
  # initial commit
  git commit --allow-empty -m "$MSG"
else
  git commit -m "$MSG" || true
fi

git push -u origin "$BRANCH"
echo "Pushed to origin/$BRANCH"
