#!/usr/bin/env bash

set -ex

# Exit code always 0
dotenv-linter \
| reviewdog -f=dotenv-linter -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
