% toolbox-unexport 1

## NAME
toolbox\-unexport - Remove an application or a binary exported from a Toolbx container

## SYNOPSIS
**toolbox unexport** [*--all* | *--app APP* | *--bin BIN*] [*--container NAME*]

## DESCRIPTION

Removes what `toolbox export` made available on the host.

Only files that record the given container are removed. Desktop entries are
matched by their `X-Toolbx-Container` key, and shell scripts by a comment in
the same form. A file that was written by hand, or that was exported from a
different container, is left alone and reported as an error.

Removing an application also removes the icons that were copied along with it,
from `$XDG_DATA_HOME/toolbx/NAME/icons/APP-ID`.

`toolbox rm` runs the equivalent of `unexport --all` on a container before
removing it, so anything exported from that container is cleaned up
automatically. Running `unexport` by hand beforehand is no longer necessary,
but is still available if only some of the exported items should be removed.

## OPTIONS ##

The following options are understood:

**--all**

Remove everything that was exported from the container. Can't be used together
with `--app` or `--bin`.

**--app** APP

Remove the application with the given desktop entry. The trailing `.desktop`
is optional. Can't be used together with `--bin`.

**--bin** BIN

Remove the binary with the given name. Can't be used together with `--app`.

**--container** NAME, **-c** NAME

Remove what was exported from the Toolbx container with the given NAME. This
is useful when multiple containers are present.

## EXAMPLES

### Remove GIMP exported from the default Toolbx container

```
$ toolbox unexport --app gimp
```

### Remove Neovim exported from a container called fedora-toolbox-42

```
$ toolbox unexport --bin nvim --container fedora-toolbox-42
```

### Remove a container along with everything exported from it

```
$ toolbox rm fedora-toolbox-42
```

## SEE ALSO

`toolbox(1)`, `toolbox-export(1)`, `toolbox-rm(1)`
