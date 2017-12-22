# How do I use opctl from Docker?

[Docker](https://docker.com) supports running containerized processes from "Image"'s.

An official docker opctl image is
[maintained on docker hub](https://hub.docker.com/r/opctl/opctl/) (as
of v0.1.15), which features a ready to use opctl node.

> the container provider leveraged by opctl in this case is an embedded
> docker daemon which requires the `--privileged` flag to run

### Examples

```shell
docker run --privileged opctl/opctl:stable opctl run github.com/opspec-pkgs/uuid.v4.generate#1.0.0
```
