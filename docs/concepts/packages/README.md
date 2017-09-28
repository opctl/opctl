# Packages

Packages define orchestrations of containerized processes [ops].

They contain:
- an `op.yml` file declaring inputs, outputs, and a call graph
  (consisting of container, op, parallel, and serial calls)
- files/dirs the call graph depends on

## Distribution

Because a package is a directory, distribution is as simple as
transferring the directory and it's contents from one place to another.

It also means you can do things like zip, tar, & version control
packages.
