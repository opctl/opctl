#!/bin/sh

cat engine.coverprofile > coverage.txt && \
cat core/core.coverprofile >> coverage.txt && \
cat core/models/models.coverprofile >> coverage.txt && \
cat core/adapters/containerengine/dockercompose/dockercompose.coverprofile >> coverage.txt && \
cat core/adapters/templatesrc/git/git.coverprofile >> coverage.txt && \
cat core/adapters/filesys/os/os.coverprofile >> coverage.txt && \
cat rest/rest.coverprofile >> coverage.txt && \
curl -s https://codecov.io/bash | bash -s
