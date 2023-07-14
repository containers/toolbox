% toolbox 1

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
access to the user's home directory, the Wayland and X11 sockets, networking
(including Avahi), removable devices (like USB sticks), systemd journal, SSH
agent, D-Bus, ulimits, /dev and the udev database, etc..

## Supported operating system distributions

By default, Toolbox tries to use an image matching the host operating system
distribution for creating containers. If the host is not supported, then it
falls back to a Fedora image. Supported host operating systems are:

* Fedora
* Red Hat Enterprise Linux >= 8.5
* CentOS Stream
* Arch Linux


However, it's possible to create containers for a different distribution
through the use of the `--distro` and `--release` options that are accepted by
the relevant commands, or their counterparts in the configuration file. The
`--distro` flag specifies the name of the distribution, and `--release`
specifies its version. Supported combinations are:

Distro |Release
-------|----------
fedora |\<release\> or f\<release\> eg., 36 or f36
rhel   |\<major\>.\<minor\> eg., 8.5

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
The following commands correspond to markdown files which will help you with execution.

| Command                                          | Description                                                                 |
| ------------------------------------------------ | --------------------------------------------------------------------------- |
| [toolbox-create(1)](toolbox-create.1.md)         | Create a new toolbox container                                              |
| [toolbox-enter(1)](toolbox-enter.1.md)           | Enter a toolbox container for interactive use                               |
| [toolbox-help(1)](toolbox-help.1.md)             | Display help information about Toolbox                                      |
| [toolbox-init-container(1)](toolbox-init-container.1.md) |Initialize a running container                                       |
| [toolbox-list(1)](toolbox-list.1.md)             | List existing toolbox containers and images                                 |
| [toolbox-rm(1)](toolbox-rm.1.md)                 | Remove one or more toolbox containers                                       |
| [toolbox-rmi(1)](toolbox-rmi.1.md)               | Remove one or more toolbox images                                           |
| [toolbox-run(1)](toolbox-run.1.md)               | Run a command in an existing toolbox container                              |

## FILES ##

**toolbox.conf(5)**

Toolbox configuration file.

## SEE ALSO

`podman(1)`, https://github.com/containers/toolbox
