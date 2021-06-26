% toolbox-create(1)

## NAME
toolbox\-create - Create a new toolbox container

## SYNOPSIS
**toolbox create** [*--distro DISTRO* | *-d DISTRO*]
               [*--image NAME* | *-i NAME*]
               [*--release RELEASE* | *-r RELEASE*]
               [*CONTAINER*]

## DESCRIPTION

Creates a new toolbox container. You can then use the `toolbox enter` command
to interact with the container at any point.

A toolbox container is an OCI container created from an OCI image. On Fedora,
the default image is known as `fedora-toolbox:N`, where N is the release of
the host. If the image is not present locally, then it is pulled from a
well-known registry like `registry.fedoraproject.org`. Other images may be
used on other host operating systems. If the host is not recognized, then the
Fedora image will be used.

The container is created with `podman create`, and its entry point is set to
`toolbox init-container`.

By default, toolbox containers are named after their corresponding images. If
the image had a tag, then the tag is included in the name of the container,
but it's separated by a hyphen, not a colon. A different name can be assigned
by using the CONTAINER argument.

Toolbox containers are primarily created in a way to be tightly integrated with
the host system. They are not meant to be secure.

### Entry Point

A key feature of toolbox containers is their entry point, the `toolbox
init-container` command.

Read more about the entry-point in `toolbox-init-container(1)`.

### Toolbox setup

`toolbox-create(1)` passes several options to `podman-create(1)` when creating
toolbox containers to provide the needed functionality. The options have the
following effects:

- Toolboxes share with the host system:
    - network stack, including dns
    - IPC (shared memory, semaphores, message queues,..)
    - PID namespace
    - ulimits
- Toolboxes have access to cherry-picked parts of host filesystem made
  available under /run/host/
- Toolboxes are privileged containers
- SELinux label separation is disabled for toolboxes
- Toolboxes use as their entry-point `toolbox-init-container(1)`

Despite being privileged, rootless containers cannot have more privileges than
the user that created them.

Thanks to these options, `toolbox-init-container(1)` can futher set up the
containers. Read more about the entry-point in `toolbox-init-container(1)`.

## OPTIONS ##

**--distro** DISTRO, **-d** DISTRO

Create a toolbox container for a different operating system DISTRO than the
host. Cannot be used with `--image`.

**--image** NAME, **-i** NAME

Change the NAME of the base image used to create the toolbox container. This
is useful for creating containers from custom-built base images. Cannot be used
used with `--release`.

If NAME does not contain a domain, the image will be pulled from
`registry.fedoraproject.org`.

**--release** RELEASE, **-r** RELEASE

Create a toolbox container for a different operating system RELEASE than the
host. Cannot be used with `--image`.

## EXAMPLES

### Create a toolbox container using the default image matching the host OS

```
$ toolbox create
```

### Create a toolbox container using the default image for Fedora 30

```
$ toolbox create --distro fedora --release f30
```

### Create a custom toolbox container from a custom image

```
$ toolbox create --image bar foo
```

## SEE ALSO

`toolbox(1)`, `toolbox-init-container(1)`, `podman(1)`, `podman-create(1)`
