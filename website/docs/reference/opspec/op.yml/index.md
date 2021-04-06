---
sidebar_label: Overview
title: op.yml
---

The `op.yml` file is a [YAML 1.2](https://yaml.org/spec/1.2/spec.html) file. It defines an operation's metadata, inputs, outputs, and call graph.

The `op.yml` file must conform to the `opspec` [JSON schema](https://github.com/opctl/opctl/blob/main/opspec/opfile/jsonschema.json). Opctl will interpret your `op.yml` dynamically and may produce runtime errors for additional reason such as invalid variables or variable coercion.

## Properties

### `description`

_(since v0.1.6)_

A human friendly description of the op, written as a [markdown string](markdown.md).

### `inputs`

An object defining input values for the op. The key of each property is a string [identifier](identifier.md) that defines the name of the input the caller will use and the name of the variable in scope within the op. The value of each property is a [parameter object](parameter/index.md) that describes the data.

### `outputs`

An object defining output values produced by the op. The key of each property is a string [identifier](identifier.md) that defines the name of the output and the name of the variable in scope within the op that stores the final outputted value. The value of each property is a [parameter object](parameter/index.md) that describes the data.

### `run`

`run` is a [call object](call/index.md) describing what and how your op runs.

## Example

```yaml
name: example
description: |
  This is an example op.yml file that returns the input as the output
inputs:
  inputName:
    string:
      description: This is a string input
outputs:
  outputName:
    string:
      description: This is a string output
run:
  container:
    image: { ref: "alpine" }
    cmd: ["sh", "-c", "echo \"$(inputName)\" > /output"]
    files:
      /output: $(outputName)
```
