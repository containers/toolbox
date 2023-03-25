% toolbox-run 1

## NAME
toolbox\-run - Run a command in an existing toolbox container

## SYNOPSIS
**toolbox run** [*--container NAME* | *-c NAME*]
            [*--distro DISTRO* | *-d DISTRO*]
            [*--preserve-fds N*]
            [*--release RELEASE* | *-r RELEASE*]
            [*COMMAND*]

## DESCRIPTION

Runs a command inside an existing toolbox container. The container should have
been created using the `toolbox create` command.

On Fedora, the default container is known as `fedora-toolbox-N`, where N is
the release of the host. A specific container can be selected using the
`--container` option.

A toolbox container is an OCI container. Therefore, `toolbox run` is analogous
to a `podman start` followed by a `podman exec`.

## OPTIONS ##

The following options are understood:

**--container** NAME, **-c** NAME

Run command inside a toolbox container with the given NAME. This is useful
when there are multiple toolbox containers created from the same image, or
entirely customized containers created from custom-built images.

**--distro** DISTRO, **-d** DISTRO

Run command inside a toolbox container for a different operating system DISTRO
than the host. Has to be coupled with `--release` unless the selected DISTRO
matches the host system.

**--preserve-fds** N

Pass down to command N additional file descriptors (in addition to 0, 1,
2). The total number of file descriptors will be 3+N.

**--release** RELEASE, **-r** RELEASE

Run command inside a toolbox container for a different operating system
RELEASE than the host.

## EXIT STATUS

The exit code gives information about why the command within the container
failed to run or why it exited.

**1** There was an internal error in Toolbox

**125** There was an internal error in Podman

**126** The run command could not be invoked

```
$ toolbox run /etc; echo $?
/bin/sh: line 1: /etc: Is a directory
/bin/sh: line 1: exec: /etc: cannot execute: Is a directory
Error: failed to invoke command /etc in container fedora-toolbox-36
126
```

**127** The run command cannot be found or the working directory does not exist

```
$ toolbox run foo; echo $?
/bin/sh: line 1: exec: foo: not found
Error: command foo not found in container fedora-toolbox-36
127
```

**Exit code** The run command exit code

```
$ toolbox run false; echo $?
1
```

## EXAMPLES

### Run ls inside the default toolbox container matching the host OS

```
$ toolbox run ls -la
```

### Run emacs inside the default toolbox container for Fedora 36

```
$ toolbox run --distro fedora --release f36 emacs
```

### Run uptime inside a toolbox container with a custom name

```
$ toolbox run --container foo uptime
```

## SEE ALSO

`toolbox(1)`, `podman(1)`, `podman-exec(1)`, `podman-start(1)`
