coreos-toolbox
---

This is a new implementation of https://github.com/debarshiray/toolbox/

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

Another obvious thing is that ctb is written in a real programming
language; bash gets problematic once one goes beyond 10-20 lines
of code.

