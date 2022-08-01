<img src="data/logo/toolbox-logo-landscape.svg" alt="Toolbox logo landscape" width="800"/>
Also doing changes in this file.
## Goals

### High Level Goals

- Provide a convenient command line interface to run containers using
  [Podman](https://podman.io/)
- Support for development, debugging and system management use cases
- Support for multiple distros
  - `toolbox` package in multiple distros
  - `toolbox` containers for multiple distros

### Non-goals

- Supporting multiple container runtimes. Toolbox will use Podman exclusively
- Adding significant features on top of Podman
  - Significant feature requests should be driven into Podman upstream
- To run containers that aren't tightly integrated with the host
  - Extremely sandboxed containers quickly become specific to the user

### Developer Use Cases

- I’m a developer hacking on source code and building/testing code
  - Most cases: user doesn't need root, rootless containers work fine
  - Some cases: user needs root for testing
- Desktop Development:
  - Developers need things like D-Bus, display, etc. to be forwarded into the
    toolbox container
- Headless Development:
  - Toolbox works properly in headless environments (no display, etc)
- Need development tools like GDB, strace, etc. to work

### Debugging and System Management Use Cases

- Inspecting host processes and the kernel
  - Typically need root access
  - Need bpftrace, strace on host processes to work
    - Ideally even do things like helping get kernel-debuginfo data for the
      host kernel
- Managing system services
  - `systemctl restart foo.service`
  - journalctl
- Managing updates to the host
  - rpm-ostree
  - dnf/yum (classic systems)

### Specific environments

- Fedora Silverblue
  - Silverblue comes with a subset of packages and discourages host software
    changes
    - Users need a toolbox container as a working environment
    - Future: use toolbox container by default when a user opens a shell
- Fedora CoreOS
  - Similar to Silverblue, but non-graphical and smaller package set
- RHEL CoreOS
  - Similar to Fedora CoreOS. Based on RHEL content and the underlying
    operating system for OpenShift
  - Need to [use default authfile on pull](https://github.com/coreos/toolbox/pull/58/commits/413f83f7240d3c31121b557bfd55e489fad24489)
  - Need to ensure compatibility with the rhel7/support-tools container
    - Currently not a toolbox image, opportunity for collaboration
  - Alignment with `oc debug node/` (OpenShift)
    - `oc debug node` opens a shell on a kubernetes node
    - Value in having a consistent environment for both Toolbox's debugging
      mode and `oc debug node`
