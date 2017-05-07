### Filesystem

The filesystem observed by nodes is described in this section with the following
semantics:

- PER_USER_APP_DATA_PATH must be interpreted as defined by
  [app data spec](https://github.com/appdataspec/spec)
- `/` must be interpreted as the native path segment
  delimiter of the platform.

#### NODE_ROOT_PATH

The root path of a nodes filesystem is `PER_USER_APP_DATA_PATH/opspec`

#### NODE_CACHED_PKG_PATH

The path at which a node caches a package is determined by:

1. substituting the last occurrence of `#` in `pkgRef` w/ `/`
2. joining NODE_ROOT_PATH, `pkg-cache`, and the result of the previous
   step w/ `/`


example file tree resulting from caching packages
`hostname1.com/nspart1/op1#1.1.1`,`hostname1.com/nspart1/op1#2.2.2`, and
`hostname2.com/nspart1/nspart2/op2#3.3.3`:

```
NODE_ROOT_PATH
  |-- pkg-cache
     |-- hostname1.com
        |-- nspart1
           |-- op1
              |-- 1.1.1
                 |-- op.yml
                 ... (pkg specific files/dirs)
              |-- 2.2.2
                 |-- op.yml
                 ... (pkg specific files/dirs)
        |-- nspart1
           |-- nspart2
              |-- op2
                 |-- 3.3.3
                    |-- op.yml
                    ... (pkg specific files/dirs)
```

