#!/usr/bin/env bash

# NOTE that this file must be executable (0755) in order for the CI to be
# able to run it.

echo "Deploying..."

janus deploy -to "isaac-tests/janus/$(janus version -format %M.%m.x)" -files "*.zip" -key gcp-key.enc.json
