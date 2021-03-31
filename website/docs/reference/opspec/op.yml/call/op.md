---
title: Op call
---

An op call is an object that runs an external op. The external op can be local or remote and will not share scope with the current op.

## Properties

### `ref`

_required_

A ref is a string identifying the external op.

For local ops, refs can be a relative path referencing an op's directory in the current local filesystem or a [variable referencing](../variable-reference.md) an op's directory.

Remote ops are referenced using git repository tags. Remote refs use the format `git-repo#SEMVER_GIT_TAG[/path]`, where `git-repo` is the hostname of the git server (checked out over `https`), `SEMVER_GIT_TAG` is a valid tag on the git repository, and `[/path]` is an optional sub-directory where the op lives (usually used by repositories that contain collections of ops).

```
ref: "github.com/opspec-pkgs/golang.build.bin#2.0.0"
```

### `pullCreds`

A [pull credentials](../pull-creds.md) object containing credentials used to pull the op from a private remote source.

### `inputs`

Inputs is an object that defines values passed from the current op into the external one, similar to passing named parameters to a function in modern programming languages. Each key is one of the [defined inputs](../index#inputs) of the op being called. Key values define the in-scope value passed to the input and can take a few forms:

|value|meaning|
|--|--|
|[value initializer](../initializer.md)|The defined value is created explicitly|
|[variable reference](../variable-reference.md)|The referenced value passed|
|null|Shorthand for a variable reference of the key/input name (`inputName:` is equivalent to `inputName: $(inputName)`). This is similar to [object initialization shorthand](../../types/object#initialization)|

Inputs defined by the external op that do not have a `default` are required, and your op will fail if not provided.

### `outputs`

Outputs is an object that defines values "returned" by the external op. Each key is one of the [defined outputs](../index#outputs) of the op being called. Key values define the variable the output's value will be stored in.

|value|meaning|
|--|--|
|[variable reference](../variable-reference.md)|The value will be stored in the referenced variable, which will be created if it doesn't exist.|
|null|Shorthand for a variable reference of the name of the key/output name (`outputName:` is equivalent to `outputName: $(outputName)`)|
