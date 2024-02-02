% toolbox-rmi 1

## NAME
toolbox\-rmi - Remove one or more Toolbx images

## SYNOPSIS
**toolbox rmi** [*--all* | *-a*] [*--force* | *-f*] [*IMAGE*...]

## DESCRIPTION

Removes one or more Toolbx images from the host. The image should have been
created using the `toolbox create` command.

A Toolbx image is an OCI image. Therefore, `toolbox rmi` can be used
interchangeably with `podman rmi`.

## OPTIONS ##

The following options are understood:

**--all, -a**

Remove all Toolbx images. It can be used in conjunction with `--force` as well.

**--force, -f**

Force the removal of Toolbx images that are used by Toolbx containers. The
dependent containers will be removed as well.

## EXAMPLES

### Remove a Toolbx image named `localhost/fedora-toolbox-gegl:36`

```
$ toolbox rmi localhost/fedora-toolbox-gegl:36
```

### Remove all Toolbx images, but not those that are used by containers

```
$ toolbox rmi --all
```

### Remove all Toolbx images and their dependent containers

```
$ toolbox rmi --all --force
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-rmi(1)`
