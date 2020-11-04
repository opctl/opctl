---
sidebar_label: self-update
title: opctl self-update
---

```sh
opctl self-update [OPTIONS]
```

Update opctl.

> if a node is running, it will be automatically killed

## Options

### `-c` or `--channel` *default: `stable`*
The release channel to update from

- `stable`
- `beta` (smoke tested alpha channel)
- `alpha` (all bets are off)

## Global Options
see [global options](global-options.md)

## Examples
get latest stable release
```sh
opctl self-update
# output: Updated to new version: 0.1.24!
```

play around w/ latest beta release
 ```sh
opctl self-update -c beta
# output: Updated to new version: 0.1.24-beta.1!
```

play times over; switch back to latest stable release
```sh
opctl self-update
# output: Updated to new version: 0.1.24!
```