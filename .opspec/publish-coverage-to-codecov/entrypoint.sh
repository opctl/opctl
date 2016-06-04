#!/bin/sh -e

# generate aggregate coverage profile
cat adapters/host/docker/docker.coverprofile > coverage.txt
cat models/models.coverprofile > coverage.txt
cat sdk.coverprofile >> coverage.txt

# strip fakes from coverage profile
sed -i '/fake/d' coverage.txt

curl -s https://codecov.io/bash | bash -s
