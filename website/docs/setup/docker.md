---
title: Docker
sidebar_label: Docker
---

Official opctl images are [published to docker hub](https://hub.docker.com/r/opctl/opctl/)

# Image Variants

## dind variant

The `opctl:0.1.29-dind` variant leverates Docker In Docker (dind) and requires a `--privileged` flag.

### Examples

```shell
docker run --privileged opctl/opctl:0.1.29-dind opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
```

## dood variant (experimental; paths might not be found)

The `opctl:0.1.29-dood` variant leverages Docker out of Docker (dood) and requires a `-v /var/run/docker.sock:/var/run/docker.sock` arg. 

### Examples

```shell
docker run -v /var/run/docker.sock:/var/run/docker.sock opctl/opctl:0.1.29-dood opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
```