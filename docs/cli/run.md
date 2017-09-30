## `run [OPTIONS] PKG_REF`

Start and wait on an op

> if a node isn't running, one will be automatically created

## Arguments

### `PKG_REF`
Package reference (either `relative/path`, `/absolute/path`, or `host/path/repo#tag` (since v0.1.19))

## Options

### `-a`
Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`

### `--arg-file` *default: `.opspec/args.yml`* (since v0.1.19)
Read in a file of args in yml format

## Examples

### local pkg ref w/out args
```shell
opctl run myop
```

### remote pkg ref w/ args (must be installed first)
```shell
opctl run -a apiToken="my-token" -a channelName="my-channel" -a msg="hello!" github.com/opspec-pkgs/slack.chat.post-message#0.1.1
```

## Notes

### pkg source username/password prompt

If auth w/ the pkg source fails the cli will (re)prompt for username & password.

> in non-interactive terminals, the cli will note that it can't prompt due to being in a
> non-interactive terminal and exit with a non zero exit code.

### input sources

Input sources are checked according to the following precedence:

- arg provided via `-a` option
- arg file (since v0.1.19)
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

```shell

-
  Please provide value for parameter.
  Name: version
  Description: version of app being compiled
-
```

#### validation (since v0.1.15)

When inputs don't meet constraints, the cli will (re)prompt for the
input until a satisfactory value is obtained.

### containers

#### networking

All containers created by opctl will be attached to a single managed
network.

> the network is visible from `docker network ls` as `opctl`.

#### cleanup

Containers will be removed as they exit.
