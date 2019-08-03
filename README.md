coreos-toolbox
===

This is a new implementation of https://github.com/debarshiray/toolbox/

Getting started
---

One time setup

```
$ coretoolbox create
<answer questions>
$
```

Now, each time you want to enter the toolbox:

```
$ coretoolbox run
```

One suggestion is to add a "profile" or configuration to your terminal
emulator that runs `coretoolbox run` by default, so that you can
easily create new tabs/windows in the toolbox.

Rationale
---

In order to disambiguate in this text we'll call this tool
"ctb", and the other one "dtb".

The main reason to introduce a new tool is that dtb too strongly
encourages true "pet" containers, where significant state is stored
inside.  We want to make it easy for people to build their own
toolbox "base images" derived from the upstream image.  For example,
rather than doing `yum install cargo` inside a toolbox container,
you use a `Dockerfile` that does:

```
FROM registry.fedoraproject.org/f30/fedora-toolbox:30
RUN yum -y install cargo
```

The `toolbox` command should ideally have at least a basic
concept of a "build" that regenerates the base container, but
at a minimum should support more easily specifying that base image.

A related problem with dtb is that it actually does create
a derived image locally with e.g. the username added; this
forces the image to be specific to one user or machine.

What "ctb" does instead is inject dynamic state (username, `HOME` path)
into the container at runtime.  This allows a lot more flexibility.

Today "dtb" has a hardcoded list of bind mounts for e.g. `HOME`
and the DBus system bus socket.
I ran into a case where I wanted e.g. the system libvirt socket.

In general, we aren't trying to confine `toolbox` - it's a privileged
container.  So "ctb" takes the approach of mounting in most
things from the host into the `/host` directory, and then uses
symlinks into `/host`.  This again makes everything a lot more
flexible as the set of things exposed can easily be changed
while the container is running.

Finally, ctb is written in a real programming language; bash
gets problematic once one goes beyond 10-20 lines
of code.

