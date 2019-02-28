% toolbox(1)

## NAME
toolbox - Unprivileged development environment

## SYNOPSIS
**toolbox** [*--verbose* | *-v*] *COMMAND* [*ARGS*]

## DESCRIPTION

Toolbox is a tool that offers a familiar RPM based environment for developing
and debugging software that runs fully unprivileged using Podman.

The toolbox container is a fully *mutable* container; when you see
`yum install ansible` for example, that's something you can do inside your
toolbox container, without affecting the base operating system.

This is particularly useful on OSTree based Fedora systems like Silverblue.
The intention of these systems is to discourage installation of software on
the host, and instead install software as (or in) containers.

However this tool doesn't *require* using an OSTree based system — it works
equally well if you're running e.g. existing Fedora Workstation or Server, and
that's a useful way to incrementally adopt containerization.

The toolbox environment is based on an OCI image. On Fedora this is the
`fedora-toolbox` image. This image is then customized for the current user to
create a toolbox container that seamlessly integrates with the rest of the
operating system.

## OPTIONS ##

The following options are understood:

**--help, -h**

Print a synopsis of this manual and exit.

**--verbose, -v**

Print debug information. This includes messages coming from the standard error
stream of internal commands.

## COMMANDS

Commands for working with toolbox containers and images:

**toolbox-create(1)**

Create a new toolbox container.

**toolbox-enter(1)**

Enter an existing toolbox container for interactive use.

**toolbox-list(1)**

List existing toolbox containers and images.

## SEE ALSO

`buildah(1)`, `podman(1)`
