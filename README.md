![README](data/gfx/README.gif)

[Toolbx](https://containertoolbx.org/) is a tool for Linux, which allows the
use of interactive command line environments for development and
troubleshooting the host operating system, without having to install software
on the host. It is built on top of [Podman](https://podman.io/) and other
standard container technologies from [OCI](https://opencontainers.org/).

Toolbx environments have seamless access to the user's home directory,
the Wayland and X11 sockets, networking (including Avahi), removable devices
(like USB sticks), systemd journal, SSH agent, D-Bus, ulimits, /dev and the
udev database, etc..

This is particularly useful on
[OSTree](https://ostreedev.github.io/ostree/) based operating systems like
[Fedora CoreOS](https://fedoraproject.org/coreos/) and
[Silverblue](https://fedoraproject.org/silverblue/). The intention of these
systems is to discourage installation of software on the host, and instead
install software as (or in) containers — they mostly don't even have package
managers like DNF or YUM. This makes it difficult to set up a development
environment or troubleshoot the operating system in the usual way.

Toolbx solves this problem by providing a fully mutable container within
which one can install their favourite development and troubleshooting tools,
editors and SDKs. For example, it's possible to do `yum install ansible`
without affecting the base operating system.

However, this tool doesn't *require* using an OSTree based system. It works
equally well on Fedora Workstation and Server, and that's a useful way to
incrementally adopt containerization.

The Toolbx environment is based on an [OCI](https://www.opencontainers.org/)
image. On Fedora this is the `fedora-toolbox` image. This image is used to
create a Toolbx container that offers the interactive command line
environment.

Note that Toolbx makes no promise about security beyond what's already
available in the usual command line environment on the host that everybody is
familiar with.


## Installation & Use

See our guides on
[installing & getting started](https://containertoolbx.org/install/) with
Toolbx and [Linux distro support](https://containertoolbx.org/distros/).


##

[![Star History Chart](https://api.star-history.com/svg?repos=containers/toolbox&type=Date)](https://star-history.com/#containers/toolbox&Date)


##

[![Zuul](https://zuul-ci.org/gated.svg)](https://softwarefactory-project.io/zuul/t/local/builds?project=containers/toolbox)
[![Daily Pipeline](https://softwarefactory-project.io/zuul/api/tenant/local/badge?project=containers/toolbox&pipeline=periodic)](https://softwarefactory-project.io/zuul/t/local/builds?project=containers%2Ftoolbox&pipeline=periodic)

[![Arch Linux package](https://img.shields.io/archlinux/v/extra/x86_64/toolbox?logo=archlinux)](https://www.archlinux.org/packages/extra/x86_64/toolbox/)
[![Fedora package](https://img.shields.io/fedora/v/toolbox/rawhide?logo=fedora)](https://src.fedoraproject.org/rpms/toolbox/)
[![Ubuntu package](https://img.shields.io/badge/ubuntu-0.0.99.3%2Bgit20230118%2B446d7bfdef6a-orange?logo=ubuntu)](https://packages.ubuntu.com/noble/podman-toolbox)
