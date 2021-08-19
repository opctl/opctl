---
title: Github
sidebar_label: Github
---

[Github](https://github.com) looks at the root of each repo under `.github/workflows/...` to identify "workflows".

Their hosted agents support starting the ci process alongside a docker daemon so running opctl is
a matter of defining a "workflow" as follows:

- downloading & untar'ing opctl
- adding "steps" with your calls to opctl

### Examples

```yaml
# .github/workflows/your_workflow.yml
name: your_workflow_name

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v1

    - name: Install Opctl
      run: curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin
    
    - name: Build
      run: opctl run build
```
