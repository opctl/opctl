## `ui [MOUNT_REF]`
Opens the opctl web UI and mounts a reference.

### Arguments

#### `MOUNT_REF` *default: `.`*
Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`).

### Examples
Open web UI to current working directory
```sh
opctl ui
```

Open web UI to remote op (github.com/opspec-pkgs/_.op.create#3.3.1)
```sh
opctl ui github.com/opspec-pkgs/_.op.create#3.3.1
```