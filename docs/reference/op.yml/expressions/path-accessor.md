Files &/or directories can be accessed via `$(dir/path)` syntax. 

Dir can be a dir from scope, or the pkg root (`/`)

```yaml
# access file1.json in pkg root directory
$(/file1.json)

# access file2.txt in scope someDir directory
$(someDir/file2.txt)
```