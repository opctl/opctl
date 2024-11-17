---
sidebar_label: create
title: opctl op create
---

```sh
opctl op create [OPTIONS] NAME
```

Create an op.

## Arguments

### `NAME`
Name of the op

## Options

### `-d` or `--description`
Description of the op

### `--path` *default: `.opctl`*
Path to create the op at

## Global Options
see [global options](../global-options.md)

## Examples
```sh
opctl op create -d "my awesome op description" --path some/path my-awesome-op-name
```