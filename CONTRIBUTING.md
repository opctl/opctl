# Dev Ops

Ops are maintained in
[![opspec 0.1.4](https://img.shields.io/badge/opspec-0.1.4-brightgreen.svg)](https://opspec.io/0.1.4/packages.html#format)
package format

They can be consumed via tools like [opctl](https://opctl.io).

# Acceptance criteria

Contributions are subject to:

- acceptance by 66% of the projects maintainers (see
  [MAINTAINERS.md](MAINTAINERS.md))
- the [build](.opspec/build) op continuing to run with a successful
  outcome
- adherence to
  [go code review comments](https://github.com/golang/go/wiki/CodeReviewComments)


# Repo organization

## /cli

CLI, distributed w/ the opctl binary

The web app is built using [mow](https://github.com/jawher/mow.cli)

## /docs

docs, hosted at [https://opctl.io/docs](https://opctl.io/docs)

## /node

daemon, distributed w/ the opctl binary

hosts the opctl web app & an opspec node

## /webapp

web app, distributed w/ the opctl binary & hosted by the opctl daemon.

It is a static web app built using
[react](https://facebook.github.io/react/) & was bootstrapped with
[Create React App](https://github.com/facebookincubator/create-react-app).

## /website

opctl website, hosted at [https://opctl.io](https://opctl.io)

It is a static website built using
[metalsmith](https://github.com/metalsmith/metalsmith)
