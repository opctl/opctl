# Apple VF (Virtualization Framework) Container Runtime

This package implements a container runtime for opctl using Apple's Virtualization Framework via vfkit.

## Overview

The Apple VF runtime provides container execution on macOS using vfkit, which is a lightweight virtualization tool that leverages Apple's Virtualization Framework. This runtime automatically downloads and caches the latest vfkit binary from GitHub releases.

## Features

- **Automatic Download**: Downloads the latest vfkit binary from GitHub if not present
- **Caching**: Caches the vfkit binary in `~/.opctl/cache/vfkit/` for future use
- **macOS Only**: Uses build constraints to ensure it only builds on macOS
- **Container Runtime Interface**: Implements the standard opctl ContainerRuntime interface

## Usage

The vfkit runtime is automatically selected as the default on macOS. Users can also explicitly specify it:

```bash
opctl node create --container-runtime=vfkit
```

## Implementation Details

### Binary Management

1. **PATH Check**: First checks if `vfkit` is available in the system PATH
2. **Cache Check**: Looks for cached vfkit binary in `~/.opctl/cache/vfkit/vfkit`
3. **Download**: If not found, downloads the latest release from `crc-org/vfkit` GitHub repository
4. **Permissions**: Makes the downloaded binary executable

### Current Status

This implementation provides a comprehensive vfkit-based container runtime with:

- **Binary Management**: Automatic download and caching of vfkit binary
- **VM Configuration**: Complete VM setup with kernel, initrd, root filesystem
- **Command Line Building**: Proper vfkit argument construction
- **Resource Management**: Memory, CPU, and network configuration
- **Process Execution**: Full vfkit command execution with stdout/stderr handling
- **Cleanup**: Automatic temporary file cleanup

### Implementation Details

The `RunContainer` method now performs these steps:

1. **Binary Resolution**: Downloads/caches vfkit if not present
2. **VM Directory Setup**: Creates temporary directory for VM files
3. **Configuration Preparation**: Sets up VM config (kernel, rootfs, networking)
4. **Argument Building**: Constructs complete vfkit command line
5. **Execution**: Runs vfkit with proper context and I/O handling
6. **Cleanup**: Removes temporary files and handles process termination

### Current Limitations

- Uses placeholder kernel/initrd files (would need real Linux boot files)
- Basic root filesystem creation (would need container image extraction)
- Simplified command execution (would need proper init system integration)

### Future Enhancements

For production use, the implementation could be extended with:

- **Real Kernel/Initrd**: Download from Linux distributions
- **Container Image Support**: OCI/Docker image pulling and extraction
- **Advanced Networking**: Multiple network configurations
- **Volume Management**: Bind mounts, named volumes, tmpfs
- **Init System Integration**: Proper process management inside VM

## Architecture

```
appleVF/
├── vfkit.go          # Main runtime implementation
├── vfkit_test.go     # Unit tests
└── README.md         # This documentation
```

## Dependencies

- `github.com/crc-org/vfkit` - The vfkit binary (downloaded automatically)
- macOS with Virtualization Framework support
- Go 1.24+ (for build constraints)

## Build Constraints

This package uses `//go:build darwin` to ensure it only builds on macOS systems where Virtualization Framework is available.
