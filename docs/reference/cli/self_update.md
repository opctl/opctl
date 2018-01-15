## `self-update [OPTIONS]`

Updates the current version of opctl.

> if a node is running, it will be automatically killed

## Options

### `-c` or `--channel` *default: `stable`*
The release channel to update from

- `stable`
- `beta` (smoke tested alpha channel)
- `alpha` (all bets are off)

## Examples

get latest stable release
```shell
opctl self-update
# output: Updated to new version: 0.1.23!
```

play around w/ latest beta release
 ```shell
opctl self-update -c beta
# output: Updated to new version: 0.1.23-beta.1!
```

play times over; switch back to latest stable release
```shell
opctl self-update
# output: Updated to new version: 0.1.23!
```
