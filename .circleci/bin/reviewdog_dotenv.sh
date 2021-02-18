#!/usr/bin/env bash

set -ex

./bin/dotenv-linter \
| ./bin/reviewdog -f=dotenv-linter -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
