---
title: Image
---

An image object defines the image to start a container with.

## Properties

### `ref`

_required_

A string referencing a local or remote image. Two types of image refs are supported:

#### OCI images

A [variable reference](../../variable-reference.md) evaluating to a [directory](../../../types/dir) conforming to a [v1.0.1 OCI `image-layout`](https://github.com/opencontainers/image-spec/blob/v1.0.1/image-layout.md).

```yaml
ref: $(myOCIImageLayoutDir)
```

#### Docker images 

A [string](../../../types/string.md#initialization) evaluating to a docker image name. Opctl will use the `latest` tag by default. These usually take the format `[host][repository]image[tag]` where the default host is `docker.io` and default tag is `latest`, but you can also reference local image tags (e.g. if running an image created by a previous call in the op).

```yaml
ref: 'ubuntu:19.10'
``` 

```yaml
ref: 'docker.io/ubuntu:19.10'
```

### `pullCreds`

A [pull credentials](../../pull-creds.md) object containing credentials used to pull the image from a private source.
