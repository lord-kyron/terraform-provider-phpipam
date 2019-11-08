#!/usr/bin/env bash

declare -A args=(
  # The GPG key to sign the binaries with.
  [keyid]="9E2D8AFF3BE44244"

  # The name of the target binary.
  [binname]="terraform-provider-phpipam"

  # The build target
  [build_target]="build"

  # The OSes to release for.
  [target_os]="windows darwin linux"

  # The arches to release for.
  [target_arch]="amd64"
)

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

# get_release fetches the last 10 releases and asks the user to pick one to
# release for. The script then switches to this branch.
get_release() {
  # gets the last ten releases.
  versions=($(git tag --sort "-v:refname" | egrep '^v[0-9]+\.[0-9]+\.[0-9]+' | head -n 10))

  if [ -z "${versions[*]}" ]; then
    message error "No non-prerelease versions available at this time. Please release a version first"
    message error "and ensure the build has a matching vMAJOR.MINOR.PATCH tag."
    exit 1
  fi

  message info "Select a build version to release with:"
  for (( n=0; n<${#versions[@]}; n+=1 )); do
    message info "${n}: ${versions[$n]}"
  done

  read -r -p "$(message info "Enter a number (Default 0, CTRL-C Cancels):")" selection_number

  if [ -z "${selection_number}" ]; then
    selection_number="0"
  fi

  if ! [ "${selection_number}" -ge "0" ] 2>/dev/null || ! [ "${selection_number}" -lt "${#versions[@]}" ] 2>/dev/null; then
    message error "Invalid selection. Please enter a valid selection from 0-$((${#versions[@]}-1))."
    exit 1
  fi

  git checkout "${versions[$selection_number]}"
  echo "${versions[$selection_number]}"
}

# build performs the build
build() {
  (
    make "${args[build_target]}"
  )
  local __status=$?
  if [ "${__status}" != "0" ]; then
    message error "ERROR: Release build exited with code ${__status}"
    exit 1
  fi
}

# move_files moves the files to a dist/ directory within the pkg/ 
# directory of the build root, with the release tag supplied to the
# function. The binaries are also compressed.
move_files() {
  local __release="$1"
  local __os=""
  local __arch=""
  set -e 
  rm -rf pkg/dist
  mkdir -p pkg/dist
  for __os in ${args[target_os]}; do
    for __arch in ${args[target_arch]}; do
      if [ "${__os}" == "windows" ]; then
        zip -j "pkg/dist/${args[binname]}_${__release}_${__os}_${__arch}.zip" "pkg/${__os}_${__arch}/${args[binname]}.exe"
      else
        zip -j "pkg/dist/${args[binname]}_${__release}_${__os}_${__arch}.zip" "pkg/${__os}_${__arch}/${args[binname]}"
      fi
    done
  done
  set +e
}

# sign_files creates a SHA256SUMS files for the releases, and creates
# a detatched GPG signature.
sign_files() {
  local __release="$1"
  (
    set -e
    cd pkg/dist
    shasum -a256 -- * > "./${args[binname]}_${__release}_SHA256SUMS"
    gpg --default-key "${args[keyid]}" --detach-sig "./${args[binname]}_${__release}_SHA256SUMS"
  )
  local __status=$?
  if [ "${__status}" != "0" ]; then
    message error "ERROR: Release signing exited with code ${__status}"
    exit 1
  fi
}

## Main
if [ "$(git rev-parse HEAD)" != "$(git rev-parse master)" ]; then
  message error "ERROR: Current HEAD does not match master."
  message error "Please switch HEAD back to master before running this script."
  exit 1
fi

release="$(get_release)"

build "${release}"
move_files "${release}"
sign_files "${release}"
git checkout master
