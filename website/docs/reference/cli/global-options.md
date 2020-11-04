## `-v` or `--version`
To print the version and exit, specify a `-v` (or `--version`) flag.

### Examples
```sh
opctl -v
```

## `--data-dir` or `OPCTL_DATA_DIR` *default: OS dependent per user app data*
To specify the path of the directory used to store opctl data, include a `--data-dir` or set an `OPCTL_DATA_DIR` env var.
to a relative or absolute path. 

### Examples
```sh
opctl --data-dir . node create
```

```sh
export OPCTL_DATA_DIR=. && opctl node create
```

## `--listen-address` or `OPCTL_LISTEN_ADDRESS` *default: 127.0.0.1:42224*
To specify the HOST:PORT on which the node will listen, include a `--listen-address` or set an `OPCTL_LISTEN_ADDRESS` env var.

### Examples
```sh
opctl --listen-address 0.0.0.0:42224
```

## `--nc`
To disable color, include a `--nc` flag w/ your command.
> this may increase readability in environments not supporting color escape codes or piping output to another program.

### Examples
```sh
opctl --no-color events
```

## `-h` or `--help`
For context specific help, include a `-h` (or `--help`) flag w/ your command.

### Examples
```sh
opctl node create -h

Usage: opctl node create

Creates a node
```
