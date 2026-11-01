#!/bin/sh
# Reads the go.sum hashes of the tagged Go op SDK from the Go module proxy.
set -eu

out="${outDir:-}"
download=$(go mod download -json "github.com/opctl/opctl/sdks/go/opspec/sdk@v$version")
printf '%s\n' "$download" | sed -n 's/^[[:space:]]*"Sum": "\([^"]*\)".*/\1/p' | tr -d '\n' > "$out/sdkSum"
printf '%s\n' "$download" | sed -n 's/^[[:space:]]*"GoModSum": "\([^"]*\)".*/\1/p' | tr -d '\n' > "$out/sdkGoModSum"
test -s "$out/sdkSum"
test -s "$out/sdkGoModSum"
