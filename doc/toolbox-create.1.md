% toolbox-create(1)

## NAME
toolbox\-create - Create a new toolbox container

## SYNOPSIS
**toolbox create** [*--container NAME* | *-c NAME*]
               [*--distro DISTRO* | *-d DISTRO*]
               [*--hostname HOSTNAME*]
               [*--image NAME* | *-i NAME*]
               [*--release RELEASE* | *-r RELEASE*]

## DESCRIPTION

Creates a new toolbox container. You can then use the `toolbox enter` command
to interact with the container at any point.

A toolbox container is an OCI container created from an OCI image. On Fedora
the base image is known as `fedora-toolbox`. If the image is not present
locally, then it is pulled from a well-known registry like
`registry.fedoraproject.org`. The base image is locally customized for the
current user to create a second image, from which the container is finally
created.

Toolbox containers and images are tagged with the version of the OS that
corresponds to the content inside them. The user-specific images and the
toolbox containers are prefixed with the name of the base image and suffixed
with the current user name.

## OPTIONS ##

The following options are understood:

**--container** NAME, **-c** NAME

Assign a different NAME to the toolbox container. This is useful for creating
multiple toolbox containers from the same base image, or for entirely
customized containers from custom-built base images.

**--distro** DISTRO, **-d** DISTRO

Create a toolbox container for a different operating system DISTRO than the
host. Cannot be used with `--image`.

**--hostname** HOSTNAME

Create the toolbox container using the specified hostname (default: toolbox).

**--image** NAME, **-i** NAME

Change the NAME of the base image used to create the toolbox container. This
is useful for creating containers from custom-built base images. Cannot be used
used with `--release`.

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
$ toolbox create --container foo --image bar
```

## SEE ALSO

`buildah(1)`, `podman(1)`
