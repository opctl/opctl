---
title: Travis
sidebar_label: Travis
---

[travis ci](https://travis-ci.org/) looks for a `travis.yml` file at the root of each repo to identify ci "stages".

Their hosted agents support starting the ci process alongside a docker daemon so running opctl is
a matter of defining your `travis.yml` as follows:

- adding opctl installation to the `before_script` array
- adding docker to the `services` array
- adding `script` array entries with your calls to opctl

### Examples

```yaml
# travis.yml
language: generic
sudo: required
before_script:
- curl -L https://github.com/opctl/opctl/releases/download/0.1.46/opctl0.1.46.linux.tgz | sudo tar -xzv -C /usr/local/bin
services:
- docker
script:
# passes an arg `gitBranch` to the op from a travis-ci variable
- opctl run -a gitBranch=$TRAVIS_BRANCH some_op
```
