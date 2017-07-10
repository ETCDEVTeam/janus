#!/usr/bin/env bash

echo "Deploying..."

janus deploy \
        # TODO: replace me with the real stuff; bucket, object, file, key
        -to isaac-tests/janus/$(janus version -format '%M.%m.x') \
        -files *.zip \
        -key isaac-test-key.enc.json
