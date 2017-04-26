## `pkg pull [OPTIONS] PKG_REF` (since v0.1.19)

Pulls a package from a remote git repo to the local pkg cache

## Arguments

### `PKG_REF`
Package reference (`host/path/repo#tag`)

## Options

### `-u` or `--username`
Username used for auth w/ source

### `-p` or `--password`
Password used for auth w/ source

## Examples

```shell
$ opctl pkg pull -u someUser -p somePass host/path/repo#tag
```
