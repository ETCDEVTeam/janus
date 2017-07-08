#!/usr/bin/env bash

./janus deploy \
        -bucket isaac-tests \
        -object janus/$(janus version -format '%M.%m.x')/janus-$TRAVIS_OS_NAME-$(janus version -format 'v%M.%m.%P+%C-%S').zip \
        -file janus-$TRAVIS_OS_NAME-$(janus -format "v%M.%m.%P+%C-%S").zip \
        -key isaac-test-key.enc.json
        # -key gcloud-travis.json.enc
