## `ls [DIR_REF]`

List ops in a directory.

## Arguments

### `DIR_REF` *default: `.opspec`* (since v0.1.24)

Reference to dir ops will be listed from (either `relative/path`, `/absolute/path`,
`host/path/repo#tag` (since v0.1.19), or `host/path/repo#tag/path`
(since v0.1.24))

## Examples

### `.opspec` dir

lists ops from `./.opspec`

```shell
opctl ls
```

### remote dir

lists ops from [github.com/opctl/opctl#0.1.24](https://github.com/opctl/opctl/tree/0.1.24)

```shell
opctl ls github.com/opctl/opctl#0.1.24/
```