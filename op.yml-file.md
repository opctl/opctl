# Name

[MUST](index.md#mustmay) be named `op.yml`

# Format

[MUST](index.md#mustmay) be
[v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

# Schema

attributes:

| name        | description           | schema                                                                | required | default                |
|:------------|:----------------------|:----------------------------------------------------------------------|:---------|:-----------------------|
| name        | Name of the op        | [FILE_SAFE_NAME](index.md#file_safe_name)                             | true     |                        |
| description | Description of the op | [string](http://yaml.org/type/str.html)                               | true     |                        |
| inputs      | Inputs to the op      | [sequence](http://yaml.org/type/seq.html) of [Parameter](#parameter)s | false    |                        |
| run         | Run instruction       | [Run Instruction](#run-instruction)                                   | false    | run docker-compose.yml |
| version     | Version of the op     | [v2.0.0 SemVer](http://semver.org/spec/v2.0.0.html)                   | false    |                        |

## Run Instruction

Types:

- [Sub Ops](#sub-ops-run-instruction)

### Sub Op Run Instruction

attributes:

| name       | description                                                                                | schema                                  | required | default |
|:-----------|:-------------------------------------------------------------------------------------------|:----------------------------------------|:---------|:--------|
| isParallel | If set this sub op and the next (if any) will be run simultaneously | [bool](http://yaml.org/type/bool.html)  | false    | false   |
| url        | Url of an op                                                                               | [string](http://yaml.org/type/str.html) | true     |         |


### Sub Ops Run Instruction

attributes

| name   | description                    | schema                                                                                          | required | default |
|:-------|:-------------------------------|:------------------------------------------------------------------------------------------------|:---------|:--------|
| subOps | Sub op run instructions to run | [sequence](http://yaml.org/type/seq.html) of [Sub Op Run Instruction](#sub-op-run-instruction)s | true     |         |

## Parameter

Types:

- [String](#string-parameter)

attributes

| name     | description                | schema                                    | required | default |
|:---------|:---------------------------|:------------------------------------------|:---------|:--------|
| type     | Type of parameter          | Name of a defined parameter type          | false    | String  |
| isSecret | If the parameter is secret | [bool](http://yaml.org/type/bool.html)    | false    | false   |
| name     | Name of the parameter      | [FILE_SAFE_NAME](index.md#file_safe_name) | true     |         |

### String Parameter

inherits [Parameter](#parameter) attributes

additional attributes:

| name    | description   | schema                                  | required | default |
|:--------|:--------------|:----------------------------------------|:---------|:--------|
| default | default value | [string](http://yaml.org/type/str.html) | false    |         |


# Examples

```YAML
name: npm-install
description: installs npm packages
inputs:
- name: NPM_CONFIG_REGISTRY
  description: Registry to use
  isSecret: true
```

```YAML
name: build-and-publish
description: builds the docker image then pushes it to docker hub
inputs:
- name: DOCKER_PASSWORD
  description: Password for docker registry
  isSecret: true
- name: DOCKER_USERNAME
  description: Username for docker registry
  isSecret: true
run:
  subOps:
  - url: build
  - url: push-to-docker-hub
```

