# Name

[MUST](index.md#mustmay) be named `op.yml`

# Format

[MUST](index.md#mustmay) be
[v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

# Schema

attributes:

| name        | type                                                     | constraints                                                                                |
|:------------|:---------------------------------------------------------|:-------------------------------------------------------------------------------------------|
| name        | [string](http://yaml.org/type/str.html)                  | Required [MUST](index.md#mustmay) be a [FILE_SAFE_NAME](index.md#file_safe_name)           |
| description | [string](http://yaml.org/type/str.html)                  | Required                                                                                   |
| version     | [string](http://yaml.org/type/str.html)                  | Optional [MUST](index.md#mustmay) be a [v2.0.0 SemVer](http://semver.org/spec/v2.0.0.html) |
| run         | [Atomic](#atomic-run) or [Composite](#composite-run) Run | Required                                                                                   |

## Atomic Run

attributes

| name      | type                                                          | constraints                                                                                         |
|:----------|:--------------------------------------------------------------|:----------------------------------------------------------------------------------------------------|
| container | [Container](#container)                                       | Required [MUST](index.md#mustmay) be a docker image reference (`[[registry/]namespace/]repo[:tag]`) |
| params    | [sequence](http://yaml.org/type/seq.html) of [Param](#param)s | Optional                                                                                            |

## Composite Run

attributes

| name   | type                                                            | constraints |
|:-------|:----------------------------------------------------------------|:------------|
| subOps | [sequence](http://yaml.org/type/seq.html) of [Sub Op](#sub-op)s | Required    |


## Param

attributes  

| name        | type                                    | constraints                                                                      |
|:------------|:----------------------------------------|:---------------------------------------------------------------------------------|
| description | [string](http://yaml.org/type/str.html) | Required                                                                         |
| isSecret    | [bool](http://yaml.org/type/bool.html)  | Optional (defaults to false)                                                     |
| name        | [string](http://yaml.org/type/str.html) | Required [MUST](index.md#mustmay) be a [FILE_SAFE_NAME](index.md#file_safe_name) |

## Container

attributes

| name       | type                                                              | constraints                                                                                         |
|:-----------|:------------------------------------------------------------------|:----------------------------------------------------------------------------------------------------|
| cmd        | [string](http://yaml.org/type/str.html)                           | Optional (defaults to the image’s CMD)                                                              |
| entrypoint | [string](http://yaml.org/type/str.html)                           | Optional (defaults to the image’s ENTRYPOINT)                                                       |
| image      | [string](http://yaml.org/type/str.html)                           | Required [MUST](index.md#mustmay) be a docker image reference (`[[registry/]namespace/]repo[:tag]`) |
| env        | [sequence](http://yaml.org/type/seq.html) of [Env Var](#env-var)s | Optional                                                                                            |

## Env Var

attributes  

| name         | type                                    | constraints                               |
|:-------------|:----------------------------------------|:------------------------------------------|
| name         | [string](http://yaml.org/type/str.html) | Required                                  |
| defaultValue | [string](http://yaml.org/type/str.html) | Optional                                  |
| valueFrom    | [string](http://yaml.org/type/str.html) | Optional (defaults to param of same name) |

## Sub Op

attributes:

| name       | type                                    | constraints                                                |
|:-----------|:----------------------------------------|:-----------------------------------------------------------|
| url        | [string](http://yaml.org/type/str.html) | Required, [MUST](index.md#mustmay) be a valid op reference |
| isParallel | [bool](http://yaml.org/type/bool.html)  | Optional (defaults to false)                               |


# Examples

```YAML
name: push-to-docker-reg
description: pushes the docker image
run:
  params:
  - name: DOCKER_PASSWORD
    description: Password for docker registry
    isSecret: true
  - name: DOCKER_USERNAME
    description: Username for docker registry
    isSecret: true
  container:
    image: docker
    entrypoint: |
      sh -c '
        docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD && \
        docker push my-image
      '
    env:
    - name: DOCKER_PASSWORD
    - name: DOCKER_USERNAME
```

```YAML
name: build-and-publish
description: builds the docker image then pushes it to docker hub
run:
  subOps
    - url: build
    - url: push-to-docker-hub
```

