## Collection Packages

encapsulate collections so they may be consumed and/or distributed.

### Dir Structure

```
pkg-name
  |-- collection.yml
  ... (op packages)
```

Constraints:

- MUST contain a
  [collection package manifest](#manifest) at their
  root.
- Name MUST match their manifest

### Manifest

Constraints:

- MUST be named `collection.yml`
- MUST be [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against
  [schema/manifest.json#definitions/collectionManifest](schema/manifest.json#definitions/collectionManifest)
