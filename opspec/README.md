# What is opspec?

Opspec is a high level language built for fully specifying operations start to finish.


# Ops

In Opspec, the highest level building block is an op. 

Ops serve a similar purpose as programs in other programming languages.

Like programs, ops can require inputs, produce outputs, and call other programs.

## Definition

Each op is defined via a file directory.

The file directory:
- MUST contain:
  - an [op file](#op-file)
- MAY contain:
  - an [icon file](#icon-file)
  - arbitrary files/dirs which MAY be referenced from the [op file](#op-file)


### Op file

Declares inputs, outputs, and call graph of the op.

Constraints:

- MUST be named (case sensitive) `op.yml`
- MUST be valid [v1.2 yaml](http://www.yaml.org/spec/1.2/spec.html)
- MUST validate against the [json schema](opfile/jsonschema.json)


### Icon file

Declares an icon for the op. 

Constraints:

- MUST be named (case sensitive) `icon.svg`
- MUST be valid [SVG 1.1](https://www.w3.org/TR/SVG11/)
- MUST have a 1:1 aspect ratio (equal height & width)


## Instantiation

Ops MAY be intantiated by a runtime such as [opctl](https://opctl.io) an arbitrary number of times, each creating a uniquely identifiable instance with it's own lifecycle.


## Distribution

Because ops are file directories, distribution is as simple as moving/copying them from one place to another.

It also means ops are natively supported by things like zip, tar, git, etc... 