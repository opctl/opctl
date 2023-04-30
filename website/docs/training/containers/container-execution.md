---
title: Container execution
---

## TLDR;
Opctl supports using [container](../../reference/opspec/op-directory/op/call/container/index.md) statements to make your op run [OCI](https://opencontainers.org/) image based containers.

> Note: a common place to obtain [OCI](https://opencontainers.org/) images is [Docker Hub](https://hub.docker.com/).

## Example
1. Start this op: 
    ```yaml
    name: containerExecution
    run:
      container:
        cmd: [echo, 'hello!']
        image: { ref: alpine }
    ```
2. Observe the container is started, `hello!` is logged, and the container exits.
