#!/bin/sh

set -e

apk add --update docker

# only continue with release if tag doesn't exist
if docker pull ghcr.io/opctl/opctl:${version}-dind; then
  echo "Opctl Image for version '${version}' already exists"
  exit 1
else
  echo "Image does not exist, proceeding with release..."
  exit 0
fi
