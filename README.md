<img src="data/logo/toolbox-logo-landscape.svg" alt="Toolbox logo landscape" width="800"/>

[Toolbox](https://github.com/containers/toolbox) is a tool that offers a
familiar package based environment for developing and debugging software that
runs fully unprivileged using [Podman](https://podman.io/).

The toolbox container is a fully *mutable* container; when you see
`yum install ansible` for example, that's something you can do inside your
toolbox container, without affecting the base operating system.

This is particularly useful on
[OSTree](https://ostree.readthedocs.io/en/latest/) based operating systems like
[CoreOS](https://coreos.fedoraproject.org/) and
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

## Distro support

By default, Toolbox creates the container using an
[OCI](https://www.opencontainers.org/) image called
`<ID>-toolbox:<VERSION-ID>`, where `<ID>` and `<VERSION-ID>` are taken from the
host's `/usr/lib/os-release`. For example, the default image on a Fedora 30
host would be `fedora-toolbox:30`.

This default can be overridden by the `--image` option in `toolbox create`,
but operating system distributors should provide an adequately configured
default image to ensure a smooth user experience.

## Image requirements

Toolbox customizes newly created containers in a certain way. This requires
certain tools and paths to be present and have certain characteristics inside
the OCI image.

Tools:
* `getent(1)`
* `id(1)`
* `ln(1)`
* `mkdir(1)`: for hosts where `/home` is a symbolic link to `/var/home`
* `passwd(1)`
* `readlink(1)`
* `rm(1)`
* `rmdir(1)`: for hosts where `/home` is a symbolic link to `/var/home`
* `sleep(1)`
* `test(1)`
* `touch(1)`
* `unlink(1)`
* `useradd(8)`

Paths:
* `/etc/host.conf`: optional, if present not a bind mount
* `/etc/hosts`: optional, if present not a bind mount
* `/etc/krb5.conf.d`: directory, not a bind mount
* `/etc/localtime`: optional, if present not a bind mount
* `/etc/resolv.conf`: optional, if present not a bind mount
* `/etc/timezone`: optional, if present not a bind mount

The image should have `sudo(8)` enabled for users belonging to either the
`sudo` or `wheel` groups, and the group itself should exist. File an
[issue](https://github.com/containers/toolbox/issues/new) if you really need
support for a different group. However, it's preferable to keep this list as
short as possible.

Since Toolbox only works with OCI images that fulfill certain requirements,
it will refuse images that aren't tagged with
`com.github.containers.toolbox="true"` and
`com.github.debarshiray.toolbox="true"` labels. These labels are meant to be
used by the maintainer of the image to indicate that they have read this
document and tested that the image works with Toolbox. You can use the
following snippet in a Dockerfile for this:
```
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true"
```
