---
sidebar_label: Overview
title: Opspec
---

Opspec (a portmanteau of operation specification) is a language designed to portably and fully define ops (operations).

An op is any directory containing a [valid opspec file with the name `op.yml`](op.yml/index). By default, the opctl cli will look for ops in the `.opspec/` directory in your current working directory, but ops can be defined in any directory.

An op directory may contain other files and directories. These files are considered "embedded" within the op and can be referenced during execution as executable code, as data, or for whatever purpose you have.

If the root of an op directory contains a `icon.svg` file (must be [SVG 1.1](https://www.w3.org/TR/SVG11) with a 1:1 aspect ratio), it will be used when displaying the op within a UI.

## Call graph

As an op runs, opctl parses [`op.yml`](op.yml/index) files to create a dynamic call graph tracking the different running containers as they are started and exit. The call graph is defined declaritively using [call objects](op.yml/call/index), which include sequencing and conditional support.

## Data flow

Opctl has first class support for data flow into, between, and out of ops and their containers. Opspecs receive data from [`inputs`](op.yml#inputs) and return data through [`outputs`](op.yml#outputs). Within an opspec, data is referenced via named variables that containing typed values.

### Types

Opspecs support a limited set of type coercion. This is important, as it bridges the gap between container data and opspec typed values. A container doesn't know how to produce a key/value object, but can write to a file which can then be coerced into an object.

The supported types are:

- [array](types/array.md)
- [boolean](types/boolean.md)
- [dir](types/dir.md)
- [file](types/file.md)
- [number](types/number.md)
- [object](types/object.md)
- [socket](types/socket.md)
- [string](types/string.md)

### Scoping

A variable is created when defined as an [`input`](op.yml#inputs) or [`output`](op.yml#inputs), bound implicitly to an `output` when making an [op call](op.yml/call/op), or declared with [loop variables](./op.yml/loop-variable). A variable name is in scope from creation through the lifetime of it's `op.yml` file. Variable scope does not extend into [child op calls](./op.yml/call/op) within the call graph.

Variables can be referenced using the syntax `$(name)`.

#### Scoping example

```yml
name: foobar
inputs:
  variable1:
    string: {}
  variable2:
    string:
      default: "hello world"
outputs:
  variable3:
    string:
      default: "goodbye!"
```

Two variables, `variable1` and `variable2` are created from inputs. The CLI will prompt the user to provide a value for `variable1`, since it doesn't have a default, and `variable2` will contain the string value `"hello world"` if the user didn't override the value with a command line argument.

One variable, `variable3` is created from the single output, which contains the string value `"goodbye!"`.

```yml
run:
  serial:
    - op:
        ref: ../other-op
        inputs:
          input1: $(variable1)
        outputs:
          output1: $(variable2)
          output2: $(variable4)
```

The variable `variable1`'s value is passed into `../other-op` where it will be stored in the variable `input1`. `../other-op` will not have access to the `variable1` or `variable2` names. After the op call completes the variable `variable2` will contain the value from `output1` that was produced by `../other-op`. Additionally, the variable `variable3` will be implicitly created from the output binding and will have the value from `output2` stored in it.

```yml
    - op:
        ref: ../third-op
        inputs:
          input1: $(variable2)
        outputs:
          output1: $(variable3)
```

`variable2` now contains the value from `../other-op`, not the original value of `"hello world"`.

The value of `output1` from `../third-op` will be stored in `variable3`, which was defined as this op's output. It will be returned to the caller through the `variable3` output.
