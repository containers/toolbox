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

Toolbx is a tool for Linux operating systems, which allows the use of
containerized command line environments. It is built on top of Podman and
other standard container technologies from OCI.

This is particularly useful on OSTree based operating systems like Fedora
CoreOS and Silverblue. The intention of these systems is to discourage
installation of software on the host, and instead install software as (or in)
containers — they mostly don't even have package managers like DNF or YUM.
This makes it difficult to set up a development environment or install tools
for debugging in the usual way.

Toolbx solves this problem by providing a fully mutable container within
which one can install their favourite development and debugging tools, editors
and SDKs. For example, it's possible to do `yum install ansible` without
affecting the base operating system.

However, this tool doesn't *require* using an OSTree based system. It works
equally well on Fedora Workstation and Server, and that's a useful way to
incrementally adopt containerization.

The Toolbx environment is based on an OCI image. On Fedora this is the
`fedora-toolbox` image. This image is used to create a Toolbx container that
seamlessly integrates with the rest of the operating system by providing
access to the user's home directory, the Wayland and X11 sockets, networking
(including Avahi), removable devices (like USB sticks), systemd journal, SSH
agent, D-Bus, ulimits, /dev and the udev database, etc..

## Supported operating system distributions

By default, Toolbx tries to use an image matching the host operating system
distribution for creating containers. If the host is not supported, then it
falls back to a Fedora image. Supported host operating systems are:

* Arch Linux
* Fedora
* Red Hat Enterprise Linux >= 8.5
* Ubuntu

However, it's possible to create containers for a different distribution
through the use of the `--distro` and `--release` options that are accepted by
the relevant commands, or their counterparts in the configuration file. The
`--distro` flag specifies the name of the distribution, and `--release`
specifies its version. Supported combinations are:

Distro |Release
-------|----------
arch   |latest or rolling
fedora |\<release\> or f\<release\> eg., 36 or f36
rhel   |\<major\>.\<minor\> eg., 8.5
ubuntu |\<YY\>.\<MM\> eg., 22.04

## USAGE

### Create a Toolbx container:

```
[user@hostname ~]$ toolbox create
Image required to create toolbox container.
Download registry.fedoraproject.org/fedora-toolbox:36 (294.1MB)? [y/N]: y
Created container: fedora-toolbox-36
Enter with: toolbox enter
[user@hostname ~]$
```

### Enter the Toolbx container:

```
[user@hostname ~]$ toolbox enter
⬢[user@toolbox ~]$
```

### Remove the Toolbx container:

```
[user@hostname ~]$ toolbox rm fedora-toolbox-36
[user@hostname ~]$
```

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

Commands for working with Toolbx containers and images:

**toolbox-create(1)**

Create a new Toolbx container.

**toolbox-enter(1)**

Enter a Toolbx container for interactive use.

**toolbox-help(1)**

Display help information about Toolbx.

**toolbox-init-container(1)**

Initialize a running container.

**toolbox-list(1)**

List existing Toolbx containers and images.

**toolbox-rm(1)**

Remove one or more Toolbx containers.

**toolbox-rmi(1)**

Remove one or more Toolbx images.

**toolbox-run(1)**

Run a command in an existing Toolbx container.

## FILES ##

**toolbox.conf(5)**

Toolbx configuration file.

## SEE ALSO

`podman(1)`, https://github.com/containers/toolbox
