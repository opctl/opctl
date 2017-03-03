## `run`

### input sources

Input sources are checked according to the following precedence:

- arg provided via `-a` option
- env var
- default
- prompt

```shell
$ export input1Name=input1Value; opctl run -a input2Name=input2Value some-op
```

### input prompt

Inputs which are invalid or missing will result in the cli prompting for
them.

> in non-interactive terminals, the cli will provide details about the
> invalid or missing input, note that it's giving up due to being in a
> non-interactive terminal and exit with a non zero exit code.

```shell
-
  Please provide value for parameter.
  Name: version
  Description: version of app being compiled
-


```

### input validation

When inputs don't meet constraints, the cli will (re)prompt for the
input until a satisfactory value is obtained.
