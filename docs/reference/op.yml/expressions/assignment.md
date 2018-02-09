Assignment expressions, when evaluated, assign a value to something.

# Examples

## Assign to file

> `rootDir` is assumed to be an in scope directory

```yaml
$(rootDir/subDir/file): hello
```

## Assign to property

> `rootObj` is assumed to be an in scope object

```yaml
$(rootObj.subProp): 2
```

