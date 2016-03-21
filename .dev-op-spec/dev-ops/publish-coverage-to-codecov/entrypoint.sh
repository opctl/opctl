#!/bin/sh

cat engine.coverprofile > coverage.txt && \
cat core/core.coverprofile >> coverage.txt && \
cat core/models.coverprofile >> coverage.txt && \
cat core/adapters/dockercompose/dockercompose.coverprofile >> coverage.txt && \
cat core/adapters/git/git.coverprofile >> coverage.txt && \
cat core/adapters/osfilesys/osfilesys.coverprofile >> coverage.txt && \
cat rest/rest.coverprofile >> coverage.txt && \
curl -s https://codecov.io/bash | bash -s
