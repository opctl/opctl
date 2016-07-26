# Name

[MUST](index.md#mustmay) be named `collection.yml`

# Format

[MUST](index.md#mustmay) be
[v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

# Schema

attributes:

| name        | type                                    | constraints                                                                                |
|:------------|:----------------------------------------|:-------------------------------------------------------------------------------------------|
| name        | [string](http://yaml.org/type/str.html) | Required [MUST](index.md#mustmay) be a [FILE_SAFE_NAME](index.md#file_safe_name)           |
| description | [string](http://yaml.org/type/str.html) | Required                                                                                   |
| version     | [string](http://yaml.org/type/str.html) | Optional [MUST](index.md#mustmay) be a [v2.0.0 SemVer](http://semver.org/spec/v2.0.0.html) |

# Example

```YAML
name: facebook-golang
description: dev ops for golang projects at facebook 
```

