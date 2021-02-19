#!/usr/bin/env bash

set -ex

dotenv-linter \
| reviewdog -f=dotenv-linter -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
