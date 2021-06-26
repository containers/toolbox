% toolbox(1)

## NAME
toolbox - Tool for containerized command line environments on Linux

## SYNOPSIS
**toolbox** [*--assumeyes* | *-y*]
        [*--help* | *-h*]
        [*--log-level LEVEL*]
        [*--log-podman*]
        [*--verbose* | *-v*]
        *COMMAND* [*ARGS*...]

## DESCRIPTION

Toolbox is a tool for Linux operating systems, which allows the use of
containerized command line environments. It is built on top of Podman and
other standard container technologies from OCI.

This is particularly useful on OSTree based operating systems like Fedora
CoreOS and Silverblue. The intention of these systems is to discourage
installation of software on the host, and instead install software as (or in)
containers — they mostly don't even have package managers like DNF or YUM.
This makes it difficult to set up a development environment or install tools
for debugging in the usual way.

Toolbox solves this problem by providing a fully mutable container within
which one can install their favourite development and debugging tools, editors
and SDKs. For example, it's possible to do `yum install ansible` without
affecting the base operating system.

However, this tool doesn't *require* using an OSTree based system. It works
equally well on Fedora Workstation and Server, and that's a useful way to
incrementally adopt containerization.

The toolbox environment is based on an OCI image. On Fedora this is the
`fedora-toolbox` image. This image is used to create a toolbox container that
seamlessly integrates with the rest of the operating system by providing
access to the user's home directory, the Wayland and X11 sockets, SSH agent,
etc..

## GLOBAL OPTIONS ##

The following options are understood:

**--assumeyes, -y**

Automatically answer yes for all questions.

**--help, -h**

Print a synopsis of this manual and exit.

**--log-level**=*level*

Log messages above specified level: debug, info, warn, error, fatal or panic
(default: error)

**--log-podman**

Show log messages of invocations of Podman based on the logging level specified
by option **log-level**.

**--verbose, -v**

Same as `--log-level=debug`. Use `-vv` to include `--log-podman`.

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

**toolbox-rm(1)**

Remove one or more toolbox containers.

**toolbox-rmi(1)**

Remove one or more toolbox images.

**toolbox-run(1)**

Run a command in an existing toolbox container.

## SEE ALSO

`podman(1)`, https://github.com/containers/toolbox
