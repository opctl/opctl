### travis.yml

```yaml
language: generic
sudo: required
before_script:
- curl -L https://bin.equinox.io/c/4fmGop7rntx/opctl-beta-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin
services:
- docker
script:
- opctl run build
```

### examples

[opctl projects travis](https://travis-ci.org/opspec-io/opctl)
