#!/bin/sh

set -e

apk add --no-cache github-cli

gh release edit "v${VERSION}" --draft=false --prerelease=false
