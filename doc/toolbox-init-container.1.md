% toolbox-init-container 1

## NAME
toolbox\-init\-container - Initialize a running container

## SYNOPSIS
**toolbox init-container** *--gid GID*
                       *--home HOME*
                       *--home-link*
                       *--media-link*
                       *--mnt-link*
                       *--shell SHELL*
                       *--uid UID*
                       *--user USER*

## DESCRIPTION

Initializes a newly created container that's running. It is primarily meant to
be used as the entry point for all toolbox containers, and must be run inside
the container that's to be initialized. It is not expected to be directly
invoked by humans, and cannot be used on the host.

A key feature of toolbox containers is their entry point, the `toolbox
init-container` command.

OCI containers are inherently immutable. Configuration options passed through
`podman create` are baked into the definition of the OCI container, and can't
be changed later. This means that changes and improvements made in newer
versions of Toolbox can't be applied to pre-existing toolbox containers
created by older versions of Toolbox. This is avoided by using the entry point
to configure the container at runtime.

The entry point of a toolbox container customizes the container to fit the
current user by ensuring that it has a user that matches the one on the host,
and grants it `sudo` and `root` access.

Crucial configuration files, such as `/etc/host.conf`, `/etc/hosts`,
`/etc/localtime`, `/etc/resolv.conf` and `/etc/timezone`, inside the container
are kept synchronized with the host. The entry point also bind mounts various
subsets of the host's file system hierarchy to their corresponding locations
inside the container to provide seamless integration with the host. This
includes `/run/libvirt`, `/run/systemd/journal`, `/run/udev/data`,
`/var/lib/libvirt`, `/var/lib/systemd/coredump`, `/var/log/journal` and others.

On some host operating systems, important paths like `/home`, `/media` or
`/mnt` are symbolic links to other locations. The entry point ensures that
paths inside the container match those on the host, to avoid needless
confusion.

## OPTIONS ##

The following options are understood:

**--gid** GID

Pass GID as the user's numerical group ID from the host to the toolbox
container.

**--home** HOME

Create a user inside the toolbox container whose login directory is HOME. This
option is required.

**--home-link**

Make `/home` a symbolic link to `/var/home`.

**--media-link**

Make `/media` a symbolic link to `/run/media`.

**--mnt-link**

Make `/mnt` a symbolic link to `/var/mnt`.

**--monitor-host**

Deprecated, does nothing.

Crucial configuration files inside the toolbox container are always kept
synchronized with their counterparts on the host, and various subsets of the
host's file system hierarchy are always bind mounted to their corresponding
locations inside the toolbox container.

**--shell** SHELL

Create a user inside the toolbox container whose login shell is SHELL. This
option is required.

**--uid** UID

Create a user inside the toolbox container whose numerical user ID is UID. This
option is required.

**--user** USER

Create a user inside the toolbox container whose login name is LOGIN. This
option is required.

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-create(1)`, `podman-start(1)`
