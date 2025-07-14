
% toolbox-upgrade 1

## NAME
toolbox\-upgrade - Upgrade packages in Toolbx containers

## SYNOPSIS
**toolbox upgrade** [*--all* | *-a*] [*--container* | *-c* *CONTAINER*]

## DESCRIPTION
Upgrades packages inside one or more Toolbx containers by automatically detecting
the package manager (dnf, apt, pacman, etc.) and running the appropriate upgrade
command. The container should have been created using the `toolbox create` command.

This command will:
1. Detect the available package manager in the container
2. Run the appropriate update/upgrade commands
3. Report any errors that occur during the process

## OPTIONS
**--all, -a**
Upgrade all Toolbx containers. Cannot be used with *--container*.

**--container, -c** *CONTAINER*
Upgrade a specific Toolbx container. Cannot be used with *--all*.

## EXAMPLES
**Upgrade packages in a specific container:**
```
$ toolbox upgrade --container fedora-toolbox-38
```

**Upgrade packages in all containers:**
```
$ toolbox upgrade --all
```

**Typical output:**
```
Detected package manager: dnf
Updating metadata...
Upgrading packages...
Complete!
```

## NOTES
Supported package managers:
- dnf (Fedora, RHEL)
- microdnf (Minimal Fedora/RHEL)
- yum (Older RHEL/Fedora)
- apt (Debian, Ubuntu)
- pacman (Arch Linux)
- xbps (Void Linux)
- zypper (openSUSE)
- apk (Alpine Linux)
- emerge (Gentoo)
- slackpkg (Slackware)
- swupd (Clear Linux)

## SEE ALSO
`toolbox(1)`, `toolbox-create(1)`, `toolbox-list(1)`
