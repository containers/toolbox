[Toolbox](https://containertoolbx.org/) is a tool for Linux, which allows the
use of interactive command line environments for development and
troubleshooting the host operating system, without having to install software
on the host. It is built on top of [Podman](https://podman.io/) and other
standard container technologies from [OCI](https://opencontainers.org/).

Toolbox environments have seamless access to the user's home directory,
the Wayland and X11 sockets, networking (including Avahi), removable devices
(like USB sticks), systemd journal, SSH agent, D-Bus, ulimits, /dev and the
udev database, etc..

This is particularly useful on
[OSTree](https://ostreedev.github.io/ostree/) based operating systems like
[Fedora CoreOS](https://coreos.fedoraproject.org/) and
[Silverblue](https://silverblue.fedoraproject.org/). The intention of these
systems is to discourage installation of software on the host, and instead
install software as (or in) containers — they mostly don't even have package
managers like DNF or YUM. This makes it difficult to set up a development
environment or troubleshoot the operating system in the usual way.

Toolbox solves this problem by providing a fully mutable container within
which one can install their favourite development and troubleshooting tools,
editors and SDKs. For example, it's possible to do `yum install ansible`
without affecting the base operating system.

However, this tool doesn't *require* using an OSTree based system. It works
equally well on Fedora Workstation and Server, and that's a useful way to
incrementally adopt containerization.

The toolbox environment is based on an [OCI](https://www.opencontainers.org/)
image. On Fedora this is the `fedora-toolbox` image. This image is used to
create a toolbox container that offers the interactive command line
environment.

Note that Toolbox makes no promise about security beyond what's already
available in the usual command line environment on the host that everybody is
familiar with.


## Installation & Use

See our guides on
[installing & getting started](https://containertoolbx.org/install/) with
Toolbox and [Linux distro support](https://containertoolbx.org/distros/).
