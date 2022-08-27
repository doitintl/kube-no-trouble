#!/usr/bin/env sh
# Set strict error checking
set -eou pipefail

BASE_REF="origin/master"
COMMIT_TYPES="build|chore|ci|dep|docs|feat|fix|ref|style|test"
COMMIT_REGEXP="^(${COMMIT_TYPES}): [A-Z]+.{5,}[^.]$"

# Needed for ci
git config --global --add safe.directory "${GITHUB_WORKSPACE:-$PWD}" || true

COMMITS_COUNT="$(git log --oneline --no-merges "${BASE_REF}..HEAD" | wc -l)"

echo "- Checking commit messages: ${COMMITS_COUNT} commits"

BAD_COMMITS="$(git log --oneline --no-merges -E --invert-grep --grep="${COMMIT_REGEXP}" "${BASE_REF}..HEAD")"
if [ -n "${BAD_COMMITS}" ]; then
  echo "${BAD_COMMITS}"
  echo "Error: commit messages do no confirm to required format: ${COMMIT_REGEXP}"
  exit 1
fi
