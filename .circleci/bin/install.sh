#!/usr/bin/env bash

set -ex

REVIEWDOG_VERSION=v0.11.0
GOLANGCILINT_VERSION=v1.37.0
HADOLINT_VERSION=v1.22.1
DOTENV_LINTER_VERSION=v3.0.0

# Install reviewdog
wget -O - -q https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh \
  | sh -s -- -b ./bin $REVIEWDOG_VERSION

# Install golangci-lint
wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $GOLANGCILINT_VERSION

# Install hadolint
wget -q https://github.com/hadolint/hadolint/releases/download/$HADOLINT_VERSION/hadolint-Linux-x86_64 \
  -O ./bin/hadolint && chmod +x ./bin/hadolint

# Install dotenv-linter
wget -q -O - https://git.io/JLbXn | sh -s -- -b bin $DOTENV_LINTER_VERSION

# For working with binaries without a relative path
cp -r bin/. /usr/local/go/bin