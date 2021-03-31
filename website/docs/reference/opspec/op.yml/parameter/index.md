---
sidebar_label: Overview
title: Parameters
---

A parameter is an object describing a value passed into or out of an op through it's inputs or outputs. Each parameter must have a single [type property](#type-properties) that declares it's type.

## Properties

### `description`

A human friendly description of the parameter, written as a [markdown string](markdown.md).

### `{type}`

Each parameter declares it's type with one of the following property keys:

- [`array`](array.md)
- [`boolean`](boolean.md)
- [`dir`](dir.md)
- [`file`](file.md)
- [`number`](number.md)
- [`object`](object.md)
- [`socket`](socket.md)
- [`string`](string.md)

## Example

```yaml
inputs:
  parameterName:
    description: An example parameter of type array
    array: {}
```
