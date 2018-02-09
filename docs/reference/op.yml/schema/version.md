Semantic version of the pkg.

MUST be a valid
[v2 semantic version](https://semver.org/spec/v2.0.0.html)

## Examples

### Major revision

Major revisions typically signify breaking changes.

```yaml
name: majorRev
version: 1.0.0
run:
  container:
    image: { ref: alpine }
```

### Minor revision

Minor revisions typically signify non breaking changes which introduce
new functionality.

```yaml
name: minorRev
version: 1.1.0
run:
  container:
    image: { ref: alpine }
```

### Patch revision

Patch revisions typically signify non breaking changes which don't
introduce new functionality

```yaml
name: PatchRev
version: 1.1.1
run:
  container:
    image: { ref: alpine }
```

