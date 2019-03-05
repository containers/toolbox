% toolbox-enter(1)

## NAME
toolbox\-enter - Enter an existing toolbox container for interactive use

## SYNOPSIS
**toolbox enter** [*--container NAME*] [*--release RELEASE*]

## DESCRIPTION

Spawns an interactive shell inside an existing toolbox container. The
container should have been created using the `toolbox create` command.

A toolbox container is an OCI container. On Fedora the toolbox containers are
tagged with the version of the OS that corresponds to the content inside them.
Their names are prefixed with the name of the base image and suffixed with the
current user name.

## OPTIONS ##

The following options are understood:

**--container** NAME

Enter a toolbox container with the given NAME. This is useful when there are
multiple toolbox containers created from the same base image, or entirely
customized containers created from custom-built base images.

**--release** RELEASE

Enter a toolbox container for a different operating system RELEASE than the
host.

## EXAMPLES

### Enter a toolbox container using the default image matching the host OS

```
$ toolbox enter
```

### Enter a toolbox container using the default image for Fedora 30

```
$ toolbox enter --release f30
```

### Enter a custom toolbox container using a custom image

```
$ toolbox enter --container foo
```

## SEE ALSO

`buildah(1)`, `podman(1)`
