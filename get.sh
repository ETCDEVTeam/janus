#!/bin/sh

# Sources cited:
# - https://raw.githubusercontent.com/goreleaser/get/master/get

# Use:
#
# Install janus executable to system PATH.
#
# curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash

set -e

TAR_FILE="/tmp/janus.tar.gz"
TAR_FILE_SIG="/tmp/janus.tar.gz.sig"
# It's really annoying that we (have to?) do this Windows workaround.
if [ "$TRAVIS_OS_NAME" = "" ]; then
        TAR_FILE_SIG="/tmp/janus.zip.sig"
fi
DOWNLOAD_URL="https://github.com/ethereumproject/janus/releases/download"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
  # local header
  # # test -z "$GITHUB_TOKEN" || header="-H \"Authorization: token $GITHUB_TOKEN\""
  # # # curl -s $header https://api.github.com/repos/ethereumproject/janus/releases/latest |
  #   grep tag_name |
  #   cut -f4 -d'"'

  # The new and improved sans-GithubAPI-rate-limited curler.
  # https://github.com/goreleaser/goreleaser/issues/157
  curl -sL -o /dev/null -w %{url_effective} https://github.com/ethereumproject/janus/releases/latest | rev | cut -f1 -d'/'| rev
}

download() {
  test -z "$VERSION" && VERSION="$(last_version)"
  test -z "$VERSION" && {
    echo "Unable to get janus version." >&2
    exit 1
  }
  echo "Version: $VERSION"
  rm -f "$TAR_FILE"
  download_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).tar.gz"
  # Check CI for AppVeyor.
  if [ "$TRAVIS_OS_NAME" = "" ]; then
          download_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).zip"
  fi
  echo "Downloading Janus: $download_target"
  curl -s -L -o "$TAR_FILE" \
    "$download_target"

  # Get and verify signature.
  rm -f "$TAR_FILE_SIG"
  sig_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).tar.gz.sig"
  # Check CI for AppVeyor.
  if [ "$TRAVIS_OS_NAME" = "" ]; then
          sig_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).zip.sig"
  fi
  echo "Downloading signature: $sig_target"
  curl -s -L -o "$TAR_FILE_SIG" \
    "$sig_target"

  # Ensure our newly downloaded files exists.
  if ! [ -f "$TAR_FILE" ]; then
          echo "Tar file not found."
          exit 1
  fi
  if ! [ -f "$TAR_FILE_SIG" ]; then
          echo "Tar sig file not found."
          exit 1
  fi
}

# TODO: we may have to download my/someone's signature to use --verify
# not happy about it.
verify() {
  # Ensure we have GPG software
  if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew install gpg2
  elif [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
    sudo apt-get install gnupg2
    # TODO: Windows appveyor
    # How great would it be if the CIs came with gpg2 pre-installed. Really great.
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
