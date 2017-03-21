### image

An official opctl image is
[maintained on docker hub](https://hub.docker.com/r/opctl/opctl/) (as
of v0.1.15).

The image features a ready to use opctl node.

> the container provider leveraged by opctl in this case is an embedded
> docker daemon which requires the `--privileged` flag to run

### examples

```shell
# attach to an sh terminal in the container
docker run --privileged -it opctl/opctl:beta /bin/sh

# install git & clone the opspec repo
apk add -U git
git clone https://github.com/opspec-io/spec.git

# cd to examples dir
cd spec/examples/nodejs

# run an op!
opctl run install-deps
```

