---
title: Bare Metal
sidebar_label: Bare Metal
---
## Installation

opctl is distributed as a self-contained executable, so installation generally consists of:

1. Downloading the OS specific binary
2. Adding it to your path

### Prerequisites
The default container runtime interface implementation relies on API access to a docker daemon to run containers.
[Install Docker for your platform](https://docs.docker.com/install/)

### OSX

```bash
curl -L https://github.com/opctl/opctl/releases/download/latest/opctl-darwin-amd64.tgz | sudo tar -xzv -C /usr/local/bin
```

### Linux

```bash
curl -L https://github.com/opctl/opctl/releases/download/latest/opctl-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin
```

### Windows

Use [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/) and the [linux install](#linux).

## Updating
to get the newest release of opctl
```bash
opctl self-update
```

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
