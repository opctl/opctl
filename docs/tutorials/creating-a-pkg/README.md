# Creating a package

> Tip: checkout
> [official opspec packages](https://github.com/opspec-pkgs) for
> lots of examples

## 1. create a pkg dir

## 2. add an op.yml

example op.yml:

```yaml
name: hello-world
run:
  container:
    image: { ref: alpine }
    cmd: [echo, hello world!]
```
