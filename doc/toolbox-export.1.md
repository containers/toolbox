# toolbox-export(1)

## NAME
toolbox-export - Export binaries or applications from a toolbox container to your host

## SYNOPSIS
**toolbox export** [--bin _binary_] [--app _application_] --container _container_

## DESCRIPTION
The **toolbox export** command allows you to expose binaries or desktop applications from a toolbox container onto your host system. This is achieved by creating wrapper scripts for binaries in `~/.local/bin` and desktop files for applications in `~/.local/share/applications`. These exported items let you launch containerized tools seamlessly from your host environment.

## OPTIONS

**--bin _binary_**
:   Export a binary from the toolbox container. The argument can be a binary name or a path inside the container.

**--app _application_**
:   Export a desktop application from the toolbox container. This will search for an appropriate `.desktop` file inside the container and adapt it for host use.

**--container _container_**
:   Name of the toolbox container from which the binary or application should be exported.

## EXAMPLES

Export the `vim` binary from the container named `arch`:
```
toolbox export --bin vim --container arch
```

Export the `firefox` application from the container named `fedora`:
```
toolbox export --app firefox --container fedora
```

## FILES

Exported binaries are placed in:
```
~/.local/bin/
```

Exported desktop files are placed in:
```
~/.local/share/applications/
```

## SEE ALSO
toolbox(1), toolbox-unexport(1)

## AUTHORS
Toolbox contributors