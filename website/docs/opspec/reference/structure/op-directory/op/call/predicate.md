---
title: Predicate [object]
---

An object defining a condition which evaluates to true or false.

## Properties
- must have exactly one of
  - [eq](#eq)
  - [exists](#exists)
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

### ne
An array defining a predicate, true when one or more items aren't equal.

Items:
- must be
  - [variable-reference [string]](../variable-reference.md)
  - [initializer](../initializer.md)

### notExists
A [variable-reference [string]](../variable-reference.md) defining a predicate, true when the referenced value doesn't exist.