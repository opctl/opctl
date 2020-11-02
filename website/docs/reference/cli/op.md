## `opctl op create [OPTIONS] NAME`
Creates an op

### Arguments

#### `NAME`
Name of the op

### Options

#### `-d` or `--description`
Description of the op

#### `--path` *default: `.opspec`*
Path to create the op at

### Examples
```sh
opctl op create -d "my awesome op description" --path some/path my-awesome-op-name
```

## `opctl op install [OPTIONS] OP_REF`
Installs an op

### Arguments

#### `OP_REF`
Op reference (`host/path/repo#tag`, or `host/path/repo#tag/path`)

### Options

#### `--path` *default: `.opspec/OP_REF`*
Path to install the op at

#### `-u` or `--username`
Username used to auth w/ the op source

#### `-p` or `--password`
Password used to auth w/ the op source

### Examples
```sh
opctl op install -u someUser -p somePass host/path/repo#tag
```

### Notes

#### op source username/password prompt
If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> and exit with a non zero exit code.

## `opctl op kill OP_ID`
Kill an op. 

### Arguments

#### `OP_ID`
Id of the op to kill

## `opctl op validate OP_REF`
Validates an op according to:

- existence of `op.yml`
- validity of `op.yml` (per
  [schema](https://opctl.io/0.1.6/op.yml.schema.json))

### Arguments

#### `OP_REF`
Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`).

### Examples
```sh
opctl op validate myop
```

### Notes

#### op source username/password prompt
If auth w/ the op source fails the cli will (re)prompt for username & password.

> in non-interactive terminals, the cli will note that it can't prompt and exit with a non zero exit code.
