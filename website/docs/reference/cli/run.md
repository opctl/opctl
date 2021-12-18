---
sidebar_label: run
title: opctl run
---

```sh
opctl run [OPTIONS] OP_REF
```

Start and wait for an op to exit.

> if a node isn't running, one will be started automatically

## Arguments

### `OP_REF`
Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)

## Options

### `-a`
Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`

### `--arg-file` *default: `.opspec/args.yml`*
Read in a file of args in yml format

### `--no-progress` *default: `false`*
Disable live call graph for the op

## Global Options
see [global options](global-options.md)

## Examples

### local op ref w/out args
```sh
opctl run myop
```

### remote op ref w/ args
```sh
opctl run -a apiToken="my-token" -a channelName="my-channel" -a msg="hello!" github.com/opspec-pkgs/slack.chat.post-message#0.1.1
```

## Notes

### op source username/password prompt
If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> due to being in a non-interactive terminal and exit with a non zero
> exit code.

### input sources
Input sources are checked according to the following precedence:

- arg provided via `-a` option
- arg file
- env var
- default
- prompt

### input prompts
Inputs which are invalid or missing will result in the cli prompting for
them.

> in non-interactive terminals, the cli will provide details about the
> invalid or missing input, note that it's giving up due to being in a
> non-interactive terminal and exit with a non zero exit code.

example:

```sh

-
  Please provide value for parameter.
  Name: version
  Description: version of app being compiled
-
```

### validation
When inputs don't meet constraints, the cli will (re)prompt for the
input until a satisfactory value is obtained.

### caching
All pulled ops/image layers will be cached

### image updates
Prior to container creation, updates to the referenced image will be
pulled and applied.

If checking for or applying updated image layers fails, graceful
fallback to cached image layers will occur

### container networking
All containers created by opctl will be attached to a single managed
network.

> the network is visible from `docker network ls` as `opctl`.

### container cleanup
Containers will be removed as they exit.
