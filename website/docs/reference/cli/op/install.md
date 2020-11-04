---
sidebar_label: install
title: opctl op install
---

```sh
opctl op install [OPTIONS] OP_REF
```

Install an op.

## Arguments

### `OP_REF`
Op reference (`host/path/repo#tag`, or `host/path/repo#tag/path`)

## Options

### `--path` *default: `.opspec/OP_REF`*
Path to install the op at

### `-u` or `--username`
Username used to auth w/ the op source

### `-p` or `--password`
Password used to auth w/ the op source

## Global Options
see [global options](../global-options.md)

## Examples
```sh
opctl op install -u someUser -p somePass host/path/repo#tag
```

## Notes

### op source username/password prompt
If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> and exit with a non zero exit code.