# prerequisites

The current implementation of opctl relies on API access to a docker
daemon to run containers

Docker has two offerings supporting osx, docker-machine and docker4mac.

Docker behaves quite differently depending on which offering you choose:

### docker-machine

| pros                                                                                                                          | cons                             |
|:------------------------------------------------------------------------------------------------------------------------------|:---------------------------------|
| bind mounts order of magnitude faster than docker4Mac via [docker-machine-nfs](https://github.com/adlogix/docker-machine-nfs) | manual creation of docker VM     |
| -                                                                                                                             | manual export of docker env vars |
| -                                                                                                                             | manual start of docker VM        |

### docker4Mac

| pros                           | cons                                                             |
|:-------------------------------|:-----------------------------------------------------------------|
| auto creation of docker VM     | bind mounts order of magnitude slower than docker-machine w/ nfs |
| auto export of docker env vars | -                                                                |
| auto start of docker VM        | -                                                                |

In the end, docker4Mac wins in convenience but falls on it's face when
it comes to bind mount performance.

We therefore strongly recommend using a
[virtualbox docker-machine](https://docs.docker.com/machine/drivers/virtualbox/)
and [docker-machine-nfs](https://github.com/adlogix/docker-machine-nfs)

# installation

curl pipe the opctl binary

```bash
curl -L https://github.com/opctl/opctl/releases/download/0.1.22/opctl0.1.22.darwin.tgz | tar -xzv -C /usr/local/bin
```

