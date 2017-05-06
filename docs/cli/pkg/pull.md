## `pkg pull [OPTIONS] PKG_REF` (since v0.1.19)

Pulls a package from a remote git repo to the local pkg cache

## Arguments

### `PKG_REF`
Package reference (`host/path/repo#tag`)

## Options

### `-u` or `--username`
Username used for pull auth

### `-p` or `--password`
Password used for pull auth

## Examples

```shell
$ opctl pkg pull -u someUser -p somePass host/path/repo#tag
```

## Notes

### username/password prompt

If auth failure occurs the cli will (re)prompt for username & password & re-attempt the pull.

> in non-interactive terminals, the cli will note that it can't prompt due to being in a
> non-interactive terminal and exit with a non zero exit code.
