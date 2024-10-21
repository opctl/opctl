## `--api-listen-address` or `OPCTL_API_LISTEN_ADDRESS` *default: 127.0.0.1:42224*
To specify the HOST:PORT on which the opctl API server will listen, include a `--api-listen-address` or set an `OPCTL_API_LISTEN_ADDRESS` env var.

### Examples
```sh
opctl --api-listen-address 0.0.0.0:42224
```

## `--container-runtime` or `OPCTL_CONTAINER_RUNTIME` *default: docker*
To specify the container runtime opctl uses to run containers, include a `--container-runtime` or set an `OPCTL_CONTAINER_RUNTIME` env var.

Allowed values are:
- `k8s` (connects to k8s via downward API)
- `docker` (connects to docker via same config as docker client)
- `qemu` (experimental)

### Examples
```sh
opctl --container-runtime qemu
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

## `--dns-listen-address` or `OPCTL_DNS_LISTEN_ADDRESS` *default: 127.0.0.1:53*
To specify the HOST:PORT on which the opctl DNS server will listen, include a `--dns-listen-address` or set an `OPCTL_DNS_LISTEN_ADDRESS` env var.

### Examples
```sh
opctl --dns-listen-address 0.0.0.0:53
```

## `-h` or `--help`
For context specific help, include a `-h` (or `--help`) flag w/ your command.

### Examples
```sh
opctl node create -h

Usage: opctl node create

Creates a node
```

## `--nc`
To disable color, include a `--nc` flag w/ your command.
> this may increase readability in environments not supporting color escape codes or piping output to another program.

### Examples
```sh
opctl --no-color events
```

## `-v` or `--version`
To print the version and exit, specify a `-v` (or `--version`) flag.

### Examples
```sh
opctl -v
```