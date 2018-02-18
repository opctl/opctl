Packages define orchestrations of containerized processes such that they
MAY be version controlled & distributed via standard, open protocols
such as `git` & `HTTP` (respectively).

In concept, they are similar to
- [NPM packages](https://docs.npmjs.com/getting-started/packages)
- [Apple bundles](https://developer.apple.com/library/content/documentation/CoreFoundation/Conceptual/CFBundles/Introduction/Introduction.html)
- and so on...

Packages are plain old directories with the following characteristics:

- MUST contain
  - an [op.yml](../../reference/op.yml/README.md) file declaring inputs,
    outputs, and call graph (container, op, parallel, and serial calls)
    of an operation
- MAY contain
  - files/dirs the call graph depends on

## Distribution

Because a package is a stateless directory, distribution is as simple as
transferring the directory and it's contents from one place to another.

It also means you can do things like zip, tar, & version control
packages.
