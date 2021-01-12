---
title: Azure Pipelines
sidebar_label: Azure Pipelines
---

The Azure Pipelines "Hosted Linux" agent queue runs dockerized linux agents. The docker socket from the host gets mounted at `/var/run/docker.sock` inside each agent container, so running opctl is a matter of:

- choosing "Hosted Linux Preview" for Agent queue
- adding a task which creates an opctl node w/ a custom data dir (see [create opctl node task](#create-opctl-node-task))
- adding tasks with your calls to opctl (see [opctl run task](#opctl-run-task))

### Examples

#### create opctl node task

Task: "Shell Script"

Type: "Inline"

Script:
```bash
curl -L https://github.com/opctl/opctl/releases/download/0.1.46/opctl0.1.46.linux.tgz | sudo tar -xzv -C /usr/local/bin

# manually create an opctl node.
# custom node data dir required because VSTS only makes build dir available to docker daemon
#
# the node will remain running throughout the build and get used by all tasks calling `opctl run ...`
nohup opctl node create --data-dir=.opctl &>/dev/null &
```

#### opctl run task

Task: "Shell Script"

Type: "Inline"

> disables color (because vsts doesn't interpret ansi escape codes), & passes an arg `someArg` to the op from a VSTS "Variable" `someArg`

Script: "opctl --nc run -a someArg=$(someArg) some_op"
