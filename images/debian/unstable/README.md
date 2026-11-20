<img src="data/logo/toolbox-logo-landscape.svg" alt="Toolbox logo landscape" width="800"/>

[Toolbox](https://github.com/debarshiray/toolbox) is a tool that offers a
familiar RPM based environment for developing and debugging software that runs
fully unprivileged using [Podman](https://podman.io/).

The toolbox container is a fully *mutable* container; when you see
`yum install ansible` for example, that's something you can do inside your
toolbox container, without affecting the base operating system.

This is particularly useful on
[OSTree](https://ostree.readthedocs.io/en/latest/) based Fedora systems like
[Silverblue](https://silverblue.fedoraproject.org/).  The intention of these
systems is to discourage installation of software on the host, and instead
install software as (or in) containers.

However, this tool doesn't *require* using an OSTree based system — it
works equally well if you're running e.g. existing Fedora Workstation or
Server, and that's a useful way to incrementally adopt containerization.

The toolbox environment is based on an [OCI](https://www.opencontainers.org/)
image. On Fedora this is the `fedora-toolbox` image. This image is used to
create a toolbox container that seamlessly integrates with the rest of the
operating system.

## Usage

### Create your toolbox container:
```
[user@hostname ~]$ toolbox create
Created container: fedora-toolbox-30
Enter with: toolbox enter
[user@hostname ~]$
```
This will create a container called `fedora-toolbox-<version-id>`.

### Enter the toolbox:
```
[user@hostname ~]$ toolbox enter
⬢[user@toolbox ~]$
```
