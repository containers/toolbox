<img src="data/logo/toolbox-logo-landscape.svg" alt="Toolbox logo landscape" width="800"/>

## Goals and Use Cases

### High Level Goals

- Provide a CLI convenience interface to run containers (via `podman`) easily
- Support for Developer and Debugging/Management use cases
- Support for multiple distros
    - toolbox package in multiple distros
    - toolbox containers for multiple distros

### Non-Goals - Anti Use Cases

- Supporting multiple container runtimes. `toolbox` will use `podman` exclusively
- Adding significant features on top of `podman`
	- Significant feature requests should be driven into `podman` upstream
- To run containers that aren't tightly integrated with the host
	- i.e. extremely sandboxed containers become specific to the user quickly

### Developer Use Cases

- I’m a developer hacking on source code and building/testing code
    - Most cases: user doesn't need root, rootless containers work fine
    - Some cases: user needs root for testing
- Desktop Development: 
    - developers need things like dbus, display, etc, to be forwarded into the toolbox
- Headless Development:
    - toolbox works properly in headless environments (no display, etc)
- Need development tools like gdb, strace, etc to work

### Debugging/System management Use Cases

- Inspecting Host Processes/Kernel
    - Typically need root access
    - Need bpftrace, strace on host processes to work
		- Ideally even do things like helping get kernel-debuginfo data for the host kernel
- Managing system services
    - systemctl restart foo.service
    - journalctl
- Managing updates to the host
    - rpm-ostree
    - dnf/yum (classic systems)

### Specific environments

- Fedora Silverblue
	- Silverblue comes with a subset of packages and discourages host software changes
		- Users need a toolbox container as a working environment
		- Future: use toolbox container by default when a user opens a shell
- Fedora CoreOS
	- Similar to silverblue, but non-graphical and smaller package set
- RHEL CoreOS
	- Similar to Fedora CoreOS. Based on RHEL content and the underlying OS for OpenShift
	- Need to [use default authfile on pull](https://github.com/coreos/toolbox/pull/58/commits/413f83f7240d3c31121b557bfd55e489fad24489)
    - Need to ensure compatibility with the rhel7/support-tools container 
		- currently not a toolbox image, opportunity for collaboration
	- Alignment with `oc debug node/` (OpenShift)
		- `oc debug node` opens a shell on a kubernetes node
		- Value in having a consistent environment for both `toolbox` in debugging mode and `oc debug node`
