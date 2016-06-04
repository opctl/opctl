A manifest file for an op.

# Name

[MUST](index.md#mustmay) be named `op.yml`

# Schema

format: [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

## Op Schema
attributes:

| name | type | constraints |
| :---- | :---- | :--- |
| name | [string](http://yaml.org/type/str.html) | Required |
| description | [string](http://yaml.org/type/str.html) | Required |
| subOps | [sequence](http://yaml.org/type/seq.html) of [Sub Op](#sub-op-schema)s | - |

example:
```YAML
name: build-and-publish
description: builds the docker image then pushes it to docker hub
subOps: 
- url: build
- url: push-to-docker-hub
```

## Sub Op Schema
attributes:

| name | type                                    | constraints                                                  |
|:-----|:----------------------------------------|:-------------------------------------------------------------|
| url  | [string](http://yaml.org/type/str.html) | Required, [MUST](index.md#mustmay) be a valid op url |

example:
```YAML
url: push-to-docker-hub
```


