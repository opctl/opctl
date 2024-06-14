#!/bin/sh

diff=$(git diff origin/main --name-only | grep CHANGELOG.md)

echo $diff

echo "---"
if [ -z $diff ]; then
  echo "WARNING: CHANGELOG.md has not been updated. As a result, no new release will be created when this change is pushed to main."
  echo "Please update CHANGELOG.md with a new version and release notes."
  exit 1
else
  echo "CHANGELOG.md has been updated."
  exit 0
fi
