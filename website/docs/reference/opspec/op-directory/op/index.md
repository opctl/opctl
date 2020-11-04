---
sidebar_label: Overview
title: Op [object]
---
An object which defines an operations inputs, outputs, call graph... etc.

## Properties
- must have
    - [name](#name)
- may have
    - [description](#description)
    - [inputs](#inputs)
    - [opspec](#opspec)
    - [outputs](#outputs)
    - [run](#run)
    - [version](#version)

### name
A string defining a human friendly identifier for the operation.

> It's considered good practice to make `name` unique by using domain
> &/or path based namespacing.

Ops MAY be network resolvable; therefore `name` MUST be a valid
[uri-reference](https://tools.ietf.org/html/rfc3986#section-4.1)

example:
```yaml
name: `github.com/opspec-pkgs/jwt.encode`
```

### description
A [markdown [string]](markdown.md) defining a human friendly description of the op (since v0.1.6).

### inputs
An object defining input parameters of the operation.

For each property:
- key is an [identifier [string]](identifier.md) defining the name of the input.
- value is a [parameter [object]](parameter/index.md) defining the output. 

### outputs
An object defining output parameters of the operation.

For each property:
- key is an [identifier [string]](identifier.md) defining the name of the output.
- value is a [parameter [object]](parameter/index.md) defining the output.

### opspec
A [semver v2.0.0 [string]](https://semver.org/spec/v2.0.0.html) which defines the version of opspec used to define the operation.

### run
A [call [object]](call/index.md) defining the ops call graph; i.e. what gets run by the operation. 

### version
A [semver v2.0.0 [string]](https://semver.org/spec/v2.0.0.html) which defines the version of the operation. 

> If the op is published remotely, this MUST correspond to a [git] tag on the containing repo.