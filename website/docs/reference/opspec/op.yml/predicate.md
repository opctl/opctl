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

### Verify a variable is in scope

In this example, `variable` has not been defined, so `../op` is not run.

```yml
run:
  if: 
    - exists: $(variable)
  op:
    ref: ../op
```

### Verify a directory or file exists

This can be use to check if a build artifact has been created, for example.

```yml
run:
  if: 
    - exists: $(./myFile.txt)
  op:
    ref: ../op
```

### Check if an object contains a key

In this example, "hello world" is printed because the key `foo` in the variable exists. Without the `exists` check in the second call, the op would crash because of an unresolvable reference.

```yml
inputs:
  value:
    object:
      default: { foo: "hello world" }
run:
  serial:
    - if:
        - exists: $(value.foo)
      container:
        image: { ref: alpine }
        cmd: ["echo", $(value.foo)]
    - if:
        - exists: $(value.bar)
      container:
        image: { ref: alpine }
        cmd: ["echo", $(value.bar)]
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

### Define a variable if it is not in scope

In this example, `variable` has not been defined, so the first container sets it before the second uses it.

```yml
serial:
  - if:
      - notExists: $(variable)
    container:
      image: { ref: alpine }
      cmd: [sh, -c, echo "hello world" > /output]
      files:
        /output: $(variable)
  - container:
      image: { ref: alpine }
      cmd: ["echo", $(variable)]
```

### Skip a build process if artifacts exist

This can be use to speed up dev ops, if a dependency rarely changes

```yml
if: 
  - notExists: $(./build)
op:
  ref: ../build
```

### Build a list of items

The following example will continually add items to the `variable` array until it contains four items. Once done, it will print the full array: `["first", "data", "data", "data"]`.

```yml
inputs:
  variable:
    array:
      default: ["first"]
run:
  serial:
    - serialLoop:
        until:
          - exists: $(variable[3])
        run:
          container:
            image: { ref: alpine }
            cmd:
              - sh
              - -c
              - |
                # this uses a mix of opctl variable expansion and sh variable expansion
                value='$(variable)'
                echo -n "${value%?}" > /output
                echo -n ', "data"]' >> /output
            files:
              /output: $(variable)
    - container:
        image: { ref: alpine }
        cmd: [echo, $(variable)]
```
