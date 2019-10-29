% toolbox-reset(1)

## NAME
toolbox\-reset - Remove all local podman (and toolbox) state

## SYNOPSIS
**toolbox reset**

## DESCRIPTION

Removes all existing podman (and toolbox) containers, images and configuration.
This can be used to factory reset your local Podman and Toolbox installations
when something has gone irrecoverably wrong with the `podman(1)` and
`toolbox(1)` commands.

This command can only be used on the host, and not from within a toolbox
container, and is only expected to be used right after a fresh boot before any
other `podman(1)` or `toolbox(1)` commands have been invoked.

## EXAMPLES

### Reset a broken Podman and Toolbox installation

```
$ toolbox reset
```
