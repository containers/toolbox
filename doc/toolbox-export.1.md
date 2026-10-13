% toolbox-export 1

## NAME
toolbox\-export - Export an application or a binary from a Toolbx container

## SYNOPSIS
**toolbox export** [*--app APP* | *--bin BIN*] [*--container NAME*] [*--force*]

## DESCRIPTION

Makes a graphical application or a command line tool from inside a Toolbx
container available on the host.

Exporting an application writes a copy of its desktop entry to
`$XDG_DATA_HOME/applications`, with the `Exec` and `TryExec` keys rewritten to
go through `toolbox run`. The basename of the desktop entry is preserved,
because Wayland compositors match the `app_id` of a window against it to find
the corresponding icon and name. Renaming the entry would break that mapping.

Icons referenced by the entry are copied to
`$XDG_DATA_HOME/toolbx/NAME/icons/APP-ID` and referenced from the exported
entry by their absolute path. The icon themes on the host are left untouched.

Exporting a binary writes a small shell script to `$XDG_BIN_HOME`, which
executes the binary through `toolbox run`. Note that this directory is not
part of `PATH` on every operating system. A warning is shown if it isn't.

Applications and binaries are looked up inside the container itself, using the
container's own `XDG_DATA_HOME` and `XDG_DATA_DIRS`. The host doesn't read the
container's file system.

Exported files record the container they came from. Desktop entries get an
`X-Toolbx-Container` key, and shell scripts get a comment in the same form.
Files that lack such a marker, or that carry the name of a different
container, aren't overwritten unless the `--force` option is used, and are
never removed by `toolbox unexport`.

Applications are resolved by their desktop file ID with the trailing
`.desktop` being optional. If no exact match is found, an entry whose ID ends
with the given name is used, so that `--app gimp` finds `org.gimp.GIMP`.

## OPTIONS ##

The following options are understood:

**--app** APP

Export the application with the given desktop entry. Can't be used together
with `--bin`.

**--bin** BIN

Export the binary with the given name. The name is looked up in the
container's `PATH`. Can't be used together with `--app`.

**--container** NAME, **-c** NAME

Export from the Toolbx container with the given NAME. This is useful when
multiple containers are present.

**--force**

Overwrite an existing file even if it wasn't exported from the same container,
or wasn't exported by Toolbx at all.

## NOTES

Two keys of an exported desktop entry are adjusted, because they're meaningless
once the `Exec` line runs on the host:

`DBusActivatable` is set to `false`. Otherwise the launcher tries to activate a
bus name that only exists inside the container.

`Path` is commented out. It's the working directory that the launcher enters
before spawning the application, and it usually doesn't exist on the host.
The working directory inside the container is unaffected.

## EXAMPLES

### Export GIMP from the default Toolbx container

```
$ toolbox export --app gimp
```

### Export Neovim from a container called fedora-toolbox-42

```
$ toolbox export --bin nvim --container fedora-toolbox-42
```

### Replace an entry that was exported from a different container

```
$ toolbox export --app gimp --container arch-toolbox-latest --force
```

## SEE ALSO

`toolbox(1)`, `toolbox-run(1)`, `toolbox-unexport(1)`,
https://specifications.freedesktop.org/desktop-entry-spec/latest/
