---
title: Hello World
---

Let's create a simple op to run the essential programming task of greeting the world

1. Create a directory named `hello-world`
2. Inside `hello-world` create a file named `op.yml` with the below contents

```yaml
# you might want to match the name to the containing directory's
name: hello-world

# describe what your op does using markdown
description: echoes hello world, because we gotta start somewhere

run:

  # run a container
  container:

    # use image resolvable via reference alpine:3.6
    image: { ref: 'alpine:3.6' }

    # invoke echo w/ arg "hello world"
    cmd: [ echo, "hello world" ]
```

3. verify you're in the right directory by running:
```bash
ls
```
you should see `hello-world`

4. run your hello-world op:
```bash
$ opctl run hello-world
```

You will see the running log of the op, and the glorious "hello world" echoed right before the `ContainerExited` line