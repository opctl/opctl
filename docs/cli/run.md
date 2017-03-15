## `run`

Run an op.

> if a node isn't running, one will be automatically created

### inputs

#### sources

Input sources are checked according to the following precedence:

- arg provided via `-a` option
- env var
- default
- prompt

#### prompt

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
