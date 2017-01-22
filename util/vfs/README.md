virtual file system interface & implementations

why does this package exist?

- Testing filesystem dependent code in isolation from the filesystem
  requires a filesystem interface, but this is not builtin to the
  standard lib.
- Filesystem development continues to be highly active. Reliant code
  should be enabled to readily swap filesystem implementations in order
  to stay current.

