## `pkg create [OPTIONS] NAME` (since v0.1.15)

Creates a package

## Arguments

### `NAME`
Name of the package

## Options

### `-d` or `--description`
Description of the package

### `--path` *default: `.opspec`*
Path to create the package at

## Examples

```shell
opctl pkg create -d "my awesome pkg description" -p some/path my-awesome-pkg-name
```
