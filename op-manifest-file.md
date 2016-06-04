A manifest file for an op.

# Name

[MUST](index.md#mustmay) be named `op.yml`

# Schema

format: [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

## Op
attributes:

| name | type | constraints |
| :---- | :---- | :--- |
| name | [string](http://yaml.org/type/str.html) | Required |
| description | [string](http://yaml.org/type/str.html) | Required |
| params | [map](http://yaml.org/type/map.html) of [Param Map Entries](#param-map-entry) | - |
| subOps | [sequence](http://yaml.org/type/seq.html) of [Sub Op](#sub-op)s | - |

## Param Map Entry
key  

| type                                    | constraints                                                  |
|:----------------------------------------|:-------------------------------------------------------------|
| [string](http://yaml.org/type/str.html) | Required [MUST](index.md#mustmay) be a [FILE_SAFE_NAME](index.md#file_safe_name) |

value

| type                                    | constraints                                                  |
|:----------------------------------------|:-------------------------------------------------------------|
| [Param Map Value](#param-map-entry-value) | Required |

## Param Map Entry Value
attributes:

| name | type                                    | constraints                                                  |
|:-----|:----------------------------------------|:-------------------------------------------------------------|
| isSecret  | [bool](http://yaml.org/type/bool.html) | Optional (defaults to false) |
| description  | [string](http://yaml.org/type/str.html) | Required |

## Sub Op
attributes:

| name | type                                    | constraints                                                  |
|:-----|:----------------------------------------|:-------------------------------------------------------------|
| url  | [string](http://yaml.org/type/str.html) | Required, [MUST](index.md#mustmay) be a valid op url |

# Example
```YAML
name: build-and-publish
description: builds the docker image then pushes it to docker hub
params:
  DOCKER_PASSWORD:
    description: Password for docker registry
    isSecure: true
  DOCKER_USERNAME:
    description: Username for docker registry
    isSecure: true
subOps: 
- url: build
- url: push-to-docker-hub
```
