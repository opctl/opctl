---
title: Image [object]
---

An object which defines the image of a container call.

## Properties
- must have
  - [ref](#ref)
- may have
  - [platform](#platform)
  - [pullCreds](#pullcreds)

### ref
A string referencing a local or remote image.

Must be one of:
- a [variable-reference [string]](../../variable-reference.md) evaluating to a [v1.0.1 OCI (Open Container Initiative) `image-layout`](https://github.com/opencontainers/image-spec/blob/v1.0.1/image-layout.md).
- a [string initializer](../../../../types/string.md#initialization) evaluating to a docker image name i.e. `[host][repository]image[tag]` where by default host is `docker.io` and tag is `latest`

### Example ref ([docker.io/ubuntu:19.10](https://hub.docker.com/_/ubuntu))
`ref: 'ubuntu:19.10'` or `ref: 'docker.io/ubuntu:19.10'`

### Example ref (variable)
`ref: $(myOCIImageLayoutDir)`

### platform
An [oci-image-platform [object]](../oci-image-platform.md) constraining the image which will be pulled from the source.

### pullCreds
A [pull-creds [object]](../pull-creds.md) defining creds used to pull the image from a private source.