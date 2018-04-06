## `op validate OP_REF` (since v0.1.15)

Validates an op according to:

- existence of `op.yml`
- validity of `op.yml` (per
  [schema](https://opspec.io/0.1.6/op.yml.schema.json))

## Arguments

### `OP_REF`

Op reference (either `relative/path`, `/absolute/path`,
`host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path`
(since v0.1.24))

## Examples

```shell
opctl op validate myop
```

## Notes

### op source username/password prompt

If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> and exit with a non zero exit code.

