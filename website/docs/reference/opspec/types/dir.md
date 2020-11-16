---
title: Dir
---

Dir typed values are a filesystem directory entry.

Dirs...
- are mutable, i.e. making changes to a directory results in the directory being changed everywhere it's referenced.
- can be passed in/out of ops via [dir parameters](../op-directory/op/parameter/dir.md).
- can be initialized via [dir initialization](#initialization)
- are not coercible to any other type.

### Initialization
Dir typed values can be constructed from a literal or templated object.
 
A templated object is an object which includes one or more [variable-reference [string]](../op-directory/op/variable-reference.md).
At runtime, each reference gets evaluated and replaced with it's corresponding value.

#### Initialization Example (literal)

```yaml
myLiteralDir:
    /singleLineFile:
      data: contents of /childFile1
    /subDir:
      /multilineFile:
        data: |
          multiline
          contents of /childFile2
    /emptySubDir: {}
```

#### Initialization Example (templated)

```yaml
myTemplatedDir:
    /childFile:
      data: $(someVariable)
```

### Entry Referencing
Dir entries (child files/directories) can be referenced via `$(ROOT/ENTRY)` syntax.

#### Entry Referencing Example (embedded)
given:
- `file1.json` exists in op

```yaml
$(./file1.json)
```

#### Entry Referencing Example (scope)
given:
- `someDir`
  - is in scope dir
  - contains `file2.txt`

```yaml
$(someDir/file2.txt)
```