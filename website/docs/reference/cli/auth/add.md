---
sidebar_label: add
title: opctl auth add
---

```sh
opctl auth add RESOURCES [ -u=<username> ] [ -p=<password> ]
```

Add auth for an OCI image registry.

## Arguments

### `RESOURCES`
Resources this auth applies to in the form of a host or host/path.

## Options

### `-u` or `--username`
Username

### `-p` or `--password`
Password

## Global Options
see [global options](../global-options.md)

### Examples 

#### [docker.io](https://hub.docker.com)
```sh
opctl auth add docker.io -u <username> -p <password>
```