Dir entries (subdirs &/or files of a dir) can be referenced via `/dirEntry` syntax.

## Examples

### Embedded json file
given:
- `/file1.json` is embedded in op

```yaml
$(/file1.json)
```

### File of in scope directory
given:
- `someDir`
  - is in scope
  - is type coercible to dir
  - contains `file2.txt`

```yaml
$(someDir/file2.txt)
```