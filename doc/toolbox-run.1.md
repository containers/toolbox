% toolbox-run(1)

## NAME
toolbox\-run - Run a command in an existing toolbox container

## SYNOPSIS
**toolbox run** [*--container NAME* | *-c NAME*]
            [*--distro DISTRO* | *-d DISTRO*]
            [*--release RELEASE* | *-r RELEASE*] [*COMMAND*]

## DESCRIPTION

Runs a command inside an existing toolbox container. The container should have
been created using the `toolbox create` command.

`toolbox run` wraps around `podman exec` and by default passes several options
to it. It allocates a tty, connects to stdin, runs the passed command as the
current user in the current directory and shares common environmental
variables.

The executed command is wrapped in `capsh` that gets rid of all extra
capabilities that could negatively affect the experience.

A toolbox container is an OCI container. Therefore, `toolbox run` is analogous
to a `podman start` followed by a `podman exec`.

## OPTIONS ##

The following options are understood:

**--container** NAME, **-c** NAME

Run command inside a toolbox container with the given NAME. This is useful
when there are multiple toolbox containers created from the same base image,
or entirely customized containers created from custom-built base images.

**--distro** DISTRO, **-d** DISTRO

Run command inside a toolbox container for a different operating system DISTRO
than the host.

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

`toolbox(1)`, `podman(1)`, `podman-exec(1)`, `podman-start(1)`, `capsh(1)`,
`sh(1)`
