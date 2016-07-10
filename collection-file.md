# Name

[MUST](index.md#mustmay) be named `collection.yml`

# Schema

format: [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)

attributes:

| name        | type                                    | constraints                                                                      |
|:------------|:----------------------------------------|:---------------------------------------------------------------------------------|
| name        | [string](http://yaml.org/type/str.html) | Required [MUST](index.md#mustmay) be a [FILE_SAFE_NAME](index.md#file_safe_name) |
| description | [string](http://yaml.org/type/str.html) | Required                                                                         |

# Example

```YAML
name: facebook-golang
description: dev ops for golang projects at facebook 
```

