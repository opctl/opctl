---
sidebar_label: Index
title: Op [directory]
---
A directory which defines an operation.

## Entries
- must have
    - [op.yml](#opyml)
- may have
    - [icon.svg](#iconsvg)
    - arbitrary files/directories embedded in the op

### icon.svg
An optional [SVG 1.1](https://www.w3.org/TR/SVG11/) file defining the icon to use when displaying the operation from a user interface. It MUST have a 1:1 aspect ratio (equal height & width)

> SVG is a vector (as opposed to raster) graphic format; it scales infinitely large/small w/out loss of quality.

### op.yml
A [YAML 1.2](https://yaml.org/spec/1.2/spec.html) file whos content is an [op](op/index) object defining the operations inputs, outputs, and call graph.