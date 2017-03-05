## collection

0-N [ops](#op) grouped together.

## node

A node is a process, typically run as a
[daemon or service](https://en.wikipedia.org/wiki/Daemon_(computing) ),
which exposes the
[opspec runtime API](https://github.com/opspec-io/spec/blob/master/runtime-oai_spec.yml),
servicing any requests it receives (since v0.1.15).

There can be only one node running at a time on a given machine.

## op

composition of containerized processes.
