---
title: Predicate
---

A predicate is a condition that evaluates to a boolean value, used for conditional logic. The predicate is an object with a single key defining the type of condition.

## `eq`

An [array](../types/array). When all items in the array are equivalent, the predicate is true.

```yml
run:
  if: 
    - eq: [true, $(variable)]
  op:
    ref: ../op
```

## `exists`

A [variable reference](variable-reference.md). The predicate is true when the referenced value exists. One use case is to verify a directory or file has been generated.

```yml
run:
  if: 
    - exists: $(variable)
  op:
    ref: ../op
```

```yml
run:
  if: 
    - exists: $(./myFile.txt)
  op:
    ref: ../op
```

## `ne`

An [array](../types/array). When one or more of the items in the array aren't equivalent, the predicate is true.

```yml
run:
  if: 
    - ne: [true, $(variable)]
  op:
    ref: ../op
```

## `notExists`

A [variable reference](variable-reference.md). The predicate is true when the referenced value does not exist. One use case is to check if a directory or file has been generated.

```yml
run:
  if: 
    - notExists: $(variable)
  op:
    ref: ../op
```

```yml
inputs:
  src:
    dir: {}
run:
  if: 
    - exists: $(src/build)
  op:
    ref: ../op
```
