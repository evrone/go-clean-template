#!/usr/bin/env bash

set -ex

git ls-files --exclude='Dockerfile*' --ignored | xargs ./bin/hadolint \
  | ./bin/reviewdog -efm="%f:%l %m" -diff="git diff ${DEFAULT_BRANCH}..HEAD" -reporter=github-pr-review
