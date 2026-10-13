% toolbox-upgrade(1)

## NAME
toolbox\-upgrade - Upgrade the packages inside one or more Toolbx containers

## SYNOPSIS
**toolbox upgrade** [*--all* | *-a*] [*CONTAINER*...]

## DESCRIPTION

Upgrades the packages inside one or more Toolbx containers. The containers
should have been created using the `toolbox create` command.

Toolbx doesn't try to detect the package manager inside a container. Instead,
the image says how its packages are meant to be upgraded through the labels
described below, and Toolbx runs that command inside the container. This way an
image that Toolbx has never heard of can support the command without changing,
or rebuilding, Toolbx itself.

The command being run is always shown before it's run, and, unless the
**--assumeyes** option was used, Toolbx asks for confirmation. Beyond that,
Toolbx doesn't interfere: it's the package manager that shows what's going to be
downloaded and installed, and asks for the final confirmation.

Note that upgrading the packages inside a container is not the same as moving a
container to a newer operating system release. The latter needs the container to
be recreated from a newer image.

## IMAGE REQUIREMENTS

Images that want to support this command must have the following labels:

**com.github.containers.toolbx.package-manager.update**

The command that upgrades the packages inside the container. It's run with
`/bin/sh -c` as the user invoking `toolbox`, and should use `sudo(8)` to gain
the necessary privileges. It's expected to be interactive, and to let the
package manager ask for confirmation before changing anything.

**com.github.containers.toolbx.package-manager.update.assumeyes**

The variant of the above command that doesn't ask any questions. It's only used
when `toolbox` is invoked with **--assumeyes**. Images that can't be upgraded
without user interaction should leave this label unset.

Example for Fedora:
```
LABEL com.github.containers.toolbx.package-manager.update="sudo dnf upgrade" \
      com.github.containers.toolbx.package-manager.update.assumeyes="sudo dnf --assumeyes upgrade"
```

Example for Ubuntu:
```
LABEL com.github.containers.toolbx.package-manager.update="sudo apt-get update && sudo apt-get upgrade" \
      com.github.containers.toolbx.package-manager.update.assumeyes="sudo sh -c 'apt-get update && DEBIAN_FRONTEND=noninteractive apt-get upgrade --assume-yes'"
```

Containers made from images without these labels, like the ones for Arch Linux,
don't support this command. Arch Linux isn't built around the assumption that
unattended upgrades work, so it deliberately opts out instead of advertising
something that can break.

## OPTIONS ##

The following options are understood:

**--all, -a**

Upgrade all Toolbx containers. Can't be used together with *CONTAINER*.
Containers that don't support the command are skipped with a warning.

The **--assumeyes**, or **-y**, option of `toolbox(1)` skips the confirmation
prompt and uses the unattended variant of the upgrade command. Containers whose
images don't offer one can't be upgraded this way.

## EXAMPLES

### Upgrade the packages inside a Toolbx container

```
$ toolbox upgrade fedora-toolbox-42
```

### Upgrade the packages inside two Toolbx containers

```
$ toolbox upgrade fedora-toolbox-42 ubuntu-toolbox-26.04
```

### Upgrade the packages inside all Toolbx containers, without any questions

```
$ toolbox --assumeyes upgrade --all
```

## SEE ALSO

`toolbox(1)`, `toolbox-create(1)`, `toolbox-list(1)`, `toolbox-run(1)`,
`podman(1)`
