A docker implementation of the
engine/pkg/containerengine/ContainerEngine interface


# Dev guide

## docker engine interaction
Do ensure all interactions w/ docker engine leverage API client rather than shelling out  
Why? Remaining agnostic about presence of docker bin on host allows
 greater flexibility/portability
