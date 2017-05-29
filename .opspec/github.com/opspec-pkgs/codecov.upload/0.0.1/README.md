# problem statement
uploads a coverage report to codecov.io

# example usage

> note: in examples, VERSION represents a version of the codecov.upload pkg

## install

```shell
opctl pkg install github.com/opspec-pkgs/codecov.upload#VERSION
```

## run

```
opctl run github.com/opspec-pkgs/codecov.upload#VERSION
```

## compose

```yaml
run:
  op:
    pkg: { ref: github.com/opspec-pkgs/codecov.upload#VERSION }
    inputs: { gitBranch, gitCommit, token, report }
```
