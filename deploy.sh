#!/usr/bin/env bash

echo "I'm deploying!"
echo "..."

./janus deploy \
        # TODO: replace me with the real stuff; bucket, object, file, key
        -bucket isaac-tests \
        -object janus/$(./janus version -format '%M.%m.x')/janus-$TRAVIS_OS_NAME-$(./janus version -format 'v%M.%m.%P+%C-%S').zip \
        -file janus-$TRAVIS_OS_NAME-$(./janus version -format 'v%M.%m.%P+%C-%S').zip \
        -key isaac-test-key.enc.json
