---
title: Javascript
sidebar_label: Javascript
---

## Examples

Run an op using the [Javascript SDKs](https://github.com/opctl/sdk-js) API client.
```typescript
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