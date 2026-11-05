% toolbox-rm 1

## NAME
toolbox\-rm - Remove one or more Toolbx containers

## SYNOPSIS
**toolbox rm** [*--all* | *-a*] [*--force* | *-f*] [*CONTAINER*...]

## DESCRIPTION

Removes one or more Toolbx containers from the host. The container should have
been created using the `toolbox create` command.

Before a container is removed, anything that was exported from it with
`toolbox export` is removed as well, the same as running
`toolbox unexport --all --container NAME` against it first. If a file
couldn't be removed for some reason, this is reported as an error, but the
container is still removed.

A Toolbx container is an OCI container. Therefore, `toolbox rm` can be used
interchangeably with `podman rm`, except for the removal of exported files,
which is a Toolbx concept `podman rm` doesn't know about.

## OPTIONS ##

The following options are understood:

**--all, -a**

Remove all Toolbx containers. It can be used in conjunction with `--force` as
well.

**--force, -f**

Force the removal of running and paused Toolbx containers.

## EXAMPLES

### Remove a Toolbx container named `fedora-toolbox-gegl:36`

```
$ toolbox rm fedora-toolbox-gegl:36
```

### Remove all Toolbx containers, but not those that are running or paused

```
$ toolbox rm --all
```

### Remove all Toolbx containers, including ones that are running or paused

```
$ toolbox rm --all --force
```

## SEE ALSO

`toolbox(1)`, `toolbox-export(1)`, `toolbox-unexport(1)`, `podman(1)`,
`podman-rm(1)`
