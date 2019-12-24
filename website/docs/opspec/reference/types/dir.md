---
title: Dir
---

Dir typed values are a filesystem directory entry.

Dirs...
- are mutable, i.e. making changes to a directory results in the directory being changed everywhere it's referenced
- can be passed in/out of ops via [dir parameters](../structure/op-directory/op/parameter/dir)
- are coerced according to [coercion](#coercion)

### Entry Referencing
Dir entries (child files/directories) can be referenced via `$(ROOT/ENTRY)` syntax.

#### Entry Referencing Example (embedded)
given:
- `/file1.json` is embedded in op

```yaml
$(/file1.json)
```

#### Entry Referencing Example (scope)
given:
- `someDir`
  - is in scope dir
  - contains `file2.txt`

```yaml
$(someDir/file2.txt)
```

### Coercion
Object typed values are coercible to:

- [boolean](#boolean) (dirs which are empty coerce to false; all else coerce to true)