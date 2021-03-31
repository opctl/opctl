---
title: How do I communicate with an opctl container?
---

## TLDR;
Adding a [ports](../../reference/opspec/op.yml/call/container/index#ports) attribute to a container binds container ports to the opctl host.

## Example
1. Start this op: 
    ```yaml
    name: curl
    run:
      parallel:
        - container:
            image: { ref: nginx:alpine }
            ports:
              # bind container port 80 to host port 8080
              80: 8080
    ```
1. On the opctl host, open a web browser to [localhost:8080](localhost:8080).
1. Observe the nginx containers default page is returned. 
