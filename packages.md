## Packages

Encapsulate the definition of an op.

### Dir Structure

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
- MUST validate against [schema/packageManifest.json](schema/packageManifest.json)
