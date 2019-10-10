#!/usr/bin/env sh

git config --global user.email "$GIT_USER@users.noreply.github.com"
git config --global user.name "$GIT_USER"

yarn run publish-gh-pages