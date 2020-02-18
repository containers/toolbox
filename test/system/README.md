# System tests

These tests are built with BATS (Bash Automated Testing System). They are
strongly influenced by the [libpod](https://github.com/containers/libpod)
project.

The tests are meant to ensure that Toolbox's functionality remains stable
throughout updates of both Toolbox and Podman/libpod.

## Structure

- **Basic Tests**
  - [] output version number (Toolbox + Podman)
  - [x] show help screen when no command is given
  - [x] create the default container
  - [x] create a container with a custom name
  - [x] create a container from a custom image
  - [x] list containers (no present)
  - [x] list default container and default image
  - [x] list containers (some present; different name patterns)
  - [x] list images (no present)
  - [x] list images (some present; different name patterns)
  - [x] remove a specific container
  - [x] try to remove nonexistent container
  - [x] try to remove a running container
  - [x] remove all containers
  - [x] try to remove all containers (running)
  - [x] remove a specific image
  - [x] remove all images
  - [x] run a command inside of an existing container

- **Advanced Tests**
  - [ ] create several containers with various configuration and then list them
  - [ ] create several containers and hop between them (series of enter/exit)
  - [ ] create a container, enter it, run a series of basic commands (id,
        whoami, dnf, top, systemctl,..)
  - [ ] enter a container and test basic set of networking tools (ping,
        traceroute,..)

The list of tests is stil rather basic. We **welcome** PRs with test
suggestions or even their implementation.

## Convention

- All tests that start with *Try to..* expect non-zero return value.

## How to run the tests

Make sure you have `bats` and `podman` with `toolbox` installed on your system.

- Enter the toolbox root folder
- Invoke command `bats ./test/system/` and the test suite should fire up

By default the test suite uses the system versions of `podman` and `toolbox`.

If you have a `podman` or `toolbox` installed in a nonstandard location then
you can use the `PODMAN` and `TOOLBOX` environmental variables to set the path
to the binaries. So the command to invoke the test suite could look something
like this: `PODMAN=/usr/libexec/podman TOOLBOX=./toolbox bats ./test/system/`.
