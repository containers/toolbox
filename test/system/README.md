# System tests

These tests are built with BATS (Bash Automated Testing System).

The tests are meant to ensure that Toolbox's functionality remains stable
throughout updates of both Toolbox and Podman/libpod.

**Warning**: The tests are not executed in an isolated environment. Running the
tests will clear all podman state (delete all containers, images, etc).

## Dependencies

- `awk`
- `bats`
- `GNU coreutils`
- `podman`
- `skopeo`
- `toolbox`

These tests use a few standard libraries for `bats` which help with clarity
and consistency. In order to use it you need to download them to the `libs`
directory:

```
# Go to the Toolbox root folder
$ git clone https://github.com/bats-core/bats-assert test/system/libs/bats-assert
$ git clone https://github.com/bats-core/bats-support test/system/libs/bats-support
```

## Convention

- All tests should follow the nomenclature: `[command]: <test description>...`
- When the test is expected to fail, start the test description with "Try
  to..."
- When the test is to give a non obvious output, it should be put in parenthesis
  at the end of the title

Examples:

* `@test "create: Create the default container"`
* `@test "rm: Try to remove a non-existent container"`

- All the tests start with a clean system (no images or containers) to make sure
  that there are no dependencies between tests and they are really isolated. Use
  the `setup()` and `teardown()` functions for that purpose.

## How to run the tests

First, make sure you have all the dependencies installed.

- Enter the toolbox root folder
- Invoke command `bats ./test/system/` and the test suite should fire up

Mocking of images is done automatically to prevent potential networking issues
and to speed up the cases.

By default the test suite uses the system versions of `podman`, `skopeo` and
`toolbox`.

If you have a `podman`, `skopeo` or `toolbox` installed in a nonstandard
location then you can use the `PODMAN`, `SKOPEO` and `TOOLBOX` environmental
variables to set the path to the binaries. So the command to invoke the test
suite could look something like this: `PODMAN=/usr/libexec/podman TOOLBOX=./toolbox bats ./test/system/`.

When running the tests, make sure the `test suite: [job]` jobs are successful.
These jobs set up the whole environment and are a strict requirement for other
jobs to run correctly.
