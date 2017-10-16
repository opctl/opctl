[![Build Status](https://travis-ci.org/opspec-io/sdk-js.svg?branch=master)](https://travis-ci.org/opspec-io/sdk-js)[![Coverage](https://codecov.io/gh/opspec-io/sdk-js/branch/master/graph/badge.svg)](https://codecov.io/gh/opspec-io/sdk-js)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

Javascript SDK for [opspec](https://opspec.io)

# Installation

```shell
npm install --save @opspec/sdk
```

# Usage

```javascript
// nodejs version:
const OpspecNodeApiClient = require('@opspec/sdk/src/node/apiClient');
// browser version: 
// import OpspecNodeApiClient from '@opspec/sdk/lib/node/apiClient';

const opspecNodeApiClient = new OpspecNodeApiClient({ baseUrl: 'https://demo.opctl.io' });

opspecNodeApiClient.liveness_get();
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
