# How do I use opctl from Visual Studio Team Services?

[Visual studio team services](https://www.visualstudio.com/team-services/) (VSTS)
allows defining "build"'s and "tasks" from their UI

Agents from their "Hosted Linux" agent queue include a docker daemon so running opctl is
a matter of:

- choosing "Hosted Linux Preview" for Agent queue
- adding a task which installs opctl (see [install opctl task](#install-opctl-task))
- adding tasks with your calls to opctl (see [opctl run task](#opctl-run-task))

### Examples

#### Install opctl task

Task: "Shell Script"

Type: "Inline"

Script: "curl -L https://bin.equinox.io/c/46sMA3YsjZ6/opctl-stable-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin"

#### Opctl run task

Task: "Shell Script"

Type: "Inline"

> passes an arg `someArg` to the op from a VSTS "Variable" `someArg`

Script: "opctl run -a someArg=$(someArg) test"
