## Packages

A package defines an orchestration of (a) containerized process(es).

> Ops are instances of a package.

Packages MUST follow the [package format](#format)

### Format

```
pkg-name
  |-- op.yml
  ... (pkg specific files/dirs)
```

Constraints:

- MUST contain a [manifest](#manifest) at their root.
- Name MUST match that defined by their manifest.

### Manifest

Constraints:

- MUST be named `op.yml`
- MUST be valid [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against the [manifest schema](#manifest-schema)

#### Manifest Schema

[include](package-manifest.schema.json)
