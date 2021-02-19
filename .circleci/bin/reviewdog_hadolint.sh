#!/usr/bin/env bash

set -ex

git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint \
| reviewdog -efm="%f:%l %m" -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
