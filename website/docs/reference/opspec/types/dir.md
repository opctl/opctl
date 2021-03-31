---
title: Dir
---

Dir typed values are a filesystem directory entry. Dirs are passed by reference, not copied, which means changing the contents or value of a directory will cause it to change in each variable that holds it.

## Initialization

Dir values can be constructed literally or through a [variable reference](../op.yml/variable-reference.md).
 
### Literal initialization

A literal dir initialization is a yaml key-value object. Keys are a single-component, absolute path name. Values are either an object with the key `data` and value a value or [variable reference](../op.yml/variable-reference.md) that can be coerced to a file, or another literal dir initialization.

```yaml
myLiteralDir:
  /singleLineFile:
    data: contents of /childFile1
  /childFile
    data: $(someVariable) # will contain the value of `someVariable`
  /subDir:
    /multilineFile:
      data: |
        multiline
        contents of /childFile2
  /emptySubDir: {}
```

## Coercion

Dir values cannot be coerced to other types.
