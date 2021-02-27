---
title: Gitlab
sidebar_label: Gitlab
---

[gitlab](https://gitlab.io) ci looks for a `.gitlab-ci.yml` file at the root of each repo to identify ci "stages".

Their hosted agents support running the ci process within a docker container so running opctl is
just a matter of defining your `.gitlab-ci.yml` as follows:

- using the official [opctl docker image](https://hub.docker.com/r/opctl/opctl/) as `image`
- adding "stages" with your calls to opctl

### Examples

```yaml
# .gitlab-ci.yml
image: opctl/opctl:0.1.48-dind
stages:
  - build
  - deploy
build:
  stage: build
  script:
    # passes args to opctl from gitlab variables
    - export gitlabUsername="$CI_REGISTRY_USER"
    - export gitlabSecret="$CI_REGISTRY_PASSWORD"
    - opctl run build
deploy:
  stage: deploy
  only:
    - master
  script:
    - opctl run deploy
```
