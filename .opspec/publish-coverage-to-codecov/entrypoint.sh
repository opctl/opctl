#!/bin/sh

# generate aggregate coverage profile
cat engine.coverprofile > coverage.txt && \
cat core/core.coverprofile >> coverage.txt && \
cat core/models/models.coverprofile >> coverage.txt && \
cat core/adapters/containerengine/dockercompose/dockercompose.coverprofile >> coverage.txt && \
cat core/adapters/filesys/os/os.coverprofile >> coverage.txt && \
cat tcp/tcp.coverprofile >> coverage.txt

# strip fakes from coverage profile
sed -i '/fake/d' coverage.txt

bash <(curl -s https://codecov.io/bash)
