#!/usr/bin/env bash

echo "Deploying..."

# TODO: replace me with the real stuff; bucket, object, file, key
janus deploy -to "isaac-tests/janus/$(janus version -format '%M.%m.x'") -files "*.zip" -key isaac-test-key.enc.json
