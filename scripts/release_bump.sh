#!/usr/bin/env bash

# The release from CHANGELOG.md.
release=$(head -n 1 CHANGELOG.md | awk '{print $2}')

# message prints text with a color, redirected to stderr in the event of
# warning or error messages.
message() {
  declare -A __colors=(
    ["error"]="31"   # red
    ["warning"]="33" # yellow
    ["begin"]="32"   # green
    ["ok"]="32"      # green
    ["info"]="1"     # bold
    ["reset"]="0"    # here just to note reset code
  )
  local __type="$1"
  local __message="$2"
  if [ -z "${__colors[$__type]}" ]; then
    __type="info"
  fi
  if [[ ! "${__type}" =~ ^(warning|error)$ ]]; then
    echo -e "\e[${__colors[$__type]}m${__message}\e[0m" 1>&2
  else
    echo -e "\e[${__colors[$__type]}m${__message}\e[0m"
  fi
}

if [[ "${release}" == *"-pre" ]]; then
  message error "This is a pre-release - release aborted." >&2
  message error "Please update the first line in CHANGELOG.md to a version without the -pre tag."
  exit 1
fi

semver=(${release//./ })

for n in 0 1 2; do 
  if ! [ "${semver[$n]}" -eq "${semver[$n]}" ]; then
    message error "${release} is not a proper semantic-versioned release." >&2
    message error "Please update the first line in CHANGELOG.md to a numeric MAJOR.MINOR.PATCH version."
    exit 1
  fi
done

changelog_status=$(git status -s | grep CHANGELOG.md)

set -e

if [ "${changelog_status}" == " M CHANGELOG.md" ]; then
  git add CHANGELOG.md
  changelog_status=$(git status -s | grep CHANGELOG.md)
fi

if [ "${changelog_status}" == "M  CHANGELOG.md" ]; then
  message begin "==> Committing CHANGELOG.md <=="
  git commit -m "$(echo -e "Release v${release}\n\nSee CHANGELOG.md for more details.")"
fi


message begin "==> Tagging Release v${release} <=="
git tag "v${release}" -m "$(echo -e "Release v${release}\n\nSee CHANGELOG.md for more details.")"

new_prerelease="${semver[0]}.${semver[1]}.$((semver[2]+1))-pre"

message begin "==> Bumping CHANGELOG.md to Release v${new_prerelease} <=="
echo -e "## ${new_prerelease}\n\nBumped version for dev.\n\n$(cat CHANGELOG.md)" > CHANGELOG.md

git add CHANGELOG.md
git commit -m "Bump CHANGELOG.md to v${new_prerelease}"

message begin "==> Pushing Commits and Tags <=="
git push origin master
git push origin --tags

message ok "\nRelease v${release} successful."
