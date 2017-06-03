## Packages

A package encapsulates the definition of an op & any requisite
artifacts.

### Resolution

unless explicitly overridden, resolution precedence MUST be:

- local path relative to `.opspec` directory
- local path relative to current
  [working directory](https://en.wikipedia.org/wiki/Working_directory)
- remote git repo where
  [URI fragment](https://tools.ietf.org/html/rfc3986#section-3.5) MUST
  be a valid
  [v2.0.0 semantic version](http://semver.org/spec/v2.0.0.html) and MUST
  resolve to a
  [tag](https://github.com/git/git/blob/v2.13.0/Documentation/git-tag.txt)


example:

```
|--some-dir
  |--.opspec
    |--op1 # could reference remote some.host/some/git/repo#0.0.0
    |--op2 # could reference local op1
    # ~ snip other pkgs
```

### Format

```
pkg-name
  |-- op.yml (manifest *required)
  ... (pkg specific files/dirs *optional)
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

[include](pkg-manifest.schema.json)
