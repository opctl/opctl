## `self-update`

Updates the current version of opctl.

> if a node is running, it will be automatically killed

### channels

The `self-update` command takes an optional `-c` (or `--channel`)
argument which allows updating from any available release channel:

- `stable`(default)
- `beta` (smoke tested alpha channel)
- `alpha` (all bets are off)

### examples

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
