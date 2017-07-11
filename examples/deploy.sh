#!/usr/bin/env bash

echo "Deploying..."

janus deploy -to "isaac-tests/janus/$(janus version -format %M.%m.x)" -files "*.zip" -key gcp-key.enc.json
