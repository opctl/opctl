---
sidebar_label: ls
title: opctl ls
---

```sh
opctl ls [DIR_REF=.opspec]
```

List ops in a local or remote directory.

### Arguments

#### `DIR_REF` *default: `.opspec`*
Reference to dir ops will be listed from (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)

## Global Options
see [global options](global-options.md)

### Examples

#### `.opspec` dir
lists ops from `./.opspec`

```sh
opctl ls
```

#### remote dir
lists ops from [github.com/opctl/opctl#0.1.24](https://github.com/opctl/opctl/tree/0.1.24)

```sh
opctl ls github.com/opctl/opctl#0.1.24/
```