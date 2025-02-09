---
title: Setup
sidebar_label: Setup
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs queryString="os">  
  <TabItem value="kubernetes" label="Kubernetes">
    Opctl can be deployed to Kubernetes in order to run ops across an arbitrary number of nodes.

    Example deployment yaml:

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
    name: opctl-in-kubernetes
    spec:
    replicas: 1
    template:
        spec:
        containers:
            - name: opctl
            image: ghcr.io/opctl/opctl:0.1.58-dind
            ports:
                # expose to other containers
                - name: http
                containerPort: 80
                protocol: TCP
            securityContext:
                privileged: true
    ```
  </TabItem>
  <TabItem value="docker" label="Docker">
    Opctl can be run in docker containers using either DinD (Docker in Docker) or DooD (Docker out of Docker).

    ### DinD

    The `-dind` variant uses DinD (Docker in Docker) as the container runtime and requires a `--privileged` flag.

    For example, to run github.com/opspec-pkgs/uuid.v4.generate#1.1.0 with DinD:

    ```sh
    docker run \
        --privileged \
        ghcr.io/opctl/opctl:0.1.58-dind \
        opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
    ```

    ### DooD

    The `-dood` variant uses DooD (Docker out of Docker) as the container runtime and requires:
    - `-v /var/run/docker.sock:/var/run/docker.sock`
    to mount the socket of the external docker daemon.
    - `-v opctl_data_dir:/root/opctl`
    to mount an external directory as opctl's data dir.

    For example, to run github.com/opspec-pkgs/uuid.v4.generate#1.1.0 with DooD:

    ```sh
    docker run \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v ~/opctl:/root/opctl \
        ghcr.io/opctl/opctl:0.1.58-dood \
        opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
    ```
  </TabItem>
  <TabItem value="github" label="Github">
    Opctl can be used from Github Actions to run ops.

    Example workflow yaml:

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
  </TabItem>
  <TabItem value="gitlab" label="Gitlab">
    Opctl can be used from Gitlab CI/CD to run ops.

    Example .gitlab-ci.yml:

    ```yaml
    image: ghcr.io/opctl/opctl:0.1.58-dind
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
        - main
    script:
        - opctl run deploy
    ```
  </TabItem>
  <TabItem value="vscode" label="VSCode">
    Intellisense for [opspec](reference/opspec/index.md) can be enabled in VSCode.

    1. install [vscode-yaml plugin](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
    2. add to your user or workspace settings
      
       ```json
       "yaml.schemas": {
         "https://raw.githubusercontent.com/opctl/opctl/main/opspec/opfile/jsonschema.json": "/op.yml"
       }
       ```

    3. edit or create an op.yml w/ your fancy intellisense.
  </TabItem>
  <TabItem value="linux" label="Linux">
    Opctl can be installed and run natively on linux.

    ## Dependencies

    Opctl currently supports multiple container runtimes. Dependencies vary based on whichever you choose:

    |container runtime|dependencies|
    |--|--|
    |docker|[docker](https://docs.docker.com/get-docker/)|
    |qemu (experimental)|[lima](https://github.com/lima-vm/lima/releases/latest)|

    ## Installation

    from terminal, run:

    ```sh
    curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-linux-amd64.tgz | sudo tar -xzv -C /usr/local/bin
    ```

    ## Upgrading

    from terminal, run:
    ```sh
    # update to latest release from https://github.com/opctl/opctl/releases
    opctl self-update
    ```
  </TabItem>
  <TabItem value="osx" label="OSX">
    Opctl can be installed and run natively on OSX.

    ## Dependencies

    Opctl currently supports multiple container runtimes. Dependencies vary based on whichever you choose:

    |container runtime|dependencies|
    |--|--|
    |docker|[docker](https://docs.docker.com/get-docker/)|
    |qemu (experimental)|[lima](https://github.com/lima-vm/lima/releases/latest)|

    ## Installation

    ### M1

    from terminal, run:

    ```sh
    curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-darwin-arm64.tgz | sudo tar -xzv -C /usr/local/bin
    ```

    ### Intel

    from terminal, run:

    ```sh
    curl -L https://github.com/opctl/opctl/releases/latest/download/opctl-darwin-amd64.tgz | sudo tar -xzv -C /usr/local/bin
    ```

    ## Upgrading

    from terminal, run:
    ```sh
    # update to latest release from https://github.com/opctl/opctl/releases
    opctl self-update
    ```
  </TabItem>
  <TabItem value="windows" label="Windows">
    Use [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/) and the [linux install](?os=linux).
  </TabItem>
</Tabs>