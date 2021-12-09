---
title: Bare Metal
sidebar_label: Bare Metal
---
## Installation

opctl is distributed as a self-contained executable, so installation generally consists of:

1. Downloading the OS specific binary
2. Adding it to your path

### Dependencies

Opctl currently supports multiple container runtimes. Dependencies vary based on whichever you choose:
|container runtime|dependencies|
|--|--|
|docker|[docker](https://docs.docker.com/get-docker/)|
|k8s|opctl must be running inside k8s (uses the downward API) |
|qemu (experimental)|[lima](https://github.com/lima-vm/lima/releases/latest)|

### OSX

1. `curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-darwin-amd64.tgz | sudo tar -xzv -C /usr/local/bin`

### Linux

1. `curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin`

### Windows

Use [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/) and the [linux install](#linux).

## Upgrading from version 0.1.48 or newer
If you're updating from version 0.1.48 or newer:
```bash
# update to latest release from https://github.com/opctl/opctl/releases
opctl self-update
```

## Upgrading from version 0.1.47 or older
In 2021 our [previous artifact hosting provider](https://equinox.io/) shut down. If using version 0.1.47 or older, this results in `self-update` finding no updates and mistakenly thinking you're up to date. 

To update you must uninstall and reinstall:

1.  Kill any running node:
    ```bash
    opctl node kill
    ```
1.  Locate your data dir:
    ```bash
    # result should include default data dir such as:
    # --data-dir            Path of dir used to store opctl data (env $OPCTL_DATA_DIR) (default "/Users/myusername/Library/Application Support/opctl")
    opctl
    ```
1.  Delete your data dir:
    ```bash
    sudo rm -rf /path/to/opctl/data/dir
    ```
1. Follow [Installation](#installation) instructions

## IDE Plugins

### VSCode

1. install [vscode-yaml plugin](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
2. add to your user or workspace settings
   ```json
   "yaml.schemas": {
    "https://raw.githubusercontent.com/opctl/opctl/main/opspec/opfile/jsonschema.json": "/op.yml"
    }
    ```
3. edit or create an op.yml w/ your fancy intellisense.
