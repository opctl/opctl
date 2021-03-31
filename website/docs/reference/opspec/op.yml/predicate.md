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

A [variable reference](variable-reference.md). The predicate is true when the referenced value exists.

```yml
run:
  if: 
    - exists: $(variable)
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

A [variable reference](variable-reference.md). The predicate is true when the referenced value does not exist.

```yml
run:
  if: 
    - notExists: $(variable)
  op:
    ref: ../op
```
