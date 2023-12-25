#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  cleanup_all
}

teardown() {
  cleanup_all
}

@test "rename: Rename a container 'bad-container' to 'good-container'" {
  create_container bad-container

  run $TOOLBOX rename bad-container good-container

  assert_success
  refute_output

  run $TOOLBOX list -c

  assert_success
  assert_line --index 1 --partial "good-container"
}

@test "rename: Rename a container using its ID" {
  local container_id

  create_container container

  container_id="$(get_container_id container)"

  run $TOOLBOX rename $container_id renamed-container

  assert_success
  refute_output

  run $TOOLBOX list -c

  assert_success
  assert_line --index 1 --partial "renamed-container"
}

@test "rename: Try to run command with no arguments" {
  run $TOOLBOX rename

  assert_failure
  assert_line --index 0 "Error: The 'rename' command takes two arguments"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "rename: Try to run command with 3 arguments" {
  run $TOOLBOX rename arg1 arg2 arg3

  assert_failure
  assert_line --index 0 "Error: The 'rename' command takes two arguments"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "rename: Try to rename a non-existent container" {
  run $TOOLBOX rename non-existent-container new-name

  assert_failure
  assert_line --index 0 "Error: Invalid argument for CONTAINER"
  assert_line --index 1 "Container non-existent-container is not a toolbox container"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "rename: Try to rename with invalid name" {
  run $TOOLBOX rename non-existent-container .bad-name

  assert_failure
  assert_line --index 0 "Error: Invalid argument for NEWNAME"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "rename: Try to rename a non-toolbox container" {
  local num_of_containers

  pull_distro_image busybox

  run $PODMAN create --name podman-container busybox

  assert_success
  num_of_containers=$(list_containers)
  assert [ $num_of_containers -eq 1 ]

  run $TOOLBOX rename podman-container toolbox-container

  assert_failure
  assert_line --index 0 "Error: Invalid argument for CONTAINER"
  assert_line --index 1 "Container podman-container is not a toolbox container"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}
