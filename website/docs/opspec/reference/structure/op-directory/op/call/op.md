---
title: Op Call [object]
---

An object defining an op call.

## Properties
- must have
  - [ref](#ref)
- may have
  - [inputs](#inputs)
  - [outputs](#outputs)
  - [pullCreds](#pullcreds)

### ref
A string referencing a local or remote operation.

Must be one of:
- an absolute path referencing an op embedded in the current op.
- a relative path referencing an op existing on the same local filesystem.
- a string in `git-repo#{SEMVER_GIT_TAG}/path` format referencing a network resolvable op.

### Example ref ([opspec-pkgs golang.build.bin 2.0.0](https://github.com/opspec-pkgs/golang.build.bin))
`ref: 'github.com/opspec-pkgs/golang.build.bin#2.0.0'`

### pullCreds
An [pull-creds](pull-creds.md) object defining creds used to pull the op from a private source.

### inputs
An object for which each key is an input of the referenced op and the value is one of:

|value|meaning|
|--|--|
|null|Bind input to variable w/ same name (equivalent to `$(INPUT_NAME)`)|
|[reference](../reference.md)|Bind referenced variable to the named input|
|[initializer](../initializer.md)|Evaluate and bind to the named input|

> This is equivalent to providing named arguments to a function in modern programming languages.

### outputs
An object for which each key is a variable to assign and the value is one of:

|value|meaning|
|--|--|
|null|Bind variable to output w/ same name (equivalent to `$(OUTPUT_NAME)`)|
|[reference](../reference.md)|Bind named output to referenced variable|

> This is equivalent to consuming named return values from a function in modern programming languages.