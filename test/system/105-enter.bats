#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  _setup_environment
  cleanup_containers
}

teardown() {
  cleanup_containers
}

@test "enter: Try to enter the default container with no containers created" {
  run $TOOLBOX enter <<< "n"

  assert_success
  assert_line --index 0 "No toolbox containers found. Create now? [y/N] A container can be created later with the 'create' command."
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter the default container with more than 1 other containers present" {
  local default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_container first
  create_container second

  run $TOOLBOX enter

  assert_failure
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the '--container' option to select a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a specific container with no containers created " {
  run $TOOLBOX enter wrong-container <<< "n"

  assert_success
  assert_line --index 0 "No toolbox containers found. Create now? [y/N] A container can be created later with the 'create' command."
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a specific non-existent container with other containers present" {
  create_container other-container

  run $TOOLBOX enter wrong-container

  assert_failure
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the '--container' option to select a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBOX --assumeyes enter --distro "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on Fedora but with wrong version" {
  run $TOOLBOX enter -d fedora -r foobar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX enter --distro fedora --release -3

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on RHEL but with wrong version" {
  run $TOOLBOX enter --distro rhel --release 8

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX enter --distro rhel --release 8.2foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX enter --distro rhel --release -2.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on non-default distro without providing a version" {
  local distro="fedora"
  local system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run $TOOLBOX enter -d "$distro"

  assert_failure
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

# TODO: Write the test
@test "enter: Enter the default toolbox" {
  skip "Testing of entering toolboxes is not implemented"
}

# TODO: Write the test
@test "enter: Enter the default toolbox when only 1 non-default toolbox is present" {
  skip "Testing of entering toolboxes is not implemented"
}

# TODO: Write the test
@test "enter: Enter a specific toolbox" {
  skip "Testing of entering toolboxes is not implemented"
}
