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

@test "enter: Try to enter a specifc non-existent container with other containers present" {
  create_container other-container

  run $TOOLBOX enter wrong-container

  assert_failure
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the '--container' option to select a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBOX -y enter -d "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  # Distro names are in a hashtable and thus the order can change
  assert_line --index 1 --regexp "Supported values are: (.?(fedora|rhel))+"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on Fedora but with wrong version" {
  run $TOOLBOX enter -d fedora -r foobar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "Supported values for distribution fedora are in format: <release>/f<release>"
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
  assert_line --index 0 "Error: release not found for non-default distribution $distro"
  assert [ ${#lines[@]} -eq 1 ]
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
