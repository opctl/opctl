---
title: What’s in an Op
sidebar_label: What’s in an Op
---

## Structure

An op is defined as a directory containing an [op.yml](#opyml) and optionally an [icon.svg](#iconsvg) at its root.
```sh
root-directory
├── op.yml # required
└── icon.svg #optional
```

Additional files/directories can be included, even other ops.

### icon.svg

An icon.svg describes the icon to use when displaying an op from a user interface. It MUST be valid [SVG 1.1](https://www.w3.org/TR/SVG11/) with 1:1 aspect ratio (equal height & width)

> SVG is a vector (as opposed to raster) graphic format; it scales infinitely large/small w/out loss of quality.

### op.yml

An `op.yml` describes the inputs, outputs, and call graph of an operation. It MUST be valid [opspec](opspec.md)