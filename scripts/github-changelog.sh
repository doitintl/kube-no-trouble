#!/usr/bin/env sh
# Set strict error checking
set -eou pipefail

GH_ORG="LeMyst"
GH_REPO="kube-no-trouble"
CHANGELOG_TEMPLATE="./scripts/changelog.tmpl"
OUTPUT_FILE="${OUTPUT_FILE:="./changelog.md"}"

# Needed for ci
git config --global --add safe.directory "${GITHUB_WORKSPACE:-$PWD}" || true

RELEASE_SHA="$(curl --silent "https://api.github.com/repos/${GH_ORG}/${GH_REPO}/releases/latest" \
	| jq -r '.target_commitish')"
MASTER_SHA="$(git rev-parse origin/master)"

echo "- Generating changelog (${OUTPUT_FILE})"

changelog-gen -changelog "${CHANGELOG_TEMPLATE}" \
	-owner "${GH_ORG}" -repo "${GH_REPO}" \
	"${RELEASE_SHA}" "${MASTER_SHA}" | tee "${OUTPUT_FILE}"
