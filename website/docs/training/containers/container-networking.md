---
title: Container networking
---

## TLDR;
Opctl attaches all containers to a virtual overlay network and bridges it to the opctl host.

Adding a [name](../../reference/opspec/op-directory/op/call/container/index.md#name) attribute to container(s) makes it resolvable from other containers and the host by that name.

If multiple containers have the same [name](../../reference/opspec/op-directory/op/call/container/index.md#name) the ips of all containers will be resolved as is standard for DNS A records.

## Example
1. Run this op:
    ```yaml
    name: containerNetworking
    run:
    parallel:
        - container:
            image: { ref: nginx }
            # some syntactically valid hostname
            name: nginx.demo.wow
        - container:
            image: { ref: alpine }
            cmd:
            - sh
            - -ce
            - |
                # wait up
                sleep 1; until ping -c1 nginx.demo.wow >/dev/null 2>&1; do :; done

                # ping forever
                ping nginx.demo.wow
    ```

1. Observe the second container succeeds in `ping`ing the `nginx.demo.wow` container.
1. On the opctl host, open a terminal and `ping nginx.demo.wow`; observe successful pings to the `nginx.demo.wow` container.
1. On the opctl host, open a browser and visit `http://nginx.demo.wow`; observe the nginx default webpage is returned.
