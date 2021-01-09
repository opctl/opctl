---
title: Docker
sidebar_label: Docker
---

Official opctl images are [published to docker hub](https://hub.docker.com/r/opctl/opctl/)

# Image Variants

## DinD variant

The `opctl:0.1.46-dind` variant leverates Docker in Docker (DinD) and requires a `--privileged` flag.

### Example run github.com/opspec-pkgs/uuid.v4.generate#1.1.0 with DinD

```shell
docker run \
    --privileged \
    opctl/opctl:0.1.46-dind \
    opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
```

## DooD variant

The `opctl:0.1.46-dood` variant leverages Docker out of Docker (DooD) and requires:
- `-v /var/run/docker.sock:/var/run/docker.sock`
  to mount the socket of the external docker daemon.
- `-v opctl_data_dir:/root/opctl`
  to mount an external directory as opctl's data dir.

### Example run github.com/opspec-pkgs/uuid.v4.generate#1.1.0 with DooD
```shell
docker run \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v ~/opctl:/root/opctl \
    opctl/opctl:0.1.46-dood \
    opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
```
