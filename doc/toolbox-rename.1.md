% toolbox-rename(1)

## NAME
toolbox\-rename - Rename an existing container

## SYNOPSIS
**toolbox rename** [*CONTAINER*] [*NEWNAME*]

## DESCRIPTION

Renames an existing toolbox container. The old name will become available for
use by new containers. This command works with toolbox containers in all basic 
states: not running & running. When a running container is renamed, the
propagation of the change will not fully complete (e.g., name used in logs)
until the container is stopped.

A container with a running session (e.g., using `toolbox-run(1)` or
`toolbox-enter(1)`) can not be renamed.

## EXAMPLES

### Rename a toolbox container called fedora-toolbox-34 to fedora-toolbox-35

```
$ toolbox rename fedora-toolbox-34 fedora-toolbox-35
```

### Rename a toolbox container using ID

```
$ toolbox rename 3e85b13c705e new-toolbox-name
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-rename(1)`
