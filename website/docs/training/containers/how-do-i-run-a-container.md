---
title: How do I run a container?
---

## TLDR;
Opctl supports running an [OCI](https://opencontainers.org/) image based container by defining a [container call](../../reference/opspec/op.yml/call/container/index).

> Note: a common place to obtain [OCI](https://opencontainers.org/) images is [Docker Hub](https://hub.docker.com/).

## Example
1. Start this op: 
    ```yaml
    name: runAContainer
    run:
      container:
        cmd: [sh, -ce, 'echo hello!']
        image: { ref: alpine }
    ```
1. Observe the container is run and `hello!` logged.
