---
title: OCI Image Platform [object]
---

An object defining the platform for an OCI image.

## Properties
- must have
  - [arch](#arch)

### arch
A [string initializer](../../../types/string.md#initialization) specifying a [v1.0.1 OCI (Open Container Initiative) `image-index`](https://github.com/opencontainers/image-spec/blob/v1.0.1/image-index.md) platform architecture.

### Example arch
`arch: amd64`

### Example arch (variable)
`arch: $(myArch)`