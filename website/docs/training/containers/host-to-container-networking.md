---
title: Host to container networking
---

## TLDR;
Opctl supports using [ports](../../reference/opspec/op-directory/op/call/container/index.md#ports) statements to bind opctl container ports to opctl node host ports.

## Example
1. Start this op: 
    ```yaml
    name: curl
    run:
      container:
        image: { ref: nginx:alpine }
        ports:
          # bind container port 80 to host port 8080
          80: 8080
    ```
2. On the opctl host, open a web browser to [localhost:8080](localhost:8080).
3. Observe the nginx containers default page is returned. 
