#!/usr/bin/env bash

set -ex
./bin/reviewdog -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
