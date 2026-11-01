#!/bin/sh
# Tags the Go op SDK at the release commit, reusing an existing tag when its
# SDK contents match. Published Go module versions must remain immutable.
set -eu

tag="sdks/go/opspec/sdk/v$version"
api=https://api.github.com/repos/opctl/opctl/git
auth="Authorization: Bearer $githubAccessToken"
ref="${TMPDIR:-/tmp}/ref.json"

status=$(curl -sS -o "$ref" -w '%{http_code}' -H "$auth" "$api/ref/tags/$tag")
case "$status" in
  200)
    taggedCommit=$(sed -n 's/.*"sha": *"\([0-9a-f]\{40\}\)".*/\1/p' "$ref")
    if [ -z "$taggedCommit" ]; then
      echo "unable to read commit for $tag" >&2
      exit 1
    fi
    if [ "$taggedCommit" = "$commit" ]; then
      echo "$tag already points to $commit"
      exit 0
    fi
    compareDir=$(mktemp -d)
    trap 'rm -rf "$compareDir"' EXIT
    git -C "$compareDir" init --quiet
    if ! git -C "$compareDir" -c "http.extraHeader=$auth" fetch --quiet --depth=1 --no-tags \
      https://github.com/opctl/opctl "$taggedCommit" "$commit"; then
      echo "unable to fetch commits for $tag" >&2
      exit 1
    fi
    if git -C "$compareDir" diff --quiet "$taggedCommit" "$commit" -- sdks/go/opspec/sdk; then
      echo "$tag has identical SDK contents; reusing $taggedCommit"
      exit 0
    else
      result=$?
    fi
    if [ "$result" -ne 1 ]; then
      echo "unable to compare SDK contents for $tag" >&2
      exit 1
    fi
    echo "$tag already exists at $taggedCommit; bump the CHANGELOG version to release again" >&2
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
