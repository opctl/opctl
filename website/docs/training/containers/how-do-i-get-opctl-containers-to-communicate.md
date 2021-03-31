---
title: How do I get opctl containers to communicate?
---

## TLDR;
Opctl attaches all containers to a virtual overlay network.  

Adding a [name](../../reference/opspec/op.yml/call/container/index#name) attribute to container(s) adds a corresponding network wide DNS A record which resolves to the assigned ip(s) of the container(s).

Whether containers are defined in the same op or not makes no difference, they can still reach each other.

If multiple containers have the same [name](../../reference/opspec/op.yml/call/container/index#name) requests will be load balanced across them.

## Example
1. Run this op:
    ```yaml
    name: ping
    run:
      parallel:
        - container:
            image: { ref: alpine }
            name: container1
            cmd: [sleep, 1000000]
        - container:
            image: { ref: alpine }
            # ping container1 by its name
            cmd: [ping, container1]
    ```

1. Observe the second container succeeds in `ping`ing `container1`. 
