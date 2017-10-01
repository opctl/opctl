# prerequisites

Opctl relies on API access to a docker daemon.

Docker has two offerings supporting windows, docker-machine and
docker4win.

Docker behaves quite differently depending on which offering you choose:

### docker-machine
> note: if you use docker-machine w/ a virtualbox driver & plan on running ops
> requiring symlinking (such as `npm install`), make sure to run as
> admin. virtualbox requires admin permissions to enable symlinks on
> windows

| pros                      | cons                             |
|:--------------------------|:---------------------------------|
| bind mounts more reliable | manual creation of docker VM     |
| -                         | manual export of docker env vars |
| -                         | manual start of docker VM        |

### docker4Win

| pros                           | cons                      |
|:-------------------------------|:--------------------------|
| auto creation of docker VM     | bind mounts less reliable |
| auto export of docker env vars | -                         |
| auto start of docker VM        | -                         |

In the end, docker4Win wins in convenience but its bind mount behavior
isn't as mature as docker-machine & can be finicky.

We therefore recommend using docker4Win over docker-machine unless you
run into issues.

# installation

download and run the
[windows installer](https://github.com/opctl/opctl/releases/download/0.1.21/opctl0.1.21.windows.msi)

