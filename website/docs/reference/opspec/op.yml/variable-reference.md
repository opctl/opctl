---
title: Variable reference
---

Variable references in an opspec file are how you reference data within your op.

A variable reference takes the form `$(reference)`, where `reference` can be the [name](identifier) of an [in scope](../index.md#scoping) variable, a newly defined variable [name](identifier), or an absolute or relative filesystem path.

If referencing a variable, `reference` can be extended with an array access, object property, or filesystem path, if the type of the variable matches or can be coerced.

## Escaping

Variable references can be escaped with `\`. Because opspecs are yaml, a second `\` is required to escape the escape when used in a normal string.

```yaml
foo: "\\$(wouldBeVariableReference)" # not treated as a variable reference
```

```yaml
foo: |
  \$(wouldBeVariableReference) # not treated as a variable reference
```

## Filesystem paths

Ops have a few filesystem paths in scope automatically.

- `./` is the the current op's directory, where the `op.yml` is located. (the current `op.yml` can be accessed with `$(./op.yml)`)
- `../` is the parent of the current op's directory, (the current `op.yml` can be also be accessed with `$(../op.yml)`)

## Array access

[Array](../types/array) items can be referenced with the syntax `$(reference[index])`, where `index` is the zero based index of the item.
If `index` is negative, indexing will take place from the end of the array.

```yaml
inputs:
  someArray:
    array:
      default: ["one", 2, [3]]
run:
  op:
    ref: ../op
    inputs:
      input1: $(someArray[0]) # "one"
      input2: $(someArray[-1][0]) # 3
```

## Object properties

[Object](../types/object) properties can be referenced with the syntax `$(reference.property)` or `$(reference[property])`, where `property` is the name of the property.

```yaml
inputs:
  someObject:
    object:
      default:
        myKey: aValue
        secondKey:
          subKey: 2
run:
  op:
    ref: ../op
    inputs:
      input1: $(someObject.myKey) # aValue
      input2: $(someObject[secondKey][subKey]) # 2
```

## Directory paths

[Directory](../types/dir) entries (child files and directories) can be referenced with the syntax `$(reference/entry)`, where `entry` is the subdirectory, file, or combination of subdirectories and files.

```yaml
inputs:
  someDir:
    dir:
      default: ./myDir
run:
  op:
    ref: ../op
    inputs:
      input1: $(someDir/file.txt) # the file at ./myDir/file.txt
      input2: $(someObject/a/b/c) # the file or directory at ./myDir/a/b/c
```

The default `./` and `../` references support these extended path references.

```yaml
run:
  op:
    ref: ../op
    inputs:
      opfile: $(./op.yml)
```
