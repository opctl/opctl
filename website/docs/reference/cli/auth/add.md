---
sidebar_label: add
title: opctl auth add
---

```sh
opctl auth add RESOURCES [ -u=<username> ] [ -p=<password> ]
```

Add default auth used to pull ops and images.

## Arguments

### `RESOURCES`
Op or image ref prefixes this auth applies to (e.g. docker.io, github.com/some-org, etc.)

## Options

### `-u` or `--username`
Username
### `-p` or `--password`
Password
## Global Options
see [global options](../global-options.md)

### Examples

#### [github.com](https://github.com)
```sh
opctl auth add github.com -u <username> -p <password>
```

#### [docker.io](https://hub.docker.com)
```sh
opctl auth add docker.io -u <username> -p <password>
```