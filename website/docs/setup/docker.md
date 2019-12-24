---
title: Docker
sidebar_label: Docker
---

An official docker opctl image is [maintained on docker hub](https://hub.docker.com/r/opctl/opctl/) (as of v0.1.15), which features a ready to use opctl node.

> The container runtime in this case will be an embedded docker daemon, which leads to the `--privileged` flag being required

### Examples

```shell
docker run --privileged opctl/opctl:0.1.25 opctl run github.com/opspec-pkgs/uuid.v4.generate#1.0.0
```
