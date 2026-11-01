#!/bin/sh
# Tags the Go op SDK at the release commit. Rerunning a release keeps a tag
# that already points there, and fails if the tag points anywhere else.
set -eu

tag="sdks/go/opspec/sdk/v$version"
api=https://api.github.com/repos/opctl/opctl/git
auth="Authorization: Bearer $githubAccessToken"
ref="${TMPDIR:-/tmp}/ref.json"

status=$(curl -sS -o "$ref" -w '%{http_code}' -H "$auth" "$api/ref/tags/$tag")
case "$status" in
  200)
    if grep -Eq "\"sha\": *\"$commit\"" "$ref"; then
      echo "$tag already points to $commit"
      exit 0
    fi
    echo "$tag already exists at another commit" >&2
    exit 1
    ;;
  404)
    curl -fsS -o /dev/null -X POST -H "$auth" -H 'Content-Type: application/json' \
      -d "{\"ref\":\"refs/tags/$tag\",\"sha\":\"$commit\"}" "$api/refs"
    echo "created $tag at $commit"
    ;;
  *)
    echo "unable to read $tag: HTTP $status" >&2
    exit 1
    ;;
esac
