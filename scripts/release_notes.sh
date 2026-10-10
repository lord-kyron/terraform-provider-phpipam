#!/usr/bin/env bash

# release_notes.sh prints the CHANGELOG.md section for a release, for use as
# the GitHub release notes. "*" bullets are rewritten as "-" bullets.
#
# Usage: scripts/release_notes.sh VERSION   (with or without a leading "v")
#
# Exits non-zero if CHANGELOG.md has no section for VERSION, or the section is
# empty, so a release is never published without notes.

set -euo pipefail

version="${1:?usage: $0 VERSION}"
version="${version#v}"
changelog="$(dirname "$0")/../CHANGELOG.md"

notes=$(awk -v ver="${version}" '
  /^## / { if (found) exit; found = ($2 == ver); next }
  found {
    if (sub(/^[[:space:]]*\* /, "- ")) { print; next }
    print
  }
' "${changelog}")

if [ -z "$(echo "${notes}" | tr -d '[:space:]')" ]; then
  echo "No CHANGELOG.md entries found for ${version}" >&2
  exit 1
fi

echo "${notes}"
