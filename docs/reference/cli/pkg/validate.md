## `pkg validate PKG_REF` (since v0.1.15)

Validates a package according to:

- manifest existence
- manifest validity (per
  [schema](https://opspec.io/0.1.5/pkg-manifest.schema.json))

## Arguments

### `PKG_REF`

Package reference (either `relative/path`, `/absolute/path`,
`host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path`
(since v0.1.24))

## Examples

```shell
opctl pkg validate myop
```

## Notes

### pkg source username/password prompt

If auth w/ the pkg source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> and exit with a non zero exit code.

