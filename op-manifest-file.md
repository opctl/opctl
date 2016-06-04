A manifest file for an op.

# criteria

## name

[MUST](index.md#requirements) be named `op.yml`

## schema

format: [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

### op
attributes:

| name | type | constraints |
| :---- | :---- | :--- |
| name | [string](http://yaml.org/type/str.html) | Required |
| description | [string](http://yaml.org/type/str.html) | Required |
| subOps | [sequence](http://yaml.org/type/seq.html) of [Sub Op](#sub-op)s | - |

example:
```YAML
name: build-and-publish
description: builds the docker image then pushes it to docker hub
subOps: 
- url: build
- url: push-to-docker-hub
```

### sub-op
attributes:

| name | type                                    | constraints                                                  |
|:-----|:----------------------------------------|:-------------------------------------------------------------|
| url  | [string](http://yaml.org/type/str.html) | Required, [MUST](index.md#requirements) be a valid op url |

example:
```YAML
url: push-to-docker-hub
```


