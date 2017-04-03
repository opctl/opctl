## Introduction

This document is a manual for maintainers old and new. It explains their
responsibilities, how they are added/removed, and how decisions are made
amongst them.

This is a living document - if you see something out of date or missing,
speak up!

## Maintainers

| Name          | Github Username |
|:--------------|:----------------|
| Chris Dostert | chrisdostert    |

### Maintainer responsibility

* 1) Deliver prompt feedback and decisions on pull requests.

* 2) Be available to anyone with questions, bug reports, criticism etc.
  about the project. This includes GitHub issues and pull requests.

* 3) Make sure pull requests respect the philosophy, design and roadmap
  of the project.

### Adding maintainers

The best maintainers have a vested interest in the project. Maintainers
are first and foremost contributors that have shown they are committed
to the long term success of the project. Contributors wanting to become
maintainers are expected to be deeply involved in contributing code,
pull request review, and triage of issues in the project for more than
two months.

### Removing maintainers

When a maintainer is unable to perform their required responsibilities
they should be removed.

### Chief Maintainer

The Chief Maintainer for the project is responsible for overall
architecture of the project to maintain conceptual integrity. Large
decisions and architecture changes should be reviewed by the chief
maintainer.

The current chief maintainer for the project is Chris Dostert
(@chrisdostert).

## Decision Making

[opspec-io/opctl](https://github.com/virtual-go/vos) is an
open-source project with an open design philosophy. This means that the
repository is the source of truth for EVERY aspect of the project,
including its philosophy, design, roadmap and APIs. *If it's part of the
project, it's in the repo.*

All decisions affecting
[opspec-io/opctl](https://github.com/virtual-go/vos), big and
small, follow the same steps:

* Step 1: A pull request is opened.

* Step 2: The pull request is discussed publicly.

* Step 3: [Maintainers](#maintainers) `Approve` or `Reject` the pull
  request by commenting as such on the pull request.

* Step 4: If 66% of the current maintainers `Approve` the pull request
  is merged (with chief maintainer holding veto power).
