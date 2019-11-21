% toolbox-init-container(1)

## NAME
toolbox\-init\-container - Initialize a running container

## SYNOPSIS
**toolbox init-container** *--home HOME*
                       *--home-link*
                       *--media-link*
                       *--monitor-host*
                       *--shell SHELL*
                       *--uid UID*
                       *--user USER*

## DESCRIPTION

Initializes a newly created container that's running. It is primarily meant to
be used as the entry point for all toolbox containers, and must be run inside
the container that's to be initialized. It is not expected to be directly
invoked by humans, and cannot be used on the host.

## OPTIONS ##

The following options are understood:

**--home** HOME

Create a user inside the toolbox container whose login directory is HOME.

**--home-link**

Make `/home` a symbolic link to `/var/home`.

**--media-link**

Make `/media` a symbolic link to `/run/media`.

**--monitor-host**

Ensure that certain configuration files inside the toolbox container are kept
synchronized with their counterparts on the host. Currently, these files are
`/etc/hosts` and `/etc/resolv.conf`.

**--shell** SHELL

Create a user inside the toolbox container whose login shell is SHELL.

**--uid** UID

Create a user inside the toolbox container whose numerical user ID is UID.

**--user** USER

Create a user inside the toolbox container whose login name is LOGIN.

## SEE ALSO

`podman(1)`, `podman-create(1)`, `podman-start(1)`
