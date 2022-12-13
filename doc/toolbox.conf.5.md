% toolbox.conf(5)

## NAME
toolbox.conf - Toolbox configuration file

## DESCRIPTION

Persistently overrides the default behaviour of `toolbox(1)`. The syntax is
TOML and the names of the options match their command line counterparts.
Currently, the only supported section is *general*.

## OPTIONS

**distro** = "DISTRO"

Create a toolbox container for a different operating system DISTRO than the
host. Cannot be used with `image`.

**image** = "NAME"

Change the NAME of the image used to create the toolbox container. This is
useful for creating containers from custom-built images. Cannot be used with
`distro` and `release`.

If NAME does not contain a registry, the local image storage will be
consulted, and if it's not present there then it will be pulled from a suitable
remote registry.

**release** = "RELEASE"

Create a toolbox container for a different operating system RELEASE than the
host. Cannot be used with `image`.

## FILES

The following locations are looked up in increasing order of priority:

**/etc/containers/toolbox.conf**

This is meant to be provided by the operating system distributor or the system
administrator, and affects all users on the host.

Fields specified here can be overridden by any of the files below.

**$XDG_CONFIG_HOME/containers/toolbox.conf**

This is meant for user-specific changes. Fields specified here override any of
the files above.

## EXAMPLES

### Override the default operating system distro:
```
[general]
distro = "fedora"
release = "36"
```

### Override the default image:
```
[general]
image = "registry.fedoraproject.org/fedora-toolbox:36"
```

## SEE ALSO

`toolbox(1)`, `toolbox-create(1)`
