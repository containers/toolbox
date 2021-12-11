% toolbox-create(1)

## NAME
toolbox\-create - Create a new toolbox container

## SYNOPSIS
**toolbox create** [*--distro DISTRO* | *-d DISTRO*]
               [*--image NAME* | *-i NAME*]
               [*--release RELEASE* | *-r RELEASE*]
               [*CONTAINER*]

## DESCRIPTION

Creates a new toolbox container. You can then use the `toolbox enter` command
to interact with the container at any point.

A toolbox container is an OCI container created from an OCI image. On Fedora,
the default image is known as `fedora-toolbox:N`, where N is the release of
the host. If the image is not present locally, then it is pulled from a
well-known registry like `registry.fedoraproject.org`. Other images may be
used on other host operating systems. If the host is not recognized, then the
Fedora image will be used.

The container is created with `podman create`, and its entry point is set to
`toolbox init-container`.

By default, a toolbox container is named after its corresponding image. If the
image had a tag, then the tag is included in the name of the container, but
it's separated by a hyphen, not a colon. A different name can be assigned by
using the CONTAINER argument.

### Container Configuration

A toolbox container seamlessly integrates with the rest of the operating
system by providing access to the user's home directory, the Wayland and X11
sockets, networking (including Avahi), removable devices (like USB sticks),
systemd journal, SSH agent, D-Bus, ulimits, /dev and the udev database, etc..

The user ID and account details from the host is propagated into the toolbox
container, SELinux label separation is disabled, and the host file system can
be accessed by the container at /run/host. The container has access to the
host's Kerberos credentials cache if it's configured to use KCM caches.

A toolbox container can be identified by the `com.github.containers.toolbox`
label or the `/run/.toolboxenv` file.

The entry point of a toolbox container is the `toolbox init-container` command
which plays a role in setting up the container, along with the options passed
to `podman create`.

### Entry Point

A key feature of toolbox containers is their entry point, the `toolbox
init-container` command.

OCI containers are inherently immutable. Configuration options passed through
`podman create` are baked into the definition of the OCI container, and can't
be changed later. This means that changes and improvements made in newer
versions of Toolbox can't be applied to pre-existing toolbox containers
created by older versions of Toolbox. This is avoided by using the entry point
to configure the container at runtime.

The entry point of a toolbox container customizes the container to fit the
current user by ensuring that it has a user that matches the one on the host,
and grants it `sudo` and `root` access.

Crucial configuration files, such as `/etc/host.conf`, `/etc/hosts`,
`/etc/localtime`, `/etc/resolv.conf` and `/etc/timezone`, inside the container
are kept synchronized with the host. The entry point also bind mounts various
subsets of the host's filesystem hierarchy to their corresponding locations
inside the container to provide seamless integration with the host. This
includes `/run/libvirt`, `/run/systemd/journal`, `/run/udev/data`,
`/var/lib/libvirt`, `/var/lib/systemd/coredump`, `/var/log/journal` and others.

On some host operating systems, important paths like `/home`, `/media` or
`/mnt` are symbolic links to other locations. The entry point ensures that
paths inside the container match those on the host, to avoid needless
confusion.

## OPTIONS ##

**--distro** DISTRO, **-d** DISTRO

Create a toolbox container for a different operating system DISTRO than the
host. Cannot be used with `--image`.

**--image** NAME, **-i** NAME

Change the NAME of the image used to create the toolbox container. This is
useful for creating containers from custom-built images. Cannot be used with
`--distro` and `--release`.

If NAME does not contain a registry, the local image storage will be
consulted, and if it's not present there then it will be pulled from a suitable
remote registry.

**--release** RELEASE, **-r** RELEASE

Create a toolbox container for a different operating system RELEASE than the
host. Cannot be used with `--image`.

**--color** COLOR, **-C** COLOR

Create a toolbox with different dot color that is displayed while you are inside toolbox. Possible values are: red, green, orange, blue, purple, white, grey

## EXAMPLES

### Create a toolbox container using the default image matching the host OS

```
$ toolbox create
```

### Create a toolbox container using the default image for Fedora 36

```
$ toolbox create --distro fedora --release f36
```

### Create a custom toolbox container from a custom image

```
$ toolbox create --image bar foo
```

## SEE ALSO

`toolbox(1)`, `toolbox-init-container(1)`, `podman(1)`, `podman-create(1)`
