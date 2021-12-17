% toolbox-run(1)

## NAME
toolbox\-run - Run a command in an existing toolbox container

## SYNOPSIS
**toolbox run** [*--container NAME* | *-c NAME*]
            [*--distro DISTRO* | *-d DISTRO*]
            [*--release RELEASE* | *-r RELEASE*]
            [*COMMAND*]

## DESCRIPTION

Runs a command inside an existing toolbox container. The container should have
been created using the `toolbox create` command.

On Fedora, the default container is known as `fedora-toolbox-N`, where N is
the release of the host. A specific container can be selected using the
`--container` option.

A toolbox container is an OCI container. Therefore, `toolbox run` is analogous
to a `podman start` followed by a `podman exec`.

## OPTIONS ##

The following options are understood:

**--container** NAME, **-c** NAME

Run command inside a toolbox container with the given NAME. This is useful
when there are multiple toolbox containers created from the same image, or
entirely customized containers created from custom-built images.

**--distro** DISTRO, **-d** DISTRO

Run command inside a toolbox container for a different operating system DISTRO
than the host. Has to be coupled with `--release` unless the selected system
matches the host system.

**--release** RELEASE, **-r** RELEASE

Run command inside a toolbox container for a different operating system
RELEASE than the host.

## EXAMPLES

### Run ls inside a toolbox container using the default image matching the host OS

```
$ toolbox run ls -la
```

### Run emacs inside a toolbox container using the default image for Fedora 36

```
$ toolbox run --distro fedora --release f36 emacs
```

### Run uptime inside a custom toolbox container using a custom image

```
$ toolbox run --container foo uptime
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-exec(1)`, `podman-start(1)`
