# Getting Started

### Understanding an 'op'

**Q.** What is an op?
**A.** An op is a containerized procedure or set of procedures that can accept many inputs and produce many outputs. A procedure can consistent of containerized commands run in serial or parallel. A procedure may also be a reference to another op. 

**Q.** Why opctl?
**A.** `opctl` cli is a tool that makes use of the opsec standards. These standards allow us to rationalize operations in an explicit way in order to create consistent implementations regardless of the tool being used. To that end, `opctl` is one of many possible ways to achieve the desired result and is not required.

**Q.** Why opspec?
**A.** Opspec was designed to make repetitive operations more easily defined and reusable. Opspec defines a way to describe operations in a modular and loosely coupled way. This language is vendor agnostic and can be used across any set of tools.

### Writing our first op

We are going to create a 'hello-world' op in this example...of course!

1. Create a folder in the location of your choice, called `hello-world`
2. In the `./hello-world` folder, create an `op.yml` file
```
name: HelloWorld
description: Just saying hello...
run:
  serial:
    - container:
        cmd:
          - echo
          - Hello, world!
        image:
          ref: alpine:3.6
```
3. Now run the cli: `opctl run .` (should be run from within the `hello-world` folder)

### Taking some input

We are going to create a 'hello-friend' op in this example...

1. Create a folder in the location of your choice, called `hello-friend`
2. In the `./hello-friend` folder, create an `op.yml` file
```
name: HelloFriend
description: Just saying hello to a friend...
inputs:
  friend:
    string:
      description: a friend's name
run:
  serial:
    - container:
        cmd:
          - echo
          - Hello, $(friend)!
        image:
          ref: alpine:3.6
```
3. Now run the cli: `opctl run .` (should be run from within the `hello-friend` folder)
	a. You should notice that you are prompted for an input from the cli
