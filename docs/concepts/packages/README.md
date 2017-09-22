# Packages

Packages define orchestrations of containerized processes [ops].

They contain:
- an `op.yml` file which declares inputs, outputs, and a call graph of
  container, op, serial, and parallel calls
- any static assets the call graph depends on

## Distribution

A package is just a directory containing an `op.yml` and optionally,
files &/or directories.

This means distribution is as simple as transferring the directory and
it's contents from one place to another.

It also means you can do things like zip, tar, & version control
packages.
