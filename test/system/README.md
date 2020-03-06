# System tests

These tests are built with BATS (Bash Automated Testing System). They are
strongly influenced by the [libpod](https://github.com/containers/libpod)
project.

The tests are meant to ensure that Toolbox's functionality remains stable
throughout updates of both Toolbox and Podman/libpod.

## Structure

- **0xx (Info)**
  - Commands that are not dependent on the presence/number of containers or
    images. eg., version, help, etc..
- **1xx (Initialization)**
  - Commands (list, create) when Toolbox has not really been used, yet.
  - It tries to list an empty list, creates several containers (default one
    and several with custom names and images).
- **2xx (Usage)**
  - The created containers are used for the first time testing the
    initialization (CMD of the container).
  - Not all containers will be used because in the *Cleanup* phase we want to
    try removing containers in both running and not running states.
- **3xx (Cleanup)**
  - In this section the containers and images from the previous *phases* are
    removed.
  - There is a difference between removing running and not running containers.
    We need to check the right behaviour.

## Convention

- All tests that start with *Try to..* expect non-zero return value.

## How to run the tests

Make sure you have `bats` and `podman` with `toolbox` installed on your system.

**Important**
Before you start the tests, you need to have present two images: the default
`fedora-toolbox` image for your version of Fedora and the `fedora-toolbox:29`
image.

- Enter the toolbox root folder
- Invoke command `bats ./test/system/` and the test suite should fire up

By default the test suite uses the system versions of `podman` and `toolbox`.

If you have a `podman` or `toolbox` installed in a nonstandard location then
you can use the `PODMAN` and `TOOLBOX` environmental variables to set the path
to the binaries. So the command to invoke the test suite could look something
like this: `PODMAN=/usr/libexec/podman TOOLBOX=./toolbox bats ./test/system/`.
