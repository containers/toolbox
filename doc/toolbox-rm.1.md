% toolbox-rm(1)

## NAME
toolbox\-rm - Remove one or more toolbox containers

## SYNOPSIS
**toolbox rm** [*--all* | *-a*] [*--force* | *-f*] [*CONTAINER*...]

## DESCRIPTION

Removes one or more toolbox containers from the host. The container should
have been created using the `toolbox create` command.

A toolbox container is an OCI container. Therefore, `toolbox rm` can be used
interchangeably with `podman rm`.

## OPTIONS ##

The following options are understood:

**--all, -a**

Remove all toolbox containers. It can be used in conjunction with `--force` as
well.

**--force, -f**

Force the removal of running and paused toolbox containers.

## EXAMPLES

### Remove a toolbox container named `fedora-toolbox-gegl:36`

```
$ toolbox rm fedora-toolbox-gegl:36
```

### Remove all toolbox containers, but not those that are running or paused

```
$ toolbox rm --all
```

### Remove all toolbox containers, including ones that are running or paused

```
$ toolbox rm --all --force
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-rm(1)`
