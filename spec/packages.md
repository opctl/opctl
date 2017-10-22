## Packages

A package is the distributable definition of an [op](ops.md).

Packages are:

- composable
- stateless
- self-describing

Packages MUST contain an `op.yml` file declaring the inputs, outputs,
and call graph (consisting of container, op, parallel, and serial calls)
of the operation.

Packages MAY contain files/dirs the operation depends on

### op.yml

Constraints:

- MUST be named (case sensitive) `op.yml`
- MUST be valid [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against the [op.yml json schema](#opyml-json-schema)

#### op.yml JSON schema

> [open in schema visualizer](https://schema-visualizer.opspec.io/?schema=https://opspec.io/0.1.5/op.yml.schema.json)

[include](op.yml.schema.json)
