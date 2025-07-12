% toolbox-unexport 1

## NAME
toolbox-unexport - Remove exported binaries and applications for a toolbox container

## SYNOPSIS
**toolbox unexport** --container _container_ [--bin _binary_] [--app _application_] [--all]

## DESCRIPTION
The **toolbox unexport** command removes exported binaries and/or desktop applications that were previously made available on the host from a specified toolbox container. This helps clean up your host from wrappers and desktop files created by the `toolbox export` command.

## OPTIONS

**--bin _binary_**
:   Remove the exported binary wrapper for the specified binary, for the given container.

**--app _application_**
:   Remove the exported desktop application for the specified app, for the given container.

**--all**
:   Remove all exported binaries and applications for the specified container.

**--container _container_**
:   The container whose exported binaries and applications should be removed.

## EXAMPLES

Remove the exported `vim` binary for container `arch`:
```
toolbox unexport --container arch --bin vim
```

Remove the exported `firefox` application for container `fedora`:
```
toolbox unexport --container fedora --app firefox
```

Remove all exported binaries and applications for container `arch`:
```
toolbox unexport --container arch --all
```

## FILES

Exported binaries are located in:
```
~/.local/bin/
```

Exported desktop files are located in:
```
~/.local/share/applications/
```

## SEE ALSO
toolbox(1), toolbox-export(1)

## AUTHORS
Toolbox contributors
