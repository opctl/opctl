# Creating a package

> Tip: checkout
> [official opspec packages](https://github.com/opspec-pkgs) for
> lots of examples

## step1
create a dir named `hello-world`

### step2
add a file named `op.yml` to `hello-world` w/ contents:
```yaml
name: hello-world
run:
  container:
    image: { ref: alpine }
    cmd: [echo, hello world!]
```

