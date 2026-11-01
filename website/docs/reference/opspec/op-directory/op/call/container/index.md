---
sidebar_label: Overview
title: Container Call [object]
---

An object defining a container call.

## Properties
- must have
  - [image](#image)
- may have
  - [cmd](#cmd)
  - [dirs](#dirs)
  - [envVars](#envvars)
  - [files](#files)
  - [name](#name)
  - [sockets](#sockets)
  - [stdout and stderr](#stdout-and-stderr)
  - [outputs](outputs.md)
  - [outputParams](outputs.md#typed-members)
  - [workDir](#workdir)

### image
An [image [object]](image.md) defining the container image run by the call.

### cmd
An [array](../../../../types/array.md) [initializer](../../../../types/array.md#initialization) or [variable-reference [string]](../../variable-reference.md) defining the path (from [workDir](#workdir)) of the binary to call and it's arguments.

> defining cmd overrides any entrypoint and/or cmd defined by the image

### dirs
An object for which each key is an absolute path in the container and each value is one of:

|value|meaning|
|--|--|
|null|Mount dir embedded in op w/ same path (equivalent to `$(./relative/path)`)|
|[dir](../../../../types/dir.md) [variable-reference [string]](../../variable-reference.md)|Mount dir|
|[dir initializer](../../../../types/dir.md#initialization)|Evaluate and mount|

### envVars
An [object initializer](../../../../types/object.md#initialization) or [variable-reference [string]](../../variable-reference.md), whos properties represent the name and value of an environment variable to be set in the container.

> upon evaluation, the key and value of each property will be coerced to a string.

### files
An object for which each key is an absolute path in the container and each value is one of:

|value|meaning|
|--|--|
|null|Mount file embedded in op w/ same path (equivalent to `$(./relative/path)`)|
|[file](../../../../types/file.md) [variable-reference [string]](../../variable-reference.md)|Mount file|
|[file initializer](../../../../types/file.md#initialization)|Evaluate and mount|

### name
A [string initializer](../../../../types/string.md#initialization) defining a name by which the container can be reached by other opctl containers and opctl host nodes.

> if multiple containers are given the same name, network requests will be distributed (load balanced) across them. 

### sockets
An object for which each key is an absolute path in the container and and each value is a [socket](../../../../types/socket.md) [variable-reference [string]](../../variable-reference.md) to mount. 

### workDir
A [string initializer](../../../../types/string.md#initialization) defining absolute path from which [cmd](#cmd) will be executed.

> defining workDir overrides any defined by the image

### stdout and stderr

Bind a process stream to a String using a whole scope reference, for example `stdout: $(version)` or `stderr: $(diagnostic)`. The target must contain only letters, digits, underscores, or hyphens. A new target cannot duplicate another output or mount binding in the same call.

Each bound stream accepts at most 1 MiB of valid UTF-8. Capture strips exactly one final LF or CRLF and preserves other whitespace. The call publishes values only after successful execution and collection. Overflow cancels the container and fails the call. Public log events retain the original bytes; a stream bound to a declared secret output or sensitive target produces no log events and remains available to later calls through private runtime results. Authors must declare derived secrets. Unbound application logs are still public.

Docker supports separate streams. Separate streams are not yet verified on Embedded and are not yet supported on Kubernetes. Calls without stream bindings keep their existing stream and TTY behavior.

```yaml
run:
  serial:
    - container:
        image: {ref: alpine:3}
        cmd: [sh, -ce, "printf 'v1.2.3\\n'"]
        stdout: $(version)
    - container:
        image: {ref: alpine:3}
        envVars: {VERSION: $(version)}
        cmd: [sh, -ce, 'test "$VERSION" = v1.2.3']
```


### outputs

A nonempty map from JSON member names to whole scope write references, such as `outputs: {version: $(buildVersion)}`. See the [outputs file reference](outputs.md).

### outputParams

Optional declarations for required outputs-file members. Validate exact JSON types and constraints before binding. See [typed members](outputs.md#typed-members).
