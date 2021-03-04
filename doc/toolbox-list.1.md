% toolbox-list(1)

## NAME
toolbox\-list - List existing toolbox containers and images

## SYNOPSIS
**toolbox list** [*--containers* | *-c*] [*--images* | *-i*]

## DESCRIPTION

Lists existing toolbox containers and images. These are OCI containers and
images, which can be managed directly with tools like `buildah` and `podman`.

## OPTIONS ##

The following options are understood:

**--containers, -c**

List only toolbox containers, not images.

**--images, -i**

List only toolbox images, not containers.

## EXAMPLES

### List all existing toolbox containers and images

```
$ toolbox list
```

### List existing toolbox containers only

```
$ toolbox list --containers
```

### List existing toolbox images only

```
$ toolbox list --images
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-ps(1)`, `podman-images(1)`
