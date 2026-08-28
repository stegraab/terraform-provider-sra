#!/usr/bin/env bash
# Copyright 2026 Stegra AB
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

upstream_remote="${UPSTREAM_REMOTE:-upstream}"
upstream_branch="${UPSTREAM_BRANCH:-main}"
mirror_branch="${MIRROR_BRANCH:-upstream-main}"

if [[ -n "$(git status --porcelain)" ]]; then
  echo "The working tree must be clean before syncing upstream." >&2
  exit 1
fi

if ! git remote get-url "$upstream_remote" >/dev/null 2>&1; then
  echo "Missing '$upstream_remote' remote." >&2
  echo "Add it with: git remote add $upstream_remote https://github.com/BeyondTrust/terraform-provider-sra.git" >&2
  exit 1
fi

if [[ "$(git branch --show-current)" == "$mirror_branch" ]]; then
  echo "Check out a branch other than '$mirror_branch' before syncing." >&2
  exit 1
fi

git fetch --prune --tags "$upstream_remote"
git fetch origin "$mirror_branch" 2>/dev/null || true
git branch --force "$mirror_branch" "refs/remotes/$upstream_remote/$upstream_branch"
git push --force-with-lease origin "refs/heads/$mirror_branch:refs/heads/$mirror_branch"

echo "Updated origin/$mirror_branch from $upstream_remote/$upstream_branch."
echo "Upstream tags were fetched locally and were not pushed to origin."
