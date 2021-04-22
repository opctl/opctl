---
sidebar_label: Index
title: Parameter [object]
---

An object defining a parameter of an operation; i.e. a value that is passed into or out of it's scope.

## Properties
- must have exactly one of
  - [array](array.md)
  - [boolean](boolean.md)
  - [dir](dir.md)
  - [file](file.md)
  - [number](number.md)
  - [object](object.md)
  - [socket](socket.md)
  - [string](string.md)

## Example

This is an example op that uses a string input, with a default value

```yaml
name: example
description: an example op
inputs:
  example-input:
    string:
      default: "a default value"
run:
  container:
    image: { ref: 'alpine' }
    cmd: ['echo', $(example-input)]
```

Using the example op:
```shell-script
opctl run example
```
The expected output is the op running and echoing "a default value".

Using the example op while overriding the default value:
```shell-script
opctl run -a example-input="hello world" example
```
The expected output is the op running and echoing "hello world".
