---
title: Introduction
sidebar_label: Introduction
---

Opctl is a distributed, event-driven runtime which visualizes and runs graphs of containerized workloads called ops (short for operations). 
Ops are defined in a language called opspec (portmanteau of operation specification).

## Use Cases
- Conventional interface to empower would be contributors and break down silos.
- Executable documentation.
- Portable dev ops (running apps and their dependencies, continuous integration/delivery pipelines, building, testing, debugging, deploying...) that live and change with the code they operate.
- Portable ML (machine learning) pipelines.

## Opspec
Opspec is a language designed to portably and fully define operations (ops). See [reference docs](reference/opspec/index.md) for full details.

It features:
- Containers as first class citizens
- Serial and parallel looping and execution
- Conditional execution
- Variables and scoping
- Explicit inputs/outputs with type specific constraints
- Composition and re-use of ops
- Versioning via standard source control e.g. Git
- Array, boolean, dir, file, number, object, socket, and string data types
- Type coercion
- Declarative dependencies between calls

## CLI
The opctl CLI allows managing ops, nodes, events (logs), and opctl updates. See [reference docs](reference/cli/index) for full details.

![opctl CLI](/img/opctl-cli.png)

## UI
The opctl UI allows visualizing local or remote ops.

When visualizing an op, it's call graph is expanded recursively and you can zoom or pan around to focus where needed.
The call graph is rendered from top to bottom in order of execution. 

![opctl UI](/img/opctl-ui.png)

## Runtime
The opctl runtime is hosted by nodes. The CLI automatically [daemonizes](https://en.wikipedia.org/wiki/Daemon_(computing)) a node on the current host the first time it's needed (nodes can be manually managed via the `opctl node` sub command).

The runtime is virtual; it leverages overlay networking, overlay filesystem, and containerization to isolate ops from the hosts on which they run.
