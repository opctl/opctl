#!/bin/sh -e

# generate aggregate coverage profile
cat sdk-golang.coverprofile > coverage.txt

# strip fakes from coverage profile
sed -i '/fake/d' coverage.txt

curl -s https://codecov.io/bash | bash -s
