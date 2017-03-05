## Op Package

Encapsulates an *op.

### Op Package Dir Structure

```
pkg-name
  |-- op.yml
  ... (pkg specific files/dirs)
```

Constraints:

- MUST contain an [op package manifest](#op-package-manifest) at their
  root.
- Name MUST match their manifest.

### Op Package Manifest

Constraints:

- MUST be named `op.yml`
- MUST be [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against
  [schema/manifest.json#definitions/opManifest](schema/manifest.json#definitions/opManifest)
