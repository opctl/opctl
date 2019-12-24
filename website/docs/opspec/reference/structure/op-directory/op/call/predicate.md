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
  - [reference](../reference) string
  - [initializer](../initializer)

### exists
A [reference](../reference) string defining a predicate, true when the referenced value exists.

### ne
An array defining a predicate, true when one or more items aren't equal.

Items:
- must be
  - [reference](../reference) string
  - [initializer](../initializer)

### notExists
A [reference](../reference) string defining a predicate, true when the referenced value doesn't exist.