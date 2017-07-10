#!/usr/bin/env bash

# Sources cited:
# - https://raw.githubusercontent.com/goreleaser/get/master/get

# Use:
#
# Install janus executable to system PATH.
#
# curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash


set -e

TAR_FILE="/tmp/janus.tar.gz"
TAR_FILE_SIG="/tmp/janus.sig"
DOWNLOAD_URL="https://github.com/ethereumproject/janus/releases/download"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
  local header
  test -z "$GITHUB_TOKEN" || header="-H \"Authorization: token $GITHUB_TOKEN\""
  curl -s $header https://api.github.com/repos/ethereumproject/janus/releases/latest |
    grep tag_name |
    cut -f4 -d'"'
}

download() {
  test -z "$VERSION" && VERSION="$(last_version)"
  test -z "$VERSION" && {
    echo "Unable to get janus version." >&2
    exit 1
  }
  rm -f "$TAR_FILE"
  curl -s -L -o "$TAR_FILE" \
    "$DOWNLOAD_URL/$VERSION/janus_$(uname -s)_$(uname -m).tar.gz"

  # Get and verify signature.
  rm -f "$TAR_FILE_SIG"
  curl -s -L -o "$TAR_FILE_SIG" \
    "$DOWNLOAD_URL/$VERSION/janus_$(uname -s)_$(uname -m).sig"
}

verify() {
  # Ensure we have GPG software
  if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew install gpg2
  elif [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    sudo apt-get install gnupg2
    # TODO: Windows appveyor
  fi

  gpg --verify "$TAR_FILE_SIG" "$TAR_FILE"
}

install() {
  tar -xf "$TAR_FILE" -C "$TMPDIR"
  # Ensure executable
  chmod +x "${TMPDIR}/janus"
  # Add to PATH
  if [[ "$TRAVIS_OS_NAME" != "" ]]; then
    export PATH=$PATH:"${TMPDIR}/janus"
  else
    set /p PATH=%PATH%;"${TMPDIR}/janus"
  fi
}

download
verify
install
