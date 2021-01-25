% toolbox-enter(1)

## NAME
toolbox\-enter - Enter a toolbox container for interactive use

## SYNOPSIS
**toolbox enter** [*--distro DISTRO* | *-d DISTRO*]
              [*--release RELEASE* | *-r RELEASE*]
              [*CONTAINER*]

## DESCRIPTION

Spawns an interactive shell inside a toolbox container. The container should
have been created using the `toolbox create` command.

When invoked without any options, `toolbox enter` will try to enter the default
toolbox container for the host, or if there's only one container available then
it will use it. On Fedora, the default container is known as
`fedora-toolbox-N`, where N is the release of the host. If there aren't any
containers, `toolbox enter` will offer to create the default one for you.

A specific container can be selected using the CONTAINER argument.

A toolbox container is an OCI container. Therefore, `toolbox enter` is
analogous to a `podman start` followed by a `podman exec`.

## OPTIONS ##

The following options are understood:

**--distro** DISTRO, **-d** DISTRO

Enter a toolbox container for a different operating system DISTRO than the
host.

**--release** RELEASE, **-r** RELEASE

Enter a toolbox container for a different operating system RELEASE than the
host.

## EXAMPLES

### Enter a toolbox container using the default image matching the host OS

```
$ toolbox enter
```

### Enter a toolbox container using the default image for Fedora 30

```
$ toolbox enter --distro fedora --release f30
```

### Enter a custom toolbox container using a custom image

```
$ toolbox enter foo
```

## SEE ALSO

`buildah(1)`, `podman(1)`, `podman-exec(1)`, `podman-start(1)`
