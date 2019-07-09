#!/usr/bin/env sh

git config --global user.email "$GIT_USER@users.noreply.github.com"
git config --global user.name "$GIT_USER"
echo "machine github.com login $GIT_USER password $GITHUB_PASSWORD" > ~/.netrc

CURRENT_BRANCH=master yarn run publish-gh-pages
exit 0;