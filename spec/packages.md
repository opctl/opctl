## Packages

A package is the distributable definition of an [op](ops.md).

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

[include](op.yml.schema.json)
