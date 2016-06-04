A manifest file for an op collection.

# name

[MUST](index.md#requirements) be named `op-collection.yml`.

# schema

format: [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

attributes:

| name         | type                                    | constraints                                                                                                 |
|:-------------|:----------------------------------------|:------------------------------------------------------------------------------------------------------------|
| spec-version | [string](http://yaml.org/type/str.html) | Required, [MUST](index.md#requirements) match a [SPEC RELEASE VERSION](index.md#spec_release_version) |

example:
```YAML
spec-version: 0.0.0
```
