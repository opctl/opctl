---
title: Inputs & Outputs
---

## Inputs
Let's modify the previous simple op to take an input and use it

```yaml
name: hello-world

description: echoes hello followed by a name you provide

# we add the inputs section
inputs:
  person: # the name of this input is "person"
    description: who to greet # the description is "who to greet"
    string: # the type of this input is string
      constraints: { minLength: 1 } # it can have a minLength of 1

run:

  container:

    image: { ref: 'alpine:3.6' }
    envVars: { person: $(person) } # we dereference our input "person" and assign its value to an environment variable called "person" inside the container
    # invoke echo w/ arg "hello $person" - shell will substitute $person with the value of environment variable "person"
    cmd:
      - sh
      - -ce
      - echo hello $person
```
if you run that, you'll be prompted for the input

```bash
$ opctl run hello-world

-
  Please provide "person".
  Description: who to greet
-

```
if you type in "you", the container will run and echo out "hello you"

Now you may not want to be prompted for the input everytime you run the op. That's why there's several ways to accept input:

1. `-a` cli flag: explicitly pass args to op. eg: `-a NAME1=VALUE1 -a NAME2=VALUE2`
2. `--arg-file` cli flag: reads in a file of args as key=value, in yml format. eg: `--arg-file="./args.yml"`. This flag has a default value of `.opspec/args.yml` i.e. opctl will automatically check for an args file at `.opspec/args.yml`
3. Environment variables: If you define an environment variable with the same name as an input on the machine you're running opctl on, its value will be supplied as the input's value
4. `default` property: You can define a `default` property for each input, containing a value to assign if no other input method was invoked (cli args or args file)

Input sources are checked according to the following precedence:

* arg provided via -a option
* arg file (since v0.1.19)
* env var
* default
* prompt

## Outputs
Let's take that simple op with 1 input and have it provide an output to be used by another op

```yaml
name: hello-world

description: echoes hello followed by a name you provide

inputs:
  person:
    description: whom to greet
    string:
      constraints: { minLength: 1 }

# we add the outputs section
outputs:
  helloperson:
    description: a string of hello $(person)
    string: {}

run:
  container:
    files:
        /output.txt: $(helloperson) # we bind our output to a file that we will create during the container run
    image: { ref: 'alpine:3.6' }
    envVars: { person: $(person) } 
    cmd:
      - sh
      - -ce
      - |
        echo hello $person > /output.txt
```

We are now producing an output, let's reference it in another op:
1. Create a new directory and call it `caddy`
2. create `op.yml` in the `caddy` directory, with the below contents

```yaml
name: caddy

description: runs a simple caddy web server that serves a welcome text at http://localthost:8080/

inputs:
# we need an input of person to pass to the hello-world op when we run it as part of the caddy op
  person:
    description: name to greet with welcome text at root of web site
    string:
      constraints: { minLength: 1 }
run:
  serial:
    - op:
        ref: ../hello-world # here we reference the other op we wrote, hello-world
        inputs: { person } # we pass our input, person, as input to hello-world
        outputs: { helloperson } # we add hello-world's output (helloperson) to the scope of this op
    - container:
        files:
            /srv/index.html: $(helloperson) # we dereference helloperson and use its value to populate an index.html file at the default root directory of the caddy image
        image: { ref: 'abiosoft/caddy' }
        ports: { '2015' : '8080' } # caddy image listens to 2015 by default, we'd like to serve on 8080 in this example
```

and run it
```bash
$ opctl run -a person=you caddy
```
Now if you navigate to http://localhost:8080 in your browser or via curl, you should see the text "hello you"
As you make requests to that web server, you should see caddy's log in your terminal

The above is an example of how ops can reference other ops, and how they can be composed. Note also how we effortlessly and implicitly coerced `helloperson`'s value from a `string` into a `file` as we mounted `/srv/index.html` in the container

