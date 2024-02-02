% toolbox-list 1

## NAME
toolbox\-list - List existing Toolbx containers and images

## SYNOPSIS
**toolbox list** [*--containers* | *-c*] [*--images* | *-i*]

## DESCRIPTION

Lists existing Toolbx containers and images. These are OCI containers and
images, which can be managed directly with a tool like `podman`.

## OPTIONS ##

The following options are understood:

**--containers, -c**

List only Toolbx containers, not images.

**--images, -i**

List only Toolbx images, not containers.

## EXAMPLES

### List all existing Toolbx containers and images

```
$ toolbox list
```

### List existing Toolbx containers only

```
$ toolbox list --containers
```

### List existing Toolbx images only

```
$ toolbox list --images
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-ps(1)`, `podman-images(1)`
