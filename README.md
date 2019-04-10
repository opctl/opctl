[![Build Status](https://travis-ci.org/opctl/sdk-js.svg?branch=master)](https://travis-ci.org/opctl/sdk-js)
[![Coverage](https://codecov.io/gh/opctl/sdk-js/branch/master/graph/badge.svg)](https://codecov.io/gh/opctl/sdk-js)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

Javascript SDK for [opctl](https://opctl.io)

# Supported runtimes

This library is isomorphic & should be consumable from either nodejs or
web browsers.

# Installation

```shell
npm install --save @opctl/sdk
```

# Usage

## API client

```javascript
import {
  eventStreamGet,
  livenessGet,
  opKill,
  opStart
} from '@opctl/sdk/dist/api/client'

// opctl api available from nodes via localhost:42224/api
const apiBaseUrl = 'localhost:42224/api'

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
[![Slack](https://opctl-slackin.herokuapp.com/badge.svg)](https://opctl-slackin.herokuapp.com/)
or [open an issue](https://github.com/opctl/sdk-js/issues)

# Releases

releases are versioned according to
[![semver 2.0.0](https://img.shields.io/badge/semver-2.0.0-brightgreen.svg)](http://semver.org/spec/v2.0.0.html)
and [tagged](https://git-scm.com/book/en/v2/Git-Basics-Tagging); see
[CHANGELOG.md](CHANGELOG.md) for release notes

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)
