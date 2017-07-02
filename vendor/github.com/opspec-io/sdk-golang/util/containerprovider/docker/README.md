A docker implementation of the
opctl/util/containerprovider/ContainerProvider interface


# Dev guide

## docker engine interaction

Do ensure all interactions w/ docker engine leverage API client rather
than shelling out

Why? Remaining agnostic about presence of docker bin on host allows
greater flexibility/portability
