#!/usr/bin/env bash

set -ex

./bin/golangci-lint run --out-format=line-number \
| ./bin/reviewdog -f=golangci-lint -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
