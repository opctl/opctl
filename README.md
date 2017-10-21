[![Build Status](https://travis-ci.org/opspec-io/sdk-js.svg?branch=master)](https://travis-ci.org/opspec-io/sdk-js)[![Coverage](https://codecov.io/gh/opspec-io/sdk-js/branch/master/graph/badge.svg)](https://codecov.io/gh/opspec-io/sdk-js)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

Javascript SDK for [opspec](https://opspec.io)

# Supported runtimes

This library is isomorphic & should be consumable from either nodejs or
web browsers.

# Installation

```shell
npm install --save @opspec/sdk
```

# Usage

## Node api client

```javascript
const OpspecNodeApiClient = require('@opspec/sdk/lib/node/apiClient');

const demoOpctlNodeBaseUrl = 'https://alpha.opctl.io/api';
// for local opctl node use
// const localOpctlNodeBaseUrl = 'localhost:42224/api';

const opspecNodeApiClient = new OpspecNodeApiClient({ baseUrl: demoOpctlNodeBaseUrl });

opspecNodeApiClient.liveness_get()
  .then(() => console.log('node alive!'))
  .catch(err => console.log(`error checking node; error was: ${err.message}`));

// start an op
const rootOpIdPromise = opspecNodeApiClient.op_start({
  args: {
    rawValue: {
      string: 'hello base64 url encoded world!',
    },
  },
  pkg: {
    ref: 'github.com/opspec-pkgs/base64url.encode#1.0.0',
  },
});

// wait for op to start then...
rootOpIdPromise.then(rootOpId => {
  
  // kill the op
  opspecNodeApiClient.op_kill({ opId: rootOpId })
  .then(() => console.log('successfully killed op!'))
  .catch(err => console.log(`error killing op; error was: ${err.message}`));
  
  // replay events via stream
  opspecNodeApiClient.event_stream_get({
    filter: {
      roots: [rootOpId],
    },
    onEvent: event => console.log(`received op event: ${JSON.stringify(event)}`),
    onError: err => console.log(`error streaming op events; error was: ${JSON.stringify(err)}`),
  });
  
});
```

# Support

join us on
[![Slack](https://opspec-slackin.herokuapp.com/badge.svg)](https://opspec-slackin.herokuapp.com/)
or [open an issue](https://github.com/opspec-io/sdk-js/issues)

# Releases

releases are versioned according to
[![semver 2.0.0](https://img.shields.io/badge/semver-2.0.0-brightgreen.svg)](http://semver.org/spec/v2.0.0.html)
and [tagged](https://git-scm.com/book/en/v2/Git-Basics-Tagging); see
[CHANGELOG.md](CHANGELOG.md) for release notes

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)
