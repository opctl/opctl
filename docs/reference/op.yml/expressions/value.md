Value expressions, when evaluated, produce values.

# Examples

## Produce sub file/dir from parent dir

```yaml
subFileOrDirValue: $(parentDir/subFileOrDir)
```

## Produce string

prepends `pre! ` to value at $(rootDir/file)

```yaml
fileValue: pre! (rootDir/file)
```

## Produce object property

```yaml
propertyValue: $(objRoot.objProp)
```
