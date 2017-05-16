### Filesystem

The filesystem observed by nodes is described in this section with the following
semantics:

- PER_USER_APP_DATA_PATH must be interpreted as defined by
  [app data spec](https://github.com/appdataspec/spec)
- `/` must be interpreted as the native path segment
  delimiter of the platform.

#### NODE_ROOT_PATH

The root path of a nodes filesystem is `PER_USER_APP_DATA_PATH/opspec`
