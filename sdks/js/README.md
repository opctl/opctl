[![Build Status](https://travis-ci.org/opctl/sdk-js.svg?branch=master)](https://travis-ci.org/opctl/sdk-js)
[![Coverage](https://codecov.io/gh/opctl/sdk-js/branch/master/graph/badge.svg)](https://codecov.io/gh/opctl/sdk-js)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

# Problem statement
Javascript SDK for [opctl](https://opctl.io).


# Supported runtimes
This library is isomorphic & should be consumable from either nodejs or
web browsers.


# Installation
```shell
npm install --save @opctl/sdk
```


# Typescript
This library is written in typescript. The package published to NPM targets ES2015 and includes type declarations.


# Usage

## Run an op using an API client.
```javascript
import {
  eventStreamGet,
  livenessGet,
  opKill,
  opStart
} from '@opctl/sdk/dist/api/client'

// by default, opctl node API available at 127.0.42.224/api
const apiBaseUrl = '127.0.42.224/api'

// get the liveness of the api
await livenessGet(
  apiBaseUrl
)

// start an op
const rootOpId = await opStart(
  apiBaseUrl,
  {
    rawValue: {
      string: 'hello base64 url encoded world!',
    },
  },
  {
    ref: 'github.com/opctl-pkgs/base64url.encode#1.0.0',
  }
)

// kill the op we started
await opKill(
  apiBaseUrl,
  rootOpId
)

// replay events from our op via stream
await eventStreamGet(
  apiBaseUrl,
  {
    filter: {
      roots: [rootOpId],
    },
    onEvent: event => console.log(`received op event: ${JSON.stringify(event)}`),
    onError: err => console.log(`error streaming op events; error was: ${JSON.stringify(err)}`),
  }
)
```


# Support
join us on
[![Slack](https://img.shields.io/badge/slack-opctl-E01563.svg)](https://join.slack.com/t/opctl/shared_invite/zt-51zodvjn-Ul_UXfkhqYLWZPQTvNPp5w)
or [open an issue](https://github.com/opctl/opctl/issues)


# Releases
releases are versioned according to
[![semver 2.0.0](https://img.shields.io/badge/semver-2.0.0-brightgreen.svg)](http://semver.org/spec/v2.0.0.html)
and [tagged](https://git-scm.com/book/en/v2/Git-Basics-Tagging); see
[CHANGELOG.md](CHANGELOG.md) for release notes


# Contributing
see [CONTRIBUTING.md](CONTRIBUTING.md)
