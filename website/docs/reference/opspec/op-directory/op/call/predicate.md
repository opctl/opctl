---
title: Predicate [object]
---

An object defining a condition which evaluates to true or false.

## Properties
- must have exactly one of
  - [eq](#eq)
  - [exists](#exists)
  - [gt](#gt)
  - [gte](#gte)
  - [lt](#lt)
  - [lte](#lte)
  - [ne](#ne)
  - [notExists](#notexists)

### eq
An array defining a predicate, true when all items are equal.

Items:
- must be one of
  - [variable-reference [string]](../variable-reference.md)
  - [initializer](../initializer.md)

### exists
A [variable-reference [string]](../variable-reference.md) defining a predicate, true when the referenced value exists.

### gt
An array defining a predicate, true when each item is greater than the next.

Items:
- must be a [number initializer](../../../types/number.md#initialization)

### gte
An array defining a predicate, true when each item is greater than or equal to the next.

Items:
- must be a [number initializer](../../../types/number.md#initialization)

### lt
An array defining a predicate, true when each item is less than the next.

Items:
- must be a [number initializer](../../../types/number.md#initialization)

### lte
An array defining a predicate, true when each item is less than or equal to the next.

Items:
- must be a [number initializer](../../../types/number.md#initialization)

### ne
An array defining a predicate, true when one or more items aren't equal.

Items:
- must be
  - [variable-reference [string]](../variable-reference.md)
  - [initializer](../initializer.md)

### notExists
A [variable-reference [string]](../variable-reference.md) defining a predicate, true when the referenced value doesn't exist.
