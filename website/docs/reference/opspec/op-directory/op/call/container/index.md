---
sidebar_label: Overview
title: Container Call [object]
---

An object defining a container call.

## Properties
- must have
  - [image](#image)
- may have
  - [dnsNames](#dnsNames)
  - [cmd](#cmd)
  - [dirs](#dirs)
  - [envVars](#envvars)
  - [files](#files)
  - [sockets](#sockets)
  - [workDir](#workdir)

### dnsNames
An [array](../../../../types/array.md) [initializer](../../../../types/array.md#initialization) or [variable-reference [string]](../../variable-reference.md) defining names opctl DNS resolves to this container.

> if a name resolves to multiple containers, network requests will be distributed (load balanced) across them. 

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
A [string initializer](../../../../types/string.md#initialization) defining a name opctl DNS resolves to this container.

> if a name resolves to multiple containers, network requests will be distributed (load balanced) across them. 

### sockets
An object for which each key is an absolute path in the container and and each value is a [socket](../../../../types/socket.md) [variable-reference [string]](../../variable-reference.md) to mount. 

### workDir
A [string initializer](../../../../types/string.md#initialization) defining absolute path from which [cmd](#cmd) will be executed.

> defining workDir overrides any defined by the image