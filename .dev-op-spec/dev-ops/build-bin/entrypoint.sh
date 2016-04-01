#!/bin/sh

rm -rf .tmp

go get ./... && \
go build -o .tmp/dev-op-spec-engine