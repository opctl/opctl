Absolute/relative path used as the default value for the parameter.

# Absolute path

MUST be interpreted from the root of the package.

# Relative path

If the operation is the root of the run graph, value MUST be interpreted
as relative to the current directory (from perspective of tooling).

OTHERWISE value MUST be ignored.

# Examples

## Package root

```yaml
dir:
  default: /
```

## Current directory

> only observed if operation is called as root of run graph

```yaml
dir:
  default: .
```

