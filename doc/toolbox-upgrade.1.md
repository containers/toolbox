
% toolbox-upgrade(1)

## NAME
toolbox-upgrade - Upgrade packages in Toolbx containers

## SYNOPSIS
**toolbox upgrade** [*--all* | *-a*] [*--container* | *-c* *CONTAINER*] [*CONTAINER*]

## DESCRIPTION
Upgrades packages inside one or more Toolbx containers. The container must have been created using `toolbox create` and must include a label specifying how to upgrade its packages.

This command will:
1. Read the container's metadata label `com.github.containers.toolbox.package-manager.update`
2. Run the specified upgrade command inside the container
3. Report any errors that occur during the process

The `--container` flag is optional when a positional container name is given.

## LABEL REQUIREMENT
Each container **must** have the following OCI label set in its metadata:

`com.github.containers.toolbox.package-manager.update="COMMAND"`

This label defines the exact package upgrade command to run inside the container. For example:

`com.github.containers.toolbox.package-manager.update="sudo dnf --assumeyes update"`

This label is typically the responsibility of the **image publisher** and should be present at container creation.

## OPTIONS

**--all, -a**  
Upgrade all Toolbx containers. Cannot be used with *--container* or a positional argument.

**--container, -c** *CONTAINER*  
Upgrade a specific Toolbx container. Optional when a positional container name is provided.

## EXAMPLES

**Upgrade a specific container (positional):**

$ toolbox upgrade fedora-toolbox-38


**Upgrade a specific container (flag):**

$ toolbox upgrade --container fedora-toolbox-38


**Upgrade all containers:**

$ toolbox upgrade --all


## NOTES

- This command doesn't perform package manager detection itself.
- It relies entirely on the container image to define the correct update mechanism.
- The `package-manager.update` label **must be set**; otherwise, the upgrade will fail.

## SEE ALSO
`toolbox(1)`, `toolbox-create(1)`, `toolbox-list(1)`
