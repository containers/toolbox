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

**--assumeyes, -y**

Automatically answer yes for all questions.

**--help, -h**

Print a synopsis of this manual and exit.

**--verbose, -v**

Print debug information including standard error stream of internal commands.
Use `-vv` for more detail.

## COMMANDS

Commands for working with toolbox containers and images:

**toolbox-create(1)**

Create a new toolbox container.

**toolbox-enter(1)**

Enter a toolbox container for interactive use.

**toolbox-help(1)**

Display help information about Toolbox.

**toolbox-init-container(1)**

Initialize a running container.

**toolbox-list(1)**

List existing toolbox containers and images.

**toolbox-reset(1)**

Remove all local podman (and toolbox) state.

**toolbox-rm(1)**

Remove one or more toolbox containers.

**toolbox-rmi(1)**

Remove one or more toolbox images.

**toolbox-run(1)**

Run a command in an existing toolbox container.

## SEE ALSO

`buildah(1)`, `podman(1)`
