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

```shell
# get latest stable release
$ opctl self-update
Updated to new version: 0.1.14!

# play around w/ latest beta release
$ opctl self-update -c beta
Updated to new version: 0.1.15-beta.122!

# play times over; switch back to latest stable release
$ opctl self-update
Updated to new version: 0.1.14!
```
