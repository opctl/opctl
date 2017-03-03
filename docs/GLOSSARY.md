## .opspec

See default collection

## collection

One or more op's, grouped together physically (via embedding) and/or
logically (via reference).

## default collection

The default collection referenced by tooling within a given directory.
It's designated by naming the collections containing dir `.opspec`
rather than the collection's name.

## node

A node is a process, typically run as a
[daemon or service](https://en.wikipedia.org/wiki/Daemon_(computing) ),
which exposes the
[opspec runtime API](https://github.com/opspec-io/spec/blob/master/runtime-oai_spec.yml),
servicing any requests it receives (since v0.1.15).

There can be only one node running at a time on a given machine.

## op

Orchestration of containerized processes which has been
[spec](https://opspec.io/)'d.
