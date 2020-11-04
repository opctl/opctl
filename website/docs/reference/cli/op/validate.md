---
sidebar_label: validate
title: opctl op validate
---

```sh
opctl op validate [OPTIONS] OP_REF
```

Validate an op according to:

- existence of `op.yml`
- validity of `op.yml` (per
  [schema](https://opctl.io/0.1.6/op.yml.schema.json))

## Arguments

### `OP_REF`
Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`).

## Examples
```sh
opctl op validate myop
```

## Global Options
see [global options](../global-options.md)

## Notes

#### op source username/password prompt
If auth w/ the op source fails the cli will (re)prompt for username & password.

> in non-interactive terminals, the cli will note that it can't prompt and exit with a non zero exit code.
