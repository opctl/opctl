---
sidebar_label: Index
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
  - [ports](#ports)
  - [sockets](#sockets)
  - [workDir](#workdir)

### image
An [image [object]](image.md) defining the container image run by the call.

### cmd
An array of [string initializers](../../../../types/string.md#initialization) defining the path (from [workDir](#workdir)) of the binary to call and it's arguments.

> defining cmd overrides any entrypoint and/or cmd defined by the image

### dirs
An object for which each key is an absolute path in the container and each value is one of:

|value|meaning|
|--|--|
|null|Mount dir embedded in op w/ same path (equivalent to `$(/absolute/path)`)|
|[dir](../../../../types/dir.md) [variable-reference [string]](../../variable-reference.md)|Mount dir|

### envVars
An [object initializer](../../../../types/object.md#initialization) or [variable-reference [string]](../../variable-reference.md), whos properties represent the name and value of an environment variable to be set in the container.

> upon evaluation, the key and value of each property will be coerced to a string.

### files
An object for which each key is an absolute path in the container and each value is one of:

|value|meaning|
|--|--|
|null|Mount file embedded in op w/ same path (equivalent to `$(/absolute/path)`)|
|[file](../../../../types/file.md) [variable-reference [string]](../../variable-reference.md)|Mount file|
|[file initializer](../../../../types/file.md#initialization)|Evaluate and mount|

### name
A [string initializer](../../../../types/string.md#initialization) defining a name by which the container can be resolved on the opctl network.

> if multiple containers are given the same name, network requests will be distributed (load balanced) across them. 

### ports
An object defining container ports exposed on the opctl host where:
- each key is a container port or range of ports (optionally including protocol) matching `[0-9]+(-[0-9]+)?(tcp|udp)`
- each value is a corresponding opctl host port or range of ports matching `[0-9]+(-[0-9]+)?`

### sockets
An object for which each key is an absolute path in the container and and each value is a [socket](../../../../types/socket.md) [variable-reference [string]](../../variable-reference.md) to mount. 

### workDir
A [string initializer](../../../../types/string.md#initialization) defining absolute path from which [cmd](#cmd) will be executed.

> defining workDir overrides any defined by the image