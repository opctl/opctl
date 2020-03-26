---
title: Parallel Loop Call [object]
---

An object defining a call loop in which all iterations happen in parallel (all at once without order).

## Properties
- must have 
  - [range](#range)
  - [run](#run)
- may have
  - [vars](#vars)

### range
A [rangeable value](rangeable-value.md) to loop over.

### run
A [call [object]](../call/index.md) defining a call run on each iteration of the loop

### vars
A [loop-vars [object]](loop-vars.md) binding iteration info to variables.
