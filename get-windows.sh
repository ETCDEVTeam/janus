#!/bin/sh

# Sources cited:
# - https://raw.githubusercontent.com/goreleaser/get/master/get

# Use:
#
# Install janus executable to system PATH.
#
# curl -sL https://raw.githubusercontent.com/ethereumproject/janus/master/get.sh | bash

set -e

TAR_FILE="/tmp/janus.zip"
TAR_FILE_SIG="/tmp/janus.zip.sig"
ISAAC_GPG_FILE="/tmp/isaac-gpg.txt"
ISAAC_GPG_URL="https://raw.githubusercontent.com/ethereumproject/volunteer/master/Volunteer-Public-Keys/isaac.ardis%40gmail.com"
DOWNLOAD_URL="https://github.com/ethereumproject/janus/releases/download"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
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
  download_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).zip"
  echo "Downloading Janus: $download_target"
  curl -s -L -o "$TAR_FILE" \
    "$download_target"

  # Get and verify signature.
  rm -f "$TAR_FILE_SIG"
  sig_target="$DOWNLOAD_URL/$VERSION/janus_${VERSION}_$(uname -s)_$(uname -m).zip.sig"
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

verify() {
  # Verify signature if gpg is available.
  if command -v gpg 2> /dev/null; then
          echo "Downloading Isaac's key file: $ISAAC_GPG_URL"
          curl -s -L -o "$ISAAC_GPG_FILE" \
            "$ISAAC_GPG_URL"

          gpg --import "$ISAAC_GPG_FILE"
          gpg --verify "$TAR_FILE_SIG" "$TAR_FILE"
  fi
}

install() {
  tar -xf "$TAR_FILE" -C "$TMPDIR"
  # Ensure executable
  chmod +x "${TMPDIR}/janus"
  # Add to PATH
  mv "$TMPDIR/janus" "%HOME%/bin/"
}

download
verify
install

echo "Janus installed to: $(which janus)"
