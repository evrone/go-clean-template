#!/usr/bin/env bash

set -ex

./integration-test/build/wait-for-it.sh app:8080 \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go test -v ./integration-test/...
