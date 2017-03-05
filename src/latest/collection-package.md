## Collection Package

Encapsulates a *collection

### Collection Package Dir Structure

```
pkg-name (or .opspec)
  |-- collection.yml
  ... (op packages)
```

Constraints:

- MUST contain a
  [collection package manifest](#collection-package-manifest) at their
  root.
- Name MUST match their manifest

### Collection Package Manifest

Constraints:

- MUST be named `collection.yml`
- MUST be [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against
  [schema/manifest.json#definitions/collectionManifest](schema/manifest.json#definitions/collectionManifest)
