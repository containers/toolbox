% toolbox-run(1)

## NAME
toolbox\-run - Run a command in an existing toolbox container

## SYNOPSIS
**toolbox run** [*--container NAME* | *-c NAME*]
            [*--distro DISTRO* | *-d DISTRO*]
            [*--no-tty* | *-T*] [*COMMAND*]
            [*--release RELEASE* | *-r RELEASE*] [*COMMAND*]

## DESCRIPTION

Runs a command inside an existing toolbox container. The container should have
been created using the `toolbox create` command.

A toolbox container is an OCI container. Therefore, `toolbox run` is analogous
to a `podman start` followed by a `podman exec`.

By default, the toolbox containers are tagged with the version of the OS that
corresponds to the content inside them. Their names are prefixed with the name
of the base image and suffixed with the current user name.

## OPTIONS ##

The following options are understood:

**--container** NAME, **-c** NAME

Run command inside a toolbox container with the given NAME. This is useful
when there are multiple toolbox containers created from the same base image,
or entirely customized containers created from custom-built base images.

**--distro** DISTRO, **-d** DISTRO

Run command inside a toolbox container for a different operating system DISTRO
than the host.

**--no-tty**, **-T**

Don't allocate pseudo-TTY.

**--release** RELEASE, **-r** RELEASE

Run command inside a toolbox container for a different operating system
RELEASE than the host.

## EXAMPLES

### Run ls inside a toolbox container using the default image matching the host OS

```
$ toolbox run ls -la
```

### Run emacs inside a toolbox container using the default image for Fedora 30

```
$ toolbox run --distro fedora --release f30 emacs
```

### Run uptime inside a custom toolbox container using a custom image

```
$ toolbox run --container foo uptime
```

## SEE ALSO

`buildah(1)`, `podman(1)`, `podman-exec(1)`, `podman-start(1)`
