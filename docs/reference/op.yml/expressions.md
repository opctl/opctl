Expressions can be used where noted.

All expressions begin w/ `$(` and end w/ `)`.

# Property accessor

Properties (of objects) can be accessed via `.propertyName` syntax

```yaml
# access property of someObj
$(someObj.someProp)
```

# Path accessor

Files &/or directories can be accessed via `/path` syntax

```yaml
# access path of directory containing op.yml
$(/file1.json)
$(/someDir/file2.txt)

# access path of someDir
$(someDir/file2.txt)
```
