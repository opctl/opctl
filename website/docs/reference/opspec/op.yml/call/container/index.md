---
sidebar_label: Overview
title: Container call
---

A container call is an object defining how to run a container. Container calls are an abstraction on top of a concrete container running tool like docker or kubernetes that have a more strictly typed interface.

## Properties

### `image`

_required_

A container [image](image.md) to run.

### `cmd`

The binary to run and it's arguments, as an array of [strings](../../../types/string). This will override the images' default command, if defined. Strings will be interpreted by opctl with variable reference replacements applied.

The first string in the command can be the name binary available in the container's `PATH`, an absolute path to the binary file, or a relative path to the binary file from the [working directory](#workdir) of the container.

The remaining strings in the command are passed directly to the binary as individual string arguments.

```yaml
cmd: [echo, my string, $(variable)] # "echo" "my string" "valueOfVariable"
```

It's important to remember that the command is not run in a shell, which means shell-specific behavior such as piping (`|`) and redirection (`>`) will not work directly. To use a shell, you can run it as your command:

```yaml
cmd:
  - bash
  - -c
  - |
    # this is a full bash script
    echo "hello world" > /file
    cat /file
```

### `dirs`

An object mapping directories to mount within the container. Each key is an absolute path within the container. Values can take several forms:

|value|meaning|
|--|--|
|`null`|Shorthand to mount a directory embedded in the op.<br />For example, `/relative/path:` is equivalent to `/relative/path: $(./relative/path)`. If the `op.yml` is at `~/ops/foo/op.yml`, this will mount the directory `~/ops/foo/relative/path`.|
|[dir initializer](../../../types/dir.md#initialization)|The directory is evaluated and created, then mounted|
|[dir](../../../types/dir.md) [variable reference](../../variable-reference.md)|The referenced directory is mounted|

### `envVars`

An [object](../../../types/object.md) or [variable reference](../../variable-reference.md) mapping environment variable names to values to be set in the container. Value are coerced to strings.

### `files`

An object mapping files to mount within the container. Each key is an absolute path within the container. Values can take several forms:

|value|meaning|
|--|--|
|`null`|Shorthand to mount a file embedded in the op.<br />For example, `/relative/path:` is equivalent to `/relative/path: $(./relative/path)`. If the `op.yml` is at `~/ops/foo/op.yml`, this will mount the file `~/ops/foo/relative/path`.|
|[file initializer](../../../types/file.md#initialization)|The file is evaluated and created, then mounted|
|[file](../../../types/file.md) [variable reference](../../variable-reference.md)|The referenced file is mounted|

### `name`

A [string](../../../types/string.md#initialization) defining a name for the container. The container will be resolvable by this name within the opctl network, which allows making http requests between containers by. If multiple containers are given the same name between _any_ ops, requests may be made to any of them.

### `ports`

An object mapping ports between the container and opctl host. Each key is a container port or range of ports (optionally including protocol) which match the regular expression `[0-9]+(-[0-9]+)?(tcp|udp)`. Each value is a corresponding opctl host port or range of ports matching the regular expression `[0-9]+(-[0-9]+)?`.

For example, you can map the port 80 within a container to your local port 8080:

```yaml
ports:
  80: 8080
```

### `sockets`

An object mapping sockets to mount within the container. Each key is an absolute path within the container. Values are [sockets](../../../types/socket.md) or [variable references](../../variable-reference) to a socket to mount.

### `workDir`

A [string](../../../types/string.md#initialization) defining the absolute path from which [cmd](#cmd) will be executed. This will override the images' default working directory, if defined.
