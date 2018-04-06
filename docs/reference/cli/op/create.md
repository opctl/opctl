## `op create [OPTIONS] NAME` (since v0.1.15)

Creates an op

## Arguments

### `NAME`
Name of the op

## Options

### `-d` or `--description`
Description of the op

### `--path` *default: `.opspec`*
Path to create the op at

## Examples

```shell
opctl op create -d "my awesome op description" --path some/path my-awesome-op-name
```
