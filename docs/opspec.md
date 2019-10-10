---
title: Opspec
sidebar_label: Opspec
---

## Introduction
Opspec is a declarative language designed for the sole purpose of describing operations.

## Syntax

- [description](#description)
- [inputs](#inputs)/[outputs](#outputs)
  - [`{NAME}`](#parameter)
    > one of...

    - [array](#array-parameter)
    - [dir](#dir-parameter)
    - [file](#file-parameter)
    - [number](#number-parameter)
    - [object](#object-parameter)
    - [socket](#socket-parameter)
    - [string](#string-parameter)
- [name](#name)
- [opspec](#opspec)
- [run](#call)
  > one of...

  - [container](#container)
    - [cmd](#container-cmd)
    - [dirs](#container-dirs)
    - [envVars](#container-env-vars)
    - [files](#container-files)
    - [image](#container-image)
      - [ref](#container-image-ref)
      - [pullCreds](#pull-creds)
    - [name](#container-name)
    - [ports](#container-ports)
    - [workDir](#container-workdir)
  - [op](#op-call)
    - [ref](#op-call-ref)
    - [pullCreds](#pull-creds)
    - [inputs](#op-call-inputs)
    - [outputs](#op-call-outputs)
  - [parallel](#parallel)
  - [parallelLoop](#parallel-loop)
    - [range](#loop-range)
    - [run](#loop-range)
    - [vars](#loop-vars)
  - [serial](#serial)
  - [serialLoop](#serial-loop)
    - [range](#loop-range)
    - [run](#loop-run)
    - [until](#loop-until)
    - [vars](#loop-vars)
  > any of...
  - [if](#if)
- [version](#version)

## Concepts

### Array
Arrays hold a list of values, i.e. they are a container for ordered values.

Arrays...
- are immutable, i.e. assigning to an array results in a copy of the original array
- can be passed in/out of ops via [array parameters](#array-parameter)
- can be initialized via [array initializers](#array-initializer)
- are coerced according to [array coercion](#array-coercion)

### Dir
Dirs hold a filesystem directory entry.

Dirs...
- are mutable, i.e. making changes to a directory results in the directory being changed everywhere it's referenced
- can be passed in/out of ops via [dir parameters](#dir-parameter)

### File
Files hold a filesystem file entry.

Files...
- are mutable, i.e. making changes to a file results in the file being changed everywhere it's referenced
- can be passed in/out of ops via [file parameters](#file-parameter)
- can be initialized via [file initializers](#file-initializer)
- are coerced according to [file coercion](#file-coercion)

### Number
Numbers hold a numeric value.

Numbers...
- are immutable, i.e. assigning to an number results in a copy of the original number
- can be passed in/out of ops via [number parameters](#number-parameter)
- can be initialized via [number initializers](#number-initializer)
- are coerced according to [number coercion](#number-coercion)

### Object
Objects hold a map of values, i.e. they are a container for unordered key value pairs.

Objects...
- are immutable, i.e. assigning to an object results in a copy of the original object
- can be passed in/out of ops via [object parameters](#object-parameter)
- can be initialized via [object initializers](#object-initializer)
- are coerced according to [object coercion](#object-coercion)

### Socket
Sockets hold one endpoint of a two way communication i.e. network/unix sockets, named-pipes etc...

Sockets...
- are immutable, i.e. assigning to an socket results in a copy of the original socket
- can be passed in/out of ops via [socket parameters](#socket-parameter)

### Reference
References allow reading or assigning values. 

References start with a reference opener `$(` and end with a reference closer `)`.

Within the reference opener and closer:
- there MUST be a [reference root](#reference-root)
- there MAY be a [reference path](#reference-path) IF the [reference root](#reference-root) is an [array](#array), [dir](#dir), [file](#file), or [object](#object)

### Reference Root
Root of the value being referenced.

one of:
- [embedded value reference](#embedded-value-reference)
- [scope value reference](#scope-value-reference)

### Reference Path
Path of the value being referenced (from [reference root](#reference-root)).
one of:
- [array path](#array-path)
- [dir path](#dir-path)
- [object path](#object-path)

#### Embedded Value Reference
Embedded value references are an absolute file path from the root of the op dir, i.e. `/op.yml` references the current op.yml, `/` references the current op dir

#### Scope Value Reference
Scope value references are an identifier either being assigned, or previously assigned, to the current scope, i.e. `myInput` references an in scope

### Array Path
Array items can be referenced via `ROOT[index]` syntax, where `index` is the zero based index of the item. 
If `index` is negative, indexing will take place from the end of the array.

#### Examples

##### First item
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[0])
```

##### Last item
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[-1])
```

### Dir Path
Dir entries (subdirs &/or files of a dir) can be referenced via `ROOT/child/path` syntax.

#### Examples

##### Embedded json file
given:
- `/file1.json` is embedded in op

```yaml
$(/file1.json)
```

##### File of in scope directory
given:
- `someDir`
  - is in scope
  - is type coercible to dir
  - contains `file2.txt`

```yaml
$(someDir/file2.txt)
```

### Object Path
Object properties can be referenced via `ROOT.child.path` syntax.

#### Examples

##### From scope
given:
- `someObject`
  - is in scope
  - is type coercible to object
  - contains property `someProperty`

```yaml
$(someObject.someProperty)
```

### Initializer
Initializes a value

one of:
- [array initializer](#array-initializer)
- [number initializer](#number-initializer)
- [object initializer](#object-initializer)
- [string initializer](#string-initializer)

### Array Initializer

Arrays can be constructed from literal arrays &/or [references](#reference).

Expressions get replaced and [type coerced](#type-coercion) (as required) on evaluation i.e. interpolation is supported. 

#### Examples

##### Literal
```yaml
- item1
- item2
```

##### Interpolated
given:
- `/someDir/file2.txt` is embedded in op
- `someObject` 
  - is in scope
  - is type coercible to object
  - has property `someProperty`

```yaml
- string $(/someDir/file2.txt)
- $(someObject.someProperty)
- [ sub, array, 2]
```

### Number Initializer
Numbers can be constructed from literal numbers &/or [references](#reference).

Expressions get replaced and [type coerced](#type-coercion) (as required) on evaluation i.e. interpolation is supported. 

#### Examples

##### Literal
```yaml
2
```

##### Interpolated
given:
- `someNumber`
  - is in scope
   - is type coercible to number

```yaml
# $(someNumber) replaced w/ someNumber
222$(someNumber)3e10
```

### Object Initializer
Objects can be constructed from literal objects &/or [references](#reference).

Expressions get replaced and [type coerced](#type-coercion) (as required) on evaluation i.e. interpolation is supported.

ES2015 style shorthand property name syntax is supported.

#### Examples

##### Literal

```yaml
myObject:
    prop1: value
```

##### Interpolated
given:
- `/someDir/file2.txt` is embedded in op
- `prop2Name` is in scope
- `someObject`
  - is in scope
  - is type coercible to object
  - has property `someProperty`
- `prop4` is in scope

```yaml
# interpolate properties
myObject:
    prop1: string $(/someDir/file2.txt)
    $(prop2Name): $(someObject.someProperty)
    prop3: [ sub, array, 2]
    # Shorthand property name; equivalent to prop4: $(prop4)
    prop4:
```

### String Initializer
Strings can be constructed from literal text &/or [references](#reference).

Expressions get replaced and [type coerced](#type-coercion) (as required) on evaluation i.e. interpolation is supported. 

#### Examples

##### Literal

```yaml
i'm a string
```

##### Interpolated
given:
- someObject
  - is in scope
  - is object

```yaml
# $(someObject) replaced w/ JSON representation of someObject
# $(dir/file.txt) replaced w/ contents of file.txt
pre $(someObject) $(dir/file.txt)
```

### Markdown

[Commonmark markdown](http://commonmark.org/) w/ table extensions can be used in descriptions.

> relative &/or absolute paths will be resolved from the root of the op

#### Examples

```markdown
checkout [this op's op.yml](op.yml)
checkout this image ![my image](/my-image.png)
# h1
## h2
### h3
#### h4
##### h5
###### h6
**bolded**
*italicized*
~~striken~~
- unordered item
  - unordered subitem
1. ordered item
  1. ordered subitem
| title 1 | title 2 | title 3 |
|:-------:|:-------:|:-------:|
| entry 1 | entry 2 | entry 3 |
```

### Pull Creds
Credentials object used to pull data (such as an image or op) from a data source (such as a git repo).

Properties:
- must have
  - [username](#string-coercible)
  - [password](#string-coercible)

### Inputs
An object defining the inputs of an operation.

For each property:
- key is a [parameter name](#parameter-name) string defining the name of the input
- value is a [parameter](#parameter) object defining the input. 

### Outputs
An object defining the outputs of an operation.

For each property:
- key is a [parameter name](#parameter-name) string defining the name of the output.
- value is a [parameter](#parameter) object defining the output. 

### Parameter Name
A string defining the name of an input or output parameter. Must match pattern `[-_.a-zA-Z0-9]+`.

### Parameter
An object defining a value that is passed in or out of an operation.

Properties:
- must have exactly one of
  - [array](#array-parameter)
  - [dir](#dir-parameter)
  - [file](#file-parameter)
  - [number](#number-parameter)
  - [object](#object-parameter)
  - [socket](#socket-parameter)
  - [string](#string-parameter)

### Array Parameter
An object defining an array typed parameter.

Properties:
- must have:
  - [description](#description)
- may have:
  - [default](#array-parameter-default)
  - [isSecret](#parameter-issecret)

### Dir Parameter
An object defining a dir typed parameter.

Properties:
- must have:
  - [description](#description)
- may have:
  - [default](#dir-parameter-default)
  - [isSecret](#parameter-issecret)

### Relative Path
A string defining a relative file path.

### Dir/File Parameter Default
A [relative](#relative-path) or [absolute](#absolute-path) path string to use as the default value of the parameter.

If the value is...
- an [absolute path](#absolute-path), the value is interpreted from the root of the op.
- a [relative path](#relative-path), the value is interpreted from the current working directory at the time the op is called.
  > relative path defaults are ignored when an op is called from an op as there is no current working directory.

### File Parameter
An object defining a file typed parameter.

Properties:
- must have
  - [description](#description)
- may have
  - [default](#file-parameter-default)
  - [isSecret](#parameter-issecret)

### Number Parameter
An object defining a number typed parameter.

Properties:
- must have
  - [description](#description)
- may have
  - [default](#number-parameter-default)
  - [isSecret](#parameter-issecret)

### Number Parameter Default
A number to use as the value of the parameter when not provided.

### Object Parameter
An object defining an object typed parameter.

Properties:
- must have
  - [description](#description)
- may have
  - [default](#object-parameter-default)
  - [isSecret](#parameter-issecret)

### Object Parameter Default
An object to use as the value of the parameter when not provided.

### Socket Parameter
An object defining a socket typed parameter.

Properties:
- must have
  - [description](#description)

### String Parameter
An object defining a string typed parameter.

Properties:
- must have
  - [description](#description)
- may have
  - [default](#string-parameter-default)
  - [isSecret](#parameter-issecret)

### String Parameter Default
A string to use as the value of the parameter when not provided.

### Parameter IsSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

### Description
Description of an op or parameter. Must be [markdown](#markdown) (since v0.1.6).

```yaml
description: |
  this description leverages:
  - [commonmark](http://commonmark.org/)
  - [yaml multiline strings](http://yaml-multiline.info/)
```

### Name
Human friendly identifier for the op.

> It's considered good practice to make `name` unique by using domain
> &/or path based namespacing.

Ops MAY be network resolvable; therefore `name` MUST be a valid
[uri-reference](https://tools.ietf.org/html/rfc3986#section-4.1)

```yaml
name: `github.com/opspec-pkgs/jwt.encode`
```

### Opspec
A string defining the version of opspec used by the op. Must be a [semver](#semver).

### Call
An object describing a call that happens as part of the operation.

Properties:
- must have exactly one of
  - [container](#container)
  - [op](#op)
  - [parallel](#parallel)
  - [parallelLoop](#parallel-loop)
  - [serial](#serial)
  - [serialLoop](#serial-loop)
- may have
  - [if](#if)

### Container
An object describing a container call.

Properties:
- must have
  - [image](#container-image)
- may have
  - [cmd](#container-cmd)
  - [dirs](#container-dirs)
  - [envVars](#container-env-vars)
  - [files](#container-files)
  - [name](#container-name)
  - [ports](#container-ports)
  - [sockets](#container-sockets)
  - [workDir](#workdir)

#### Container Cmd
Array of [string coercible](#string-coercible) items defining which binary to call and it's arguments.

> defining cmd overrides any entrypoint and/or cmd defined by the image

#### Container Dirs
Object where each key is an [absolute path](#absolute-path) in the container and the value is one of:
|value|behavior|
|--|--|
|null|Shorthand for reference to dir path|
|[reference](#reference)|On call start, reference will be de-referenced at specified path. On call end, reference will be set to dir at specified path|

#### Container Env Vars
[Object](#object) or [reference](#reference) to an object, where each key and value of the object represent the name and value of an environment variable in the container.

> note, keys and values will be coerced to strings. Refer to [coercion](#coercion) docs for details.

##### Examples

```
name: envVars
inputs:
  customEnvVars:
    object:
      default: 
        key1: 1
        key2: two
  greeting:
    string:
      default: hello
run:
  serial:
    - container:
        image: { ref: alpine }
        envVars: $(customEnvVars)
        cmd: [env]
        # prints:
        # ...
        # key1=1
        # key2=two
    - container:
        image: { ref: alpine }
        envVars:
          key1: $(customEnvVars.key1)
          key2:
            key2.2: 2.2
          key3: [1,2,3]
          greeting:
        cmd: [env]
        # prints:
        # ...
        # key1=1
        # key2={"key2.2":2.2}
        # key3=[1,2,3]
        # greeting=hello
```

#### Container Files
Object where each key is an [absolute path](#absolute-path) in the container and the value is one of:
|value|behavior|
|--|--|
|null|Shorthand for reference to file path|
|[reference](#reference)|On call start, reference will be de-referenced and coerced to file at specified path. On call end, reference will be set to file at specified path|
|[initializer](#initializer)|Initializer will be initialized and coerced to file at specified path|

#### Container Ports
Object defining container ports exposed on the opctl host where:
- each key is a container port or range of ports (optionally including protocol) matching `[0-9]+(-[0-9]+)?(tcp|udp)`
- each value is a corresponding opctl host port or range of ports matching `[0-9]+(-[0-9]+)?`

#### Container Sockets
Object where each key is an [absolute path](#absolute-path) in the container and the value is a [reference](#reference) to a socket to be mounted. 

#### Absolute Path
String defining an absolute path. Must match pattern `^([a-zA-Z]:)?[-_.\\/a-zA-Z0-9]+$`

#### Container Image
Object defining a container image

Properties:
- must have
  - [ref](#container-image-ref)
- may have
  - [pullCreds](#pull-creds)

#### Container Image Ref
String referencing a container image.

#### Container Name
Name the container can be referenced by from other containers. Must be [string-coercible](#string-coercible)

#### Container WorkDir
Defines the [absolute path](#absolute-path) in the container from which [cmd](#container-cmd) is executed.
> defining workDir overrides any defined by the image

#### Examples

##### Kitchen sink

```yaml
name: kitchenSink
inputs:
  mySocket: 
    string: {}
  registryCreds:
    object:
      constraints:
        properties:
          username: { type: string }
          password: { type: string }
        required: [username, password]
run:
  container:
    dirs:
      /:
    envVars:
      MY_ENV_VAR: my value
    files:
      /op.yml:
      /hello.txt: hello!
    image:
      ref: customBase:1.0.0
      pullCreds:
        username: $(registryCreds.username)
        password: $(registryCreds.password)
    ports: { 80: 80 }
    sockets:
      /mySocket: $(mySocket)
    name: my-kitchen-sink
    workDir: /root
```

### Version
Version of the op. Must be a [semver](#semver).

### SemVer
Defines a [v2 semantic version](https://semver.org/spec/v2.0.0.html).

```text
1.0.0-alpha
```

### If
An array of [predicates](#predicate) which must be true for the call to take place.

### Predicate
Defines a condition which evaluates to true or false.

Properties
- must have exactly one of
  - [eq](#eq-predicate)
  - [exists](#exists-predicate)
  - [ne](#ne-predicate)
  - [notExists](#not-exists-predicate)

### Eq Predicate
An array defining a predicate, true when all items are equal.

Items:
- must be one of
  - [reference](#reference)
  - [initializer](#initializer)

### Exists Predicate
A [reference](#reference) string defining a predicate, true when the referenced value exists.

### Ne Predicate
An array defining a predicate, true when one or more items aren't equal.

Items:
- must be
  - [reference](#reference)
  - [initializer](#initializer)

### Not Exists Predicate
A [reference](#reference) string defining a predicate, true when the referenced value doesn't exist.

### Op Call
An object defining an operation call.

Properties
- must have
  - [ref](#op-call-ref)
  - [pullCreds](#pull-creds)
  - [inputs](#op-call-inputs)
  - [outputs](#op-call-outputs)

### Op Call Ref
A string referencing a local or remote operation.

Must be:
- a relative path to a local op
- a network git-repo#[semver tag](#semver)/path to a remote op

### Op Call Inputs
An object where each key is an input to the op and the value is one of:
|value|behavior|
|--|--|
|null|Shorthand for reference to input name|
|[reference](#reference)|Reference will be de-referenced and passed as the value|
|[initializer](#initializer)|Initializer will be initialized and passed as the value|

This is equivalent to passing arguments to a function. 

### Op Call Outputs
An object where each key is a variable to assign and the value is one of:
|value|behavior|
|--|--|
|null|Shorthand for reference to output name|
|[reference](#reference)|Output will be de-referenced and assigned as the value|

This is equivalent to consuming return values from a function.

### Parallel
Array defining [calls](#call) which happen in parallel (all at once without order).

Items:
- must be
  - [calls](#call)

### Parallel Loop
Object defining a call loop in which all iterations happen in parallel (all at once without order).

Properties:
- may have
  - [range](#loop-range)
  - [run](#loop-run)
  - [vars](#loop-vars)

### Serial
Array defining [calls](#call) which happen in serial (one after another in order).

Items:
- must be
  - [calls](#call)

### Serial Loop
Object defining a call loop in which each iteration happens in serial (one after another in order)

Properties:
- may have
  - [range](#loop-range)
  - [run](#loop-run)
  - [until](#loop-until)
  - [vars](#loop-vars)

### Loop Range
A [rangeable value](#rangeable-value) defining the loops range.

### Loop Run
A [call](#call) object defining what's run on each iteration of the loop

### Loop Until
An array of [predicates](#predicate) which must be true for the loop to exit.

### Loop Vars
An object defining variable names iteration info (index, key, value) is available through in each iteration.

Properties:
- may have
  - [index](#loop-index-var)
  - [key](#loop-key-var)
  - [value](#loop-value-var)

### Loop Index Var
String defining the name of a variable to set equal to the current loop index.

### Loop Key Var
String defining the name of a variable to set equal to the current loop key.

Behavior varies based on the [range](#loop-range) value:
|range|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item index|
|object|Variable set to current property name|

### Loop Value Var
String defining the name of a variable to set equal to the current loop value.

Behavior varies based on the [range](#loop-range) value:
|range value|behavior|
|--|--|
|null|Variable not set|
|array|Variable set to current item|
|object|Variable set to current property value|

### Rangeable Value
An [array](#array-coercible) or [object](#object-coercible) defining a range of values.

### String Coercible
An [initializer](#initializer) or [reference](#reference) which is [coercible](#type-coercion) to a string.

### Array Coercible
An [initializer](#initializer) or [reference](#reference) which is [coercible](#type-coercion) to an array.

### Object Coercible
An [initializer](#initializer) or [reference](#reference) which is [coercible](#type-coercion) to an object.

### Type Coercion
Type coercion takes place automatically when necessary/possible.

one of:
- [array coercion](#array-coercion)
- [file coercion](#file-coercion)
- [number coercion](#number-coercion)
- [object coercion](#object-coercion)
- [string coercion](#string-coercion)
- [number coercion](#number-coercion)

#### Array coercion
Array typed values are coercible to:

- [file](#file) (will be serialized to JSON)
- [string](#string) (will be serialized to JSON)

#### File coercion
File typed values are coercible to:

- [array](#array) (if value of file is an array in JSON format)
- [number](#number) (if value of file is numeric)
- [object](#object) (if value of file is an object in JSON format)
- [string](#string)

#### Number coercion
Number typed values are coercible to:

- [file](#file)
- [string](#string)

#### Object coercion
Object typed values are coercible to:

- [file](#file) (will be serialized to JSON)
- [string](#string) (will be serialized to JSON)

#### String coercion
String typed values are coercible to:

- [file](#file)
- [number](#number) (if value of string is numeric)
- [object](#object) (if value of string is an object in JSON format)

#### Examples

##### Object to string
```yaml
name: objAsString
inputs:
  obj:
    object:
      default:
        prop1: prop1Value
        prop2: [ item1 ]
run:
  container:
    image: { ref: alpine }
    cmd:
    - echo
    - $(obj)
```

##### Number to file
```yaml
name: numAsFile
run:
  container:
    image: { ref: alpine }
    cmd:
    - sh
    - -ce
    - cat /numCoercedToFile
    files:
      /numCoercedToFile: 2.2
```
