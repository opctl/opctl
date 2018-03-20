A docker implementation of the
node/core/containerruntime/ContainerRuntime interface


# Dev guide

## docker engine interaction

Do ensure all interactions w/ docker engine leverage API client rather
than shelling out

Why? Remaining agnostic about presence of docker bin on host allows
greater flexibility/portability
