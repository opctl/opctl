Files &/or directories can be accessed via `$(dir/path)` syntax. 

Dir can be a dir from scope, or the op (`/`)

```yaml
# access file1.json embedded in op
$(/file1.json)

# access file2.txt in scope someDir directory
$(someDir/file2.txt)
```