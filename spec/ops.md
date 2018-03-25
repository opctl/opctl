# Ops

Ops define orchestrations of containerized processes.

Ops are:

- composable
- stateless
- self-describing

## Definition

An op MAY be defined via a file directory.

The file directory:
- MUST contain:
  - an [op.yml file](#opyml-file)
- MAY contain:
  - an [icon.svg file](#iconsvg-file)
  - arbitrary files/dirs which MAY be referenced from the [op.yml file](#opyml)

### op.yml file

Declares inputs, outputs, and call graph of the op.

Constraints:

- MUST be named (case sensitive) `op.yml`
- MUST be valid [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against the [op.yml json schema](op.yml.schema.json)

### icon.svg file

Declares an icon for the op. 

Constraints:

- MUST be valid [SVG 1.1](https://www.w3.org/TR/SVG11/)
- MUST have a 1:1 aspect ratio (equal height & width)

## Instantiation

Ops MAY be intantiated an arbitrary number of times, each creating a uniquely identifiable instance with it's own lifecycle.

## Distribution

Because ops are just stateless directories, distribution is as simple as moving/copying them from one place to another.

It also means ops are natively supported by things like zip, tar, git, etc... 